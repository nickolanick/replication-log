#!/bin/bash

curl -X POST http://127.0.0.1:8082/commit \
   -H 'Content-Type: application/json' \
   -d '{"message":"my_great_message","write_consistency":2, "total_order":2}'

curl -X POST http://127.0.0.1:8082/commit \
   -H 'Content-Type: application/json' \
   -d '{"message":"my_great_message","write_consistency":2, "total_order":3}'

curl -X GET http://127.0.0.1:8082/

curl -X POST http://127.0.0.1:8082/commit \
   -H 'Content-Type: application/json' \
   -d '{"message":"my_great_message","write_consistency":2, "total_order":1}'


curl -X GET http://127.0.0.1:8082/
