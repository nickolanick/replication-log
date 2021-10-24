from app import app
from flask import Flask, request,jsonify,Response

messages = []
@app.route('/', methods=['GET'])
def read_messages():
    return jsonify({"messages":messages})
    
@app.route('/', methods=['POST'])
def write_message():
    content = request.json
    print(content)
    new_message = content["message"]

    messages.append(new_message)
    return {"message":new_message}
