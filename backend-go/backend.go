package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type UserList struct {
	List []string `json:"user_list"`
}

type ChannelMessageList struct {
	List []ChannelMessage `json:"message_list"`
}

type ChannelMessage struct {
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
	User      string `json:"username"`
}

type Message struct {
	Message string `json:"message"`
}

type ChannelList struct {
	List []string `json:"channel_list"`
}

type User struct {
	cookieVal string
}

type OnlineUsers struct {
	userMap map[User]time.Time
	lock    sync.Mutex
}

func NewUserMap(ttl time.Duration) (users *OnlineUsers) {
	users = &OnlineUsers{userMap: make(map[User]time.Time)}
	go func() {
		for now := range time.Tick(time.Second) {
			users.lock.Lock()
			for u, hb := range users.userMap {
				if now.Sub(hb) > ttl {
					delete(users.userMap, u)
					_ = updateUserStatus(u, false)
				}
			}
			users.lock.Unlock()
		}
	}()

	return
}

func (users *OnlineUsers) Put(user User, heartbeat time.Time) {
	users.lock.Lock()
	users.userMap[user] = heartbeat
	users.lock.Unlock()
}

var db *sql.DB
var users *OnlineUsers

func main() {
	var err error
	postgresUser := os.Getenv("POST_USER")
	postgresPass := os.Getenv("POST_PASS")
	postgresHost := os.Getenv("POST_HOST")
	postgresDB := os.Getenv("POST_DB")
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s", postgresUser, postgresPass, postgresHost, postgresDB)
	// Connect to database
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	users = NewUserMap(5 * time.Second)

	ws := new(restful.WebService)

	ws.Route(ws.GET("/channel/{channelName}").
		Produces(restful.MIME_JSON).
		To(handleChannelGET))

	ws.Route(ws.POST("/channel/{channelName}").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		To(handleChannel))

	ws.Route(ws.GET("/users/{channelName}").
		Produces(restful.MIME_JSON).
		To(handleChannelUsers))

	ws.Route(ws.GET("/channels").
		Produces(restful.MIME_JSON).
		To(handleChannelList))

	ws.Route(ws.POST("/heartbeatz").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON).
		To(handleHeartbeatz))

	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedDomains: []string{},
		AllowedMethods: []string{"POST"},
		CookiesAllowed: true,
		Container:      restful.DefaultContainer}
	restful.DefaultContainer.Filter(cors.Filter)
	restful.Add(ws)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func getCookie(req *restful.Request, cookieName string) *http.Cookie {
	if req.Request == nil || req.Request.Cookies() == nil {
		return nil
	}
	cookies := req.Request.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == cookieName {
			return cookie
		}
	}

	return nil
}

func updateUserStatus(user User, online bool) error {
	_, err := db.Exec("UPDATE users SET online = $1 where cookie = $2", online, user.cookieVal)
	return err
}

func handleHeartbeatz(req *restful.Request, resp *restful.Response) {
	cookie := getCookie(req, "ChatUserAuth")
	if cookie == nil {
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookieVal := cookie.Value
	user := User{cookieVal: cookieVal}
	users.Put(user, time.Now())

	err := updateUserStatus(user, true)
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
	}

	return
}

func handleChannelList(req *restful.Request, resp *restful.Response) {
	cookie := getCookie(req, "ChatUserAuth")
	if cookie == nil {
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookieVal := cookie.Value
	user := User{cookieVal: cookieVal}
	res, err := db.Query("SELECT id FROM users where cookie=$1", user.cookieVal)
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Next()
	var userID string
	err = res.Scan(userID)
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = res.Close()
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err = db.Query("SELECT channel_name FROM listChannels WHERE id = (SELECT channelID from allowedChannel where userID = $1) OR public = true", userID)
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	var channel string
	var channels []string

	for res.Next() {
		err = res.Scan(channel)
		if err != nil {
			log.Printf("DB error %v", err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		channels = append(channels, channel)
	}

	channelList := ChannelList{List: channels}

	marshalled, err := json.Marshal(channelList)
	if err != nil {
		log.Printf("Could not marshall channel list")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = resp.Write(marshalled)
	if err != nil {
		log.Printf("Writing channel list to response error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}

func handleChannel(req *restful.Request, resp *restful.Response) {
	cookie := getCookie(req, "ChatUserAuth")
	if cookie == nil {
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookieVal := cookie.Value
	channelName := req.PathParameter("channelName")

	byteArr, err := io.ReadAll(req.Request.Body)
	if err != nil {
		log.Printf("Error on reading req body %v", req)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	var msg Message

	err = json.Unmarshal(byteArr, &msg)
	if err != nil {
		log.Printf("Could not unmarshall games request")
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	user := User{cookieVal: cookieVal}
	res, err := db.Query("SELECT id FROM users where cookie=$1", user.cookieVal)
	res.Next()
	var userID string
	err = res.Scan(userID)
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = res.Close()
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("INSERT INTO $1 (userID, msg) VALUES ($2, $3)", channelName, userID, msg.Message)
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}

func handleChannelUsers(req *restful.Request, resp *restful.Response) {
	cookie := getCookie(req, "ChatUserAuth")
	if cookie == nil {
		resp.WriteHeader(http.StatusUnauthorized)
		return
	}

	cookieVal := cookie.Value
	channelName := req.PathParameter("channelName")

	user := User{cookieVal: cookieVal}
	res, err := db.Query("SELECT id FROM users where cookie=$1", user.cookieVal)
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	res.Next()
	var userID string
	err = res.Scan(userID)
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	err = res.Close()
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err = db.Query("SELECT id FROM listChannels where channel_name=$1", channelName)
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Next()
	var channelID string
	err = res.Scan(channelName)
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = res.Close()
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err = db.Query("SELECT username FROM users WHERE online = true AND lastChannel = $1", channelID)
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	var username string
	var usersInChannel []string

	for res.Next() {
		err = res.Scan(username)
		if err != nil {
			log.Printf("DB error %v", err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		usersInChannel = append(usersInChannel, username)
	}

	userList := UserList{List: usersInChannel}

	marshalled, err := json.Marshal(userList)
	if err != nil {
		log.Printf("Could not marshall user list")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = resp.Write(marshalled)
	if err != nil {
		log.Printf("Writing user list to response error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}

func handleChannelGET(req *restful.Request, resp *restful.Response) {
	channelName := req.PathParameter("channelName")

	res, err := db.Query("SELECT username, extract(epoch from stamp), msg FROM $1 LIMIT 100", channelName)
	if err != nil {
		log.Printf("DB error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	var username string
	var unix int64
	var msg string
	var messages []ChannelMessage

	for res.Next() {
		err = res.Scan(username, unix, msg)
		if err != nil {
			log.Printf("DB error %v", err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		messages = append(messages, ChannelMessage{Message: msg, Timestamp: unix, User: username})
	}

	messageList := ChannelMessageList{List: messages}

	marshalled, err := json.Marshal(messageList)
	if err != nil {
		log.Printf("Could not marshall message list")
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = resp.Write(marshalled)
	if err != nil {
		log.Printf("Writing message list to response error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}
