#!/bin/bash


tag=$1

PG="registry.gitlab.com/axolotl.team/pikachu/demo/api:$tag"


docker build -t $PG .

docker push $PG
