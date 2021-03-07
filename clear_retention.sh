#! /bin/sh

export $(egrep -v '^#' .env | xargs)
# * * * * * cd /path/to.dir/&& ./clear_retention.sh
echo "clearing old files"
find $APP_DIR/assets/*.tar.gz -mtime +7 -exec rm -rf {} \;
