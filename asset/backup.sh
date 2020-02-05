#!/bin/bash

. /root/env.sh

echo "${PG_HOST}:${PG_PORT}:${PG_DATABASE}:${PG_USER}:${PG_PASSWORD}" 
#echo "${PG_HOST}:${PG_PORT}:${PG_DATABASE}:${PG_USER}:${PG_PASSWORD}" > ~/.pgpass
echo "*:*:*:*:${PG_PASSWORD}" > ~/.pgpass
#sudo su postgres
cat ~/.pgpass

chmod 0600 ~/.pgpass

dir="/backup"

mkdir -p $dir

filename="backup.tar"

pg_dump  ${PG_DATABASE} \
		-h ${PG_HOST} \
		-U ${PG_USER} \
		-w -Ft -f ${dir}/${filename} \
		-t yj -t hot -t custom_yj 

filename2=`date '+%Y%m%d%H%M'`.tar

cp ${dir}/${filename} ${dir}/${filename2}

echo "backup:${filename2}"
