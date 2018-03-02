#!/bin/bash

releasePath="build/conf"

if [ ! -d "$releasePath" ]; then
  echo "create " + $releasePath
  mkdir -p $releasePath
fi

go build -o build/cron main.go
GOOS=linux GOARCH=386 go build -o build/cron-linux main.go
GOOS=windows GOARCH=386 go build -o build/cron.exe main.go
echo "buid done"
cp -R templates build/
echo "copy config"
cp conf/config.yaml build/conf/config.yaml
echo "copy templates complete"
cp -R public build/
echo "copy static complete"
