package main

import (
	"command/service/windows/host"
	"fmt"
)

var DOMAXINS = []string{
	"sellercentral.amazon.com",
	"sellercentral.amazon.co.uk",
	"sellercentral.amazon.de",
	"sellercentral.amazon.es",
	"sellercentral.amazon.fr",
	"sellercentral.amazon.com.au",
	"sellercentral.amazon.co.jp",
	"sellercentral.amazon.ca",
	"sellercentral.amazon.com.mx",
	"sellercentral.amazon.nl",
	"sellercentral.amazon.pl",
	"sellercentral.amazon.se",
	"sellercentral.amazon.sg",
	"sellercentral.amazon.in",
	"sellercentral.amazon.com.br",
	"sellercentral.amazon.it",
	"amazon.com",
	"amazon.co.uk",
	"amazon.de",
	"amazon.es",
	"amazon.fr",
	"amazon.com.au",
	"amazon.co.jp",
	"amazon.ca",
	"amazon.com.mx",
	"amazon.nl",
	"amazon.pl",
	"amazon.se",
	"amazon.sg",
	"amazon.in",
	"amazon.com.br",
	"amazon.it",
}

func main() {
	fmt.Println("开始执行：")
	if host.HasHostCheck("./host") {
		fmt.Println("当前文件夹存在名为 host 的文件，请清理后再执行！")
		finish()
		return
	}
	//cmd.Execute()
	host.AutomaticModificationByDomainArray(DOMAXINS)
	finish()
	//webStart()
}

func finish() {
	var systemIn string
	fmt.Print("输入任意字符结束：")
	scanln, err := fmt.Scanln(&systemIn)
	if err != nil {
		return
	}
	fmt.Println(scanln)
}
