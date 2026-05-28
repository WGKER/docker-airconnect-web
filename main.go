package main

import (
	"encoding/xml"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
)

// XML 结构体定义 与配置文件完全匹配
type AirUPnP struct {
	XMLName    xml.Name `xml:"airupnp"`
	Common     Common   `xml:"common"`
	MainLog    string   `xml:"main_log"`
	UpnpLog    string   `xml:"upnp_log"`
	UtilLog    string   `xml:"util_log"`
	RaopLog    string   `xml:"raop_log"`
	LogLimit   int      `xml:"log_limit"`
	MaxPlayers int      `xml:"max_players"`
	Binding    string   `xml:"binding"`
	Ports      string   `xml:"ports"`
	Devices    []Device `xml:"device"`
}

type Common struct {
	Enabled    int    `xml:"enabled"`
	MaxVolume  int    `xml:"max_volume"`
	HttpLength int    `xml:"http_length"`
	UpnpMax    int    `xml:"upnp_max"`
	Codec      string `xml:"codec"`
	Metadata   int    `xml:"metadata"`
	Flush      int    `xml:"flush"`
	Artwork    string `xml:"artwork"`
	Latency    string `xml:"latency"`
	Drift      int    `xml:"drift"`
}

type Device struct {
	UDN     string `xml:"udn"`
	Name    string `xml:"name"`
	Mac     string `xml:"mac"`
	Enabled int    `xml:"enabled"`
}

// 配置文件路径（容器内路径）
const configPath = "/config/config.xml"

// 页面模板
const htmlTemplate = `
<!DOCTYPE html>
<html lang="zh-CN">
<head>
    <meta charset="UTF-8">
    <title>AirConnect 管理面板</title>
    <style>
        * {box-sizing: border-box; margin: 0; padding: 0; font-family: Arial, sans-serif;}
        body {background: #f5f5f5; padding: 20px; max-width: 800px; margin: 0 auto;}
        .card {background: white; padding: 20px; border-radius: 10px; margin-bottom: 20px; box-shadow: 0 2px 8px rgba(0,0,0,0.1);}
        h1 {color: #2c3e50; margin-bottom: 20px; text-align: center;}
        h2 {color: #34495e; margin: 15px 0; font-size: 18px;}
        .switch {margin: 12px 0; display: flex; justify-content: space-between; align-items: center; padding: 10px; background: #fafafa; border-radius: 6px;}
        label {font-size: 16px; color: #2c3e50; font-weight: 500;}
        button {background: #3498db; color: white; border: none; padding: 10px 20px; border-radius: 6px; cursor: pointer; font-size: 16px; margin-top: 10px;}
        button:hover {background: #2980b9;}
        .status {color: #27ae60; margin-top: 15px; text-align: center; font-weight: bold;}
    </style>
</head>
<body>
    <div class="card">
        <h1>AirConnect 音箱管理</h1>
        {{if .Msg}}
        <div class="status">{{.Msg}}</div>
        {{end}}

        <!-- 总开关 -->
        <form method="post">
            <h2>全局总开关</h2>
            <div class="switch">
                <label>AirConnect 总开关</label>
                <select name="global_enabled">
                    <option value="1" {{if eq .Config.Common.Enabled 1}}selected{{end}}>开启</option>
                    <option value="0" {{if eq .Config.Common.Enabled 0}}selected{{end}}>关闭</option>
                </select>
            </div>

            <h2>音箱独立开关</h2>
            {{range $index, $device := .Config.Devices}}
            <div class="switch">
                <label>{{$device.Name}}</label>
                <select name="device_{{$index}}">
                    <option value="1" {{if eq $device.Enabled 1}}selected{{end}}>开启</option>
                    <option value="0" {{if eq $device.Enabled 0}}selected{{end}}>关闭</option>
                </select>
            </div>
            {{end}}

            <button type="submit">保存配置</button>
        </form>
    </div>
</body>
</html>
`

// 页面数据
type PageData struct {
	Config *AirUPnP
	Msg    string
}

// 读取XML配置
func loadConfig() (*AirUPnP, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
	var config AirUPnP
	if err := xml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

// 保存XML配置
func saveConfig(config *AirUPnP) error {
	data, err := xml.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	// 写入文件
	return os.WriteFile(configPath, append([]byte(xml.Header), data...), 0644)
}

// Web处理函数
func handler(w http.ResponseWriter, r *http.Request) {
	config, err := loadConfig()
	if err != nil {
		http.Error(w, "加载配置失败: "+err.Error(), 500)
		return
	}

	msg := ""
	// POST 提交保存
	if r.Method == http.MethodPost {
		// 全局开关
		globalVal := r.PostFormValue("global_enabled")
		globalEnabled, _ := strconv.Atoi(globalVal)
		config.Common.Enabled = globalEnabled

		// 遍历设备
		for i := range config.Devices {
			key := fmt.Sprintf("device_%d", i)
			val := r.PostFormValue(key)
			enabled, _ := strconv.Atoi(val)
			config.Devices[i].Enabled = enabled
		}

		// 保存
		if err := saveConfig(config); err != nil {
			msg = "保存失败：" + err.Error()
		} else {
			msg = "✅ 配置已保存！重启 AirConnect 容器生效"
		}
	}

	// 渲染页面
	tpl, _ := template.New("webui").Parse(htmlTemplate)
	tpl.Execute(w, PageData{Config: config, Msg: msg})
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("WebUI 启动在 http://0.0.0.0:8087")
	// 监听8087端口
	http.ListenAndServe(":8087", nil)
}
