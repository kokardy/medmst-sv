#!/bin/bash

. /root/env.sh

echo "${PG_HOST}:${PG_PORT}:${PG_DATABASE}:${PG_USER}:${PG_PASSWORD}" 
#echo "${PG_HOST}:${PG_PORT}:${PG_DATABASE}:${PG_USER}:${PG_PASSWORD}" > ~/.pgpass
echo "*:*:*:*:${PG_PASSWORD}" > ~/.pgpass
#sudo su postgres
cat ~/.pgpass

chmod 0600 ~/.pgpass

dir="/backup"
filename="backup.tar"


psql -d medmst \
    -h ${PG_HOST} \
	-U ${PG_USER} \
    -w \
    -c 'DELETE FROM yj; DELETE FROM hot; DELETE FROM custom_yj'


pg_restore -d ${PG_DATABASE} \
		-h ${PG_HOST} \
		-U ${PG_USER} \
		--data-only \
		-t yj -t hot -t custom_yj \
		-w -Ft ${dir}/${filename}

echo "restore:${filename}"
