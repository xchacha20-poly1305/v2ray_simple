package utils

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
)

func OpenFile(name string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", name).Start()

	case "windows":
		err = exec.Command("cmd", "/C start "+name).Start()
	case "darwin":
		err = exec.Command("open", name).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err

}

// https://gist.github.com/hyg/9c4afcd91fe24316cbf0
func Openbrowser(url string) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	return err

}

func GetSystemKillChan() <-chan os.Signal {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, os.Interrupt, syscall.SIGTERM) //os.Kill cannot be trapped
	return osSignals
}

func GetSystemProxyState(isSocks5 bool) (ok, enabled bool, addr, port string) {

	switch runtime.GOOS {
	case "darwin":
		var out string
		if isSocks5 {
			var e error
			out, e = LogRunCmd("networksetup", "-getsocksfirewallproxy", "Wi-Fi")
			if e != nil {
				return
			}
			ok = true

		} else {
			var e error
			out, e = LogRunCmd("networksetup", "-getwebproxy", "Wi-Fi")
			if e != nil {
				return
			}
			ok = true
		}
		strs := strings.Split(out, "\n")
		if len(strs) < 1 {
			return
		}

		if strings.Contains(strs[0], "Yes") {
			enabled = true
		}
		if len(strs) < 3 {
			return
		}
		if strings.Contains(strs[1], "Server: ") {
			addr = strings.TrimPrefix(strs[1], "Server: ")
		}
		if strings.Contains(strs[2], "Port: ") {
			port = strings.TrimPrefix(strs[2], "Port: ")
		}

	}
	return
}

func ToggleSystemProxy(isSocks5 bool, addr, port string, enable bool) {
	//我们使用命令行方式。

	//todo: 还可以参考 https://github.com/getlantern/sysproxy ， 这里用了另一种实现，还用到elevate

	switch runtime.GOOS {
	case "darwin":
		if isSocks5 {
			if enable {
				LogRunCmd("networksetup", "-setsocksfirewallproxy", "Wi-Fi", addr, port)

			} else {
				LogRunCmd("networksetup", "-setsocksfirewallproxystate", "Wi-Fi", "off")
			}
		} else {
			if enable {
				LogRunCmd("networksetup", "-setwebproxy", "Wi-Fi", addr, port)
				LogRunCmd("networksetup", "-setsecurewebproxy", "Wi-Fi", addr, port)

			} else {
				LogRunCmd("networksetup", "-setwebproxystate", "Wi-Fi", "off")
				LogRunCmd("networksetup", "-setsecurewebproxystate", "Wi-Fi", "off")
			}
		}
	case "windows":
		const inetSettings = `HKCU\Software\Microsoft\Windows\CurrentVersion\Internet Settings`
		if enable {

			LogRunCmd("reg", "add", inetSettings, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "1", "/f")
			addr = addr + ":" + port

			if isSocks5 {
				LogRunCmd("reg", "add", inetSettings, "/v", "ProxyServer", "/d", "socks="+addr, "/f")

			} else {

				LogRunCmd("reg", "add", inetSettings, "/v", "ProxyServer", "/d", "http="+addr+";https="+addr, "/f")

			}

			LogRunCmd("reg", "add", inetSettings, "/v", "ProxyOverride", "/t", "REG_SZ", "/d", "<-loopback>", "/f")

		} else {
			LogRunCmd("reg", "add", inetSettings, "/v", "ProxyEnable", "/t", "REG_DWORD", "/d", "0", "/f")

			LogRunCmd("reg", "add", inetSettings, "/v", "ProxyServer", "/d", "", "/f")

			LogRunCmd("reg", "delete", inetSettings, "/v", "ProxyOverride", "/f")
		}

	}

}