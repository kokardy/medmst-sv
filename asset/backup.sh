#!/bin/bash

. /root/env.sh

cat /root/env.sh

echo "${PG_HOST}:${PG_PORT}:${PG_DATABASE}:${PG_USER}:${PG_PASSWORD}" > ~/.pgpass

cat ~/.pgpass

chmod 0600 ~/.pgpass

dir="/backup"
filename="backup.sql"

mkdir -p ${dir}

pg_dump ${PG_DATABASE} \
		-h ${PG_HOST} \
		-U ${PG_USER} \
		-w -Fp -f ${dir}/${filename} \
		-t yj -t hot -t costom_yj 

cat ~/.pgpass
filename2=`date '+%Y%m%d%H%M'`.sql

cp ${dir}/${filename} ${dir}/${filename2}

echo "backup:${filename2}"
