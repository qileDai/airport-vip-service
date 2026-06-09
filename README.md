# 机场贵宾厅预约权益核验批量对账与审计服务

## 项目简介

本项目是一个纯后端 REST API 服务，专注于机场贵宾厅预约权益核验、批量对账与审计功能。目标用户包括权益运营、机场服务台、客服支持等角色，帮助他们完成核验权益资格、安排候补入场、校验同行人规则和复盘使用率等业务操作。

## 技术栈

- **语言**: Go 1.21+
- **Web框架**: Echo v4
- **数据库**: SQLite
- **架构**: 分层架构 (Models -> Repositories -> Services -> Routes)

## 核心功能模块

1. **会员权益数据模型与 CRUD API** - 管理会员权益信息、配额、有效期
2. **预约记录批量导入预检 API** - 支持批量导入预约记录并进行预检验证
3. **核验权益资格领域服务** - 校验会员权益、预约记录、航班时段一致性
4. **安排候补入场计算服务** - 根据优先级算法安排候补入场
5. **权益过期异常判定 API** - 自动检测过期权益并生成异常事件
6. **权益核验通过率聚合统计 API** - 按日/批次/角色聚合统计
7. **航班时段归档与快照 API** - 归档已完成航班并生成快照
8. **入场优先级计算 API** - 根据会员等级、剩余额度、等待时间计算优先级
9. **审计日志和操作追踪 API** - 记录所有操作并支持查询
10. **接口测试脚本与 HTTP 样本** - 提供完整的API测试样本

## 核心数据实体

- **会员权益 (MemberBenefit)**: 编号、会员姓名、等级、配额、有效期、状态等
- **预约记录 (ReservationRecord)**: 编号、航班信息、VIP厅、同行人数量、状态等
- **航班时段 (FlightSchedule)**: 航班号、起降机场、时间、VIP厅容量等
- **同行人 (Companion)**: 姓名、证件、关系、VIP资格等
- **使用凭证 (UsageVoucher)**: 凭证号、二维码、有效期、使用状态等
- **候补名单 (WaitlistEntry)**: 会员信息、优先级分数、等待时间等
- **核验结果 (VerificationResult)**: 核验类型、结果、失败原因等
- **状态流转记录 (StatusTransition)**: 实体类型、状态变更、操作人等
- **规则配置 (RuleConfig)**: 规则名称、类型、阈值、适用等级等
- **异常事件 (ExceptionEvent)**: 事件类型、触发字段、处理状态等

## 领域业务规则

1. **入场优先级计算**: 根据会员等级(40/30/20/10分)、权益剩余额度(最高20分)、候补等待时间(最高40分)计算综合优先级
2. **权益过期异常处理**: 权益过期时自动进入异常队列，记录触发字段、阈值、处理人和处理时限
3. **核验一致性校验**: 核验权益资格前必须校验会员权益、预约记录和航班时段之间的一致性，缺少关键字段时只允许保存草稿
4. **聚合统计追溯**: 权益核验通过率、候补转入率和入场使用率按日/批次/责任角色聚合，支持查看明细来源
5. **状态流转限制**: 状态流转必须限制下一步动作，驳回必须填写原因

## 快速开始

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 初始化数据库并导入种子数据

数据库会在首次运行时自动初始化，种子数据也会自动导入。

### 3. 启动服务

```bash
go run src/main.go
```

服务默认运行在 `http://localhost:8080`

### 4. 验证服务

```bash
curl http://localhost:8080/api/v1/health
```

预期响应:
```json
{
  "status": "ok",
  "version": "1.0.0"
}
```

## API 端点列表

### 健康检查
- `GET /api/v1/health` - 服务健康检查

### 会员权益
- `GET /api/v1/member-benefits` - 获取会员权益列表
- `GET /api/v1/member-benefits/:id` - 获取单个会员权益
- `POST /api/v1/member-benefits` - 创建会员权益
- `PUT /api/v1/member-benefits/:id` - 更新会员权益
- `DELETE /api/v1/member-benefits/:id` - 删除会员权益
- `POST /api/v1/member-benefits/check-expiry` - 检查权益过期

### 预约记录
- `GET /api/v1/reservations` - 获取预约记录列表
- `GET /api/v1/reservations/:id` - 获取单个预约记录
- `POST /api/v1/reservations` - 创建预约记录
- `PUT /api/v1/reservations/:id` - 更新预约记录
- `DELETE /api/v1/reservations/:id` - 删除预约记录
- `POST /api/v1/reservations/batch-import/preview` - 批量导入预检
- `POST /api/v1/reservations/batch-import` - 执行批量导入
- `POST /api/v1/reservations/:id/status` - 状态流转
- `GET /api/v1/reservations/:id/transitions` - 获取允许的状态流转

### 航班时段
- `GET /api/v1/flight-schedules` - 获取航班时段列表
- `GET /api/v1/flight-schedules/:id` - 获取单个航班时段
- `POST /api/v1/flight-schedules` - 创建航班时段
- `PUT /api/v1/flight-schedules/:id` - 更新航班时段
- `DELETE /api/v1/flight-schedules/:id` - 删除航班时段
- `POST /api/v1/flight-schedules/:id/archive` - 归档航班时段
- `POST /api/v1/flight-schedules/batch-archive` - 批量归档

### 候补名单
- `GET /api/v1/waitlist` - 获取候补名单列表
- `GET /api/v1/waitlist/:id` - 获取单个候补记录
- `POST /api/v1/waitlist` - 创建候补记录
- `PUT /api/v1/waitlist/:id` - 更新候补记录
- `DELETE /api/v1/waitlist/:id` - 删除候补记录
- `POST /api/v1/waitlist/calculate-priority` - 计算入场优先级
- `POST /api/v1/waitlist/arrange` - 安排候补入场

### 核验结果
- `GET /api/v1/verifications` - 获取核验结果列表
- `GET /api/v1/verifications/:id` - 获取单个核验结果
- `POST /api/v1/verifications` - 创建核验结果
- `POST /api/v1/verifications/verify-eligibility` - 核验权益资格

### 异常事件
- `GET /api/v1/exceptions` - 获取异常事件列表
- `GET /api/v1/exceptions/:id` - 获取单个异常事件
- `POST /api/v1/exceptions` - 创建异常事件
- `POST /api/v1/exceptions/handle` - 处理异常事件
- `GET /api/v1/exceptions/open` - 获取待处理异常事件

### 审计日志
- `GET /api/v1/audit-logs` - 获取审计日志列表
- `GET /api/v1/audit-logs/:id` - 获取单个审计日志
- `GET /api/v1/audit-logs/entity/:type/:id` - 获取实体相关审计日志

### 统计
- `GET /api/v1/statistics/verification` - 核验通过率统计
- `GET /api/v1/statistics/waitlist` - 候补转入率统计
- `GET /api/v1/statistics/usage` - 入场使用率统计

## 运行测试

```bash
go test ./tests/... -v
```

## 使用 HTTP 样本文件

项目提供了 `tests/http_examples.http` 文件，包含所有 API 的请求样本。可以使用 VS Code 的 REST Client 插件或 IntelliJ IDEA 的 HTTP Client 来执行这些请求。

## 数据库文件

## 构建项目

```bash
go build -o airport-vip-service ./src
```
