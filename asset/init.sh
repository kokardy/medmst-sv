go get github.com/kokardy/medmst
go get github.com/lib/pq
go get github.com/jmoiron/sqlx
cd /bootstrap
/go/bin/medmst -f
cd /bootstrap/save/hot
jlha xif *.lzh
cd /bootstrap/save/y
unzip -jo y.zip
