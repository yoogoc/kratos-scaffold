# About

kratos-layout脚手架，可以生成项目，proto，biz，service，data。mono库的生成未完全依照beer-shop，直接生成到app/xx而不是app/xx/service。

# Install

```shell
go install github.com/yoogoc/kratos-scaffold@latest
```

# Command

脚手架的基本格式

```
kratos-scaffold [proto | service | biz | data] [model] [field_name:field_type[, 'a/array/slice'[, 'n,nil,null']:predicate1,predicate2]...
```

- field_type:

|   类型    | go实体类型           |  go参数类型   |       proto实体类型        |    proto参数类型    |    数据库类型     |
|:-------:|------------------|:---------:|:----------------------:|:---------------:|:------------:|
| float64 | float64          | *float64  |         double         | optional double |   numeric    |
| float32 | float32          | *float32  |         float          | optional float  |   numeric    |
|  int32  | int32            |  *int32   |         int32          | optional int32  |     int      |
|  int64  | int64            |  *int64   |         int64          | optional int64  |    bigint    |
| uint32  | uint32           |  *uint32  |         uint32         | optional uint32 |     int      |
| uint64  | uint64           |  *uint64  |         uint64         | optional uint64 |    bigint    |
|  bool   | bool             |   *bool   |          bool          |  optional bool  |   tinyint    |
| string  | string           |  *string  |         string         | optional string | varchar(255) |
|  text   | string           |  *string  |         string         | optional string |     text     |
|  time   | time.Time        | time.Time |         string         | optional string |  timestamp   |
|  date   | time.Time        | time.Time |         string         | optional string |  timestamp   |
|  json   | *structpb.Struct |     x     | google.protobuf.Struct |        x        |    jsonb     |

- field_type: 用于描述数据类型，除了基础类型外，还支持json，array符合类型

- predicate:谓语最终用于sql query时需要的where条件，目前支持：
  - eq 等于
  - cont like
  - gt 大于
  - gte 大于等于
  - lt 小于
  - lte 小于等于
  - in 数组

# Usage

1. 生成项目

- 生成mono主库

```shell
kratos-scaffold new --mono demo
```

 - 生成单体库

```shell
kratos-scaffold new user
```

 - 生成mono子服务库

```shell
kratos-scaffold new user
```

> new 生成单体库与mono子服务库的区别是通过当前目录下是否存在`go.mod`文件来判断

2. 生成proto文件
```shell
kratos-scaffold proto -o api/user/v1/user.proto user id:int64:eq,in name:string:contains age:int32:gte,lte
```

3. 生成biz。可用flag:
- -n --namespace 指定子服务,如果不指定则默认此库为单体库,直接生成到{{project_dir}}/internal/biz目录下

  如果指定了子服务,则会生成到生成到{{project_dir}}/app/{{namespace}}/internal/biz目录下
  (data,service同)

```shell
kratos-scaffold biz -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte
```

4. 生成data
```shell
kratos-scaffold data -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte
```

5. 生成service
```shell
kratos-scaffold service -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte
```

6. 一键生成proto, biz, data, service
```shell
kratos-scaffold g -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte
```

7. 复杂用例字段
```shell
service user id:int64:eq,in name:string:cont age:int32:gte,lte birthday:time:eq,in,gte extra:json last_login_at:time,nil:gte,lte array:string,a -n user
```

extra: json类型
last_login_at: time类型，允许为nil
array: string数组类型

# Roadmap

- [ ] 灵活的生成proto客户端
- [ ] 丰富配置，可以使用配置文件来约定配置，更轻量的使用cli
- [ ] biz，service，data可以通过proto文件生成
- [ ] data: 支持生成proto client和gorm
- [x] proto 生成可以指定proto风格: aa_bb, aaBb, AaBb
- [ ] i18n
- [x] 一次生成biz, service, data
- [ ] 完善文档
