#!/bin/sh
SERVER=120.77.213.111

time=$(date "+%Y-%m-%d_%H:%M:%S")

DUMP_DIR=/Users/hejialin/Documents/aliens/gok/backup/gok_dump_${time}

echo ${DUMP_DIR}

#clean redis
redis-cli -h ${SERVER} -p 20101 flushall

#clean mongo
collections="gok gok_log gok_event gok_mail gok_passport gok_star gok_community gok_search gok_center"

for collection in $collections; do
	mongodump -h ${SERVER} --port 27017 -d ${collection} -u aliens001 -p aliens -o ${DUMP_DIR}
done

for collection in $collections;
do
	mongo ${SERVER}:27017/${collection} -u aliens -p aliens --eval "db.dropDatabase()"
done