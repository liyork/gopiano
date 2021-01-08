package logger

type myLogger struct {
}

var Logger = &myLogger{}

func (l *myLogger) Info(fmt string, s ...string) {

}

func (l *myLogger) Warnln(s string) {

}

func (l *myLogger) Warnf(s string, param ...interface{}) {

}

func (l *myLogger) Fatal(s string) {

}

func (l *myLogger) Infof(s string, param ...interface{}) {

}

func (l *myLogger) Fatalf(s string, e ...interface{}) {

}

func (l *myLogger) Errorf(s string, e ...interface{}) {

}

func (l *myLogger) Infoln(s string) {

}

func (l *myLogger) Errorln(s string) {

}

func (l *myLogger) Fatalln(e error) {

}
