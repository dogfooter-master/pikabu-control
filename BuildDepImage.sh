#!/bin/sh
EXE=docker
unameOut="$(uname -s)"
if [ 0 -ne 0 ]; then
	EXE=${EXE}.exe
fi

${EXE} build --no-cache -t flowork/pikabu-control:latest .
docker push flowork/pikabu-control:latest
