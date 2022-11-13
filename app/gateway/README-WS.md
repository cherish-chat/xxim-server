# websocket gateway

## 路径以及参数

- 路径：/ws?token=${token}&userId=${userId}&networkUsed=${networkUsed}&platform=${platform}&deviceId=${deviceId}
- 参数：

| 参数名 | 类型 | 是否必须 | 说明 | 示例 |
| :--- | :--- | :--- | :--- | :--- |
| token | string | 是 | 用户token | ey.x.x |
| userId | string | 是 | 用户id | 123456 |
| networkUsed | string | 是 | 网络类型 | wifi |
| platform | string | 是 | 平台类型 | ios |
| deviceId | string | 是 | 设备id | xxx-xxx-xxx |

## 连接说明

    连接失败后，请首先先查是否是网络原因。其次检查token是否有效。
    连接成功后，服务端会发送一条 `text` 类型消息，内容为 `connected`，表示连接成功。

## 部署说明

    如果使用云平台提供的负载均衡，需要配置支持websocket，一般云平台支持websocket最大连接时间为3600s，超过这个时间，连接会断开，需要重新连接。

## 开发说明

    如果需要多地区部署，请修改 `internal/handler/handler.go` 中的 `// 自定义负载均衡` 处代码，实现自己的负载均衡策略。