cd /asset
go build -o server .
chmod +x server
cd /bootstrap
/go/bin/medmst -f
cd /bootstrap/save/hot
jlha xif *.lzh
cd /bootstrap/save/y
unzip -jo y.zip

cp /asset/cron_backup /etc/cron.d/
echo "0 18 * * *  /asset/backup.sh" | crontab -

cd /bootstrap
python3 /asset/register.py -DCI
python3 /asset/import_drug_code.py
