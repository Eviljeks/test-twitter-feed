#!/bin/bash

docker build -t twitter-feed .

docker network create -d bridge roachnet
docker volume create roach1
docker volume create roach2

docker run -d \
  --env COCKROACH_DATABASE=test\
  --env COCKROACH_USER=test \
  --name=roach1 \
  --hostname=roach1 \
  --net=roachnet \
  -p 26257:26257 \
  -p 8080:8080 \
  -v "roach1:/cockroach/cockroach-data" \
  cockroachdb/cockroach:v23.1.3 start \
    --advertise-addr=roach1:26357 \
    --http-addr=roach1:8080 \
    --listen-addr=roach1:26357 \
    --sql-addr=roach1:26257 \
    --insecure \
    --join=roach1:26357,roach2:26357

docker run -d \
  --env COCKROACH_DATABASE=test\
  --env COCKROACH_USER=test \
  --name=roach2 \
  --hostname=roach2 \
  --net=roachnet \
  -p 26258:26258 \
  -p 8081:8081 \
  -v "roach2:/cockroach/cockroach-data" \
  cockroachdb/cockroach:v23.1.3 start \
    --advertise-addr=roach2:26357 \
    --http-addr=roach2:8081 \
    --listen-addr=roach2:26357 \
    --sql-addr=roach2:26258 \
    --insecure \
    --join=roach1:26357,roach2:26357

docker exec -it roach1 ./cockroach --host=roach1:26357 init --insecure

docker run -v $PWD/migrations/:/migrations --network roachnet migrate/migrate -path=/migrations/ -database cockroachdb://root@roach1:26257/defaultdb?sslmode=disable up

docker compose up -d