#!/bin/bash

. /root/env.sh

date
#echo "${PG_HOST}:${PG_PORT}:${PG_DATABASE}:${PG_USER}:${PG_PASSWORD}" 
#echo "${PG_HOST}:${PG_PORT}:${PG_DATABASE}:${PG_USER}:${PG_PASSWORD}" > ~/.pgpass
echo "*:*:*:*:${PG_PASSWORD}" > ~/.pgpass
#sudo su postgres
#cat ~/.pgpass

chmod 0600 ~/.pgpass

dir="/backup"

mkdir -p $dir

filename="backup"

pg_dump  ${PG_DATABASE} \
		-h ${PG_HOST} \
		-U ${PG_USER} \
		-w -Fc -f ${dir}/${filename} \
		-t yj -t hot -t custom_yj 

filename2=`date '+%Y%m%d%H%M'`

cp ${dir}/${filename} ${dir}/${filename2}

echo "backup:${dir}/${filename2}"
