#!/bin/bash
curl -X POST http://127.0.0.1:8080 \
   -H 'Content-Type: application/json' \
   -d '{"message":"my_great_message_6", "write_consistency":1}';
#curl -X POST http://127.0.0.1:8080 \
#   -H 'Content-Type: application/json' \
#   -d '{"message":"my_great_message_7","write_consistency":2}'
#curl -X POST http://127.0.0.1:8080 \
#   -H 'Content-Type: application/json' \
#   -d '{"message":"my_great_message_3","write_consistency":3}'


curl -X GET http://127.0.0.1:8080/

curl -X GET http://127.0.0.1:8080/health
