FROM golang:latest as binder

COPY . /app
WORKDIR /app


# Input the origin package path
#RUN go mod init $PG

ENV GO111MODULE=on

RUN CGO_ENABLED=0 GOOS=linux go build -v -a -installsuffix cgo \
    -o app

# 
FROM alpine:latest 
RUN apk add --no-cache tzdata
ENV TZ Asia/Taipei

EXPOSE 1314

COPY --from=binder /app/app /app

WORKDIR /
ENTRYPOINT ["./app"]


