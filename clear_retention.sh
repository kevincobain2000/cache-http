#! /bin/sh

export $(egrep -v '^#' .env | xargs)

echo "clearing old files"
find $APP_DIR/assets* -mtime +7 -exec rm {} \;
