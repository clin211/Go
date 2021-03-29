package mongo_connection

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
)

type Student struct {
	Name string `bson: "name"`
	Age int `bson: "age"`
	Sid string `bson: "sid"`
	Status bool `bson: "status"`
}

type Person struct {
	person []Student
}

func main() {
	// 创建MongoDB连接
	mongo, err := mgo.Dial("mongodb://119.3.48.150:27000/test")
	if err != nil {
		log.Fatal(err)
		return
	}
	// 最后执行关闭操作
	defer mongo.Close()
	fmt.Println("connect to mongodb success")
	
	// 设置集合
	client := mongo.DB("test").C("cc_student")
	
	// 数据
	data := Student{
		Name: "小王",
		Age: 18,
		Sid: "s24",
		Status: true,
	}
	
	// 插入文档
	err = client.Insert(&data)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("insert document success")
	
	// 根据条件查找一条记录
	student := Student{}
	err = client.Find(bson.M{"sid": "s24"}).One(&student)
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(student)
	fmt.Println("find 1 document success")
	
	// 根据条件更新文档
	err = client.Update(bson.M{"status": true}, bson.M{"$set": bson.M{"age": 19}})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("update document success")
	
	// 根据条件返回所有记录
	persons := Person{}
	iter := client.Find(bson.M{"status": true}).Sort("_id").Skip(0).Limit(15).Iter()
	for iter.Next(&student) {
		persons.person = append(persons.person, student)
	}
	if err = iter.Close(); err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println(persons)
	fmt.Println("find all document sucess")
	
	// 删除文档
	err = client.Remove(bson.M{"sid": "s24"})
	if err != nil {
		log.Fatal(err)
		return
	}
	fmt.Println("delete document success")
}