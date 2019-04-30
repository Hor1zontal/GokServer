#!/bin/sh


#导出服务器数据到本地
workdir=/Users/hejialin/Documents/aliens/gok/backup
cd ${workdir}

DUMP_DIR=gok_sdw_dump_20181114

#clean redis
redis-cli -h 127.0.0.1 -p 6379 -a aliens flushall


local_collections=(gok_game gok_passport gok_star gok_community gok_search gok_mail gok_trade gok_log)


#clean and import
for i in "${!local_collections[@]}"; do
	mongo ${local_collections[$i]} --eval "db.dropDatabase()" 
	mongorestore -d ${local_collections[$i]} --dir ${DUMP_DIR}/${local_collections[$i]}
done

