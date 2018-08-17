package gorester

import (
	"reflect"
	"github.com/gin-gonic/gin"
	"net/url"
	"strings"
)



type JsonSuccess struct {
	Data interface{} `json:"data"`
}

func ReturnSuccess(c *gin.Context, code int, data interface{}) {
	js := new(JsonSuccess)
	js.Data = data
	c.JSON(code, js)
}

func MergeUrlCondition(condition map[string]interface{}, query url.Values, obj interface{}){
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
		t = t.Elem()
	}
	for i := 0; i < v.NumField(); i++ {
		f := v.Field(i)
		if f.Kind() == reflect.Struct {
			for j := 0; j < f.NumField(); j++ {
				nm := snakeString(t.FieldByIndex([]int{i, j}).Name)
				if p := query.Get(nm);p != ""{
					condition[nm] = p
				}
			}
			continue
		}
		s := t.Field(i)
		nm := snakeString(s.Name)
		if p := query.Get(nm);p != ""{
			condition[nm] = p
		}
	}
}

// 驼峰式写法转为下划线写法
func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	num := len(s)
	if num<3 {
		return strings.ToLower(s)
	}
	for i := 0; i < num; i++ {
		d := s[i]
		if i > 0 && d >= 'A' && d <= 'Z'  {
			data = append(data, '_')
		}

		data = append(data, d)
	}
	return strings.ToLower(string(data[:]))
}


