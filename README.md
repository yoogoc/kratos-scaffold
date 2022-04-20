# Usage
1. 生成proto文件
```shell
kratos-scaffold proto -o api/user/v1/user.proto user id:int64:eq,in name:string:contains age:int32:gte,lte
# 暂时是要kratos或直接生成proto client
kratos proto client api/user/v1/user.proto
```
2. 生成biz。可用flag:
- -n --namespace 指定子服务,如果不指定则默认此库为单体库,直接生成到{{project_dir}}/internal/biz目录下

  如果指定了子服务,则会生成到生成到{{project_dir}}/app/{{namespace}}/internal/biz目录下
  (data,service同)

```shell
kratos-scaffold biz -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte
```
3. 生成data
```shell
kratos-scaffold data -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte
```
4. 生成service
```shell
kratos-scaffold service -n user-service user id:int64:eq,in name:string:contains age:int32:gte,lte
```

5. 生成项目

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

# Roadmap

## New Feature
- [ ] 在kratos proto client 的基础上更加灵活的生成proto客户端
- [ ] 丰富配置，可以使用配置文件来约定配置，更轻量的使用cli
- [ ] biz，service，data可以通过proto文件生成
- [ ] data: 支持生成proto client和gorm
- [ ] 支持更多类型, 目前只有字符型和数字型, 需支持time等动态类型,如time类型在proto中会体现为timestamp类型,bool和数字类型的eq谓语需要使用wrappers类型包装
- [ ] proto 生成可以指定proto风格: aa_bb, aaBb, AaBb

# 已知问题

- [ ] 如果命名为xxx.com/xx/xx会有部分模块无法正常生成
