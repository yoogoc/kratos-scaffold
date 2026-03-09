---
name: kratos-scaffold-usage
description: kratos-scaffold 脚手架使用指南。当用户询问如何使用 kratos-scaffold 工具、生成 kratos-layout 项目结构、生成 proto/biz/data/service 文件，或询问脚手架命令和字段类型定义时使用此 skill。
---

# kratos-scaffold 使用指南

kratos-scaffold 是一个 kratos-layout 风格的脚手架工具，用于快速生成项目结构、proto 文件、biz 层、data 层和 service 层代码。

## 安装

```shell
go install github.com/yoogoc/kratos-scaffold@latest
```

## 命令格式

```
kratos-scaffold [command] [model] [field_name:field_type[, 'a/array/slice'[, 'n,nil,null']:predicate1,predicate2]...
```

## 支持的命令

### 1. new - 生成项目

生成 mono 主库（多服务仓库）：
```shell
kratos-scaffold new --mono demo
```

生成单体库或 mono 子服务库：
```shell
kratos-scaffold new user
```

> 区别：通过当前目录下是否存在 `go.mod` 文件来判断是单体库还是 mono 子服务库。

### 2. proto - 生成 proto 文件

```shell
kratos-scaffold proto -o api/user/v1/user.proto user id:int64:eq,in name:string:cont age:int32:gte,lte
```

参数说明：
- `-o, --output`: 指定输出路径
- `--http`: 同时生成 xx.http.pb.go

### 3. biz - 生成业务逻辑层

```shell
kratos-scaffold biz -n user-service user id:int64:eq,in name:string:cont age:int32:gte,lte
```

参数说明：
- `-n, --namespace`: 指定子服务名称
  - 不指定：生成到 `{{project_dir}}/internal/biz`
  - 指定：生成到 `{{project_dir}}/app/{{namespace}}/internal/biz`

### 4. data - 生成数据访问层

```shell
kratos-scaffold data -n user-service user id:int64:eq,in name:string:cont age:int32:gte,lte
```

### 5. service - 生成服务层

```shell
kratos-scaffold service -n user-service user id:int64:eq,in name:string:cont age:int32:gte,lte
```

### 6. g (generate) - 一键生成所有层

同时生成 proto、biz、data、service：

```shell
kratos-scaffold g -n user-service user id:int64:eq,in name:string:cont age:int32:gte,lte
```

## 字段类型定义

### 基础类型

| 类型 | Go 实体类型 | Go 参数类型 | Proto 实体类型 | Proto 参数类型 | 数据库类型 |
|:----:|:-----------:|:-----------:|:--------------:|:--------------:|:----------:|
| float64 | float64 | *float64 | double | optional double | numeric |
| float32 | float32 | *float32 | float | optional float | numeric |
| int32 | int32 | *int32 | int32 | optional int32 | int |
| int64 | int64 | *int64 | int64 | optional int64 | bigint |
| uint32 | uint32 | *uint32 | uint32 | optional uint32 | int |
| uint64 | uint64 | *uint64 | uint64 | optional uint64 | bigint |
| bool | bool | *bool | bool | optional bool | tinyint |
| string | string | *string | string | optional string | varchar(255) |
| text | string | *string | string | optional string | text |
| time | time.Time | time.Time | string | optional string | timestamp |
| date | time.Time | time.Time | string | optional string | timestamp |
| json | *structpb.Struct | x | google.protobuf.Struct | x | jsonb |

### 数组类型

在字段类型后添加 `,a` 或 `,array` 或 `,slice` 表示数组：

```shell
# string 数组
kratos-scaffold proto user tags:string,a

# int64 数组
kratos-scaffold proto user ids:int64,a
```

### 可空类型

在字段类型后添加 `,n` 或 `,nil` 或 `,null` 表示可空：

```shell
# 可空的 time 类型
kratos-scaffold proto user deleted_at:time,nil
```

## 谓词（Predicate）

用于生成 SQL 查询条件，支持以下操作：

- `eq` - 等于 (=)
- `cont` - 包含/模糊查询 (like)
- `gt` - 大于 (>)
- `gte` - 大于等于 (>=)
- `lt` - 小于 (<)
- `lte` - 小于等于 (<=)
- `in` - 在数组中 (in)

多个谓词用逗号分隔：

```shell
kratos-scaffold proto user id:int64:eq,in name:string:cont age:int32:gte,lte
```

## 完整示例

```shell
# 复杂字段示例
kratos-scaffold service user \
  id:int64:eq,in \
  name:string:cont \
  age:int32:gte,lte \
  birthday:time:eq,in,gte \
  extra:json \
  last_login_at:time,nil:gte,lte \
  tags:string,a \
  -n user
```

字段说明：
- `id`: int64 类型，支持 eq 和 in 查询
- `name`: string 类型，支持 contains 模糊查询
- `age`: int32 类型，支持 gte 和 lte 范围查询
- `birthday`: time 类型，支持 eq、in、gte 查询
- `extra`: json 类型
- `last_login_at`: 可空的 time 类型，支持 gte 和 lte 查询
- `tags`: string 数组类型

## 输出路径规则

### 单体库结构

```
project/
├── api/              # proto 文件
├── internal/
│   ├── biz/         # 业务逻辑
│   ├── data/        # 数据访问
│   └── service/     # 服务层
```

### Mono 仓库结构

```
project/
├── api/              # proto 文件
└── app/
    └── {{namespace}}/
        └── internal/
            ├── biz/  # 业务逻辑
            ├── data/ # 数据访问
            └── service/ # 服务层
```

## 常用命令速查

| 场景 | 命令 |
|------|------|
| 生成单体项目 | `kratos-scaffold new myproject` |
| 生成 mono 主库 | `kratos-scaffold new --mono myproject` |
| 生成 proto | `kratos-scaffold proto -o api/v1/user.proto user id:int64 name:string` |
| 生成 biz | `kratos-scaffold biz user id:int64 name:string` |
| 生成 data | `kratos-scaffold data user id:int64 name:string` |
| 生成 service | `kratos-scaffold service user id:int64 name:string` |
| 一键生成所有 | `kratos-scaffold g user id:int64 name:string` |
| 指定命名空间 | `kratos-scaffold g -n user-service user id:int64 name:string` |
