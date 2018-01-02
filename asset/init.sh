go get github.com/kokardy/medmst
go get github.com/lib/pq
go get github.com/jmoiron/sqlx
cd /bootstrap
DATE=`date +%Y%m%d`
/go/bin/medmst -d $DATE -p $http_proxy -f
mv $DATE save
cd /bootstrap/save/hot
jlha xif *.lzh
cd /bootstrap/save/y
unzip -jo y.zip
