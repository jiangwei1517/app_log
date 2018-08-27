package LogController

import(
	"github.com/astaxie/beego"
	"app_log_beego_config/model"
	"github.com/astaxie/beego/logs"
)

type LogController struct{
	beego.Controller
}

func (p *LogController) LogApply()  {
	logs.Debug("enter LogApply page success...")
	p.Layout = "layout/layout.html"
	p.TplName = "log/apply.html"
}

func (p *LogController) LogCreate()  {
	logs.Debug("enter log create page...")
	p.Layout = "layout/layout.html"
	loginfo := &model.LogInfo{}
	loginfo.AppName = p.GetString("app_name")
	loginfo.LogPath = p.GetString("log_path")
	loginfo.Topic = p.GetString("topic")
	if (len(loginfo.AppName) == 0 || len(loginfo.LogPath) == 0 || len(loginfo.Topic) == 0) {
		logs.Error("LogCreate 参数异常，请重新输入..")
		p.TplName = "log/error.html"
		p.Data["Error"] = "参数异常，请重新输入..."
		return
	} 
	err := model.LogCreate(loginfo)
	if (err != nil) {
		logs.Error("LogCreate failed LogCreate,", err)
		return
	}
	logs.Debug("Log create success...")
	p.Redirect("/log/index", 302)
}

func (p *LogController) LogList()  {
	logs.Debug("enter log list page success...")
	p.Layout = "layout/layout.html"
	logList,err:=model.GetAllLogInfo()
	if (err != nil) {
		logs.Error("model.GetAllLogInfo error,", err)
		p.TplName = "log/error.html"
		p.Data["Error"] = "创建数据库失败，请稍后再试...s"
		return
	}
	logs.Debug("model.GetAllLogInfo success...,")
	p.Data["loglist"] = logList
	p.TplName = "log/index.html"
}

func (p *LogController) LogDelete()  {
	p.Layout = "layout/layout.html"
	info := &model.LogInfo{}
	info.AppName = p.GetString("app_name")
	info.LogPath = p.GetString("log_path")
	info.Topic = p.GetString("topic")
	if (len(info.AppName) == 0 || len(info.LogPath) == 0 || len(info.Topic) == 0) {
		logs.Error("参数错误，请稍后再试")
		p.TplName = "log/error.html"
		p.Data["Error"] = "参数错误，请稍后再试"
		return
	}
	err := model.LogDelete(info)
	if (err != nil) {
		logs.Error("model.LogDelete(info) error,", err)
		p.TplName = "log/error.html"
		p.Data["Error"] = "数据库异常， 请稍后再试"
		return
	}
	p.Redirect("/log/index", 302)
}

func (p *LogController) LogDeleteApply()  {
	logs.Debug("enter delete  page success...")
	p.Layout = "layout/layout.html"
	p.TplName = "log/delete.html"
}