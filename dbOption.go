package main

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 链接数据库
var (
	client     *mongo.Client
	err        error
	db         *mongo.Database
	collection *mongo.Collection
)

// 1.建立连接
func DbInit() {
	defer func() {
		if dbErr := recover(); dbErr != nil {
			fmt.Println(dbErr)
		}
	}()
	if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://107.172.86.242:27017").SetConnectTimeout(5*time.Second)); err != nil {
		panic("链接数据库失败")
	}
	//2.选择数据库 my_db
	db = client.Database("michichong")

	//3.选择表 my_collection
	collection = db.Collection("res")
}
