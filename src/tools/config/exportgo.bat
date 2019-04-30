#!/bin/sh
workdir=$(cd $(dirname $0); pwd)
configdir=G:\Client\GokClient\config\table

java -jar tool.jar -d go -i ${configdir} -o ../../gok/constant/tableconstant.go -t ${workdir}/template/go_enum.template
java -jar tool.jar -d go -i ${configdir} -o test.go -t ${workdir}/template/go_struct.template