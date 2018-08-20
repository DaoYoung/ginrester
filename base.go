package gorester

import (
	"github.com/gin-gonic/gin"
	"errors"
	"strconv"
)

type BaseController struct{}

func (this BaseController) Rester() {
	panic(errors.New("can't find func:Rester in your controller"))
}
func (this BaseController) modelSlice() interface{} {
	panic(errors.New("can't find func:modelSlice in your controller"))
}
func (this *BaseController) model() ResourceInterface {
	panic(errors.New("can't find func:model in your controller"))
}
func (this *BaseController) parentController() ControllerInterface                           { return nil }
func (this *BaseController) beforeCreate(c *gin.Context, m ResourceInterface)                     {}
func (this *BaseController) afterCreate(c *gin.Context, m ResourceInterface)                      {}
func (this *BaseController) beforeUpdate(c *gin.Context, old ResourceInterface, new ResourceInterface) {}
func (this *BaseController) afterUpdate(c *gin.Context, old ResourceInterface, new ResourceInterface)  {}
func (this *BaseController) beforeDelete(c *gin.Context, m ResourceInterface, id int)             {}
func (this *BaseController) afterDelete(c *gin.Context, m ResourceInterface, id int)              {}
func (this *BaseController) listCondition(c *gin.Context) map[string]interface{} {
	return make(map[string]interface{})
}
func (this *BaseController) updateCondition(c *gin.Context, pk string) map[string]interface{} {
	condition := make(map[string]interface{})
	id, err := strconv.Atoi(c.Param(pk))
	if err != nil {
		panic(errors.New("can't Update without ID"))
	}
	condition["id"] = id
	return condition
}
