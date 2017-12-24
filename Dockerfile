FROM ubuntu:16.04

COPY asset /asset

RUN apt-get update && \
    apt-get install -y \
    golang \
    python \
    python-psycopg2 \
    git \
    jlha-utils \
    unzip \
    && apt-get clean

ENV GOPATH=/go
ENV http_proxy=${http_proxy}
ENV https_proxy=${https_proxy}
RUN mkdir /go && mkdir /bootstrap
RUN go get github.com/kokardy/medmst
RUN sh /asset/routine.sh 

