#!/bin/bash

run=$1

# remove stopped containers
docker rm $(docker ps -a -q)
# remove untagged images
docker rmi $(docker images | grep "^<none>" | awk "{print $3}")

killall -9 CompileDaemon
killall -9 api
killall -9 news

docker rm -f datadoor
docker-compose up -d  # creates and start the postgres container

if [[ ${run} == "init" ]];
then

echo "#***********Building Postgres News DB and Installing Go Libs on HOST******#"
#*************** Run below step only for first time********************#
# Below Golang Docker app not used currently, app run directly on host
#docker build -t pinnaman/datadoor . # creates the go app container

docker volume rm pasha_news_data
docker volume create pasha_news_data

docker stop datadoor
docker rm -f datadoor
docker-compose up -d  # creates and start the postgres container

sleep 20

echo "create and set up news database"
~/pasha/scripts/load_pgdb.sh

source ~/.bashrc
# Install Go version Centrally (/usr/local/bin)
go version

cd ~/ajay/pasha
echo "Installing golang libraries"
./go_install.sh

fi

sleep 30
echo "Building API Server"
CompileDaemon -log-prefix=false -build="go build -o ./bin ./cmd/api" -command="./bin/api" &
echo "Building NEWS Server"
CompileDaemon -log-prefix=false -build="go build -o ./bin ./cmd/news/" -command="./bin/news" &

# confirm that apps started on ports
netstat -antp|grep tcp6
ps -ef|grep Compile|grep -v grep

echo "##*************TESTS****************##"
echo "#****Running a Test API Routes*****#"
curl localhost:8090?name=ajay