package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// 메뉴를 추가할 때 body에 해당하는 구조체
type AddMenuStruct struct {
	NewOrder OrderedList
	NewItem primitive.ObjectID
	OrderId primitive.ObjectID
}