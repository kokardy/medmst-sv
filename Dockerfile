FROM ubuntu:18.04

COPY asset /asset
COPY supervisord.conf /etc/supervisord.conf
COPY asset/routine.sh /etc/cron.daily/routine.sh

ENV GOPATH=/go \
    http_proxy=${http_proxy} \
    https_proxy=${https_proxy} \
    YJ_REDIRECTER=pmda-kv

RUN apt-get update && \
    apt-get install -y \
    golang \
    python3 \
    python3-psycopg2 \
    git \
    jlha-utils \
    unzip \
    wget \
    supervisor \
    cron \
    && apt-get clean \
	mkdir /go && mkdir /bootstrap 
    
RUN go get github.com/kokardy/medmst \
    github.com/lib/pq \
    github.com/jmoiron/sqlx \
	github.com/gin-gonic/gin

RUN cd /bootstrap && \
    git clone https://github.com/riot/riot && \
    git clone https://github.com/taylorhakes/promise-polyfill

RUN mkdir -p /bootstrap/fetch && \
    cd /bootstrap/fetch && \
    wget https://github.com/github/fetch/releases/download/v3.0.0/fetch.umd.js

RUN sh /asset/init.sh

CMD supervisord -c /etc/supervisord.conf
