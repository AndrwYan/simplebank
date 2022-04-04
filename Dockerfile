# build stage
FROM golang:1.16-alpine AS builder
WORKDIR /app
COPY . .
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

COPY app.env .
RUN go build -o main main.go

#Run stage
FROM alpine
WORKDIR /app
#这条命令的意思就是将第一阶段构建的二进制可执行文件复制到第二阶段的路径也就是命令中的.（.代表着当前路径也就是WORKDIR中/app）
COPY --from=builder /app/main .

EXPOSE 8080
CMD [ "/app/main" ]


