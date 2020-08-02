FROM golang:1.13-alpine

MAINTAINER Apinnamaneni <ajai10@gmail.com>

RUN apk add git bash

WORKDIR /opt/ddoor

COPY . /opt/ddoor

ENV PORT blergh

# installing our golang dependencies
RUN go get github.com/fsnotify/fsnotify && \
  go get github.com/githubnemo/CompileDaemon && \
  go get github.com/gorilla/mux && \
  go get -u github.com/gocolly/colly && \
  go get -u github.com/fatih/color && \
  go get -u github.com/PuerkitoBio/goquery && \
  go get -u github.com/lib/pq

EXPOSE 8000

#ENTRYPOINT go run main.go
ENTRYPOINT CompileDaemon -log-prefix=false -directory="./server/api/" -build="go build ./server/api/" -command="./api"