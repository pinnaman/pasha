#!/bin/bash

run=$1

# remove stopped containers
docker rm $(docker ps -a -q)
# remove untagged images
docker rmi $(docker images | grep "^<none>" | awk "{print $3}")

killall -9 CompileDaemon
killall -9 api
killall -9 news

docker stop datadoor
docker rm -f datadoor
docker-compose up -d  # creates and start the postgres container

if [[ ${run} == "init" ]];
then

echo "#***********Building Postgres News DB and Installing Go Libs on HOST******#"
#*************** Run below step only for first time********************#
# Below Golang Docker app not used currently, app run directly on host
#docker build -t pinnaman/datadoor . # creates the go app container

docker volume rm news_data
docker volume create news_data

###################################
# Below not required, created and started using docker compose below
#docker run -d --name datadoor -v news_data:/var/lib/postgresql/data -v $(pwd)/scripts:/scripts -p 54320:5432 postgres:12-alpine
#docker run -d --rm --name datadoor \
#-e POSTGRES_PASSWORD=postgres \
#-v news_data:/var/lib/postgresql/data \
#-v $(pwd)/scripts:/scripts \
#-p 54320:5432 postgres:12-alpine
#############################
docker rm -f datadoor
docker-compose up -d  # creates and start the postgres container

sleep 20

echo "create news database"
docker exec -it datadoor psql -h localhost -U postgres -f ./scripts/db_setup.sql
echo "loading news database...."
docker exec -i datadoor psql -U postgres -d ddoor_db < pg_data/pg_backup.sql

# Add below to ~/.bashrc
#echo 'export GOPATH=$HOME/go' >> ~/.bashrc
#echo 'export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin' >> ~/.bashrc
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

echo "#****Running a Test API Routes*****#"
curl localhost:8090?name=ajay