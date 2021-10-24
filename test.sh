#!/bin/bash

curl -X POST http://127.0.0.1:5000/
   -H 'Content-Type: application/json'
   -d '{"message":"my_great_message"}'
