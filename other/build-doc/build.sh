#!/bin/zsh

# shellcheck disable=SC2034
TAG=$(date +%Y%m%d%H%M%S)

docker build --platform linux/x86_64 -t "registry.cn-hangzhou.aliyuncs.com/xxim-dev/xxim-doc:$TAG" . --no-cache

docker push "registry.cn-hangzhou.aliyuncs.com/xxim-dev/xxim-doc:$TAG"

echo "编译并推送完成："
echo "registry.cn-hangzhou.aliyuncs.com/xxim-dev/xxim-doc:$TAG"