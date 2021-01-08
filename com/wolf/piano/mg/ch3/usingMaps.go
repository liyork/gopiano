package main

import "fmt"

// go run usingMaps.go
func main() {
	iMap := make(map[string]int)
	iMap["k1"] = 12
	iMap["k2"] = 13
	fmt.Println("iMap:", iMap)

	anotherMap := map[string]int{
		"k1": 12,
		"k2": 13,
	}

	fmt.Println("anotherMap:", anotherMap)
	// 多次调用delete()去删除同一个元素并没有导致异常或者警告。
	delete(anotherMap, "k1")
	delete(anotherMap, "k1")
	delete(anotherMap, "k1")
	fmt.Println("anotherMap:", anotherMap)

	// 是否存在
	_, ok := iMap["doseItExist"]
	if ok {
		fmt.Println("Exist!")
	} else {
		fmt.Println("dose NOT exist")
	}

	for key, value := range iMap {
		fmt.Println(key, value)
	}
}
