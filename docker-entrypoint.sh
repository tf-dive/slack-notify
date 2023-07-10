#!/bin/sh
set -e

if [ "$1" != "./slack-notify" ]; then
  exec "$@"
  exit $?
fi

exec "$@" -logtostderr=true -stderrthreshold=${LOGGING_LEVEL}
