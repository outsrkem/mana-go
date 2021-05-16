#!/bin/sh
set -e
if [ -z "$1" ];then
  # mana --db-username=abc --db-password=123456 --db-server=10.10.10.10 --db-name=mana
  set -- mana
fi

exec "$@"
