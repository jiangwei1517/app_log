package others

import(
	"time"
)

const (
	CONFIG_TYPE = "ini"
	CONFIG_FILE_PATH = "config/config.conf"
	ETCD_DIAL_TIME_OUT = 5*time.Second
	ETCD_GET_TIME_OUT = 5*time.Second
	ETCD_PUT_TIME_OUT = 5*time.Second
	ETCD_TYPE_DELETE = 1
	ETCD_TYPE_NORMAL = 2
)