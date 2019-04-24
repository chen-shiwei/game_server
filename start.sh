#!/bin/bash

EXEC_ROOT=$(cd `dirname $0`; pwd)
CFG_FILE=${EXEC_ROOT}/config.toml
PID_FILE=${EXEC_ROOT}/game_server.pid
NAME=$0

#while getopts "Dk" arg
#do
#    case $arg in
#	D)
#	;;
#	k)
#	;;
#   esac
#done

function usage() {
    echo ${NAME} "[start | stop | restart | reload]"
}

function stop() {
    echo Stopping Game Server
    if [[ ! -e ${PID_FILE} ]]; then
	echo game server was not running.
	return 2
    fi
    kill -QUIT $(cat ${PID_FILE})
    rm ${PID_FILE}
    return 0
}

function reload() {
    echo Reloading Game Server Config
    if [[ ! -e ${PID_FILE} ]]; then
	echo game server was not running.
	return 2
    fi
    kill -USR1 $(cat ${PID_FILE})
    return 0
}

function start() {
    if [[ -e ${PID_FILE} ]]; then
	stop
    fi
    echo -n "Starting Game Server "
    if [[ ! -e ${CFG_FILE} ]]; then
	echo Config file does not exists.
	exit 1
    fi
    ./game_server 2>&1 >> error.log &
    if [[ "$?" == "0" ]]; then
	echo $! > ${PID_FILE}
	echo Done
	return 0
    fi
    echo Failed
    return 0
}

if [[ "$#" != "1" ]]; then
    usage
    exit 1
fi

case $1 in
    start)
	start
	;;
    stop)
	stop
	;;
    reload)
	reload
	;;
    restart)
	stop
	start
	;;
esac


