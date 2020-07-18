#!/bin/bash -eu

echo 'init process'

MYSQL="mysql -u isucon -pisucon isucon"
$MYSQL <<EOF
ALTER TABLE memos ADD INDEX index_user(user);
ALTER TABLE memos ADD INDEX index_created_at(created_at);
ALTER TABLE memos ADD summary text NULL;
EOF
export ISUCON_ENV=prod
pushd /opt/isucon3-mod/app/tools/init_summary/src
./main
popd
