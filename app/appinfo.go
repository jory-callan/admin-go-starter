package app

import (
	"fmt"
	"runtime"
	"strings"
)

// AppInfo 应用基本信息（全局变量，编译时或启动时设置）
// 用于标记应用名称、版本号等元信息，不来自配置文件
var AppInfo = Info{
	Name:    "aicode",
	Version: "0.1.0",
	Desc:    "Go Web Framework with RBAC",
}

// Info 应用元信息结构体
type Info struct {
	Name      string // 应用名称
	Version   string // 版本号
	GoVersion string // Go 版本（自动获取）
	Desc      string // 应用描述
}

// GoVersion 返回当前 Go 运行时版本
func GoVersion() string {
	return runtime.Version()
}

// String 返回格式化的应用信息字符串
func (i Info) String() string {
	goVer := i.GoVersion
	if goVer == "" {
		goVer = GoVersion()
	}
	return fmt.Sprintf("%s v%s (Go %s)", i.Name, i.Version, goVer)
}

// Banner 返回启动时的 ASCII Banner
func (i Info) Banner() string {
	return `
 ██████╗ ██╗   ██╗ █████╗ ███╗   ██╗██╗   ██╗ ██████╗
██╔═══██╗██║   ██║██╔══██╗████╗  ██║██║   ██║██╔═══██╗
██║   ██║██║   ██║███████║██╔██╗ ██║██║   ██║██║   ██║
██║▄▄ ██║██║   ██║██╔══██║██║╚██╗██║██║   ██║██║   ██║
╚██████╔╝╚██████╔╝██║  ██║██║ ╚████║╚██████╔╝╚██████╔╝
 ╚══▀▀═╝  ╚═════╝ ╚═╝  ╚═╝╚═╝  ╚═══╝ ╚═════╝  ╚═════╝`
}

// PrintBanner 漂亮打印启动信息
func (i Info) PrintBanner() string {
	goVer := GoVersion()
	info := strings.TrimSpace(fmt.Sprintf(`
 %s  v%s
 %s
 Go Runtime : %s
 Platform   : %s/%s
`,
		i.Name, i.Version,
		i.Desc,
		goVer,
		runtime.GOOS, runtime.GOARCH,
	))
	return i.Banner() + info
}
