package main

import(
	"app_log_consume/config"
	"app_log_consume/log"
	"app_log_consume/kafkaconsumer"
	"app_log_consume/es"
	"github.com/astaxie/beego/logs"
	"fmt"
)

func main()  {
	// init config
	err := config.InitConfig("ini", "config/config.conf")
	if (err != nil){
		panic(err)
	}
	fmt.Println("init config success...")

	// initlog
	err = log.InitLog(config.AppConfig.LogPath, config.AppConfig.LogLevel)
	if (err != nil) {
		panic(err)
	}
	fmt.Println("init log success...")

	// initES
	err = es.InitEs(config.AppConfig.ESAddr)
	if (err != nil) {
		logs.Warn("InitEs failed...", err)
		return
	}

	// initKafka
	err = kafkaconsumer.InitKafka(config.AppConfig.KafkaAddr)
	if (err != nil) {
		logs.Warn("InitKafka failed:", err)
		return
	}
	logs.Debug("InitKafka success...:")

	// start consume
	logs.Debug("Kafka start consume...:")
	for _,topic := range config.AppConfig.Topics{
		kafkaconsumer.WG.Add(1)
		go kafkaconsumer.Consume(topic)
	}
	kafkaconsumer.WG.Wait()
}