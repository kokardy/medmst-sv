echo "input password for postgresql(10.26.61.131)"
psql -W -U postgres -h 10.26.61.131  -c "COPY (select * from yj_status) TO STDOUT;" > /tmp/yj_status
cat /tmp/yj_status | psql medmst -h postgres -U postgres -c "DELETE from yj; COPY yj FROM STDIN;"
