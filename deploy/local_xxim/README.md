# local_xxim 本地调试环境

## 问题

### 1. 怎么连接docker内网

#### MacOS

> 原文章 https://www.haoyizebo.com/posts/fd0b9bd8/#docker-connector

-
    1. 安装 docker-connector

> brew install wenjunxiao/brew/docker-connector

-
    2. docker-connector 配置文件

> docker network ls --filter driver=bridge --format "{{.ID}}" | xargs docker network inspect --format "route {{range
> .IPAM.Config}}{{.Subnet}}{{end}}" >> /opt/homebrew/etc/docker-connector.conf

-
    3. 启动服务

> sudo brew services start docker-connector

-
    4. docker启动mac-docker-connector

> docker run -it -d --restart always --net host --cap-add NET_ADMIN --name connector wenjunxiao/mac-docker-connector
