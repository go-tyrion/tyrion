# go-tyrion
Tyrion 多功能 Golang 服务框架

### 关于命名

提利昂·兰尼斯特（Tyrion Lannister）是美国作家乔治·R·R·马丁所著长篇史诗奇幻小说《冰与火之歌》中重要的POV角色之一，前五卷pov数量排行榜第一。他是西境守护、凯岩城公爵泰温·兰尼斯特最小的孩子，是个容貌丑陋的侏儒。提利昂非常喜爱读书，善于思考，富有谋略，但是由于他天生畸形，出生时还导致母亲乔安娜·兰尼斯特难产死亡，所以父亲对他极其厌恶。虽然出身高贵并且富有权势，畸形的身材仍然给他带来了许多问题和困扰。

### 功能

**命令行**
- cli

**网络功能**
- http
- websocket
- socket

**框架**
- MySQL
- Redis(Cluster)
- ActiveMQ
- Kafka
- Memcache

**扩展功能**
- 日志处理
- 配置文件解析处理
- util / tool公共函数
- Http 客户端
- Consul
- Apollo 配置或 consul 配置
- Proto 协议
- RPC 通信
- Aes 加密

**后期**
- 前端框架支持

### 目录结构

```$xslt
|- project/
    |- dist
        |- name/
            |- bin/
            |- config/
        |- log_dir/
            |- access.log
            |- error.log
            |- info.log
    |- src/
        |- name/
            |- main.go
        |- lib[submodule]/
```
