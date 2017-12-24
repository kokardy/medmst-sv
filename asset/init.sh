cd /bootstrap
DATE=`date +%Y%m%d`
/go/bin/medmst -d $DATE -p $http_proxy -f
mv $DATE save
cd /bootstrap/save/hot
jlha xif *.lzh
cd /bootstrap/save/y
unzip -jo y.zip
