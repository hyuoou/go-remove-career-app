package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"
)

func main() {
	if adb_check() != nil {
		fmt.Println("adbコマンドが見つかりません。")
		os.Exit(1)
	}

	if device_check() != nil {
		fmt.Println("USBデバッグが有効なデバイスが見つかりません。")
		fmt.Println("Android端末が正しく接続されているか確認してください")
		os.Exit(1)
	}

	result, _ := exec.Command("adb", "shell", "pm", "list", "package").Output()
	apps := strings.Split(string(result), "\n")
	var remove_list []string
	for _, v := range apps {
		if app_check(v) {
			remove_list = append(remove_list, strings.Replace(v, "package:", "", 1))
		}
	}
	for _, v := range remove_list {
		fmt.Println(v)
	}
	fmt.Printf("%d個のアプリが削除されます [Y/n]: ", len(remove_list))

	if input() {
		for _, v := range remove_list {
			exec.Command("adb", "shell", "pm", "uninstall", "-k", "--user", "0", v).Run()
		}
		fmt.Println("アプリの消去を実行しました")
	} else {
		fmt.Println("処理を中止しました")
	}
}

func adb_check() error {
	return exec.Command("adb", "--version").Run()
}

func device_check() error {
	exec.Command("adb", "kill-server").Run()
	return exec.Command("adb", "shell", "exit").Run()
}

func app_check(app string) bool {
	reg := regexp.MustCompile(`docomo|ntt|kddi|auone|rakuten|softbank`)
	if reg.MatchString(app) {
		return true
	} else {
		return false
	}
}

func input() bool {
	scan := bufio.NewScanner(os.Stdin)
	for {
		scan.Scan()
		in := scan.Text()
		switch in {
		case "", "y", "Y", "Yes", "yes":
			return true
		default:
			return false
		}
	}
}
