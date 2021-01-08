package main

/*
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#define RET_OK 0

int GetUserName(char *userName)
{	// 拷贝xx123到userName，一共5字节.
    // dst,src,len
    memcpy(userName, "xx123", 5);
    return RET_OK;
}

int GetKey(char *key, int *keyLen)
{
    memcpy(key, "\x00\x11\x22\x33\x44\x55\x66\x77\x00\x11\x22\x33\x44\x55\x66\x77", 16);
    *keyLen = 16;
    return RET_OK;
}
*/
import "C"
import (
	"fmt"
	"unsafe"
)

//注意cgo的注释和 import "C"之间不能有空行！！！ 注意cgo的注释和 import "C"之间不能有空行！！！ 注意cgo的注释和 import "C"之间不能有空行！！！

const (
	MAX_BUFFER = 0x5000 // 留个足够大缓冲区，方便以后通用
)

// 截取到0的
func byteString(p []byte) string {
	for i := 0; i < len(p); i++ {
		if p[i] == 0 {
			return string(p[0:i])
		}
	}
	return string(p)
}

// 转换成Char指针
func getCharPointer(val []byte) *C.char {
	return (*C.char)(unsafe.Pointer(&val[0]))
}

func main() {
	userNameBuff := make([]byte, MAX_BUFFER)

	c_userName := getCharPointer(userNameBuff)

	// 调用c
	ret := C.GetUserName(c_userName)
	fmt.Println(ret)

	go_userName := C.GoString(c_userName)
	fmt.Println(go_userName)

	go_userName = C.GoStringN(c_userName, C.int(len(go_userName)))
	fmt.Println(go_userName)

	str := byteString(userNameBuff[:])
	fmt.Println(str)

	keyBuff := make([]byte, MAX_BUFFER)
	keyLen := 0

	c_key := getCharPointer(keyBuff)

	ret = C.GetKey(c_key, (*C.int)(unsafe.Pointer(&keyLen)))

	fmt.Println(ret)
	fmt.Println(keyLen)
	fmt.Println(keyBuff[:keyLen])

	// 通过C.GoStringN转换测试
	go_key := C.GoStringN(c_key, C.int(keyLen))
	fmt.Println([]byte(go_key))
}
