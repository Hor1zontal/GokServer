#!/bin/sh
SET SERVER=120.77.213.111

DUMP_DIR=./gok_test_dump

#clean redis
redis-cli -h %SERVER% -p 20101 -a 55d1JRzYLwnD4Zhh flushall

#clean mongo
collections="gok gok_log gok_mail gok_passport gok_star gok_community gok_search "

for %%a IN (%remote_collections%) DO (
	mongodump -h %SERVER% --port 27017 -d "%%a" -u aliens001 -p eZkZ6pMstAm8MY7Y -o %DUMP_DIR%
)

for %%a IN (%remote_collections%) DO (
	mongo %SERVER%:27017/%%a -u aliens -p 1h37l6foHtQ7GzzH --eval "db.dropDatabase()"
)