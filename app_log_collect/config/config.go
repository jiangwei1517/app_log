package config

import(
	"github.com/astaxie/beego/config"
	"fmt"
	"net"
	"strings"
)

type Config struct{
	EtcdAddr string
	KafkaAddr string
	LogPath	string
	LogLevel string
	TailCollectConfigs []*TailConfig
	MessageChanSize int
	EtcdTailKey string
	TopicNormal string
	AddrNormal string
	TopicError string
	AddrError string
}

type TailConfig struct{
	TailAddr string
	Topic string
}

var AppConfig *Config = &Config{}
var localIps []string

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
}

func InitConfig(typ string, filePath string) (err error) {
	conf,err := config.NewConfig(typ, filePath)
	if (err != nil) {
		fmt.Println("config.NewConfig error:",err)
		return
	}
	err = parseConfigFile(conf)
	if (err != nil) {
		fmt.Println("parseConfigFile error:",err)
		return
	}
	return
}

func parseConfigFile(conf config.Configer) (err error) {
	config := &Config{}
	config.EtcdAddr = conf.String("etcd::listen_ip")
	if (len(config.EtcdAddr) == 0) {
		config.EtcdAddr = "localhost:2379"
	}
	config.KafkaAddr = conf.String("kafka::listen_ip")
	if (len(config.KafkaAddr) == 0) {
		config.KafkaAddr = "localhost:9092"
	}
	config.LogPath = conf.String("etcd::log_path")
	if (len(config.LogPath) == 0) {
		config.LogPath = "log/app_collect.log"
	}
	config.LogLevel = conf.String("etcd::log_level")
	if (len(config.LogLevel) == 0) {
		config.LogLevel = "debug"
	}
	config.EtcdTailKey = conf.String("tail::etcd_key")
	if (len(config.EtcdTailKey) == 0) {
		config.EtcdTailKey = "/jiangwei18/app_log_manage/tail/"
	}
	if (!strings.HasSuffix(config.EtcdTailKey, "/")) {
		config.EtcdTailKey = config.EtcdTailKey + "/"
	}
	config.EtcdTailKey = config.EtcdTailKey + localIps[0]
	config.MessageChanSize,err = conf.Int("tail::message_chan_size")
	if (err != nil) {
		fmt.Println("config.MessageChanSize init error...")
		config.MessageChanSize = 1000
	}
	config.TopicNormal = conf.String("collect::topic_normal")
	if (len(config.TopicNormal) == 0) {
		config.TopicNormal = "log_normal"
	}
	config.AddrNormal = conf.String("collect::addr_normal")
	if (len(config.AddrNormal) == 0) {
		config.AddrNormal = "log/normal.log"
	}
	config.TopicError = conf.String("collect::topic_error")
	if (len(config.TopicError) == 0) {
		config.TopicError = "log_error"
	}
	config.AddrError = conf.String("collect::addr_error")
	if (len(config.AddrError) == 0) {
		config.AddrError = "log/error.log"
	}
	AppConfig = config
	return
}