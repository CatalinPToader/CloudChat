package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/emicklei/go-restful-openapi/v2"
	restfulspec "github.com/emicklei/go-restful-openapi/v2"
	"github.com/emicklei/go-restful/v3"
	"github.com/go-openapi/spec"
	_ "github.com/go-openapi/spec"
	_ "github.com/lib/pq"
	"io"
	"log"
	"net/http"
	"os"
	"path"
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
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		postgresHost, 5432, postgresUser, postgresPass, postgresDB)
	// Connect to database
	db, err = sql.Open("postgres", psqlInfo)
	if err != nil {
		log.Fatal(err)
	}

	users = NewUserMap(5 * time.Second)

	ws := new(restful.WebService)

	tags := []string{""}

	ws.Route(ws.GET("/channel/{channelName}").
		Produces(restful.MIME_JSON).
		Doc("Gets messages from channel").
		Param(ws.PathParameter("channelName", "name of the channel").DataType("string").DefaultValue("public_1")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(ChannelMessageList{}). // on the response
		Returns(200, "OK", ChannelMessageList{}).
		Returns(500, "Internal Server Error", nil).
		To(handleChannelGET))

	ws.Route(ws.POST("/channel/{channelName}").
		Consumes(restful.MIME_JSON).
		Doc("Posts messages from channel").
		Param(ws.PathParameter("channelName", "name of the channel").DataType("string").DefaultValue("public_1")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Reads(Message{}). // on the response
		Returns(200, "OK", nil).
		Returns(500, "Internal Server Error", nil).
		To(handleChannel))

	ws.Route(ws.GET("/users/{channelName}").
		Produces(restful.MIME_JSON).
		Doc("Gets online users from channel").
		Param(ws.PathParameter("channelName", "name of the channel").DataType("string").DefaultValue("public_1")).
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(UserList{}). // on the response
		Returns(200, "OK", UserList{}).
		Returns(500, "Internal Server Error", nil).
		To(handleChannelUsers))

	ws.Route(ws.GET("/channels").
		Produces(restful.MIME_JSON).
		Doc("Gets channels user has access to").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(ChannelList{}). // on the response
		Returns(200, "OK", ChannelList{}).
		Returns(500, "Internal Server Error", nil).
		To(handleChannelList))

	ws.Route(ws.POST("/heartbeatz").
		Produces(restful.MIME_JSON).
		Doc("Allows user to send a heartbeat").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Returns(200, "OK", nil).
		Returns(500, "Internal Server Error", nil).
		To(handleHeartbeatz))

	ws.Route(ws.GET("/chat/{subpath:*}").
		To(handleSubFile))
	ws.Route(ws.GET("/chat").
		To(handleFile))

	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedDomains: []string{},
		AllowedMethods: []string{"POST"},
		CookiesAllowed: true,
		Container:      restful.DefaultContainer}
	restful.DefaultContainer.Filter(cors.Filter)
	restful.Add(ws)

	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(), // you control what services are visible
		APIPath:                       "/swagger.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleSubFile(req *restful.Request, resp *restful.Response) {
	actual := path.Join(req.PathParameter("subpath"))
	fmt.Printf("serving %s ... (from %s)\n", actual, req.PathParameter("subpath"))
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		actual)
}

func handleFile(req *restful.Request, resp *restful.Response) {
	http.ServeFile(
		resp.ResponseWriter,
		req.Request,
		path.Join(req.QueryParameter("resource")))
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
	if online {
		_, err := db.Exec("UPDATE users SET online = $1 where cookie = $2", online, user.cookieVal)
		return err
	} else {
		_, err := db.Exec("UPDATE users SET online = $1, lastchannel=NULL where cookie = $2", online, user.cookieVal)
		return err
	}
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
		log.Printf("DB error update user status %v", err)
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

	userID, err := getUserID(user)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	var res *sql.Rows
	res, err = db.Query("SELECT channel_name FROM listChannels WHERE id in (SELECT channelID from allowedChannel where userID = $1) OR public = true", userID)
	if err != nil {
		log.Printf("DB error query channel list %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	var channel string
	var channels []string

	for res.Next() {
		err = res.Scan(&channel)
		if err != nil {
			log.Printf("DB error scan channel name %v", err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		channels = append(channels, channel)
	}

	channelList := ChannelList{List: channels}

	err = resp.WriteAsJson(channelList)
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
		log.Printf("Could not unmarshall message")
		resp.WriteHeader(http.StatusBadRequest)
		return
	}

	user := User{cookieVal: cookieVal}

	userID, err := getUsername(user)
	if err != nil {
		return
	}

	_, err = db.Exec(fmt.Sprintf("INSERT INTO %s (username, msg) VALUES ($1, $2)", channelName), userID, msg.Message)
	if err != nil {
		log.Printf("DB error insert user message %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	_, err = db.Exec("UPDATE users SET lastChannel = (SELECT id from listChannels where channel_name=$1) WHERE cookie=$2", channelName, cookieVal)
	if err != nil {
		log.Printf("DB error update user channel %v", err)
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

	_, err := getUserID(user)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := db.Query("SELECT id FROM listChannels where channel_name=$1", channelName)
	if err != nil {
		log.Printf("DB error query channel id %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}
	res.Next()
	var channelID string
	err = res.Scan(&channelID)
	if err != nil {
		log.Printf("DB error scan channel id %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = res.Close()
	if err != nil {
		log.Printf("DB error close %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err = db.Query("SELECT username FROM users WHERE online = true AND lastChannel = $1", channelID)
	if err != nil {
		log.Printf("DB error query online users %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	var username string
	var usersInChannel []string

	for res.Next() {
		err = res.Scan(&username)
		if err != nil {
			log.Printf("DB error scan username for online users %v", err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		usersInChannel = append(usersInChannel, username)
	}

	userList := UserList{List: usersInChannel}

	err = resp.WriteAsJson(userList)
	if err != nil {
		log.Printf("Writing user list to response error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}

func getUserID(user User) (string, error) {
	res, err := db.Query("SELECT id FROM users where cookie=$1", user.cookieVal)
	if err != nil {
		log.Printf("DB error query user id %v", err)
		return "", err
	}

	res.Next()
	var userID string
	err = res.Scan(&userID)
	if err != nil {
		log.Printf("DB error scan user id %v", err)
		return "", err
	}
	err = res.Close()
	if err != nil {
		log.Printf("DB error close %v", err)
		return "", err
	}
	return userID, nil
}

func getUsername(user User) (string, error) {
	res, err := db.Query("SELECT username FROM users where cookie=$1", user.cookieVal)
	if err != nil {
		log.Printf("DB error query username %v", err)
		return "", err
	}

	res.Next()
	var userID string
	err = res.Scan(&userID)
	if err != nil {
		log.Printf("DB error scan username %v", err)
		return "", err
	}
	err = res.Close()
	if err != nil {
		log.Printf("DB error close %v", err)
		return "", err
	}
	return userID, nil
}

func handleChannelGET(req *restful.Request, resp *restful.Response) {
	channelName := req.PathParameter("channelName")

	res, err := db.Query(fmt.Sprintf("SELECT username, extract(epoch from stamp)::bigint, msg FROM %s LIMIT 100", channelName))
	if err != nil {
		log.Printf("DB error query channel messages %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	var username string
	var unix int64
	var msg string
	var messages []ChannelMessage

	for res.Next() {
		err = res.Scan(&username, &unix, &msg)
		if err != nil {
			log.Printf("DB error scan user, stamp, message %v", err)
			resp.WriteHeader(http.StatusInternalServerError)
			return
		}

		messages = append(messages, ChannelMessage{Message: msg, Timestamp: unix, User: username})
	}

	messageList := ChannelMessageList{List: messages}

	err = resp.WriteAsJson(messageList)
	if err != nil {
		log.Printf("Writing message list to response error %v", err)
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	return
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "UserService",
			Description: "Resource for managing Users",
			Version:     "1.0.0",
		},
	}
	swo.Tags = []spec.Tag{{TagProps: spec.TagProps{
		Name:        "",
		Description: "Everything"}}}
}
