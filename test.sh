#!/bin/bash

curl -X POST http://127.0.0.1:8080/ \
   -H 'Content-Type: application/json' \
   -d '{"message":"my_great_message"}'

curl -X POST http://127.0.0.1:8080/ \
   -H 'Content-Type: application/json' \
   -d '{"message":"another_great_message"}'

curl -X GET http://127.0.0.1:8080/
