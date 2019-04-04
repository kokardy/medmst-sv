cd /asset
go build -o server .
cd /bootstrap
/go/bin/medmst -f
cd /bootstrap/save/hot
jlha xif *.lzh
cd /bootstrap/save/y
unzip -jo y.zip
cd /bootstrap

python3 register.py -CI