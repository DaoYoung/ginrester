package gorester

import (
	"github.com/gin-gonic/gin"
	"errors"
	"strconv"
	"reflect"
	"strings"
)

type EmptyController struct{}

func (this EmptyController) Rester() {
	panic(NOContentError(errors.New("can't find func:Rester in your controller")))
}
func (this EmptyController) modelSlice() interface{} {
	panic(NOContentError(errors.New("can't find func:modelSlice in your controller")))
}
func (this *EmptyController) model() ResourceInterface {
	panic(NOContentError(errors.New("can't find func:model in your controller")))
}
func (this *EmptyController) parentController() ControllerInterface                           { return nil }
func (this *EmptyController) beforeCreate(c *gin.Context, m ResourceInterface)                     {}
func (this *EmptyController) afterCreate(c *gin.Context, m ResourceInterface)                      {}
func (this *EmptyController) beforeUpdate(c *gin.Context, old ResourceInterface, new ResourceInterface) {}
func (this *EmptyController) afterUpdate(c *gin.Context, old ResourceInterface, new ResourceInterface)  {}
func (this *EmptyController) beforeDelete(c *gin.Context, m ResourceInterface, id int)             {}
func (this *EmptyController) afterDelete(c *gin.Context, m ResourceInterface, id int)              {}
func (this *EmptyController) listCondition(c *gin.Context) map[string]interface{} {
	return make(map[string]interface{})
}
func (this *EmptyController) updateCondition(c *gin.Context, pk string) map[string]interface{} {
	condition := make(map[string]interface{})
	id, err := strconv.Atoi(c.Param(pk))
	if err != nil {
		panic(NOContentError(errors.New("can't Update without ID")))
	}
	condition["id"] = id
	return condition
}

func BuildRoute(controller ControllerInterface) (path, resourceName, routeId string) {
	if controller.ParentNode() != nil {
		pp, pr, pi := BuildRoute(controller.ParentNode())
		path = pp + "/" + pr + "/:" + pi
	}
	resourceName = controller.RouteName() + "s"
	routeId = GetRouteID(controller)
	return
}
func GetRouteID(controller ControllerInterface) (routeId string) {
	routeId = "id"
	if controller.IsRestRoutePk() {
		routeId = controller.RouteName() + "_id"
	}
	return
}
func CheckupdateCondition(m ResourceInterface, condition map[string]interface{}) {
	v := reflect.ValueOf(m)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for key, val := range condition {
		old := v.FieldByName(camelString(key))
		switch old.Kind() {
		case reflect.String:
			if old.String() != val {
				panic(ForbidError(errors.New("forbid update by field:" + key)))
			}
			break
		case reflect.Int:
			if old.Int() != int64(val.(int)) {
				panic(ForbidError(errors.New("forbid update by field:" + key)))
			}
			break
		default:
			panic(ForbidError(errors.New("forbid update by field type:" + old.Kind().String())))
		}
	}
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
