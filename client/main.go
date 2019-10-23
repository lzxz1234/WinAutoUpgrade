package main

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/viper"
)

var steps = make(chan string)

func main() {
	viper.AddConfigPath(".")
	viper.SetConfigName(".upgrade")
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("未加载到配置文件 .upgrade.yml")
		} else {
			fmt.Println("配置文件格式错误", err)
		}
		os.Exit(1)
	}

	go run()
	go check()
	steps <- "start"

	time.Sleep(time.Hour * 24 * 36500)
}
