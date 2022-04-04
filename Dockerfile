# syntax=docker/dockerfile:1
FROM golang:1.16-alpine

WORKDIR /app
COPY . .

# 启用go module
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

RUN go build -o main main.go
EXPOSE 8080
CMD [ "/app/main" ]


