package AppController

import(
	"github.com/astaxie/beego"
	"app_log_beego_config/model"
	"github.com/astaxie/beego/logs"
	"strings"
)

type AppController struct{
	beego.Controller
}

func (p *AppController) AppList() {
	logs.Debug("enter index page success...")
	p.Layout = "layout/layout.html"
	appinfos,err := model.GetAllAppInfo()
	if (err != nil) {
		logs.Error("model.GetAllAppInfo() failed", err)
		p.TplName = "app/error.html"
		p.Data["Error"] = "获取数据库资料异常，请稍后再试"
		return
	}
	p.Data["applist"] = appinfos
	p.TplName = "app/index.html"
}

func (p *AppController) AppApply()  {
	logs.Debug("enter apply page success...")
	p.Layout = "layout/layout.html"
	p.TplName = "app/apply.html"
}

func (p *AppController) AppCreate()  {
	logs.Debug("enter AppCreate page success...")
	p.Layout = "layout/layout.html"
	appinfo := &model.AppInfo{}
	appinfo.AppName = p.GetString("app_name")
	appinfo.AppType = p.GetString("app_type")
	appinfo.DevelopPath = p.GetString("develop_path")
	appinfo.IPs = strings.Split(p.GetString("iplist"), ",")
	if (len(appinfo.AppName) == 0 || len(appinfo.AppType) == 0 || len(appinfo.DevelopPath) == 0 || len(appinfo.IPs) == 0) {
		logs.Error("app create params error")
		p.TplName = "app/error.html"
		p.Data["Error"] = "参数异常，请检查后重新输入"
		return
	}
	err := model.CreateAppInfo(appinfo)
	if (err != nil) {
		p.TplName = "app/error.html"
		p.Data["Error"] = "创建appinfo异常，请稍后再试"
		return
	}
	p.Redirect("/index", 302)
}