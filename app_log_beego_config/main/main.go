package main

import(
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
	etcd_client "github.com/coreos/etcd/clientv3"
	"time"
	"github.com/astaxie/beego"
	_ "app_log_beego_config/router"
	"app_log_beego_config/model"
)

var logsConfig map[string]interface{} = make(map[string]interface{})

func main()  {
	err := initLog()
	if (err != nil) {
		panic(err)
	}
	logs.Debug("init Log success...")

	err = initDB()
	if (err != nil) {
		logs.Error("initDB failed:", err)
		return
	}
	logs.Debug("init initDB success...")

	err = initEtcd()
	if (err != nil) {
		logs.Error("initEtcd failed:", err)
		return
	}
	logs.Debug("init initEtcd success...")

	beego.Run()
}

func initEtcd() (err error) {
	etcdCli,err := etcd_client.New(etcd_client.Config{
		Endpoints : []string{"localhost:2379"},
		DialTimeout : 5*time.Second,
	})
	if (err != nil) {
		logs.Warn("initEtcd failed:", err)
		return
	}
	model.InitEtcd(etcdCli)
	return
}

func initDB() (err error) {
	db, err := sqlx.Open("mysql", "root:Ww18519303997#@!@tcp(localhost:3306)/test")
	if (err != nil) {
		logs.Error("sqlx.Open mysql failed...", err)
		return
	}
	model.SetDB(db)
	return
}

func initLog() (err error) {
	logsConfig["filename"] = "./log/app_beego_config.log"
	logsConfig["level"] = logs.LevelDebug
	logsConfigJson, err := json.Marshal(logsConfig)
	if (err != nil) {
		fmt.Println("initLog json.Marshal error:", err)
		return
	}
	err = logs.SetLogger(logs.AdapterFile, string(logsConfigJson))
	if (err != nil) {
		fmt.Println("logs.SetLogger error:", err)
		return
	}
	return
}