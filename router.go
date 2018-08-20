package ginrester

import (
	"github.com/gin-gonic/gin"
)

func CreateRoutes(group *gin.RouterGroup, controller ControllerInterface, actions ...string) {
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
func BuildRoute(controller ControllerInterface) (path, resourceName, routeId string) {
	if controller.ParentNode() != nil {
		pp, pr, pi := BuildRoute(controller.ParentNode())
		path = pp + "/" + pr + "/:" + pi
	}
	resourceName = controller.RouteName() + "s"
	routeId = GetRouteID(controller)
	return
}



