#!/bin/bash

TEST_DIR="./e2e/build"
ENV_FILE_NAME=".env.example"
EXECUTABLE_NAME="lethimcook"

set -x
if [ -d $TEST_DIR ]; then
  rm -r $TEST_DIR
fi

cd app
make

cp -r build ../e2e

cp $ENV_FILE_NAME ../e2e/build/.env

cd ../e2e/build

nohup ./$EXECUTABLE_NAME --init-admin admin &

TIMEOUT=5
while ! curl -s "http://localhost:8080" > /dev/null; do
  ((TIMEOUT--))
  if [ $TIMEOUT -eq 0 ]; then
    echo "Timed out while waiting for server to start"
    exit 1
  fi
  sleep 1
done

cd ../
npx playwright test

killall $EXECUTABLE_NAME
set +x

