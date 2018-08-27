package model

import(
	// "github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/logs"
	etcd_client "github.com/coreos/etcd/clientv3"
	"fmt"
	"net"
	"strings"
	"context"
	"time"
	"encoding/json"
	"app_log_collect/config"
)

var etcdCli *etcd_client.Client
var localIps []string 
var etcd_key = "/jiangwei18/app_log_manage/tail/"

func init()  {
	addrs,err := net.InterfaceAddrs()
	if (err != nil) {
		fmt.Println("net.InterfaceAddrs error,",err)
		return
	}
	for _,addr := range addrs{
		if ipnet,ok := addr.(*net.IPNet);ok && !ipnet.IP.IsLoopback(){
			if (ipnet.IP.To4() != nil) {
				localIp := ipnet.IP.String()
				localIps = append(localIps, localIp)
			}
		}
	}
	if (!strings.HasSuffix(etcd_key, "/")) {
		etcd_key = etcd_key + "/"
	}
	etcd_key = etcd_key + localIps[0]
}

type LogInfo struct{
	AppId int `db:"app_id"`
	AppName string `db:"app_name"`
	CreateTime string `db:"create_time"`
	Topic string `db:"topic"`
	LogPath string `db:"log_path"`
	LogId int `db:"log_id"`
	Status int `db:"status"`
}

func InitEtcd(cli *etcd_client.Client)  {
	etcdCli = cli
}

func LogCreate(loginfo *LogInfo) (err error) {
	conn,err:=DB.Begin()
	if (err != nil) {
		logs.Error("LogCreate DB.Begin() failed,", err)
		return
	}
	defer func ()  {
		if (err != nil) {
			conn.Rollback();
			logs.Error("LogCreate failed start rollback,", err)
			return
		}
		err := conn.Commit()
		if (err != nil) {
			logs.Error("LogCreate() conn create failed,", err)
			return
		}
	}()
	var loginfos []*LogInfo
	err = DB.Select(&loginfos, "select app_id from tbl_app_info where app_name = ?", loginfo.AppName)
	if (err != nil) {
		logs.Error("LogCreate DB.Select failed...", err)
		return
	}
	for _,info := range loginfos{
		res,err1:=conn.Exec("insert into tbl_log_info(app_id, app_name, log_path, topic)values(?,?,?,?)", info.AppId, loginfo.AppName, loginfo.LogPath, loginfo.Topic)
		if (err1 != nil) {
			err = err1
			logs.Error("LogCreate conn.Exec insert failed,", err)
			return
		}
		_, err1 = res.LastInsertId()
		if (err1 != nil) {
			err = err1
			logs.Error("LogCreate LastInsertId failed,", err)
			return
		}
	}
	logs.Debug("LogCreate success...")
	err = putToEtcd(etcd_key, loginfo)
	if (err != nil) {
		logs.Error("sendToEtcd failed,", err)
		return
	}
	return
}

/**
	将新的配置更新到etcd服务器
*/
var etcdConfigs []*config.TailConfig

func putToEtcd(key string, loginfo *LogInfo) (err error) {
	etcdConfig := &config.TailConfig{}
	etcdConfig.TailAddr = loginfo.LogPath
	etcdConfig.Topic = loginfo.Topic
	// Get
	ctx,cancel := context.WithTimeout(context.Background(), 2*time.Second)
	resp,err := etcdCli.Get(ctx, key)
	cancel()
	if (err != nil) {
		logs.Error("etcdCli.Get failed,", err)
		return
	}
	var value []byte
	for _,v:=range resp.Kvs{
		if (string(v.Key) == etcd_key) {
			value = v.Value
		}
	}
	err = json.Unmarshal(value, &etcdConfigs)
	if (err != nil) {
		logs.Error("etcdCli.Get json.Unmarshal failed,", err)
		return
	}
	var isContains bool = false
	for _,conf := range etcdConfigs{
		if (conf.TailAddr == etcdConfig.TailAddr) {
			isContains = true
		}
	}
	if (!isContains) {
		etcdConfigs = append(etcdConfigs, etcdConfig)
	}

	// PUT
	etcdConfigJson, err := json.Marshal(etcdConfigs)
	if (err != nil){
		logs.Error("json.Marshal(etcdConfig) failed,", err)
		return
	}
	ctx,cancel = context.WithTimeout(context.Background(), 2*time.Second)
	_, err = etcdCli.Put(ctx, key, string(etcdConfigJson))
	cancel()
	if (err != nil){
		logs.Error("etcdCli.Put failed,", err)
		return
	}
	logs.Debug("etcdCli put loginfo config success...")
	etcdConfigs = nil
	return
}

func deleteFromEtcd(key string, loginfo *LogInfo) (err error) {
	etcdConfig := &config.TailConfig{}
	etcdConfig.TailAddr = loginfo.LogPath
	etcdConfig.Topic = loginfo.Topic

	// Get
	ctx,cancel := context.WithTimeout(context.Background(), 2*time.Second)
	resp,err := etcdCli.Get(ctx, key)
	cancel()
	if (err != nil) {
		logs.Error("etcdCli.Get failed,", err)
		return
	}
	var value []byte
	for _,v:=range resp.Kvs{
		if (string(v.Key) == etcd_key) {
			value = v.Value
		}
	}
	err = json.Unmarshal(value, &etcdConfigs)
	if (err != nil) {
		logs.Error("etcdCli.Get json.Unmarshal failed,", err)
		return
	}
	var k int
	for index,conf := range etcdConfigs{
		if (conf.TailAddr == etcdConfig.TailAddr && conf.Topic == etcdConfig.Topic) {
			k = index
			break
		}
	}
	etcdConfigs = append(etcdConfigs[:k], etcdConfigs[k+1:]...)

	// Delete
	etcdConfigJson, err := json.Marshal(etcdConfigs)
	if (err != nil){
		logs.Error("json.Marshal(etcdConfig) failed,", err)
		return
	}
	ctx,cancel = context.WithTimeout(context.Background(), 2*time.Second)
	_, err = etcdCli.Put(ctx, key, string(etcdConfigJson))
	cancel()
	if (err != nil){
		logs.Error("etcdCli.Delete failed,", err)
		return
	}
	logs.Debug("etcdCli delete loginfo config success...")
	etcdConfigs = nil
	return
}

func GetAllLogInfo() (infos []*LogInfo, err error) {
	err = DB.Select(&infos, "select a.app_id, a.app_name, a.log_id, a.create_time, a.log_path, a.topic from tbl_log_info a, tbl_app_info b where a.app_id = b.app_id")
	if (err != nil) {
		logs.Error("GetAllLogInfo DB select Error,", err)
		return
	}
	return
}

func LogDelete(loginfo *LogInfo) (err error) {
	conn,err := DB.Begin()
	if (err != nil) {
		logs.Error("LogDelete DB.Begin() failed,", err)
		return
	}
	defer func ()  {
		if (err != nil) {
			logs.Error("LogDelete failed start rollback,", err)
			conn.Rollback()
			return
		}
		err = conn.Commit()
		if (err != nil) {
			logs.Error("LogDelete conn.Commit() failed,", err)
			return
		}
	}()
	var logIds []int
	err = DB.Select(&logIds, "select log_id from tbl_log_info where app_name = ? and log_path = ? and topic = ?", loginfo.AppName, loginfo.LogPath, loginfo.Topic)
	if (err != nil) {
		logs.Error("LogDelete DB.Select error,", err)
		return
	}
	for _,logId := range logIds{
		_, err1 := conn.Exec("Delete from tbl_log_info where log_id = ?", logId)
		if (err1 != nil) {
			logs.Error("Delete from tbl_log_info error", err)
			return
		}
	}
	err = deleteFromEtcd(etcd_key, loginfo)
	if (err != nil){
		logs.Error("Delete from etcd error", err)
		return
	}
	return
}