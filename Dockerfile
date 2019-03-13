FROM ubuntu:18.04

COPY asset /asset

ENV GOPATH=/go \
    http_proxy=${http_proxy} \
    https_proxy=${https_proxy}

RUN apt-get update && \
    apt-get install -y \
    golang \
    python3 \
    python3-psycopg2 \
    git \
    jlha-utils \
    unzip \
    && apt-get clean \
	mkdir /go && mkdir /bootstrap 
    
RUN go get github.com/kokardy/medmst \
    github.com/lib/pq \
    github.com/jmoiron/sqlx \
	github.com/gin-gonic/gin

RUN cd /bootstrap && \
    git clone https://github.com/riot/riot && \
    git clone https://github.com/github/fetch && \
    git clone https://github.com/taylorhakes/promise-polyfill

ENTRYPOINT sh /asset/routine.sh && /asset/server

