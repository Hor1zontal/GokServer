::导出服务器数据到本地
::workdir=$(cd $(dirname $0); pwd)
::echo ${workdir}

::set SERVER=120.77.213.111
::set DUMP_DIR=./gok_test_dump

::clean redis
redis-cli -h 127.0.0.1 -p 6379 -a aliens flushall

set DUMP_DIR=./gok_test_dump

set DEV_SERVER=120.77.254.228
set LOCAL_SERVER=127.0.0.1

set remote_collections=gok_game gok_passport gok_star gok_community gok_search gok_log gok_trade gok_center gok_gm gok_mail
set local_collections=gok_game gok_passport gok_star gok_community gok_search gok_log gok_trade gok_center gok_gm gok_mail


::dump mongo
for %%a IN (%remote_collections%) DO (
	mongodump -h %DEV_SERVER% --port 17017 -d "%%a" -u root -p 1h37l6foHtQ7GzzH -o %DUMP_DIR% --authenticationDatabase admin
)

::clean
for %%b IN (%local_collections%) DO (
    mongo mongodb://%LOCAL_SERVER%/%%b --eval "db.dropDatabase()"
)

::import

mongorestore -h %LOCAL_SERVER% --port 27017 -d gok_game --dir %DUMP_DIR%/gok_game

mongorestore -h %LOCAL_SERVER% --port 27017 -d gok_passport --dir %DUMP_DIR%/gok_passport

mongorestore -h %LOCAL_SERVER% --port 27017 -d gok_star --dir %DUMP_DIR%/gok_star

mongorestore -h %LOCAL_SERVER% --port 27017 -d gok_community --dir %DUMP_DIR%/gok_community

mongorestore -h %LOCAL_SERVER% --port 27017 -d gok_search --dir %DUMP_DIR%/gok_search

mongorestore -h %LOCAL_SERVER% --port 27017 -d gok_log --dir %DUMP_DIR%/gok_log