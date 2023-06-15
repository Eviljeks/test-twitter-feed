#!/bin/bash

docker build -t twitter-feed .

docker network create -d bridge roachnet

docker compose up -d

docker compose  exec -it roach1 ./cockroach --host=roach1:26357 init --insecure

docker run --rm -v $PWD/migrations/:/migrations --network roachnet migrate/migrate -path=/migrations/ -database cockroachdb://root@roach1:26257/defaultdb?sslmode=disable up