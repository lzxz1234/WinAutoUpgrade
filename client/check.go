package main

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"time"

	"github.com/lzxz1234/PowerUpgrade/util"

	"github.com/spf13/viper"
)

func check() {
	for ; ; time.Sleep(time.Minute) {
		apiURL, err := url.Parse(viper.GetString("upgrade_url"))
		if err != nil {
			fmt.Println("  --> 地址格式错误 "+viper.GetString("upgrade_url"), err)
			continue
		}
		apiURL.Query().Set("major", viper.GetString("major"))
		apiURL.Query().Set("minor", viper.GetString("minor"))
		apiURL.Query().Set("ip", util.GetIP())
		if resp, err := http.Get(apiURL.String()); err != nil {
			fmt.Println("  --> 校验版本信息失败 "+apiURL.String(), err)
		} else {
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println("  --> 校验版本信息失败", err)
				continue
			}
			var result CheckResult
			json.Unmarshal(body, &result)
			if result.Code != 1 {
				fmt.Println("  --> 未发现新版本")
				continue
			}
			fmt.Println("  --> 发现新版本 " + result.URL)
			dstResp, err := http.Get(result.URL)
			if err != nil {
				fmt.Println("  --> 升级包下载失败", err)
				continue
			}
			fmt.Println("  --> 新版本下载完成")

			dstBytes, err := ioutil.ReadAll(dstResp.Body)
			if err != nil {
				fmt.Println("  --> 升级包下载失败", err)
				continue
			}
			zr, err := zip.NewReader(bytes.NewReader(dstBytes), int64(len(dstBytes)))
			if err != nil {
				fmt.Println("  --> 压缩包读取失败", err)
				continue
			}
			// 停机替换文件了
			steps <- "stop"
			select {
			case <-steps: //wait for ready
				for _, f := range zr.File {
					root, _ := filepath.Abs(filepath.Dir(os.Args[0]))
					fpath := filepath.Join(root, f.Name)
					if f.FileInfo().IsDir() {
						os.Mkdir(fpath, os.ModePerm)
					} else {
						os.MkdirAll(filepath.Dir(fpath), os.ModePerm)
						if fr, err := f.Open(); err == nil {
							if fc, err := ioutil.ReadAll(fr); err == nil {
								ioutil.WriteFile(fpath, fc, f.Mode())
							}
						}
						fmt.Println("  --> 替换资源文件：", f.Name)
					}
				}
				fmt.Println("  --> 升级完成")
				viper.Set("major", result.Major)
				viper.Set("minor", result.Minor)
				viper.WriteConfig()
				steps <- "start"
			case <-time.After(5 * time.Second): //超时5s
				fmt.Println("  --> 停机失败")
				continue
			}
		}
	}
}

type CheckResult struct {
	Code    int    `json:"code"` // 0 无更新， 1 有更新
	Msg     string `json:"msg"`
	Major   string `json:"major"`
	Minor   string `json:"minor"`
	TarType string `json:"tar_type"` // 打包类型，只支持 zip
	URL     string `json:"url"`
}
