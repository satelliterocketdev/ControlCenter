#!/bin/bash

/app/bin/migrate -source file:///app/migrations/ -database "mysql://$MYSQL_USER:$MYSQL_PASSWORD@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE" version
/app/bin/migrate -source file:///app/migrations/ -database "mysql://$MYSQL_USER:$MYSQL_PASSWORD@tcp($MYSQL_HOST:$MYSQL_PORT)/$MYSQL_DATABASE" up

if [[ $NODE_ENV != "production" ]]; then
    echo "Development watcher"
    cd /app/src/i69social/static && ./node_modules/webpack/bin/webpack.js -d --watch &
fi

exec /app/bin/i69social
