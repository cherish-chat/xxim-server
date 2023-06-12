#!/bin/zsh

#services=(appmgmt conn group im mgmt msg notice relation user xos)
#配置要编译的服务
services=(conn appmgmt mgmt)

# shellcheck disable=SC2034
SCRIPT_DIR=$(pwd)
TAG=$(date +%Y%m%d%H%M%S)
# 定义writeDockerFile方法
function writeDockerFile() {
cat > Dockerfile << EOF
FROM showurl/debian:go-service
WORKDIR /app
COPY bin /app
CMD ["/app/bin"]
EOF
}
echo "docker TAG = ${TAG}"
cd ../../ && echo "进入项目根目录"
ROOT_DIR=$(pwd)
echo "项目根目录为：${ROOT_DIR}"
# shellcheck disable=SC2128
for service in $services; do
  cd "${ROOT_DIR}/app/$service" && echo "进入$service目录"
  GOOS=linux GOARCH=amd64 go build -o "$SCRIPT_DIR/tmp/$service/bin" "./$service.go"
  cd "$SCRIPT_DIR/tmp/$service" && echo "进入$service build目录"
  writeDockerFile
  docker build --platform linux/x86_64 -t "registry.cn-hangzhou.aliyuncs.com/xxim-dev/xxim-$service:$TAG" .
  docker push "registry.cn-hangzhou.aliyuncs.com/xxim-dev/xxim-$service:$TAG"
  cd "$SCRIPT_DIR" && echo "进入脚本目录"
done
echo "编译并推送完成："
# shellcheck disable=SC2128
for service in $services; do
  echo "registry.cn-hangzhou.aliyuncs.com/xxim-dev/xxim-$service:$TAG"
done
