package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strings"
)

var chatbotName string

func init() {
	flag.StringVar(&chatbotName, "chatbo", "simple.en", "The chatbot's name for dialogue")
}

// go run simple.go
// a
// bye
func main() {
	inputReader := bufio.NewReader(os.Stdin)
	fmt.Println("Please input your name:")
	input, err := inputReader.ReadString('\n')
	if err != nil {
		fmt.Printf("An error occured: %s\n", err)
		os.Exit(1)
	} else {
		name := input[:len(input)-1]
		fmt.Printf("Hello, %s! What can I do for you?\n", name)
	}
	for {
		input, err = inputReader.ReadString('\n')
		if err != nil {
			fmt.Printf("An error occured: %s\n", err)
			continue
		}
		input = input[:len(input)-1]
		input = strings.ToLower(input)
		switch input {
		case "":
			continue
		case "nothing", "bye":
			fmt.Println("Bye!")
			os.Exit(0)
		default:
			fmt.Println("Sorry, I didn't catch you.")
		}
	}

	// 可以使用chatbotName从map找到然后进行使用
	// go run simple.go -chatbot simple.cn
}

type simpleCN struct {
	name string
	talk Talk
}

func NewSimpleCN(name string, talk Talk) Chatbot {
	return &simpleCN{
		name: name,
		talk: talk,
	}
}

// *simpleCN实现了Chatbot和Talk接口
func (robot *simpleCN) Name() string {
	return robot.name
}

func (robot *simpleCN) Begin() (string, error) {
	return "begin", nil
}

func (robot *simpleCN) Hello(userName string) string {
	if robot.talk != nil {
		return robot.talk.Hello(userName)
	}
	return "hello robot"
}

func (robot *simpleCN) Talk(heard string) (saying string, end bool, err error) {
	if robot.talk != nil {
		return robot.talk.Talk(heard)
	}
	return "robot saying", true, nil
}
func (robot *simpleCN) ReportError(err error) string {
	return ""
}

func (robot *simpleCN) End() error {
	return nil
}

// 自定义注册
var chatbotMap = map[string]Chatbot{}

func Register(chatbot Chatbot) error {
	chatbotMap[chatbot.Name()] = chatbot
	return nil
}

func Get(name string) Chatbot {
	return chatbotMap[name]
}
