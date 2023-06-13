FROM golang:1.20-buster AS build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

COPY cmd cmd
COPY internal internal
COPY migrations migrations
COPY pkg pkg

RUN go mod download

RUN go build -o /twitter-feed ./cmd/

FROM gcr.io/distroless/base-debian11

WORKDIR /

COPY --from=build /twitter-feed /twitter-feed

EXPOSE 3000
EXPOSE 3001

USER nonroot:nonroot