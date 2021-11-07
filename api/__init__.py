import threading
import os
import requests
import json
import time

from flask import Flask, request, jsonify, Response

app = Flask(__name__)

class AppDB():
    def __init__(self):
        self.role = os.environ.get("ROLE", "follower")
        self.followers = os.environ.get("FOLLOWERS").split(",")
        self.delay = int(os.environ.get("DELAY", "0"))
        self.messages = []

    def write_message(self, message):
        self.message = message
        time.sleep(self.delay)
        self.threads = [threading.Thread(target=self.write_secondaries, args=(follower,)) for follower in self.followers]
        if self.role == "leader":
            for thread in self.threads:
                thread.start()
            for thread in self.threads:
                thread.join()

        self.messages.append(message)

    def read_messages(self):
        return self.messages

    def write_secondaries(self, follower):
        data = {'message': self.message}
        headers = {'Content-type': 'application/json', 'Accept': 'text/plain'}

        r = requests.post(follower, data=json.dumps(data), headers=headers)

app_db = AppDB()

@app.route('/', methods=['GET'])
def read_messages():
    return {"messages": app_db.read_messages()}
    
@app.route('/', methods=['POST'])
def write_message():
    content = request.json
    new_message = content["message"]
    app_db.write_message(new_message)

    return {"message": new_message, "followers": app_db.followers}

if __name__ == '__main__':
    app.run(debug=True, use_reloader=True)
