# About

kratos-layout脚手架，可以生成项目，proto，biz，service，data。mono库的生成未完全依照beer-shop，直接生成到app/xx而不是app/xx/service。

# Install

```shell
go install github.com/yoogoc/kratos-scaffold@lastest
```

# Command

脚手架的基本格式

```
kratos-scaffold [proto | service | biz | data] [model] [field_name:field_type:predicate1,predicate2]...
```

- field_type: 

| 类型    | go实体类型 | go参数类型             |             proto实体类型             | proto参数类型 |
| :-----: | ------ |:------------------:|:---------------------------------:|:---------:|
| float64 |float64| *wrapperspb.DoubleValue |              double               |     google.protobuf.DoubleValue      |
| float32 |float32| *wrapperspb.FloatValue |               float               |      google.protobuf.FloatValue     |
| int32 |int32| *wrapperspb.Int32Value |               int32               |          google.protobuf.Int32Value |
| int64 |int64| *wrapperspb.Int64Value |               int64               |          google.protobuf.Int64Value |
| uint32  |uint32| *wrapperspb.UInt32Value |              uint32               |      google.protobuf.UInt32Value     |
| uint64  |uint64| *wrapperspb.UInt64Value |              uint64               |      google.protobuf.UInt64Value     |
| bool  |bool| *wrapperspb.BoolValue |               bool                |         google.protobuf.BoolValue  |
| string|string| *wrapperspb.StringValue |              string               |      google.protobuf.StringValue     |
| time |time.Time| time.Time          |       google.protobuf.Timestamp   |  google.protobuf.Timestamp         |


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

# Roadmap

- [ ] 灵活的生成proto客户端
- [ ] 丰富配置，可以使用配置文件来约定配置，更轻量的使用cli
- [ ] biz，service，data可以通过proto文件生成
- [ ] data: 支持生成proto client和gorm
- [x] proto 生成可以指定proto风格: aa_bb, aaBb, AaBb
- [ ] i18n
- [ ] 一次生成biz, service, data
- [ ] 完善文档
