# miniblog project

## 数据表生成结构体

- db2struct 工具

```shell
db2struct --gorm --json -H 127.0.0.1 -d miniblog -t 表名 --package 包名 --struct 结构体名 -u 数据库用户名 -p '数据库密码' --target=生成的文件名.go
```

示例：
```shell
$ mkdir -p internal/pkg/model
$ cd internal/pkg/model
$ db2struct --gorm --no-json -H 127.0.0.1:13306 -d miniblog -t user --package model --struct UserM -u miniblog -p 'miniblog1234' --target=user.go
$ db2struct --gorm --no-json -H 127.0.0.1:13306 -d miniblog -t post --package model --struct PostM -u miniblog -p 'miniblog1234' --target=post.go

```