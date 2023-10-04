package dboption

import (
	"context"
	"time"

	"github.com/nsplnp/michichong/logutil"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 链接数据库
var (
	client        *mongo.Client
	err           error
	db            *mongo.Database
	collection    *mongo.Collection
	annCollection *mongo.Collection
)

// 1.建立连接
func DbInit() *mongo.Collection {
	log := logutil.GetLog()

	if client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://107.172.86.242:27017").SetConnectTimeout(5*time.Second)); err != nil {
		log.Err(err).Msg("链接数据库失败")
		panic("链接数据库失败")
	}
	//2.选择数据库 my_db
	db = client.Database("michichong")

	//3.选择表 res
	collection = db.Collection("res")

	//4.选择表 ann
	annCollection = db.Collection("ann")

	return collection
}

func GetCollection() *mongo.Collection {
	return collection
}

func GetAnnCollection() *mongo.Collection {
	return annCollection
}

func GetError() error {
	return err
}
