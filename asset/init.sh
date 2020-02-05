cd /asset
go build -o server .
chmod +x server
cd /bootstrap
/go/bin/medmst -f
cd /bootstrap/save/hot
jlha xif *.lzh
cd /bootstrap/save/y
unzip -jo y.zip

#cron
cp /asset/cron_backup /etc/cron.d/
echo "30 6 * * *  /asset/routine.sh" | crontab -
#cron reload
/etc/init.d/cron reload


sleep 5
cd /bootstrap
python3 /asset/register.py -DCI
#python3 /asset/import_drug_code.py
cd /asset
bash restore.sh

