#! /bin/sh
export $(egrep -v '^#' .env | xargs)

########################################################################
# 1.GO_SERVER_START_PID
########################################################################
port=$1
pid=$2
if [ -z "$pod" ]; then
    echo "Please provide port (first arg)"
    exit
fi
if [ -z "$pid" ]; then
    echo "Please provide pid dir as (second arg)"
    exit
fi

GO_SERVER_START_PID=`cat $APP_DIR/$1.pid`
GO_SERVER_START_IsAlive=`ps ${GO_SERVER_START_PID} | wc -l`
isAlive=`ps ${GO_SERVER_START_PID} | wc -l`

if [ ${isAlive} = 2 ]; then
    echo "${TARGET} $1 is alive now"
else
    echo "${TARGET} $1 is dead now"
    echo "Start Application on Port: ${1}";

    $APP_DIR/main $1 $2 > /dev/null 2>&1 &
fi
########################################################################
