package ginrester

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"net/http"
	"reflect"
	"strings"
	"errors"
	"github.com/jinzhu/gorm"
	"net/url"
)
var Db *gorm.DB
var PerPage =20

type ControllerInterface interface {
	Update(c *gin.Context)
	Create(c *gin.Context)
	Info(c *gin.Context)
	List(c *gin.Context)
	Delete(c *gin.Context)
	IsRestRoutePk() bool //false id
	RouteName() string //rewrite resource name in route url
	ParentNode() ControllerInterface

	Init(r ControllerInterface)
	model() ResourceInterface
	modelSlice() interface{}
	parentController() ControllerInterface
	beforeDelete(c *gin.Context, m ResourceInterface, id int)
	afterDelete(c *gin.Context, m ResourceInterface, id int)
	beforeCreate(c *gin.Context, m ResourceInterface)
	afterCreate(c *gin.Context, m ResourceInterface)
	beforeUpdate(c *gin.Context, old ResourceInterface, new ResourceInterface)
	updateCondition(c *gin.Context, pk string) map[string]interface{}
	afterUpdate(c *gin.Context, old ResourceInterface, new ResourceInterface)
	listCondition(c *gin.Context) map[string]interface{}
}

type Controller struct {
	ParentController ControllerInterface
	Rester           ControllerInterface
	RestModel        func() ResourceInterface
	RestModelSlice   func() interface{} //https://golang.org/doc/faq#convert_slice_of_interface
	*BaseController
}
func (action *Controller) Init(r ControllerInterface){
	if r == nil {
		panic(errors.New("param r: is not a controller"))
	}
	action.Rester = r
	action.RestModel = r.model
	action.RestModelSlice = r.modelSlice
	action.ParentController = r.parentController()
}
func (action *Controller) ParentNode() ControllerInterface {
	return action.ParentController
}
func (action *Controller) IsRestRoutePk() bool {
	return false
}
func (action *Controller) RouteName() string {
	obj := action.RestModel()
	f := reflect.TypeOf(obj)
	if f.Kind() == reflect.Ptr {
		f = f.Elem()
	}
	return strings.ToLower(f.Name())
}

func (action *Controller) Create(c *gin.Context) {
	obj := action.RestModel()
	err := c.BindJSON(obj)
	if err != nil {
		panic(err)
	}
	action.Rester.beforeCreate(c, obj)
	info := Create(obj)
	action.Rester.afterCreate(c, info)
	ReturnSuccess(c, http.StatusCreated, info)
}

func (action *Controller) Info(c *gin.Context) {
	obj := action.RestModel()
	id, _ := strconv.Atoi(c.Param(GetRouteID(action.Rester)))
	info := FindOneByID(obj, id)
	ReturnSuccess(c, http.StatusOK, info)
}

func (action *Controller) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		panic(err)
	}
	obj := action.RestModelSlice()
	condition := action.Rester.listCondition(c)
	MergeUrlCondition(condition, c.Request.URL.Query(), action.RestModel())
	FindListByMap(obj, condition, "id desc", page, PerPage)
	ReturnSuccess(c, http.StatusOK, obj)
}

func (action *Controller) Update(c *gin.Context) {
	obj := action.RestModel()
	err := c.BindJSON(obj)
	if err != nil {
		panic(err)
	}
	condition := action.Rester.updateCondition(c, GetRouteID(action.Rester))
	if val, ok := condition["id"]; ok {
		old := FindOneByID(action.RestModel(), val.(int))
		CheckUpdateCondition(old, condition)
		action.Rester.beforeUpdate(c, old, obj)
		info := UpdateByID(val.(int), obj)
		action.Rester.afterUpdate(c, old, info)
		ReturnSuccess(c, http.StatusOK, info)
	}else {
		panic(errors.New("can't find data to update"))
	}
}
func (action *Controller) Delete(c *gin.Context) {
	obj := action.RestModel()
	id, _ := strconv.Atoi(c.Param(GetRouteID(action.Rester)))
	action.Rester.beforeDelete(c,obj, id)
	DeleteByID(obj, id)
	action.Rester.afterDelete(c,obj, id)
	ReturnSuccess(c, http.StatusOK, gin.H{"id": id})
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
func CheckUpdateCondition(m ResourceInterface, condition map[string]interface{}) {
	v := reflect.ValueOf(m)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	for key, val := range condition {
		old := v.FieldByName(camelString(key))
		switch old.Kind() {
		case reflect.String:
			if old.String() != val {
				panic(errors.New("forbid update by field:" + key))
			}
			break
		case reflect.Int:
			if old.Int() != int64(val.(int)) {
				panic(errors.New("forbid update by field:" + key))
			}
			break
		default:
			panic(errors.New("forbid update by field type:" + old.Kind().String()))
		}
	}
}