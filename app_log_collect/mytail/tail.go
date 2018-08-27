package mytail

import(
	"app_log_collect/config"
	"github.com/astaxie/beego/logs"
	"fmt"
	"github.com/hpcloud/tail"
	"sync"
	"app_log_collect/kafkaserver"
	"encoding/json"
	"app_log_collect/others"
)

type TailObj struct{
	tail *tail.Tail
	conf *config.TailConfig
	status int
	exitchan chan int
}

type TailManager struct{
	tails []*TailObj
	msgschan chan *kafkaserver.Message
	lock sync.Mutex
}

var tailMgr *TailManager
var wg sync.WaitGroup

func InitTailAndCollect(tailConfigs []*config.TailConfig) (err error) {
	tailMgr = &TailManager{}
	tailMgr.msgschan = make(chan *kafkaserver.Message, config.AppConfig.MessageChanSize)
	if (len(tailConfigs) == 0) {
		logs.Warn("InitTail error...")
		err = fmt.Errorf("InitTail error:%v",err)
		return
	}
	makeTailsByConfigs(tailConfigs)
	return
}

func makeTailsByConfigs(tailConfigs []*config.TailConfig) {
	for _,tailConfig := range tailConfigs {
		tailone,err := tail.TailFile(tailConfig.TailAddr, tail.Config{
			ReOpen : true,
			Follow : true,
			MustExist : false,
			Poll : true,
		})
		if (err != nil) {
			logs.Warn("tail.TailFile error...",err)
			continue
		}
		tailobj := &TailObj{}
		tailobj.tail = tailone
		tailobj.conf = tailConfig
		tailobj.exitchan = make(chan int, 1)
		tailMgr.tails = append(tailMgr.tails, tailobj) 
		go collectLog(tailobj)
	}
	wg.Add(1)
	go sendMsgToKafkaServer()
	wg.Wait()
}

func UpdateTailConfig(configJson string) (err error) {
	tailMgr.lock.Lock()
	defer tailMgr.lock.Unlock()
	var tailConfigs []*config.TailConfig
	err = json.Unmarshal([]byte(configJson), &tailConfigs)
	if (err != nil) {
		logs.Warn("UpdateTailConfig json.Unmarshal() failed:", err)
		return
	}
	if (len(tailConfigs) == 0) {
		logs.Warn("UpdateTailConfig json.Unmarshal tailConfigs length == 0")
		err = fmt.Errorf("InitTail error:%v",err)
		return
	}
	for _,tailConfig := range tailConfigs{
		isRunning := false
		for _,tailObj := range tailMgr.tails{
			if (tailConfig.TailAddr == tailObj.conf.TailAddr) {
				isRunning = true
				break
			}
		}
		if (isRunning) {
			continue
		}
		makeTailsByConfigs(tailConfigs)
	}

	// 删除的线程关闭
	for _,tailObj := range tailMgr.tails{
		tailObj.status = others.ETCD_TYPE_DELETE
		for _,tailConfig := range tailConfigs{
			if (tailObj.conf.TailAddr == tailConfig.TailAddr) {
				tailObj.status = others.ETCD_TYPE_NORMAL
				break
			}
		}
		if (tailObj.status == others.ETCD_TYPE_DELETE) {
			// 需要关闭那个指定的线程
			tailObj.exitchan<-1
		}
	}
	var objs []*TailObj
	for _,tailObj := range tailMgr.tails{
		if (tailObj.status == others.ETCD_TYPE_NORMAL) {
			objs = append(objs, tailObj)
		}
	}
	tailMgr.tails = objs
	return
}

func collectLog(tailObj *TailObj)  {
	line := tailObj.tail.Lines
	for {
		select{
			case msg,ok:=<-line:
				if (!ok) {
					logs.Debug("collectLog tail file finished...")
					return
				}
				msgOne := &kafkaserver.Message{}
				msgOne.Text = msg.Text
				msgOne.Topic = tailObj.conf.Topic
				tailMgr.msgschan <- msgOne
			case <-tailObj.exitchan:
				return
		}
	}
}

func sendMsgToKafkaServer()  {
	for msg := range tailMgr.msgschan{
		err := kafkaserver.SendToKafakServer(msg)
		if (err != nil) {
			logs.Warn("sendMsgToKafkaServer error", err)
			continue
		}
	}
	wg.Done()
}