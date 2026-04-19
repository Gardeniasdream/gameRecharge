# 游戏充值接口demo README

# 游戏充值服务

基于 Go \+ Gin \+ PostgreSQL \+ Redis \+ RocketMQ 实现的游戏钻石充值系统，支持订单创建、支付回调、异步发钻、订单取消、自动关单，保障幂等安全、数据一致性，可直接部署上线。

## 一、技术栈

- Web 框架：Gin

- ORM 框架：GORM

- 数据库：PostgreSQL（存储订单、流水、幂等数据）

- 缓存：Redis（分布式锁，防止重复下单）

- 消息队列：RocketMQ（异步处理支付回调、30分钟自动关单）

- 配置管理：\.env 环境变量

## 二、核心功能

- 订单管理：创建待支付订单、查询订单状态、手动取消未支付订单

- 支付处理：接收支付平台回调、异步写入支付流水

- 钻石发放：异步发放钻石、记录发放流水、防重复发放

- 自动关单：30分钟未支付订单自动取消（RocketMQ 延迟消息）

- 安全保障：Redis 防重复下单、PostgreSQL 幂等校验、数据库事务一致性

## 三、核心流程

1. 前端调用接口创建订单，Redis 锁防止重复提交，生成待支付订单（状态 0）

2. 系统返回支付链接，用户完成支付

3. 支付平台回调接口，发送 RocketMQ 消息异步处理

4. 消费者监听消息：幂等校验 → 更新订单为已支付（状态 1）→ 写入支付流水 → 写入钻石流水 → 调用游戏服发钻

5. 未支付订单：手动调用取消接口，或30分钟后自动取消（状态 2）

## 四、接口说明

### 1\. 创建充值订单

**请求方式**：POST /recharge/create

**Content\-Type**：application/json

**请求参数**：

```json
{
  "user_id": 10001,       // 用户ID（必填）
  "role_id": "role_666888",// 游戏角色ID（必填）
  "server_id": 101,        // 游戏服务器ID（必填）
  "product_id": "1",       // 充值商品ID（必填）
  "pay_type": 1            // 支付渠道（1=微信，2=支付宝，必填）
}
```

**返回示例**：

```json
{
  "code": 0,
  "msg": "下单成功",
  "order_no": "GM1776590908444908",
  "pay_url": "https://xxx/pay"
}
```

### 2\. 支付回调（模拟/真实）

**请求方式**：POST /pay/notify

**Content\-Type**：application/x\-www\-form\-urlencoded

**请求参数**：

```text
out_trade_no=GM1776590908444908  // 订单号（必填）
transaction_id=TEST888888888888  // 支付平台交易号（必填）
```

**返回要求**：成功返回 \&\#34;success\&\#34;（支付平台要求）

## 五、部署步骤

1. 环境准备：安装 PostgreSQL、Redis、RocketMQ（JDK 1\.8\+）

2. 配置修改：复制 \.env 模板，修改对应数据库、Redis、RocketMQ 地址

3. 依赖安装：go mod tidy

4. 启动服务：go run main.go（程序自动通过 GORM 生成数据表）

## 六、数据表说明

|表名|说明|核心字段|
|---|---|---|
|game\_order|订单主表|order\_no（订单号）、user\_id、role\_id、order\_status（状态）、price、diamond|
|pay\_log|支付流水表（资金记录）|order\_no、transaction\_id（交易号）、pay\_channel、amount|
|reward\_flow|钻石发放流水表（资产记录）|order\_no、role\_id、server\_id、diamond、status（发放状态）|
|idempotent|幂等表（防重复处理）|unique\_key（支付交易号）、created\_at|

## 七、安全说明

- 重复下单：Redis 分布式锁 \+ 订单号唯一索引，防止用户重复提交

- 重复发放：PostgreSQL 幂等表（交易号唯一索引），避免支付回调重复处理

- 数据一致：数据库事务保证订单、支付流水、钻石流水、幂等记录原子性（要么全成功，要么全失败）

- 异步安全：RocketMQ 消息可重试、可堆积，避免支付回调阻塞

- 订单安全：仅待支付订单可取消，已支付订单禁止取消，避免资损
