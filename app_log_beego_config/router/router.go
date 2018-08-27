package router

import(
	"github.com/astaxie/beego"
	"app_log/app_log_beego_config/controller/AppController"
	"app_log/app_log_beego_config/controller/LogController"
)

func init()  {
	// app
	beego.Router("/index", &AppController.AppController{}, "*:AppList")
	beego.Router("/app/apply", &AppController.AppController{}, "*:AppApply")
	beego.Router("/app/create", &AppController.AppController{}, "*:AppCreate")
	beego.Router("/app/list", &AppController.AppController{}, "*:AppList")
	
	// log
	beego.Router("/log/apply", &LogController.LogController{}, "*:LogApply")
	beego.Router("/log/create", &LogController.LogController{}, "*:LogCreate")
	beego.Router("/log/list", &LogController.LogController{}, "*:LogList")
	beego.Router("/log/index", &LogController.LogController{}, "*:LogList")
	beego.Router("/log/delete_apply", &LogController.LogController{}, "*:LogDeleteApply")
	beego.Router("/log/delete", &LogController.LogController{}, "*:LogDelete")
}