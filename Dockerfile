FROM golang:1.15.3-alpine3.12

ADD ./app /app
WORKDIR /app

RUN go install -v ./...

ADD ./resources /resources
CMD ["go", "run", "./cmd/api/main.go"]
