package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nsplnp/michichong/dboption"
	"go.mongodb.org/mongo-driver/bson"
)

type Res struct {
	TourTime        string `bson:"tourtime"`
	GroupName       string `bson:"groupname"`
	ComuName        string `bson:"comuname"`
	ComuPhoneNumber string `bson:"comuphonenumber"`
	ResTime         string `bson:"restime"`
	IsDone          bool   `bson:"isdone"`
}

// 返回全部预约
func DoGetRes(c *gin.Context) {
	selectCount := 0
	collection := dboption.GetCollection()
	ress := []*Res{}
	results, findErr := collection.Find(context.Background(), bson.M{})
	defer results.Close(context.TODO())
	if findErr != nil {
		log.Println("全数据查询失败")
		return
	}
	for results.Next(context.TODO()) {
		temp := Res{}
		decodeErr := results.Decode(&temp)
		if decodeErr != nil {
			continue
		}
		ress = append(ress, &temp)
		selectCount++
	}
	c.JSON(http.StatusOK, gin.H{
		"list": ress,
	})
	log.Println("查询方式：全数据   查询记录:" + string(rune(selectCount)))
}

// 根据团体名返回模糊查询预约
func DoGetSpeRes(c *gin.Context) {
	collection := dboption.GetCollection()
	var data Res
	if err := c.ShouldBind(&data); err != nil {
		c.Error(err)
		return
	}
	ress := []*Res{}
	results, findErr := collection.Find(context.Background(), bson.M{"groupname": bson.M{"$regex": data.GroupName}})
	defer results.Close(context.TODO())
	if findErr != nil {
		log.Println("查询：" + data.GroupName + "团体失败")
		c.Error(findErr)
		return
	}
	for results.Next(context.TODO()) {
		temp := Res{}
		if decodeErr := results.Decode(&temp); decodeErr != nil {
			c.Error(decodeErr)
			continue
		}
		ress = append(ress, &temp)
	}
	c.JSON(http.StatusOK, gin.H{
		"list": ress,
	})
	log.Println("查询：" + data.GroupName)
}

// 插入
func DoInsertRes(c *gin.Context) bool {
	collection := dboption.GetCollection()
	err := dboption.GetError()
	if !NilCheck(c) {
		c.Error(err)
	}
	results, findErr := collection.Find(context.Background(), bson.M{"tourtime": bson.M{"$regex": c.PostForm("TourTime")}})
	if findErr != nil {
		log.Println("插入失败:" + c.PostForm("GroupName"))
		c.Error(findErr)
		return false
	}
	if results.Next(context.TODO()) {
		log.Println("插入失败:" + c.PostForm("GroupName") + "  时间已被预约")
		return false
	}
	res := Res{c.PostForm("TourTime"), c.PostForm("GroupName"), c.PostForm("ComuName"), c.PostForm("ComuPhoneNumber"), c.PostForm("ResTime"), false}
	if _, err = collection.InsertOne(context.TODO(), res); err != nil {
		log.Println("插入失败:" + c.PostForm("GroupName"))
		return false
	}
	results.Close(context.TODO())
	log.Println("插入成功：" + res.GroupName)
	return true
}

// 确认已参观
func Done(c *gin.Context) bool {
	collection := dboption.GetCollection()
	if !NilCheck(c) {
		return false
	} else {
		results, seleErr := collection.Find(context.TODO(), bson.M{"tourtime": bson.M{"$eq": c.PostForm("TourTime")}})
		if seleErr != nil {
			log.Println("确认参观失败:" + c.PostForm("TourTime"))
			return false
		}
		if !results.Next(context.TODO()) {
			log.Println("没有预约记录：" + c.PostForm("TourTime"))
			return false
		} else {
			temp := Res{}
			if decErr := results.Decode(&temp); decErr != nil {
				log.Println("预约记录解码失败：" + c.PostForm("TourTime"))
				return false
			}
			if temp.IsDone {
				log.Println("确认参观失败，参观已经完成：" + c.PostForm("TourTime"))
				return false
			}
		}
		_, modierr := collection.UpdateOne(context.TODO(), bson.M{"tourtime": bson.M{"$eq": c.PostForm("TourTime")}}, bson.M{"$set": bson.M{"isdone": true}})
		if modierr != nil {
			log.Println("确认参观失败:" + c.PostForm("TourTime"))
			return false
		} else {
			log.Println("确认参观:" + c.PostForm("TourTime"))
			return true
		}
	}
}

// 删除预约
func DoDeleteRes(c *gin.Context) bool {
	collection := dboption.GetCollection()
	if c.PostForm("TourTime") == "200000000000" {
		_, delerr := collection.DeleteMany(context.TODO(), bson.M{})
		if delerr != nil {
			log.Println("全数据删除失败")
			return false
		}
		return true
	}
	inputTime := c.PostForm("TourTime")
	_, delerr := collection.DeleteMany(context.TODO(), bson.M{"tourtime": inputTime})
	if delerr != nil {
		log.Println("删除失败:" + inputTime)
		return false
	}
	return true
}

func NilCheck(c *gin.Context) bool {
	tourTime := c.PostForm("TourTime")
	if tourTime == "" {
		return false
	}
	groupName := c.PostForm("GroupName")
	if groupName == "" {
		return false
	}
	comuName := c.PostForm("ComuName")
	if comuName == "" {
		return false
	}
	comuPhoneNumber := c.PostForm("ComuPhoneNumber")
	if comuPhoneNumber == "" {
		return false
	}
	resTime := c.PostForm("ResTime")
	if resTime == "" {
		return false
	}
	return true
}
