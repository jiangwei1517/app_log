package es

import(
	elastic "gopkg.in/olivere/elastic.v2"
	"github.com/astaxie/beego/logs"
)

var esClient *elastic.Client

type Message struct{
	Text string
	Topic string
}

func InitEs(addr string) (err error) {
	esClient, err = elastic.NewClient(elastic.SetSniff(false), elastic.SetURL(addr))
	if (err != nil) {
		logs.Warn("InitEs elastic.NewClient failed:", err)
		return
	}
	return
}



func SendToES(data, id, topic string) (err error) {
	msg := &Message{}
	msg.Text = data
	msg.Topic = topic
	_, err = esClient.Index().Index(topic).Type(topic).Id(id).BodyJson(msg).Do()
	if (err != nil) {
		logs.Warn("SendToES failed:", err)
		return
	}
	logs.Debug("Send to ES Server Success. Data = %v  Topic = %v", data, topic)
	return
}

