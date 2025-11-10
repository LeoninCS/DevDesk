# 启动 Redis 服务
sudo service redis-server start && \
redis-cli ping && \

# 设置 Redis 地址（解决数据不共享问题）
export REDIS_ADDR=127.0.0.1:6379

# 启动 Moment 服务
sudo kill -9 $(sudo lsof -t -i:8080)
sudo kill -9 $(sudo lsof -t -i:8081)
sudo kill -9 $(sudo lsof -t -i:8082)
sudo kill -9 $(sudo lsof -t -i:8083)
go run cmd/main.go & \
go run cmd/CodeShare/main.go 8081 & \
go run cmd/CodeShare/main.go 8082 & \
go run cmd/CodeShare/main.go 8083 &

trap "echo; echo '停止所有服务...'; pkill -f 'go run cmd'; exit 0" SIGINT
wait

# 运行方式
# chmod +x start.sh
# ./start.sh (在项目根目录下运行，即Moment目录)
# Ctrl+C 组合键可停止所有服务