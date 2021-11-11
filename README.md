# replication-log

### Run application
docker-compose up -d

### Test application
8080 port is the leader port for write,read
8081,8082 follower port for read (Write to followers will be redirected to leader)
#### Example
test.sh includes some basic requets
bash test.sh # test app behaviour

curl -d "@payload.json" -X POST http://localhost:8080/write # write message
curl http://localhost:8080 # reade messsages
