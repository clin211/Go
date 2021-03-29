package connection

import (
	"context"
	"fmt"
	"log"

	"github.com/mongodb/mongo-go-driver"
	"go.mongodb.org/mongo-driver"
)

func main() {
	connection()
	fmt.Println("程序执行了！！！！")
}

//定义student结构,变量大写
type student struct {
	Name string
	Age  int
}

func connection() {
	// 设置客户端连接配置
	clientOptions := options.Client().ApplyURI("mongodb://119.3.48.150:27017")

	// 连接到MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// 检查连接
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}
