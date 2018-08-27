package config

import(
	"github.com/astaxie/beego/config"
	"fmt"
)

type Config struct{
	KafkaAddr string
	LogPath string
	LogLevel string
	Topics []string
	ESAddr string
}

var AppConfig = &Config{}

func InitConfig(typ string, filePath string) (err error) {
	conf,err := config.NewConfig(typ, filePath)
	if (err != nil) {
		fmt.Println("InitConfig config.NewConfig error,", err)
		return
	}
	err = parseConfig(conf)
	if (err != nil) {
		fmt.Println("InitConfig parseConfig error,")
		return
	}
	return
}

func parseConfig(conf config.Configer) (err error) {
	AppConfig.KafkaAddr = conf.String("kafka::listen_ip")
	if (len(AppConfig.KafkaAddr) == 0) {
		AppConfig.KafkaAddr = "localhost:9092"
	}
	AppConfig.LogLevel = conf.String("logs::log_level")
	if (len(AppConfig.LogLevel) == 0) {
		AppConfig.LogLevel = "debug"
	}
	AppConfig.LogPath = conf.String("logs::log_path")
	if (len(AppConfig.LogPath) == 0) {
		AppConfig.LogPath = "log/kafka_consumer.log"
	}
	AppConfig.ESAddr = conf.String("es::listen_ip")
	if (len(AppConfig.ESAddr) == 0) {
		AppConfig.ESAddr = "http://localhost:9200"
	}
	AppConfig.Topics = append(AppConfig.Topics, "log_normal")
	AppConfig.Topics = append(AppConfig.Topics, "log_error")
	// 测试etcd的watch用的
	AppConfig.Topics = append(AppConfig.Topics, "new_topic")
	return
}