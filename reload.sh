#! /bin/sh

export $(egrep -v '^#' .env | xargs)

kill -HUP `cat $APP_DIR/3000.pid`


