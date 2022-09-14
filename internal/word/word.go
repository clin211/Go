package word

import (
	"strings"
	"unicode"
)

// 将全部单词转换为大写
func ToUpper(s string) string {
	return strings.ToUpper(s)
}

// 将单词全部转为小写
func ToLower(s string) string {
	return strings.ToLower(s)
}

// 下划线单词转大写驼峰单词
func UnderscoreToUpperCamelCase(s string) string {
	// 将下划线替换成空格
	s = strings.Replace(s, "_", " ", -1)
	// 将所有的首字母改为大写形式
	s = strings.Title(s)
	// 将空格替换为空
	return strings.Replace(s, " ", "", -1)
}

// 下划线单词转为小写驼峰单词
func UnderscoreToLowerCamelCase(s string) string {
	// 通过大写驼峰转换方法获取数据
	s = UnderscoreToUpperCamelCase(s)
	// 对其首字母进行小写处理
	return string(unicode.ToLower(rune(s[0]))) + s[1:]
}

// 驼峰单词转下划线单词
func CamelCaseToUnderscore(s string) string {
	var output []rune
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
		}

		if unicode.IsUpper(r) {
			output = append(output, '_')
		}
		output = append(output, unicode.ToLower(r))
	}

	return string(output)[1:]
}
