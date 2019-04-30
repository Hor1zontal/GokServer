::clean local redis
redis-cli -h 127.0.0.1 -p 6379 -a aliens flushall

set DUMP_DIR=./gok_test_dump

set LOCAL_SERVER=127.0.0.1

set local_collections=gok_game gok_passport gok_star gok_community gok_search gok_log gok_trade gok_center gok_gm gok_mail

::clean local MongoDB
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

mongorestore -h %LOCAL_SERVER% --port 27017 -d gok_trade --dir %DUMP_DIR%/gok_trade

mongorestore -h %LOCAL_SERVER% --port 27017 -d gok_center --dir %DUMP_DIR%/gok_center

mongorestore -h %LOCAL_SERVER% --port 27017 -d gok_gm --dir %DUMP_DIR%/gok_gm

mongorestore -h %LOCAL_SERVER% --port 27017 -d gok_mail --dir %DUMP_DIR%/gok_mail