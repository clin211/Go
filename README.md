# miniblog project

## 数据表生成结构体

- db2struct 工具

```shell
db2struct --gorm --json -H 127.0.0.1 -d miniblog -t 表名 --package 包名 --struct 结构体名 -u 数据库用户名 -p '数据库密码' --target=生成的文件名.go
```
