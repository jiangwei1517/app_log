package main

import(
	"app_log_collect/config"
	"app_log_collect/log"
	"app_log_collect/others"
	"fmt"
	"github.com/astaxie/beego/logs"
	"app_log_collect/etcd"
	"app_log_collect/kafkaserver"
	"app_log_collect/mytail"
)

func main()  {
	// initConfig
	err := config.InitConfig(others.CONFIG_TYPE, others.CONFIG_FILE_PATH)
	if (err != nil) {
		fmt.Println("init Config failed...")
		panic(err)
	}
	fmt.Println("init Config success...")

	// initLog
	err = log.InitLog(config.AppConfig.LogPath, config.AppConfig.LogLevel)
	if (err != nil) {
		fmt.Println("init Log failed...")
		panic(err)
	}
	fmt.Println("init log success...")
	fmt.Println("start logging...")
	logs.Debug("init log success...")

	// initKafka
	err = kafkaserver.InitKafka(config.AppConfig.KafkaAddr)
	if (err != nil) {
		logs.Warn("init kafka server failed...")
		return
	}
	logs.Debug("init kafka server success...")

	// initEtcd
	err = etcd.InitEtcd(config.AppConfig.EtcdAddr)
	if (err != nil) {
		logs.Warn("init etcd failed...")
		return
	}
	logs.Debug("init etcd success...")

	// put tailConfig to Etcd   only for test
	var tailconfigs []*config.TailConfig
	logNormal := &config.TailConfig{
		Topic : config.AppConfig.TopicNormal,
		TailAddr : config.AppConfig.AddrNormal,
	}
	tailconfigs = append(tailconfigs, logNormal)
	logError := &config.TailConfig{
		Topic : config.AppConfig.TopicError,
		TailAddr : config.AppConfig.AddrError,
	}
	tailconfigs = append(tailconfigs, logError)

	/** 测试etcd的watch函数
		logNew := &config.TailConfig{
			Topic : "new_topic",
			TailAddr : "log/new.log",
		}
		tailconfigs = append(tailconfigs, logNew)
	*/

	logs.Debug("config.AppConfig.EtcdTailKey = %v", config.AppConfig.EtcdTailKey)
	err = etcd.PutTailCollect(config.AppConfig.EtcdTailKey, tailconfigs)
	if (err != nil) {
		logs.Warn("put tailConfig to Etcd error,", err)
		return
	}
	logs.Debug("put tailConfig to Etcd success...")

	// get tailconfig from Etcd  only for test
	tailconfigs,err = etcd.GetTailCollect(config.AppConfig.EtcdTailKey)
	if (err != nil) {
		logs.Warn("get tailconfig from etcd error,", err)
		return
	}
	logs.Debug("get tailConfig from Etcd success...")

	go etcd.WatchEtcd()

	// init tail and collect logs
	err = mytail.InitTailAndCollect(tailconfigs)
	if (err != nil) {
		logs.Warn("init tail and collect logs error...", err)
		return
	}
	logs.Debug("init tail and collect logs success...")
}