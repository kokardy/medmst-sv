#!/bin/bash
day=`date '+%Y%m%d'`
echo "$day routine.sh runnning"
mv /bootstrap/save /bootstrap/${day}backup
cd /bootstrap
/go/bin/medmst -f
cd /bootstrap/save/hot
unzip -jo *.zip
cd /bootstrap/save/y
unzip -jo y.zip
cd /bootstrap
python3 /asset/register.py -DI
