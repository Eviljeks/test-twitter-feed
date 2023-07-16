# Task

Twitter feed back end
Use docker-compose and any programming language

1. Implement an endpoint to add message
2. Implement an endpoint to get feed (get existing messages and stream new ones - use HTTP streaming)
3. Implement back pressure for message creation (use RabbitMQ/Kafka)
4. Use Cockroachdb(at least 2-node cluster) as a database
5. Implement a bot to generate messages (at configurable speed)

CRITICAL - Project must start with one command (bash file) without installing anything except docker
Result is a link to a git project

# Usage

## Run

To startup application use `./run.sh` (CRITICAL requirement).

** If you see something like this in the console
```
ERROR: cluster has already been initialized
Failed running "init"
no change
```
startup was completed successfully.

Type 

```
docker compose ps -a
```
to see your containers are running.

Type 

```
docker compose logs -f
```
to see containers logs.

## Stop

To stop application use `./stop.sh`

## Dev

Visit `http://localhost:3000/messages` for feed.

You can open multiple tabs and you will receive on each tab updates from bot / other users.

Bot'll begin creating messages after 15 seconds (configurable in `docker-compose.yaml` argument `--delay` for `bot` service command).

To increase bot's post frequency, increase `--reqs-per-min` argument for `bot` service command in `docker-compose.yaml`.

# Implementation

Look into `docker-compose.yaml` for better usage userstanding. 

Entrypoints for containers are in `cmd/*`.

Three (bot, api, sse) containers are using same image, so you need to build it before (`docker build -t twitter-feed .`). This command is already in use in `./run.sh`.

I decided divide endoints for list messsages and SSE fetching, because as for me, these endoints should be accepted by defferent containers, so in the future they may need independent scaling, configuration, etc. 



# TODO

- Better configuration
- Tests
- Graceful shutdown for all services

# Other
Please contact me if you have other questions.
