package gorester

import (
	"github.com/gin-gonic/gin"
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
func isContain(strSlice []string, searchStr string) bool {
	if len(strSlice) == 0 {
		return true
	}
	str := strings.Join(strSlice,",")
	return strings.Contains(str, searchStr)
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
func GetRouteID(controller ControllerInterface) (routeId string) {
	routeId = "id"
	if controller.IsRestRoutePk() {
		routeId = controller.RouteName() + "_id"
	}
	return
}

// 下划线写法转为驼峰写法
func camelString(s string) string {
	data := make([]byte, 0, len(s))
	j := false
	k := false
	num := len(s) - 1
	if num < 2 {
		return strings.ToUpper(s)
	}
	for i := 0; i <= num; i++ {
		d := s[i]
		if k == false && d >= 'A' && d <= 'Z' {
			k = true
		}
		if d >= 'a' && d <= 'z' && (j || k == false) {
			d = d - 32
			j = false
			k = true
		}
		if k && d == '_' && num > i && s[i+1] >= 'a' && s[i+1] <= 'z' {
			j = true
			continue
		}
		data = append(data, d)
	}
	return string(data[:])
}

