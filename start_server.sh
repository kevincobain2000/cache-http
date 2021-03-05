#! /bin/sh

########################################################################
# 1.GO_SERVER_START_PID
########################################################################
GO_SERVER_START_PID=`cat $APP_DIR/$1.pid`
GO_SERVER_START_IsAlive=`ps ${GO_SERVER_START_PID} | wc -l`

if [ ${GO_SERVER_START_IsAlive} = 2 ]; then
    echo "${TARGET} $1 is alive now"
else
    echo "${TARGET} $1 is dead now"
    echo "Start Application on Port: ${1}";

    start_server --pid-file $APP_DIR/$1.pid -- $APP_DIR/main $1 > /dev/null 2>&1 &
fi
########################################################################
