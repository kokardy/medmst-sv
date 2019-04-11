#/bin/bash

echo "${PG_HOST}:${PG_PORT}:${PG_DATABASE}:${PG_USER}:${PG_PASSWORD}" > ~/.pgpass

chmod 0600 ~/.pgpass

dir="/backup"
filename="backup.sql"

pg_dump -w \ 
        -Fp -f ${dir}/${filename}
        -t yj -t hot -t costom_yj 

filename2=`date '+%Y%m%d%H%M'`.sql

cp ${dir}/${filename} ${dir}/${filename2}

echo "backup:${filename2}"
