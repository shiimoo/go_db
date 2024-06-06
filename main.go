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
	log.Println("connect succ")
}
