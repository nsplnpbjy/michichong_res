package main

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
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
	ress := []*Res{}
	results, findErr := collection.Find(context.Background(), bson.M{})
	if findErr != nil {
		return
	}
	for results.Next(context.TODO()) {
		temp := Res{}
		decodeErr := results.Decode(&temp)
		if decodeErr != nil {
			continue
		}
		ress = append(ress, &temp)
	}
	c.JSON(http.StatusOK, gin.H{
		"list": ress,
	})
	results.Close(context.TODO())

}

// 根据团体名返回模糊查询预约
func DoGetSpeRes(c *gin.Context) {
	var data Res
	if err := c.ShouldBind(&data); err != nil {
		c.Error(err)
		return
	}
	ress := []*Res{}
	results, findErr := collection.Find(context.Background(), bson.M{"groupname": bson.M{"$regex": data.GroupName}})
	if findErr != nil {
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
	results.Close(context.TODO())
}

// 插入
func DoInsertRes(c *gin.Context) bool {
	if !NilCheck(c) {
		c.Error(err)
	}
	results, findErr := collection.Find(context.Background(), bson.M{"tourtime": bson.M{"$regex": c.PostForm("TourTime")}})
	if findErr != nil {
		c.Error(findErr)
		return false
	}
	if results.Next(context.TODO()) {
		return false
	}
	res := Res{c.PostForm("TourTime"), c.PostForm("GroupName"), c.PostForm("ComuName"), c.PostForm("ComuPhoneNumber"), c.PostForm("ResTime"), false}
	if _, err = collection.InsertOne(context.TODO(), res); err != nil {
		log.Println("插入失败:", err)
		return false
	}
	results.Close(context.TODO())
	return true
}

// 确认已参观
func Done(c *gin.Context) bool {
	if !NilCheck(c) {
		return false
	} else {
		_, modierr := collection.UpdateOne(context.TODO(), bson.M{"tourtime": bson.M{"$eq": c.PostForm("TourTime")}}, bson.M{"$set": bson.M{"isdone": true}})
		return modierr == nil
	}
}

// 删除预约
func DoDeleteRes(c *gin.Context) bool {
	if c.PostForm("TourTime") == "200000000000" {
		_, delerr := collection.DeleteMany(context.TODO(), bson.M{})
		if delerr != nil {
			log.Println("删除失败", delerr)
			return false
		}
		return true
	}
	inputTime := c.PostForm("TourTime")
	_, delerr := collection.DeleteMany(context.TODO(), bson.M{"tourtime": inputTime})
	if delerr != nil {
		log.Println("删除失败", delerr)
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
