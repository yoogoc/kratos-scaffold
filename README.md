# About

kratos-layout脚手架，可以生成项目，proto，biz，service，data。mono库的生成未完全依照beer-shop，直接生成到app/xx而不是app/xx/service。

## Features

- 生成完整项目（mono主库、mono子服务、单体库、BFF）并自带 demo Greeter 服务，开箱即用
- 生成 biz/service/data 时自动将构造器注入到对应的 `wire.NewSet` ProviderSet
- 支持排序（`order_by`），采用 Go Option 可变参数模式
- 时间类型：gRPC 使用 `google.protobuf.Timestamp`，HTTP 使用 `int64`（unix timestamp）
- CRUD 方法：Create、Update、Patch（部分更新）、DestroyBy、List、FindBy、ExistBy
- SubMono 项目的 wire/proto 编译使用 Makefile target（`make wire-xx`、`make proto-xx`）

# Install

```shell
go install github.com/yoogoc/kratos-scaffold@latest
```

# Command

脚手架的基本格式

```
kratos-scaffold [proto | service | biz | data] [model] [field_name:field_type[, 'a/array/slice'[, 'n,nil,null']:predicate1,predicate2]...
```

## field_type

### gRPC 模式（默认）

|   类型    | go实体类型           |  go参数类型    |          proto实体类型          |       proto参数类型        |    数据库类型     |
|:-------:|------------------|:----------:|:---------------------------:|:----------------------:|:------------:|
| float64 | float64          |  *float64  |           double            |    optional double     |   numeric    |
| float32 | float32          |  *float32  |           float             |    optional float      |   numeric    |
|  int32  | int32            |   *int32   |           int32             |    optional int32      |     int      |
|  int64  | int64            |   *int64   |           int64             |    optional int64      |    bigint    |
| uint32  | uint32           |  *uint32   |           uint32            |    optional uint32     |     int      |
| uint64  | uint64           |  *uint64   |           uint64            |    optional uint64     |    bigint    |
|  bool   | bool             |   *bool    |           bool              |    optional bool       |   tinyint    |
| string  | string           |  *string   |           string            |    optional string     | varchar(255) |
|  text   | string           |  *string   |           string            |    optional string     |     text     |
|  time   | time.Time        | *time.Time | google.protobuf.Timestamp   | google.protobuf.Timestamp |  timestamp   |
|  date   | time.Time        | *time.Time | google.protobuf.Timestamp   | google.protobuf.Timestamp |  timestamp   |
|  json   | *structpb.Struct |     x      |   google.protobuf.Struct    |           x            |    jsonb     |

### HTTP 模式（`--http` flag）

时间字段在 HTTP 模式下使用 `int64`（unix timestamp）：

|  类型  | proto实体类型 | proto参数类型    |
|:----:|:--------:|:------------:|
| time |  int64   | optional int64 |
| date |  int64   | optional int64 |

## predicate

谓语最终用于sql query时需要的where条件，目前支持：
  - eq 等于
  - cont like
  - gt 大于
  - gte 大于等于
  - lt 小于
  - lte 小于等于
  - in 数组

## 生成的 API 方法

| 方法 | 说明 | HTTP | gRPC |
|------|------|------|------|
| Create | 创建 | POST | rpc Create |
| Update | 全量更新 | PUT | rpc Update |
| Patch | 部分更新（专用 Patch 结构体，非 id 字段均为指针） | PATCH | rpc Patch |
| Destroy/DestroyBy | 删除 | DELETE（HTTP）/ rpc DestroyBy（gRPC） | rpc DestroyBy |
| Get/FindBy | 查询单条 | GET（HTTP）/ rpc FindBy（gRPC） | rpc FindBy |
| ExistBy | 判断是否存在 | - | rpc ExistBy |
| List | 分页列表（支持过滤、排序） | GET | rpc List |

### List 排序

List 接口支持 `order_by` 参数，格式为 `"field1 desc, field2 asc"`：

- **Proto**: `string order_by = 4;` 在 `ListXxxRequest` 中
- **Biz 层**: 使用 Go Option 模式 `...ListOption`，内部 apply 为 `*ListOptions`
- **Data 层**: `applyListOptions(opts)` 通过 `Modify` 统一处理分页和排序

# Usage

## 1. 生成项目

生成的项目自带 demo Greeter 服务（proto + biz + service + data），wire 注入完整可直接运行。

### 生成 mono 主库

```shell
kratos-scaffold new --mono demo
```

### 生成单体库

```shell
kratos-scaffold new user
```

### 生成 mono 子服务库

```shell
# 在 mono 主库目录下执行
kratos-scaffold new user
```

> new 生成单体库与mono子服务库的区别是通过当前目录下是否存在`go.mod`文件来判断

### 生成 BFF 服务

```shell
# 在 mono 主库目录下执行
kratos-scaffold new --bff bff
```

BFF 项目仅生成 HTTP server，data 层使用 gRPC client 调用下游服务。

## 2. 生成 proto 文件
```shell
kratos-scaffold proto -o api/user/v1/user.proto user id:int64:eq,in name:string:contains age:int32:gte,lte
```

HTTP 模式（生成 HTTP annotation）：
```shell
kratos-scaffold proto --http user id:int64:eq,in name:string:contains age:int32:gte,lte
```

## 3. 生成 biz

可用 flag:
- `-n --namespace` 指定子服务,如果不指定则默认此库为单体库

```shell
kratos-scaffold biz -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte
```

> 生成后会自动将 `NewUserUsecase` 注入到 `biz.go` 的 `ProviderSet` 中

## 4. 生成 data
```shell
kratos-scaffold data -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte
```

> 生成后会自动将 `NewUserRepo` 注入到 `data.go` 的 `ProviderSet` 中

## 5. 生成 service
```shell
kratos-scaffold service -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte
```

> 生成后会自动将 `NewUserService` 注入到 `service.go` 的 `ProviderSet` 中

## 6. 一键生成 proto, biz, data, service
```shell
kratos-scaffold g -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte
```

## 7. 复杂用例字段
```shell
service user id:int64:eq,in name:string:cont age:int32:gte,lte birthday:time:eq,in,gte extra:json last_login_at:time,nil:gte,lte array:string,a -n user
```

- extra: json类型
- last_login_at: time类型，允许为nil
- array: string数组类型

# 依赖

生成的项目包含以下核心依赖：

- `github.com/go-kratos/kratos/v2` — Kratos 框架
- `entgo.io/ent` — ORM
- `github.com/google/wire` — 依赖注入
- `github.com/samber/lo` — Go 工具库（用于类型转换）
- `go.opentelemetry.io/otel` — OpenTelemetry 链路追踪
- `go.uber.org/zap` — 日志
- `github.com/spf13/cobra` — CLI

# Roadmap

- [ ] 灵活的生成proto客户端
- [ ] 丰富配置，可以使用配置文件来约定配置，更轻量的使用cli
- [ ] biz，service，data可以通过proto文件生成
- [ ] data: 支持生成gorm
- [x] proto 生成可以指定proto风格: aa_bb, aaBb, AaBb
- [ ] i18n
- [x] 一次生成biz, service, data
- [x] 生成项目时自带 demo 服务
- [x] 自动注入 ProviderSet
- [x] List 排序支持
- [x] Patch 部分更新
- [x] ExistBy 存在性检查
- [x] gRPC 时间类型使用 protobuf Timestamp
- [x] HTTP 时间类型使用 int64
- [ ] 完善文档
