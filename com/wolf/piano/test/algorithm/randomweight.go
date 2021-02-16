package main

import (
	"fmt"
	"math/rand"
)

type invoker struct {
	weight int
}

func main() {
	invokers := []invoker{{weight: 80}, {weight: 20}}
	// Number of invokers
	length := len(invokers)
	// Every invoker has the same weight?
	sameWeight := true
	// the maxWeight of every invokers, the minWeight = 0 or the maxWeight of the last invoker
	weights := make([]int, length) // 第一个值为第一个，第二个值为当前值+前面的汇总值，第三个值为当前值+前面的汇总值
	// The sum of weights
	totalWeight := 0
	for i := 0; i < length; i++ {
		weight := invokers[i].weight
		// Sum
		totalWeight += weight
		// save for later use
		weights[i] = totalWeight
		if sameWeight && totalWeight != weight*(i+1) { // 若是都等，则总的值，等于当前值*(i+1)，i+1表示第几个+1即一共个数
			sameWeight = false
		}
	}
	if totalWeight > 0 && !sameWeight {
		// If (not every invoker has the same weight & at least one invoker's weight>0), select randomly based on totalWeight.
		offset := rand.Intn(totalWeight)
		// Return a invoker based on the random value.
		for i := 0; i < length; i++ {
			if offset < weights[i] {
				fmt.Println("get result in not sameWeight ", invokers[i])
				return
			}
		}
	}
	// If all invokers have the same weight value or totalWeight=0, return evenly.
	fmt.Println("get result in sameWeight ", rand.Intn(length))
}
