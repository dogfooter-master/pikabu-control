#!/bin/sh
docker service rm flowork_pikabu_control 2>/dev/null
docker stack deploy -c docker-stack.deploy.yml flowork
