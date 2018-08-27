package etcd

import(
	etcd_client "github.com/coreos/etcd/clientv3"
	"app_log_collect/others"
	"github.com/astaxie/beego/logs"
	"app_log_collect/config"
	"context"
	"encoding/json"
	"time"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"app_log_collect/mytail"
)

var EtcdClient *etcd_client.Client

func InitEtcd(addr string) (err error) {
	EtcdClient,err=etcd_client.New(etcd_client.Config{
		Endpoints : []string{addr},
		DialTimeout : others.ETCD_DIAL_TIME_OUT,
	})
	if (err != nil) {
		logs.Warn("etcd_client.New error:",err)
		return
	}
	return
}

func WatchEtcd()  {
	for {
		watchChan := EtcdClient.Watch(context.Background(), config.AppConfig.EtcdTailKey)
		watchMsg := <-watchChan
		for _,value := range watchMsg.Events{
			if (value.Type == mvccpb.DELETE) {
				logs.Debug("detect etcdlib has deleted key = %v  value = %v", string(value.Kv.Key), string(value.Kv.Value))
			} else if (value.Type == mvccpb.PUT) {
				logs.Debug("detect etcdlib has put key = %v  value = %v", string(value.Kv.Key), string(value.Kv.Value))
			}
			go mytail.UpdateTailConfig(string(value.Kv.Value))
			// if (err != nil) {
			// 	logs.Warn("mytail.UpdateTailConfig failed, key = %v  value = %v  error = %v", string(value.Kv.Key), string(value.Kv.Value), err)
			// }
		}
		time.Sleep(time.Second)
	}
}

func GetTailCollect(key string) (tailConfigs []*config.TailConfig, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), others.ETCD_GET_TIME_OUT)
	resp, err := EtcdClient.Get(ctx, key)
	defer cancel()
	if (err != nil) {
		logs.Warn("GetTailCollect error,",err)
		return
	}
	for _,res := range resp.Kvs{
		if (string(res.Key) == key) {
			err = json.Unmarshal(res.Value, &tailConfigs)
			if (err != nil) {
				logs.Warn("GetTailCollect json.Unmarshal error,",err)
				return
			}
		}
	}
	return
}

func PutTailCollect(key string, tailConfigs []*config.TailConfig) (err error) {
	ctx,cancel := context.WithTimeout(context.Background(), others.ETCD_PUT_TIME_OUT);
	defer cancel()
	tailsConfig,err := json.Marshal(tailConfigs)
	if (err != nil) {
		logs.Warn("PutTailCollect json.Marshal error,",err)
		return
	}
	_,err = EtcdClient.Put(ctx, key, string(tailsConfig))
	if (err != nil) {
		logs.Warn("PutTailCollect EtcdClient.Put error,",err)
		return
	}
	return
}