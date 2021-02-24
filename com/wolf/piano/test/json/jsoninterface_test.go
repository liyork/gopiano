package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"testing"
)

type Shape interface {
	Name() string
	Area() float64
	Circumference() float64
}

type Square struct {
	Size int `json:"size"`
}

func (s Square) Name() string {
	return "square"
}

func (s Square) Area() float64 {
	return float64(s.Size * s.Size)
}

func (s Square) Circumference() float64 {
	return float64(4 * s.Size)
}

type Circle struct {
	Radius int `json:"radius"`
}

func (c Circle) Name() string {
	return "circle"
}

func (c Circle) Area() float64 {
	return math.Pi * float64(c.Radius*c.Radius)
}

func (c Circle) Circumference() float64 {
	return math.Pi * float64(c.Radius*2)
}

type ShapeMap map[string]Shape

func (sm *ShapeMap) UnmarshalJSON(data []byte) error {
	// 解析key
	shapes := make(map[string]json.RawMessage)
	err := json.Unmarshal(data, &shapes)
	if err != nil {
		return err
	}
	result := make(ShapeMap)
	// 再解析value
	for k, v := range shapes {
		switch k {
		case "square":
			s := Square{}
			err := json.Unmarshal(v, &s)
			if err != nil {
				return err
			}
			result[k] = s
		case "circle":
			c := Circle{}
			err := json.Unmarshal(v, &c)
			if err != nil {
				return err
			}
			result[k] = c
		default:
			return errors.New("Unrecognized shape")
		}
	}
	*sm = result
	return nil
}

type Scene struct {
	Shapes ShapeMap `json:"shapes"`
}

func TestInterfacejson(t *testing.T) {
	data := []byte(`{"shapes":{
		"square": {
			"size": 2
		},
		"circle": {
			"radius": 1
		}
	}}`)
	s := Scene{}
	err := json.Unmarshal(data, &s)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(s.Shapes)
}
