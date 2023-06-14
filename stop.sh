#!/bin/bash

 docker compose down

 docker stop -t 300 roach1 roach2
 
 docker rm roach1 roach2
 
 docker volume rm roach1 roach2

 docker network rm roachnet

