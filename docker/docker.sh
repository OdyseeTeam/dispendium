#!/usr/bin/env bash

docker build --tag odyseeteam/dispendium:$TRAVIS_BRANCH ./
echo "$DOCKER_PASSWORD" | docker login --username "$DOCKER_USERNAME" --password-stdin
docker push odyseeteam/dispendium:$TRAVIS_BRANCH