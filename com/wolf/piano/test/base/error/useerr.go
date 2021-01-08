package main

import (
	"errors"
	"log"
)

//使用
// 子数据层
func InsertA(id int) (err error) {
	if id == 0 {
		// 非500, 则是"错误"
		err = NewCoder(400, "id不能为空")
		return
	}

	err = Insert("insert data")
	if err != nil {
		// 500则是"异常"
		err = NewCoder(500, err.Error()+",插入数据库错误")
		return
	}
	return
}

func Insert(i interface{}) error {
	return errors.New("insert error")
}

func InsertB(id int) (err error) {
	if id == 0 {
		err = NewCoder(400, "id不能为空")
		return
	}
	return
}

// 数据层
func Data(aid int, bid int) (err error) {
	err = InsertA(aid)
	if err != nil {
		// 使用warp方法返回更详细的错
		err = Wrap(err, "Data 插入A错误")
		return
	}
	err = InsertB(aid)
	if err != nil {
		err = Wrap(err, "Data 插入B错误")
		return
	}

	return
}

// controll层
// 区分"错误"与"异常"并返回对应的响应
func Api() {
	var aid = 1
	var bid = 1
	err := Data(aid, bid)
	errorCode := err.(*ErrorCode)
	if errorCode.code == 500 {
		// 将"异常"位置和信息都打印方便排除
		log.Println("log err:", errorCode.where, errorCode.Error())
		log.Println("to client response ", 500, "服务器错误")
		return
	}

	log.Println("to client response ", errorCode.code, errorCode.Error())
}

// 最终的错误提示会是这样:
//
//(错误) 插入A错误: id不能为空
//(错误) 插入B错误: id不能为空  --这是包装，错误
//(异常) code.go:154 插入A错误: 插入数据库错误: sql: Scan error on column index 1: unsupported Scan, storing driver.Value type into type *string  --这是根音，异常
func main() {
	Api()
}
