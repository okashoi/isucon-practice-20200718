#!/bin/bash -eu

echo 'init process'

MYSQL="mysql -u isucon -pisucon isucon"
$MYSQL <<EOF
ALTER TABLE memos ADD INDEX index_user(user);
ALTER TABLE memos ADD INDEX index_created_at(created_at);
EOF