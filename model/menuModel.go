package model

import (
	"go.mongodb.org/mongo-driver/bson"
	"fmt"
	"github.com/codestates/WBABEProject-08/commits/main/util"
	"encoding/json"
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MenuModel struct {
	Client *mongo.Client
	Menucollection *mongo.Collection
}

type Review struct {
	Score int `bson:"score,omitemepty" json:"score"`
	Review string `bson:"review,omitemepty" json:"review"`
	OrderId primitive.ObjectID `bson:"orderid"`
}

type Menu struct {
	ID primitive.ObjectID `bson:"_id,omitempty"`
	IsVisible bool `bson:"isvisible" json:"isvisible"`
	Name string `bson:"name" json:"name"`
	IsOrderable bool `bson:"orderable" json:"orderable"`
	Limit int `bson:"limit" json:"limit"`
	Price int `bson:"price" json:"price"`
	From string `bson:"from" json:"from"`
	Orderedcount int `bson:"orderedcount" json:"orderedcount"`
	Avg float64 `bson:"avg" json:"avg"`
	Suggestion bool `bson:"suggestion" json:"suggestion"`
	Reviews []Review `bson:"reviews,omitemepty" json:"reviews"`
}

func GetMenuModel(db, host, model string) (*MenuModel, error) {
	m := &MenuModel{}
	var err error

	if m.Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(host)); err != nil {
		return nil, err
	} else if err := m.Client.Ping(context.TODO(), nil); err != nil {
		return nil, err
	} else {
		m.Menucollection = m.Client.Database(db).Collection(model)
	}
	return m, nil
}

// DB에 메뉴 data를 추가하는 메서드
func (m *MenuModel) AddMenu(data []byte) (primitive.ObjectID, error) {
	// 새 메뉴를 추가할 때의 default값 생성
	newMenu := &Menu{
		IsVisible : true,
		Orderedcount : 0,
		Avg : 0,
		Suggestion : false,
	}
	err := json.Unmarshal(data, newMenu)
	
	util.ErrorHandler(err)

	result, err := m.Menucollection.InsertOne(context.TODO(), newMenu)
	
	return result.InsertedID.(primitive.ObjectID), err
}


// DB 메뉴 data를 업데이트하는 메서드
func (m *MenuModel) UpdateMenu(data []byte) (interface{}, error) {
	id, key, value := util.GetJsonIdKeyValue(data)
	
	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: key, Value: value}}}}

	result, err := m.Menucollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return result.MatchedCount, err
	} else {
		return result.MatchedCount, nil
	}
}


// DB 메뉴 data를 삭제하는 메서드
func (m *MenuModel) DeleteMenu(menuId primitive.ObjectID) (interface{}, error) {
	filter := bson.D{{Key: "_id", Value: menuId}}
	option := bson.D{{Key: "$set", Value: bson.D{{Key: "isvisible", Value: false}}}}
	result, err := m.Menucollection.UpdateOne(context.TODO(), filter, option)

	if err != nil {
		return nil, err
	} else {
		return result.MatchedCount, nil
	}
}


// 메뉴 리스트를 조회하는 메서드
func (m *MenuModel) GetMenuList(category string, page int64) []Menu {
	menus := []Menu{}
	filter := bson.D{{Key: "isvisible", Value: true}}

	if category == "suggestion" {
		filter = bson.D{{Key: "suggestion", Value: true}, {Key: "isvisible", Value: true}}
	}
	opt := options.Find().SetSort(bson.D{{Key: category, Value: -1}}).SetLimit(5).SetSkip((page - 1) * 5)
	cursor, err := m.Menucollection.Find(context.TODO(), filter, opt)
	util.ErrorHandler(err)

	err = cursor.All(context.TODO(), &menus)
	util.ErrorHandler(err)

	return menus
}


// 메뉴별 평점/리뷰 데이터 추가하는 메서드
func (m *MenuModel) AddReview(sfoodId string, review *Review) {
	foodId := util.ConvertStringToObjectId(sfoodId)

	// 음식에 리뷰 추가하기
	food, _ := m.GetOneMenu(foodId)
	fmt.Println(food)
	food.Reviews = append(food.Reviews, *review)
	filter := bson.D{{Key: "_id", Value: foodId}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "reviews", Value: food.Reviews}}}}
	_, err := m.Menucollection.UpdateOne(context.TODO(), filter, update)
	util.ErrorHandler(err)
}


// 메뉴의 limit을 amount만큼 줄여주고 count를 amount만큼 올려주는 함수
func (m *MenuModel) LimitAndCountUpdate(id primitive.ObjectID, limit, count, amount int) {
	var update bson.D
	filter := bson.D{{Key: "_id", Value: id}}
	if limit - amount > 0 {
		update = bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "limit", Value : limit - amount}, 
				{Key: "orderedcount", Value: count + amount}}}}
	} else {
		update = bson.D{
			{Key: "$set", Value: bson.D{
				{Key: "limit", Value : limit - amount}, 
				{Key: "orderedcount", Value: count + amount}, 
				{Key: "orderable", Value: false}}}}
	}
	_, err := m.Menucollection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		util.ErrorHandler(err)
	}
}


// 메뉴를 하나만 찾아주는 함수
func (m *MenuModel) GetOneMenu(id primitive.ObjectID) (*Menu, error) {
	menu := &Menu{}
	filter := bson.D{{Key: "_id", Value: id}}
	err := m.Menucollection.FindOne(context.TODO(), filter).Decode(menu)
	if err != nil {
		return nil, err
	} else {
		return menu, nil
	}
}


// 음식 평점의 평균을 구해주는 함수
func (m *MenuModel) CalcAvg(sid string) {
	id := util.ConvertStringToObjectId(sid)
	menu, _ := m.GetOneMenu(id)

	avg := 0.0;

	for _, review := range menu.Reviews {
		avg += float64(review.Score)
	}

	avg /= float64(len(menu.Reviews))

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{
		{Key: "$set", Value: bson.D{
			{Key: "avg", Value: avg},
		},
	}}
	_, err := m.Menucollection.UpdateOne(context.TODO(), filter, update)

	util.ErrorHandler(err)
}


// 추천메뉴를 업데이트하는 메서드
func (m *MenuModel) SuggestionUpdate(suggestion *SuggestionType) error {
	// 먼저, 원래 suggestion 되어있던 메뉴들을 초기화해준다.
	filter := bson.D{
		{Key : "suggestion", Value : true},
	}
	update := bson.D{{
		Key : "$set", Value : bson.D{
			{Key: "suggestion", Value : false},
		},
	}}
	_, err := m.Menucollection.UpdateMany(context.TODO(), filter, update)
	util.ErrorHandler(err)

	// 이후 새로운 아이디의 음식들을 업데이트해준다.
	// 배열 핸들링
	filter = bson.D{{Key: "_id", Value: bson.D{{Key: "$in", Value: suggestion.Ids}}}}
	update = bson.D{{
		Key : "$set", Value : bson.D{
			{Key: "suggestion", Value : true},
		},
	}}

	_, err = m.Menucollection.UpdateMany(context.TODO(), filter, update)

	return err
}


// 모든 메뉴를 가져오는 메서드
func (m *MenuModel) GetAllMenu(page int64) []Menu {
	list := []Menu{}
	filter := bson.D{}
	option := options.Find().SetLimit(5).SetSkip((page - 1) * 5)
	cursor, err := m.Menucollection.Find(context.TODO(), filter, option)
	util.ErrorHandler(err)

	err = cursor.All(context.TODO(), &list)
	util.ErrorHandler(err)

	return list
}