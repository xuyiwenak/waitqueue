package queue

import (
	"reflect"
	"strings"
)

type Queue struct {
	queue chan interface{}
}

// 初始化一个队列

func NewQueue(len int32) *Queue{
	return &Queue{queue: make(chan interface{}, len)}
}

func (q * Queue) QPOP()  interface{}{
	if q.QLEN() != 0{
		return <-q.queue
	}
	return nil
}

func (q * Queue) QPUSH(in interface{}){
	q.queue <- in
}
// 返回队列的长度
func(q * Queue) QLEN() int{
	return len(q.queue)
}


func StructToMap(s interface{},key string) map[string]interface{} {
	var ret = make(map[string]interface{})
	modelReflect := reflect.ValueOf(s)

	if modelReflect.Kind() == reflect.Ptr {
		modelReflect = modelReflect.Elem()
	}

	modelRefType := modelReflect.Type()
	fieldsCount := modelReflect.NumField()

	var fieldData interface{}

	for i := 0; i < fieldsCount; i++ {
		field := modelReflect.Field(i)
		switch field.Kind() {
		case reflect.Struct:
			fieldData = StructToMap(field.Interface(),key)
		case reflect.String:
			fieldData = field.Interface().(string)
		case reflect.Int:
			fieldData = field.Interface().(int)
		case reflect.Float64:
			fieldData = field.Interface().(float64)
		default:
			fieldData = field.Interface()
		}
		ret[strings.ToLower(modelRefType.Field(i).Tag.Get(key))] = fieldData
	}
	return ret
}
//只保存有效数据即 整形>0;字符串类型不为空
func StructToMapOnlyValidData(st interface{}) map[string]interface{} {
	var info = make(map[string]interface{},0)
	modelReflect := reflect.ValueOf(st)
	if modelReflect.Kind() == reflect.Ptr {
		modelReflect = modelReflect.Elem()
	}

	modelRefType := modelReflect.Type()
	fieldsCount := modelReflect.NumField()

	for i := 0; i < fieldsCount; i++ {
		field := modelReflect.Field(i)
		switch field.Kind() {
		case reflect.String:
			value:= field.Interface().(string)
			if value != ""{
				info[modelRefType.Field(i).Tag.Get("json")] = value
			}
		case reflect.Int64:
			value := field.Interface().(int64)
			if value != 0 {
				info[modelRefType.Field(i).Tag.Get("json")] = value
			}
		case reflect.Int:
			value := field.Interface().(int)
			if value != 0 {
				info[modelRefType.Field(i).Tag.Get("json")] = value
			}
		case reflect.Float64:
			value := field.Interface().(float64)
			if value != 0 {
				info[modelRefType.Field(i).Tag.Get("json")] = value
			}
		case reflect.Float32:
			value := field.Interface().(float32)
			if value != 0{
				info[modelRefType.Field(i).Tag.Get("json")] = value
			}
		case reflect.Bool:
			value := field.Interface().(bool)
			if value != false {
				info[modelRefType.Field(i).Tag.Get("json")] = value
			}
		}
	}
	return info
}