# 工具库

## 单词转换
>
> 该子命令支持各种单词格式转换，模式如下：
1: 全部转换为大写
2: 全部转换为小写
3: 下划线单词转为大写驼峰单词
4: 下划线单词转为小写驼峰单词
5: 驼峰单词转为下划线单词

Usage:
  clin word [flags]

Flags:
  -h, --help         help for word
  -m, --mode int8    请输入单词转换的模式：
  -s, --str string   请输入单词内容

```go
const (
 ModeUpper                      = iota + 1 // 全部转大写
 ModeLower                                 // 全部转小写
 ModeCamelCaseToUnderscore                 // 驼峰转下划线
 ModeUnderscoreToUpperCamelCase            // 下划线转大写驼峰
 ModeUnderscoreToLowerCamelCase            // 下线线转小写驼峰
)
```

### 使用示例

```shell
$ go run main.go word -s=word -m=1
2022/09/15 00:30:41 result: WORD


$ go run main.go word -s=WORD -m=2
2022/09/15 00:31:12 result: word


$ go run main.go word -s=keyWord -m=3
2022/09/15 00:34:58 result: key_word


$ go run main.go word -s=key_word -m=4
2022/09/15 00:44:29 result: KeyWord


$ go run main.go word -s=key_word -m=4
2022/09/15 00:40:27 result: keyWord
```

## 时间转换
>
> 时间格式处理

Usage:
  clin time [flags]
  clin time [command]

Available Commands:
  calc        计算所需时间
  format      时间戳格式化工具
  now         获取当前时间

Flags:
  -h, --help   help for time

### 使用示例

```shell
# 获取当前时间
$ go run main.go time now
2022/09/15 08:32:31 result: 2022-09-15 08:32:31, 1663201951


# 传入一个时间，获取指定分钟后的时间；比如：获取5分钟的时间
$ go run main.go time calc -c="2022-09-15 08:32:31" -d=5m
2022/09/15 08:38:43 result: 2022-09-15 08:37:31, 1663231051


# 传入一个时间，获取指定小时后的时间；比如：获取2小时前的时间
$ go run main.go time calc -c="2022-09-15 08:32:31" -d=-2h
2022/09/15 08:43:02 result: 2022-09-15 06:32:31, 1663223551

# 传入一个时间戳，格式化成当前时区的标注时间
$ go run main.go time format 1697513074411
2023/10/17 11:53:10 格式化后的时间: 2023-10-17 11:24:34
```

## SQL 语句转结构体

### 使用示例

> sql 转换和处理

Usage:
  clin sql [flags]
  clin sql [command]

Available Commands:
  struct      sql 转换

Flags:
  -h, --help   help for sql

Use "clin sql [command] --help" for more information about a command.

```shell
$ go run main.go sql struct --username root --password 123456 --host localhost:3306 --database service --table article
type Article struct {
 // id
 Id     int32   `json:id`
 // 文章标题
 Title  string  `json:title`
 // 文章简述
 Desc   string  `json:desc`
 // 封面图片地址
 CoverImageUrl  string  `json:cover_image_url`
 // 文章内容
 Content        string  `json:content`
 // 新建时间
 CreatedOn      int32   `json:created_on`
 // 创建人
 CreatedBy      string  `json:created_by`
 // 修改时间
 ModifiedOn     int32   `json:modified_on`
 // 修改人
 ModifiedBy     string  `json:modified_by`
 // 删除时间
 DeletedOn      int32   `json:deleted_on`
 // 是否删除 0为未删除、1为已删除
 IsDel  int8    `json:is_del`
 // 状态 0为禁用、1为启用
 State  int8    `json:state`
}
func (model Article) TableName() string {
 return "article"
}
```
