## compile project
FROM golang:1.25 AS build

RUN apt update && apt install -y libasound2-dev

WORKDIR /project

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o nyaweria ./cmd/main.go

# build final image
FROM ubuntu:25.04

RUN apt update && apt install -y libasound2-data

COPY --from=build /project/nyaweria /usr/bin/nyaweria

CMD [ "nyaweria" ]
