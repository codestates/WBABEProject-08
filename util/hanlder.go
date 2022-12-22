package util

import (
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func PanicHandler (err error) {
	if err != nil {
		panic(err)
	}
}

func GetJsonIdKeyValue(data []byte) (primitive.ObjectID, string, interface{}) {
	var unMarshared map[string]interface{}
	json.Unmarshal(data, &unMarshared)

	var id primitive.ObjectID
	var key string
	var value interface{}

	// id로 구분하려면, primitive.ObjectIDFromHex()함수를 사용해 형변환을 해줘야한다.
	if _, exist := unMarshared["id"]; exist {
		id, _ = primitive.ObjectIDFromHex(unMarshared["id"].(string))
	}
	if _, exist := unMarshared["key"]; exist {
		key = unMarshared["key"].(string)
	}
	if _, exist := unMarshared["value"]; exist {
		value = unMarshared["value"]
	}

	return id, key, value
}