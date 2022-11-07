package api

import (
	"encoding/json"
	"fmt"
	"vpn/center/serializer"
)

func ErrorResponse(err error) serializer.Response {
	//if ve, ok := err.(validator.ValidationErrors); ok {
	//	for _, e := range ve {
	//		field := conf.T(fmt.Sprintf("Field.%s", e.Field))
	//		tag := conf.T(fmt.Sprintf("Tag.Valid.%s", e.Tag))
	//		return serializer.Response{
	//			Code:  -1,
	//			Msg:   fmt.Sprintf("%s%s", field, tag),
	//			Error: fmt.Sprint(err),
	//		}
	//	}
	//}
	if _, ok := err.(*json.UnmarshalTypeError); ok {
		return serializer.Response{
			Code:  -1,
			Msg:   "JSON类型不匹配",
			Error: fmt.Sprint(err),
		}
	}
	return serializer.Response{
		Code:  -1,
		Msg:   "参数错误",
		Error: fmt.Sprint(err),
	}
}
