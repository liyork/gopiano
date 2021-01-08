package loadgen

import "log"

// 限定符log代表gopcp.v2/helper/log代码包。
// 日志记录器
var logger = myLog{}

type myLog log.Logger

func (m myLog) Infoln(s string) {

}

func (m myLog) Infof(format string, param ...interface{}) {

}

func (m myLog) Errorln(s string) {

}
