#!/bin/bash
day=`date '+%Y%m%d'`
mv /bootstrap/save /bootstrap/${day}backup
cd /bootstrap
/go/bin/medmst -f
cd /bootstrap/save/hot
jlha xif *.lzh
cd /bootstrap/save/y
unzip -jo y.zip
cd /bootstrap
cd /asset
python3 register.py -DI
