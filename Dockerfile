FROM golang:1.20.5-alpine3.18 AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY migrations migrations
COPY pkg pkg
COPY cmd cmd
COPY internal internal
COPY templates templates

RUN go mod download

RUN go build -o /twitter-feed ./cmd/

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /twitter-feed /twitter-feed
COPY --from=build /app/templates /templates

EXPOSE 3000

USER nonroot:nonroot