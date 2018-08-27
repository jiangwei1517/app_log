package kafkaserver

import(
	"github.com/Shopify/sarama"
	"github.com/astaxie/beego/logs"
)

type Message struct{
	Topic string
	Text string
}

var producer sarama.SyncProducer

func InitKafka(addr string) (err error) {
	config:=sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true
	config.Producer.Return.Errors = true
	producer,err = sarama.NewSyncProducer([]string{addr}, config)
	if (err != nil) {
		logs.Warn("sarama.NewClient error",err)
		return
	}
	return
}

func SendToKafakServer(msg *Message) (err error) {
	msgOne := &sarama.ProducerMessage{}
	msgOne.Topic = msg.Topic
	msgOne.Value = sarama.StringEncoder(msg.Text) 
	partition,offset,err := producer.SendMessage(msgOne)
	if (err != nil) {
		logs.Warn("SendToKafakServer producer.SendMessage error", err)
		return
	}
	logs.Debug("SendToKafakServer Partition = %v, Offset = %v Value = %v Topic = %v", partition, offset, msg.Text, msg.Topic)
	return
}