package main

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 获取数据库配置
	dsn := "gorm:gorm123456@tcp(127.0.0.1:3306)/gormgen?charset=utf8mb4&parseTime=True&loc=Local"

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		fmt.Println(`failed to connect database:`, err)
		panic(err)
	}

	// 创建生成器实例
	g := gen.NewGenerator(gen.Config{
		// 输出路径
		OutPath: "./dal/query",
		// 输出模式
		Mode: gen.WithDefaultQuery | gen.WithQueryInterface | gen.WithoutContext,
		// 表字段可为空值时，对应结构体字段使用指针类型
		FieldNullable: true,
		// 生成字段类型标签
		FieldWithTypeTag: true,
		// 生成字段DB标签
		FieldWithIndexTag: true,
	})

	// 使用数据库
	g.UseDB(db)

	// 生成所有表的模型和查询代码
	// 也可以用 g.GenerateModel("users") 指定表
	g.ApplyBasic(g.GenerateAllTable()...)

	// 执行代码生成
	g.Execute()
}
