package info

import (
	"fmt"
	"runtime"
)

// AppInfo 应用基本信息（全局变量，编译时或启动时设置）
// 用于标记应用名称、版本号等元信息，不来自配置文件
var AppInfo = Info{
	Name:      "aicode",
	Version:   "0.1.0",
	Desc:      "Go Web Framework with RBAC",
	BuildTime: "", // 编译时通过 -ldflags 注入
}

// Info 应用元信息结构体
type Info struct {
	Name      string // 应用名称
	Version   string // 版本号
	GoVersion string // Go 版本（自动获取）
	Desc      string // 应用描述
	BuildTime string // 构建时间（编译时注入）
}

// GoVersion 返回当前 Go 运行时版本
func GoVersion() string {
	return runtime.Version()
}

// Banner 返回启动时的 ASCII Banner
func (i Info) Banner() string {
	return `
 ██████╗ ██╗   ██╗ █████╗ ███╗   ██╗██╗   ██╗ ██████╗
██╔═══██╗██║   ██║██╔══██╗████╗  ██║██║   ██║██╔═══██╗
██║   ██║██║   ██║███████║██╔██╗ ██║██║   ██║██║   ██║
██║▄▄ ██║██║   ██║██╔══██║██║╚██╗██║██║   ██║██║   ██║
╚██████╔╝╚██████╔╝██║  ██║██║ ╚████║╚██████╔╝╚██████╔╝
 ╚══▀▀═╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝  ╚═════╝
`
}

// PrintBanner 漂亮打印启动信息
func (i Info) PrintBanner() string {
	goVer := GoVersion()
	buildTimeStr := ""
	if i.BuildTime != "" {
		buildTimeStr = fmt.Sprintf("\n Build Time : %s", i.BuildTime)
	}
	info := fmt.Sprintf(`
  %s  %s
  %s
  Go Runtime : %s
  Platform   : %s/%s%s
`,
		i.Name, i.Version,
		i.Desc,
		goVer,
		runtime.GOOS, runtime.GOARCH, buildTimeStr,
	)
	return i.Banner() + info
}
