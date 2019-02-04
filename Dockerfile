FROM ubuntu:18.04

COPY asset /asset

ENV GOPATH=/go \
    http_proxy=${http_proxy} \
    https_proxy=${https_proxy}

RUN apt-get update && \
    apt-get install -y \
    golang \
    python \
    python-psycopg2 \
    git \
    jlha-utils \
    unzip \
    && apt-get clean \
	mkdir /go && mkdir /bootstrap 
    
RUN go get github.com/kokardy/medmst \
    github.com/lib/pq \
    github.com/jmoiron/sqlx

ENTRYPOINT sh /asset/routine.sh && go run /asset/server.go

