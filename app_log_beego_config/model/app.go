package model

import(
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/logs"
)

var DB *sqlx.DB

type AppInfo struct{
	AppId int	`db:"app_id"`
	AppName string	`db:"app_name"`	
	AppType string	`db:"app_type"`
	CreateTime string	`db:"create_time"`
	DevelopPath string	`db:"develop_path"`
	IPs	[]string
}

func GetAllAppInfo() (appInfos []*AppInfo, err error) {
	err = DB.Select(&appInfos, "select app_id, app_name, app_type, create_time, develop_path from tbl_app_info")
	if (err != nil) {
		logs.Error("GetAllAppInfo failed,", err)
		return
	}
	return
}

func CreateAppInfo(appinfo *AppInfo) (err error) {
	conn,err := DB.Begin()
	if (err != nil) {
		logs.Error("DB.Begin error,", err)
		return
	}
	defer func ()  {
		if (err != nil){
			conn.Rollback()
			logs.Error("CreateAppInfo failed start rollback,", err)
			return
		}
		err = conn.Commit()
		if (err != nil) {
			logs.Error("LogCreate() conn create failed,", err)
			return
		}
	}()
	result,err := conn.Exec("insert into tbl_app_info(app_name, app_type, develop_path)values(?,?,?)", appinfo.AppName, appinfo.AppType, appinfo.DevelopPath)
	if (err != nil) {
		logs.Error("CreateAppInfo insert into tbl_app_info failed", err)
		return
	}
	insertId,err := result.LastInsertId()
	for _,ip:=range appinfo.IPs{
		_,err1 := conn.Exec("insert into tbl_app_ip(app_id, app_name, ip)values(?,?,?)", insertId, appinfo.AppName, ip)
		if (err1 != nil) {
			err = err1
			logs.Error("insert into tbl_app_ip error,",err)
			return
		}
	}
	return
}

func SetDB(db *sqlx.DB)  {
	DB = db
}