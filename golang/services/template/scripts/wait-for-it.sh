#!/bin/sh
# wait-for-it.sh

set -e

host="$1"
shift

until nc -vz "$host"; do
  >&2 echo "$host is unavailable - sleeping"
  sleep 1
done

>&2 echo "$host is up - executing command"
exec "$@"
