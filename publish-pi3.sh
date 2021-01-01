#!/bin/bash

echo "Builds app. go build is running..."
go build -o live-streamer.bin

cd ./deploy

echo "build the zip package"
./deploy.bin -target pi3 -outdir ~/app/live-streamer/zips/
cd ~/app/live-streamer/

echo "update the service"
./update-service.sh

echo "Ready to fly"