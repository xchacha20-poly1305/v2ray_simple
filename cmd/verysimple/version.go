/*Package main 读取配置文件，将其内容转化为 proxy.Client和 proxy.Server，然后进行代理转发.

命令行参数请使用 --help查看详情。

如果一个命令行参数无法在标准配置中进行配置，那么它就属于高级选项，或者不推荐的选项，或者正在开发中的功能.


main.go 中，读取配置文件，生成 Dail、Listen 、 RoutePolicy 和 Fallback等 对象后，开始监听，并顺便选择性打开交互模式和 apiServer；

*/
package main

import (
	"fmt"
	"runtime"

	"github.com/e1732a364fed/v2ray_simple/netLayer"
)

const delimiter = "===============================\n"

var Version string = "[version_undefined]" //版本号可由 -ldflags "-X 'main.Version=v1.x.x'" 指定, 本项目的Makefile就是用这种方式确定版本号

func versionStr() string {
	return fmt.Sprintf("verysimple %s, %s %s %s", Version, runtime.Version(), runtime.GOOS, runtime.GOARCH)
}

func printVersion_simple() {
	fmt.Printf(versionStr())
	fmt.Printf("\n")
}

func printVersion() {

	fmt.Printf(delimiter)
	printVersion_simple()
	fmt.Printf(delimiter)

	const desc = "A very simple implementation of V2Ray with some innovation"
	fmt.Printf(desc)
	fmt.Printf("\n")

	if netLayer.HasEmbedGeoip() {
		fmt.Printf("Contains embedded Geoip file\n")
	}
	fmt.Printf(delimiter)

}