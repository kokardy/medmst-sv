FROM ubuntu:18.04

COPY asset /asset

ENV GOPATH=/go \
    http_proxy=${http_proxy} \
    https_proxy=${https_proxy}

RUN apt-get update && \
    apt-get install -y \
    golang \
    python \
    python3 \
    python-psycopg2 \
    git \
    jlha-utils \
    unzip \
    && apt-get clean \
	mkdir /go && mkdir /bootstrap 
    
RUN go get github.com/kokardy/medmst \
    github.com/lib/pq \
    github.com/jmoiron/sqlx \
	github.com/gin-gonic/gin

ENTRYPOINT sh /asset/routine.sh && /asset/server

