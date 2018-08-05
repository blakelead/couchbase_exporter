FROM golang:1.10-alpine
LABEL maintainer="Blake Lead"

ENV LISTEN_ADDR=:9191
ENV TELEMETRY_PATH=/metrics
ENV CB_URI=http://localhost:8091
ENV CB_ADMIN_USER=
ENV CB_ADMIN_PASSWORD=

RUN apk add --update git && mkdir /app 
ADD . /app/ 
WORKDIR /app 

RUN go get -d -v
RUN go build -o main && ls 

CMD "/app/main -web.listen-address=$LISTEN_ADDR -db.url=$CB_URI -db.user=$CB_ADMIN_USER -db.pwd=$CB_ADMIN_PASSWORD"