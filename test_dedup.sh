#!/bin/bash
while true;
do
curl -X POST http://127.0.0.1:8080 \
   -H 'Content-Type: application/json' \
   -d '{"message":"my_great_message_2", "write_consistency":3}';
done;
curl -X POST http://127.0.0.1:8080 \
   -H 'Content-Type: application/json' \
   -d '{"message":"my_great_message_2","write_consistency":2}'
curl -X POST http://127.0.0.1:8080 \
   -H 'Content-Type: application/json' \
   -d '{"message":"my_great_message_3","write_consistency":2}'


curl -X GET http://127.0.0.1:8080/
