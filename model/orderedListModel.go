package model

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderedListModel struct {
	Client *mongo.Client
	Collection *mongo.Collection
}

type BuyerInfo struct {
	Pnum int `bson:"pnum" json:"pnum"`
	Address string `bson:"address" json:"address"`
}

type OrderedList struct {
	Id primitive.ObjectID `bson:"id" json:"id"`
	IsReviewed bool `bson:"id" json:"id"`
	Status string `bson:"status" json:"status"`
	BuyerInfo BuyerInfo `bson:"buyerinfo" json:"buyerinfo"`
	Orderedmenus []primitive.ObjectID `bson:"orderedmenus"`
}

func GetOrderedListModel(host, db, model string) (*OrderedListModel, error) {
	m := &OrderedListModel{}
	var err error
	
	if m.Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(host)); err != nil {
		
	}

}