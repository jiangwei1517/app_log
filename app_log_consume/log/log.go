package log

import(
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego/logs"
)

var logConfig map[string]interface{} = make(map[string]interface{})

func InitLog(logPath, level string) (err error) {
	logConfig["filename"] = logPath
	intLevel := convertLevelStringToInt(level)
	logConfig["level"] = intLevel
	logConfigJson,err := json.Marshal(logConfig)
	if (err != nil) {
		fmt.Println("InitLog json.Marshal error", err)
		return
	}
	logs.SetLogger(logs.AdapterFile, string(logConfigJson))
	return
}

func convertLevelStringToInt(level string) (intLevel int) {
	switch level {
	case "debug":
		intLevel = logs.LevelDebug
	case "warn":
		intLevel = logs.LevelWarn
	case "error":
		intLevel = logs.LevelError
	case "info":
		intLevel = logs.LevelInfo
	default:
		intLevel = logs.LevelDebug
	}
	return
}