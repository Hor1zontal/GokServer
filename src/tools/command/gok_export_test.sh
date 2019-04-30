#!/bin/sh

#导出测试服务器数据到本地
workdir=$(cd $(dirname $0); pwd)
echo ${workdir}

SERVER=120.77.213.111
DUMP_DIR=gok_test_dump

#clean redis
redis-cli -h 127.0.0.1 -p 6379 -a aliens flushall



local_collections=(gok_game gok_passport gok_event gok_star gok_community gok_search gok_mail gok_trade)
remote_collections=(gok gok_passport gok_event gok_star gok_community gok_search gok_mail gok_trade)

#dump mongo
for i in "${!remote_collections[@]}"; do 
	mongodump -h ${SERVER} --port 27017 -d ${remote_collections[$i]} -u aliens001 -p eZkZ6pMstAm8MY7Y -o ${DUMP_DIR}
done


#clean and import
for i in "${!remote_collections[@]}"; do 
	mongo ${local_collections[$i]} --eval "db.dropDatabase()" 
	mongorestore -d ${local_collections[$i]} --dir ${DUMP_DIR}/${remote_collections[$i]}
done

