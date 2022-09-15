# 工具库

## 单词转换

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
```