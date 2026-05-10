#!/bin/bash

EXEC_NAME="webbot"
USER=carpet
HOST=10.0.0.185

GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o $EXEC_NAME .
scp $EXEC_NAME $USER@$HOST:~/$EXEC_NAME

rm $EXEC_NAME
