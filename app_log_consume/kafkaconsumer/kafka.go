package kafkaconsumer

import(
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
	"sync"
	"app_log_consume/es"
	"sync/atomic"
	"fmt"
)

var ConsumerCli sarama.Consumer
var WG sync.WaitGroup
var WGchild sync.WaitGroup
var id int32 = 0

func InitKafka(addr string) (err error) {
	config:=sarama.NewConfig()
	config.Consumer.Return.Errors = true
	ConsumerCli, err = sarama.NewConsumer([]string{addr}, config)
	if (err != nil) {
		logs.Warn("InitKafka sarama.NewClient error:", err)
		return
	}
	return
}

func Consume(topic string) (err error) {
	partitions,err := ConsumerCli.Partitions(topic)
	if (err != nil) {
		logs.Warn("InitKafka Consumer.Partitions error:", err)
	}
	for _,partition := range partitions{
		pc, err := ConsumerCli.ConsumePartition(topic, partition, sarama.OffsetOldest)
		defer pc.AsyncClose()
		if (err != nil) {
			logs.Warn("ConsumerCli.ConsumePartition error:",err)
			continue
		}
		WGchild.Add(1)
		go consumeProcess(pc)
	}
	WGchild.Wait()
	WG.Done()
	return
}

func consumeProcess(pc sarama.PartitionConsumer)  {
	logs.Debug("START Consume......")
	for {
		msg,ok := <-pc.Messages()
		if (!ok) {
			logs.Debug("consumeProcess consume finished...")
			break
		}
		// 传给ES服务器
		logs.Debug("Topic = %v    Value = %v", msg.Topic, string(msg.Value))
		newId := atomic.AddInt32(&id, 1)
		es.SendToES(string(msg.Value), fmt.Sprintf("%v", newId), msg.Topic)
	}
	WGchild.Done()
}
