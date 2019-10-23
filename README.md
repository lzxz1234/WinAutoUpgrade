## PowerUpgrade 

Windows 平台下的自动升级程序

## 配置文件

upgrade.exe 同级目录下保存 .upgrade.yml:

```yaml
upgrade_url: http://127.0.0.1:80/check.json
run_exe: test.exe
run_cmd: cmd /C test.exe
major: 1
minor: 1
```

## upgrade_url 

版本校验地址，会添加三个参数,

- major、minor 代指主副版本号，
- ip 当前网卡 ip，用于服务端控制版本

返回报文格式应当如下：

```json
{
	"code": 1,
	"url": "http://127.0.0.1:80/test.zip", 
	"major": "1", 
	"minor": "2"
}
```

- code: 1 有新版本，0 无
- url: 资源下载路径，格式应当为 zip