#/bin/bash

echo "${PG_HOST}:${PG_PORT}:${PG_DATABASE}:${PG_USER}:${PG_PASSWORD}" 
#echo "${PG_HOST}:${PG_PORT}:${PG_DATABASE}:${PG_USER}:${PG_PASSWORD}" > ~/.pgpass
echo "*:*:*:*:${PG_PASSWORD}" > ~/.pgpass
#sudo su postgres
cat ~/.pgpass

chmod 0600 ~/.pgpass

dir="/backup"

mkdir -p $dir

filename="backup.sql"

pg_dump ${PG_DATABASE} \
	-h ${PG_HOST} \
	-U ${PG_USER} \
        -Fp -f ${dir}/${filename} \
        -t yj -t hot -t costom_yj \
	-w

filename2=`date '+%Y%m%d%H%M'`.sql

cp ${dir}/${filename} ${dir}/${filename2}

echo "backup:${filename2}"
