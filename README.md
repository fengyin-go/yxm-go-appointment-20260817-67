# 预约挂号系统（Appointment）

一个纯 Go 标准库实现的预约挂号 REST API 服务，采用标准 Go 工程目录结构，内存存储，零第三方依赖。

## 目录结构

```
origin/
├── cmd/server/          # 程序入口
├── internal/
│   ├── app/             # 依赖装配
│   ├── config/          # 配置加载
│   ├── model/           # 领域模型与校验
│   ├── store/           # 数据访问接口 + 内存实现
│   ├── service/         # 业务逻辑层
│   └── handler/         # HTTP 处理器层
└── pkg/
    ├── httpx/           # HTTP 响应工具
    ├── idgen/           # ID 生成
    └── logger/          # 分级日志
```

## 运行

```bash
go run ./cmd/server
PORT=8081 go run ./cmd/server
```

默认监听 `:8080`。

## 测试

```bash
go test ./...
```

## 预约状态机

`booked → checked_in → completed`，`booked → cancelled`（取消自动释放号源）。

## API 接口

### 科室

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/departments` | 创建科室 |
| GET | `/api/departments` | 科室列表 |
| GET | `/api/departments/{id}` | 科室详情 |
| PUT | `/api/departments/{id}` | 更新科室 |
| DELETE | `/api/departments/{id}` | 删除科室 |

### 医生

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/doctors` | 创建医生 |
| GET | `/api/doctors?department_id=` | 医生列表 |
| GET | `/api/doctors/{id}` | 医生详情 |
| PUT | `/api/doctors/{id}` | 更新医生 |
| DELETE | `/api/doctors/{id}` | 删除医生 |

### 排班

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/schedules` | 创建排班 `{"doctor_id","date":"2006-01-02","time_slot","total_slots"}` |
| GET | `/api/schedules?doctor_id=&date=&page=&size=` | 排班列表 |
| GET | `/api/schedules/{id}` | 排班详情 |
| DELETE | `/api/schedules/{id}` | 删除排班 |

### 患者

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/patients` | 创建患者 |
| GET | `/api/patients` | 患者列表 |
| GET | `/api/patients/{id}` | 患者详情 |
| PUT | `/api/patients/{id}` | 更新患者 |
| DELETE | `/api/patients/{id}` | 删除患者 |

### 预约

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/api/appointments` | 预约挂号 `{"schedule_id","patient_id"}` |
| GET | `/api/appointments?patient_id=&doctor_id=&schedule_id=&status=&page=&size=` | 预约列表 |
| GET | `/api/appointments/{id}` | 预约详情 |
| PATCH | `/api/appointments/{id}/status` | 状态流转 |
| POST | `/api/appointments/{id}/check-in` | 签到 |
| POST | `/api/appointments/{id}/cancel` | 取消（释放号源） |
| POST | `/api/appointments/{id}/complete` | 完成就诊 |
| DELETE | `/api/appointments/{id}` | 删除预约 |

### 统计

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/stats/overview` | 全局总览 |
| GET | `/api/stats/appointments` | 预约状态分布 |
| GET | `/api/stats/doctors` | 医生排班汇总 |
