package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/lzxz1234/PowerUpgrade/util"

	"github.com/spf13/viper"
)

func run() {
	var cmd *exec.Cmd
	for {
		step := <-steps
		if step == "stop" && cmd != nil {
			if err := util.RecursiveKill(uint32(cmd.Process.Pid)); err != nil {
				fmt.Println("==> 停止失败", err)
			} else {
				fmt.Println("==> 停止成功")
				steps <- "ready"
			}
		} else if step == "start" {
			cmds := strings.Split(viper.GetString("run_cmd"), " ")
			cmd = exec.Command(cmds[0], cmds[1:]...)
			cmd.Stdout = os.Stdout
			cmd.Stdin = os.Stdin
			cmd.Stderr = os.Stderr
			if err := cmd.Start(); err != nil {
				fmt.Println("==> 启动失败", err)
			} else {
				fmt.Println("==> 启动成功")
			}
		}
	}
}
