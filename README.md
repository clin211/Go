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
