UP
```
docker volume create roach-single

docker run -d \
  --env COCKROACH_DATABASE=test\
  --env COCKROACH_USER=test \
  --name=roach-single \
  --hostname=roach-single \
  --net=roachnet \
  -p 26257:26257 \
  -p 8080:8080 \
  -v "roach-single:/cockroach/cockroach-data"  \
  cockroachdb/cockroach:v23.1.2 start-single-node \
  --http-addr=localhost:8080 \
  --insecure 


docker exec -it roach-single grep 'node starting' /cockroach/cockroach-data/logs/cockroach.log -A 11 

docker exec -it roach-single ./cockroach sql --url="postgresql://root@127.0.0.1:26257/defaultdb?sslmode=disable"
```

DOWN
```
docker stop -t 20 roach-single

docker rm roach-single

docker volume rm roach-single
```

```
docker build -t twitter-feed .
```