package main

import (
	"context"
	"log"

	"github.com/shiimoo/godb/mgo"
)

func main() {
	ctx := context.Background()
	if err := mgo.Connect(ctx, mgo.NewConnCfg("127.0.0.1", 27017)); err != nil {
		log.Fatalln(err)
	}
	log.Println("数据库连接成功!")
	mgo.SetDatabase("test") // 指定链接数据库

	log.Println("判定数据库[test]中的集合[account]是否存在!", mgo.HasCollection("test1"))
}
