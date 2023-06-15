# user 用户服务

- 该服务负责用户的相关逻辑

## logic 目录结构

```bash
logic
├── README.md
├── userRegisterLogic.go # 用户注册逻辑
├── x_consumerLogic.go # 消息队列消费者逻辑 比如用户注册成功后逻辑
