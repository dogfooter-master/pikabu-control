#!/bin/sh
EXE=docker
unameOut="$(uname -s)"
if [ 0 -ne 0 ]; then
	EXE=${EXE}.exe
fi

${EXE} build --no-cache -t dermaster/pikabu-control:latest .
docker push dermaster/pikabu-control:latest
