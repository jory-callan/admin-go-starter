package echoutil

import (
	"fmt"
	"os"
	"sort"
	"text/tabwriter"

	"github.com/labstack/echo/v4"
)

// PrintRoutes 格式化并打印所有 Echo 路由
func PrintRoutes(e *echo.Echo) {
	// 获取所有路由
	routes := e.Routes()

	// 按照 Path 进行排序
	sort.Slice(routes, func(i, j int) bool {
		if routes[i].Path != routes[j].Path {
			return routes[i].Path < routes[j].Path
		}
		return routes[i].Method < routes[j].Method
	})

	// 初始化 tabwriter
	// 参数说明：输出目标, 最小单元格宽度, 制表符宽度, 填充空格数, 填充字符, 标志
	w := tabwriter.NewWriter(os.Stdout, 0, 8, 1, ' ', tabwriter.Debug)

	fmt.Fprintln(w, "\n [ROUTE TABLE]")
	fmt.Fprintln(w, " METHOD\t PATH\t HANDLER")
	fmt.Fprintln(w, " ------\t ----\t -------")

	for _, r := range routes {
		// 排除 Echo 自动生成的路由（可选）
		if r.Path == "/*" {
			continue
		}

		line := fmt.Sprintf(" %s\t %s\t %s", r.Method, r.Path, r.Name)
		fmt.Fprintln(w, line)
	}

	w.Flush() // 必须调用 Flush 才能写入 stdout
	fmt.Println()
}
