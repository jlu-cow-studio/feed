# 基础镜像
FROM golang:1.17-alpine

# 设置工作目录
WORKDIR /app

# 将宿主机中的当前目录的所有文件拷贝到镜像中的 /app 目录下
COPY . .

# 暴露服务端口
EXPOSE 8080
EXPOSE 8081

# 执行命令
ENTRYPOINT ["pwd", "ls -la /app/build", "/app/build/run.sh"]
ENTRYPOINT ["ls -la /app/build"]
ENTRYPOINT ["sh", "/app/build/run.sh"]
