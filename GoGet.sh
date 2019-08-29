#!/bin/sh
if [ $# -eq 1 ]; then
        GOPATH=$1
fi

cd control/cmd
go get -v 
