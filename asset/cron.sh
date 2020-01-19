#!/bin/bash
printenv | awk '{print "export " $1}' > /root/env.sh
echo "0 7,18 * * *  /asset/backup.sh >> /var/log/backup.log 2>&1" | crontab -
/usr/sbin/cron -f -L 15
