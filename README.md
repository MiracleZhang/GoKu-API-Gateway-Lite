## eoLinker AGW（Lite）- 开源网关轻量版

### 简介
***
**eoLinker AGW（Lite）是eoLinker旗下轻量的开源版API网关，同时也是国内唯一的Go语言开源API网关，性能优异，提高API业务安全性。**

通过Go语言开发的高性能API网关，支持私有云部署，实现API转发、请求参数转换、数据校验、IP黑白名单等功能，帮助减轻网络攻击对后端造成的影响。eoLinker AGW（Lite）提供完全图形化的界面管理，能够快速管理多个API网关。

### 特性
***

1. **免费且开源**：eoLinker AGW（Lite）秉承开源精神，是国内第一个企业开源的API接口网关，为广大的开发、运维以及管理人员提供专业的产品。

2. **鉴权方式**：支持token认证。

3. **访问控制**：可设置的流量控制策略，支持QPS。

4. **请求转发**：默认支持http rest路由。

5. **IP黑白名单**：支持IP白名单、IP黑名单机制。

6. **数据整形**：支持参数的转换与绑定。

7. **动态数据更新**：支持在管理平台内进行配置，服务器不用重启就可直接生效。

8. **UI界面管理**：支持图形化界面管理，方便API网关管理员对业务API进行管理。

9. **快速部署**：支持手动部署与Docker部署。

** 如需更强大的网关功能，欢迎使用eoLinker AGW：**https://github.com/eolinker/eoLinker-API-Gateway.git


### 部署要求
***
* go 1.8及以上版本

* python2.7.x

* redis2.8.x及以上版本

* python拓展库：MySQLdb、ConfigParser、redis

### 快速入门
***

1. [安装指南](http://help.eolinker.com/agw/?target=/md/%E9%83%A8%E7%BD%B2/%E9%83%A8%E7%BD%B2%E6%8C%87%E5%8D%97 "安装指南") 

2. [Docker安装指南](http://help.eolinker.com/agw/?target=/md/%E9%83%A8%E7%BD%B2/Docker%E9%83%A8%E7%BD%B2%E6%8C%87%E5%8D%97 "Docker安装指南")

3. [使用指南](http://help.eolinker.com/agw/?target=/md/index "使用指南")

4. 官方交流Q群：[用户交流1群](https://jq.qq.com/?_wv=1027&k=5ikfC2S)（群号：725853895）

### 相关链接
***
* 开源支持：https://www.eolinker.com/#/os/default#agw

* Github：https://github.com/eolinker/eoLinker-API-Gateway

* Docker：https://hub.docker.com/r/eolinker/eolinker-api-gateway

* 教程文档：http://help.eolinker.com/agw
