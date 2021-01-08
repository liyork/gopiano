package main

type Talk interface {
	Hello(userName string) string
	Talk(heard string) (saying string, end bool, err error)
}

type myTalk string

func (talk *myTalk) Hello(userName string) string                           { return "" }
func (talk *myTalk) Talk(heard string) (saying string, end bool, err error) { return "", false, nil }

type Chatbot interface {
	Name() string
	Begin() (string, error)
	Talk
	ReportError(err error) string
	End() error
}
