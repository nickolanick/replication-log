#!/bin/bash

curl -X POST http://127.0.0.1:8081/ \
   -H 'Content-Type: application/json' \
   -d '{"message":"my_great_message","write_consistency":2}'

curl -X POST http://127.0.0.1:8080/ \
   -H 'Content-Type: application/json' \
   -d '{"message":"my_great_message", "write_consistency": 1}'

curl -X POST http://127.0.0.1:8080/ \
   -H 'Content-Type: application/json' \
   -d '{"message":"another_great_message", "write_consistency": 2}'

curl -X GET http://127.0.0.1:8082/

#./vegeta attack -targets=vegeta.txt -duration=10s | tee results.bin | vegeta plot
