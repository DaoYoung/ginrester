package gorester

import (
	"github.com/gin-gonic/gin"
	"strings"
)

func CreateRestRoutes(group *gin.RouterGroup, controller ControllerInterface, actions ...string) {
	path, resourceName, routeId := BuildRoute(controller)
	if isContain(actions, "list") {
		group.GET(path+"/"+resourceName, controller.List)
	}
	if isContain(actions, "info") {
		group.GET(path+"/"+resourceName+"/:"+routeId, controller.Info)
	}
	if isContain(actions, "update") {
		group.PUT(path+"/"+resourceName+"/:"+routeId, controller.Update)
	}
	if isContain(actions, "create") {
		group.POST(path+"/"+resourceName, controller.Create)
	}
	if isContain(actions, "delete") {
		group.DELETE(path+"/"+resourceName+"/:"+routeId, controller.Delete)
	}
}

func isContain(strSlice []string, searchStr string) bool {
	if len(strSlice) == 0 {
		return true
	}
	str := strings.Join(strSlice,",")
	return strings.Contains(str, searchStr)
}


