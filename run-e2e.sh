#!/bin/bash

docker build -t lethimcook ./app
docker build -t lethimcook-e2e ./e2e

docker run -p 8080:8080 -d -t lethimcook
docker run --network=host -t lethimcook-e2e npm run test

docker rm $(docker stop $(docker ps -q --filter ancestor="lethimcook"))
docker rm $(docker ps -q -a --filter ancestor="lethimcook-e2e")

