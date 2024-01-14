#Imports
import firebase_admin
import pyrebase
import json
import requests
import psycopg2
import os
import jwt
from firebase_admin import credentials, auth
from flask import Flask, request, render_template, make_response, url_for
from functools import wraps
import logging

#App configuration
app = Flask(__name__)

#Connect to firebase
cred = credentials.Certificate('fbAdminConfig_private_key.json')
firebase = firebase_admin.initialize_app(cred)
pb = pyrebase.initialize_app(json.load(open('firebaseConfig.json')))

postgres_user = os.getenv("POST_USER")
postgress_pass = os.getenv("POST_PASS")
postgres_host = os.getenv("POST_HOST")
postgres_db = os.getenv("POST_DB")

#Wrapper for token verification
def check_token(f):
    @wraps(f)
    def wrap(*args,**kwargs):
        if not request.headers.get('authorization'):
            return {'message': 'No auth token provided'}, 400
        try:
            user = auth.verify_id_token(request.headers['authorization'])
            request.user = user
        except:
            return {'message':'Invalid token provided.'}, 400
        return f(*args, **kwargs)
    return wrap
    
    
def check_db(username):
    conn = psycopg2.connect(
        host=postgres_host,
        database=postgres_db,
        user=postgres_user,
        password=postgress_pass
    )

    cur = conn.cursor()

    try:
        cur.execute(f"SELECT * from users where username='{username}'")
        if cur.fetchone() == None:
            return False, None
    except psycopg2.errors.DatabaseError as e:
        cur.close()
        conn.rollback()
        logging.error(f"Error at database: {e}")
        return None, e

    cur.close()
    return True, None  

def add_to_db(db_string):
    conn = psycopg2.connect(
        host=postgres_host,
        database=postgres_db,
        user=postgres_user,
        password=postgress_pass
    )

    cur = conn.cursor()

    try:
        cur.execute(db_string)
        conn.commit()
    except psycopg2.errors.DatabaseError as e:
        cur.close()
        conn.rollback()
        logging.error(f"Error at database: {e}")
        return e

    cur.close()
    return None

#Api route for signup
@app.route("/signup", methods=['GET', 'POST'])
def signup():
    if request.method == 'GET':
        return render_template('signup_page.html')
    else:
        email = request.form.get('email')
        password = request.form.get('password')
        username = request.form.get('username')

        if any(not field for field in [email, password, username]):
            return {'message': 'Incomplete request'}, 400

        if request.form.get('password') != request.form.get('password2'):
            return {'message': 'Passwords do not match'}, 400

        try:
            #Add user in firebase
            user = auth.create_user(
                email=email,
                password=password,
                display_name=username
            )
            
            logging.log(logging.INFO, f"username={username}")

            #Get token
            try:
                auth_resp = requests.get(url_for('token', _external=True), data={'email': email, 'password': password})
            except Exception as e:
                return {'message': f'There was an error signing up: {str(e)}'}, 400
    
            auth_json = auth_resp.json()
            if 'token' not in auth_json:
                return auth_json, auth_resp.status_code

            token = auth_json['token']

            token_split = token.split('.')
            cookieVal = token_split[1]
            
            #Set cookie
            resp = make_response({})
            resp.set_cookie('ChatUserAuth', cookieVal)

            e = add_to_db(f"INSERT INTO users (id, username, cookie) VALUES (gen_random_uuid(), '{username}', '{cookieVal}');")
            
            if e != None:
                return {'message': f'Error creating user: {str(e)}'}, 400

            return {'message': f'Successfully created user {user.uid}'}, 200
        except Exception as e:
            print(f'Error creating user: {e}')
            return {'message': f'Error creating user: {str(e)}'}, 400
        
#Api route to get a new token for a valid user
@app.route('/api/token')
def token():
    email = request.form.get('email')
    password = request.form.get('password')
    try:
        user = pb.auth().sign_in_with_email_and_password(email, password)
        jwt = user['idToken']
        return {'token': jwt}, 200
    except Exception as e:
        return {'message': f'There was an error logging in: {str(e)}'}, 400

#Sign in page that the client will see
@app.route("/", methods=['GET'])
def loginrender():
    return render_template('login_page.html')

#Login logic
@app.route("/login", methods=['POST'])
def login():
    email = request.form.get('email')
    password = request.form.get('password')

    #Get token
    try:
        auth_resp = requests.get(url_for('token', _external=True), data={'email': email, 'password': password})
    except Exception as e:
        return {'message': f'There was an error logging in: {str(e)}'},400

    auth_json = auth_resp.json()
    if 'token' not in auth_json:
        return auth_json, auth_resp.status_code

    token = auth_json['token']
    token_split = token.split('.')
    cookieVal = token_split[1]
    
    username = auth.get_user_by_email(email=email).display_name
    
    #set cookie
    resp = make_response({})
    resp.set_cookie('ChatUserAuth', cookieVal)
    
    exists, e = check_db(username)
    if e != None:
        return {'message': f'Error creating user: {str(e)}'}, 400
    if exists:
        e = add_to_db(f"UPDATE users set cookie = '{cookieVal}' WHERE username = '{username}';")
        if e != None:
            return {'message': f'Error creating user: {str(e)}'}, 400
    else:
        e = add_to_db(f"INSERT INTO users (id, username, cookie) VALUES (gen_random_uuid(), '{username}', '{cookieVal}');")
        if e != None:
            return {'message': f'Error creating user: {str(e)}'}, 400

    return resp

#Api route to get users
@app.route('/api/userinfo')
@check_token
def userinfo():
    return {'data': users}, 200

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=9000, debug=True)
