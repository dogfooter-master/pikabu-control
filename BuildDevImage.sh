#!/bin/sh
EXE=docker-compose
unameOut="$(uname -s)"
if [ 0 -ne 0 ]; then
	EXE=${EXE}.exe
fi

${EXE} build --no-cache
docker push flowork/pikabu-control-dev:latest
