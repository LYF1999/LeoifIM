package routers

import (
	"LeoifIM/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/api/im/", &controllers.WebSocketController{})
	beego.Router("/api/message/", &controllers.WebSocketController{}, "get:Join")
}
