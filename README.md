## eoLinker AGW（eoLinker-API Gateway开源版）

### 简介
***
**eoLinker AGW 是eoLinker旗下的开源版API Gateway，同时也是国内唯一的Go语言开源API Gateway，性能优异、功能强大，提高API业务安全性。**

通过Go语言开发的高性能API网关，支持私有云部署，帮助企业对外封装API，实现API转发、请求参数转换、数据校验、请求过滤等功能，帮助减轻网络攻击对后端造成的影响。eoLinker AGW 提供完全图形化界面的网关管理网站，能够快速管理多个API网关。

### 特性
***
1. **免费且开源**：eoLinker-API Gateway 秉承开源精神，是国内第一个企业开源的API接口网关，为广大的开发、运维以及管理人员提供专业的产品。

2. **多种鉴权方式**：支持Basic 认证、API Key授权、IP认证、放行认证等方式。

3. **支持Open API**：不同账户拥有独立的访问密钥。

4. **权限管理**：可设置不同策略组设置流量控制策略，包括访问QPS、访问总次数、访问IP、访问时间段等

5. **请求转发**：默认支持http rest路由。

6. **IP黑白名单**：支持用户的IP白名单、黑名单机制。

7. **数据整形**：支持参数的转换与绑定。

8. **性能控制**：包括超时设置、熔断设置等。

9. **插件系统**：支持丰富的插件系统，能够自由搭配满足不同业务需求，如数据缓存、负载均衡等。

10. **动态数据更新**：API、组件等都支持在管理平台进行配置，服务器不用重启就可直接生效。

11. **UI界面管理**：eoLinker AGW 拥有清晰的监控与操作界面，方便API网关管理员了解系统主要运行情况，对API进行管理。

12. **快速部署**：支持手动部署与Docker部署。

13. **兼容eoLinker-AMS**：可与国内最大的接口管理平台打通。

### 部署要求
***
* go 1.8及以上版本

* python2.7.x

* redis2.8.x及以上版本

* python拓展库：MySQLdb、ConfigParser、redis

### 相关链接
***
* 开源支持：https://www.eolinker.com/#/os/default#agw

* Github：https://github.com/eolinker/eoLinker-API-Gateway

* Docker：https://hub.docker.com/r/eolinker/eolinker-api-gateway

* 教程文档：http://help.eolinker.com/agw

* 官方交流Q群：[用户交流1群](https://jq.qq.com/?_wv=1027&k=5ikfC2S)
