for f in `find . -type f`;
do
        if [ -s $f ] && [ "${f##*.}"x = "proto"x ];then
              echo copy $f
              cp -rf $f /Users/hejialin/git/Gok/tools/copyFile/protocol/
        fi
done

cd /Users/hejialin/git/Gok/tools/copyFile/protocol/
python protoToJs.py