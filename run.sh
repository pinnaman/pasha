killall -9 CompileDaemon
killall -9 api
echo "Building API Server"
#cd ~/pasha/cmd
#CompileDaemon -command="./cmd/api" &
#CompileDaemon -log-prefix=false -build="go build -o /home/ajay/pasha/cmd" -command="./api" &
CompileDaemon -log-prefix=false -build="go build ./cmd/api" -command="./api" &

ps -ef|grep Compile|grep -v grep

echo "#****Running a Test API Routes*****#"
curl localhost:8090?name=ajay