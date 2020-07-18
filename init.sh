#!/bin/bash -eu

echo 'init process'

MYSQL="mysql -u isucon -pisucon"
$MYSQL <<EOF
USE isucon;
ALTER TABLE memos ADD INDEX index_user(user);
ALTER TABLE memos ADD INDEX index_is_private(created_at);
EOF