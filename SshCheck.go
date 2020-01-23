package main

import (
	"JustGo/config"
	"fmt"
)

func main() {
	vipConfig, error := config.Init() //vipConfig是配置
	fmt.Printf("config.init error是%v\n", error)
	for key, val := range vipConfig.(map[string]interface{}) { //循环接口类型，获取配置信息
		fmt.Printf("vipConfig 的key是%v val是%v\n", key, val)

		switch val.(type) { //判断val的类型
		case map[string]interface{}: //如果是 interface接口类型
			for ke, va := range val.(map[string]interface{}) { //循环接口类型，获取配置信息
				fmt.Printf("vipConfig 的ke是%v va是%v\n", ke, va)
			}
		}
	}
}
