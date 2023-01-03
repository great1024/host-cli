package host

import (
	"bytes"
	"command/service"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var DNS_SERVER = [...]string{
	"1.1.1.1",
	"8.26.56.26",
	"8.8.8.8",
	"208.67.222.222",
	"9.9.9.9",
	"64.6.64.6",
}

/*
 *https://www.whatsmydns.net/api/details?server=325&type=A&query=amazon.com
 *有用的网站，
 *这是一个没用的程序
 * 默认索引 20 行
 * 首先查询 包含所属域名的 行 有的话就修改没有的话 从最后一行追加
 */
func editHost(delays ...Delay) {
	f, err := os.Open("C:\\Windows\\System32\\drivers\\etc\\hosts")
	if err != nil {
		fmt.Println("read fail")
		//return ""
	}
	defer f.Close()
	var chunk []byte
	buf := make([]byte, 1024)
	for {
		//从file读取到buf中
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("read buf fail", err)
			//return ""
		}
		//说明读取结束
		if n == 0 {
			break
		}
		//读取到最终的缓冲区中
		chunk = append(chunk, buf[:n]...)
	}
	oldHosts := service.ConvertByte2String(chunk, "GB18030")
	hostLineList := strings.Split(oldHosts, "\n")
	index := -1
	for _, delay := range delays {
		for i := 0; i < len(hostLineList); i++ {
			if strings.Contains(hostLineList[i], delay.domain) {
				index = i
			}
		}
		fmt.Println(delay.ip + " " + delay.domain)
		if index > -1 {
			hostLineList[index] = delay.ip + " " + delay.domain
		} else {
			hostLineList = append(hostLineList, "\n"+delay.ip+" "+delay.domain)
		}
	}
	newHost := ""
	for _, s := range hostLineList {
		newHost += s
	}
	newHostBytes := []byte(newHost)
	//fileName :="hosts_"+strconv.FormatInt(time.Now().Unix(),10)
	//filePath := "./"+fileName
	//filePath := "C:\\Windows\\System32\\drivers\\etc\\hosts"
	err1 := ioutil.WriteFile("./hosts", newHostBytes, os.ModeDevice)
	if err1 != nil {
		fmt.Println(err1)
	}
}

type Delay struct {
	domain   string
	ip       string
	delayInt int    // 20
	delay    string // 20ms
	describe string // 描述
}

// 获取IP 地址 的延迟
func networkDelay(ip string) Delay {
	reg := regexp.MustCompile("平均 = \\d+ms{1}")
	regNumber := regexp.MustCompile("\\d+")
	result := runCmd("ping " + ip)
	result = service.ConvertByte2String([]byte(result), "GB18030")
	//fmt.Println(result)
	networkDelays := reg.FindAllStringSubmatch(result, -1)
	var networkDelay = networkDelays[0][0]
	//fmt.Println(networkDelay)
	networkDelayString := networkDelay[9:12]
	p := regNumber.FindAllStringSubmatch(networkDelayString, -1)
	networkDelayInt, err := strconv.Atoi(p[0][0])
	if err != nil {
		fmt.Println(err)
	}
	d := Delay{
		ip:       ip,
		delayInt: networkDelayInt,
		delay:    networkDelay,
		describe: result,
	}
	return d
}
func ipSearch(host string) []Delay {
	reg := regexp.MustCompile("((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})(\\.((2(5[0-5]|[0-4]\\d))|[0-1]?\\d{1,2})){3}")
	var delays []Delay
	set := make(map[string]int)
	for _, dns := range DNS_SERVER {
		result := runCmd("nslookup " + host + " " + dns)
		results := strings.Split(result, " ")
		ips := results[10:]
		for _, ip := range ips {
			regResult := reg.FindStringSubmatch(ip)
			if regResult == nil {
				continue
			}
			ip = regResult[0]
			ip = strings.ReplaceAll(ip, " ", "")
			if ip != "" {
				_, inSet := set[ip]
				if ip != "" && len(ip) != 0 && !inSet {
					set[ip] = 0
					defer func() {
						err := recover()
						if err != nil {
							//fmt.Println(err)
							fmt.Println(host + " " + ip + "无法访问!")
							return
						}
					}()
					delay := networkDelay(ip)
					delay.domain = host
					delays = append(delays, delay)
				}
			}
		}
	}
	//fmt.Println(delays)
	return delays
}

// 根据域名选择最小IP修改Host文件
func AutomaticModification(domain string) {
	ipList := ipSearch(domain)
	sort.Slice(ipList, func(i, j int) bool {
		return ipList[i].delayInt < ipList[j].delayInt
	})
	if len(ipList) > 0 {
		editHost(ipList[0])
	} else {
		fmt.Println("没有找到域名对应的IP地址！")
	}

}

// 根据域名选择最小IP修改Host文件
func AutomaticModificationByDomainArray(domains []string) {
	var all []Delay
	for _, domain := range domains {
		fmt.Println("正在测试：" + domain)
		ipList := ipSearch(domain)
		sort.Slice(ipList, func(i, j int) bool {
			return ipList[i].delayInt < ipList[j].delayInt
		})
		if len(ipList) > 0 {
			all = append(all, ipList[0])
		}

	}
	if len(all) > 0 {
		editHostByArray(all)
	} else {
		fmt.Println("没有找到域名对应的IP地址！")
	}

}

func runCmd(cmdOrder string) string {
	list := strings.Split(cmdOrder, " ")
	cmd := exec.Command(list[0], list[1:]...)
	//c := exec.Command("cmd", "/C", "nslookup", "amazon.com")
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stderr.String()
	} else {
		return out.String()
	}
}

func editHostByArray(delays []Delay) {
	f, err := os.Open("C:\\Windows\\System32\\drivers\\etc\\hosts")
	if err != nil {
		fmt.Println("read fail")
		//return ""
	}
	defer f.Close()
	var chunk []byte
	buf := make([]byte, 1024)
	for {
		//从file读取到buf中
		n, err := f.Read(buf)
		if err != nil && err != io.EOF {
			fmt.Println("read buf fail", err)
			//return ""
		}
		//说明读取结束
		if n == 0 {
			break
		}
		//读取到最终的缓冲区中
		chunk = append(chunk, buf[:n]...)
	}
	oldHosts := service.ConvertByte2String(chunk, "GB18030")
	hostLineList := strings.Split(oldHosts, "\n")
	index := -1
	for _, delay := range delays {
		for i := 0; i < len(hostLineList); i++ {
			if strings.Contains(hostLineList[i], delay.domain) {
				index = i
			}
		}
		fmt.Println(delay.ip + " " + delay.domain)
		if index > -1 {
			hostLineList[index] = delay.ip + " " + delay.domain
		} else {
			hostLineList = append(hostLineList, "\n"+delay.ip+" "+delay.domain)
		}
	}
	newHost := ""
	for _, s := range hostLineList {
		newHost += s
	}
	newHostBytes := []byte(newHost)
	//fileName :="hosts_"+strconv.FormatInt(time.Now().Unix(),10)
	//filePath := "./"+fileName
	//filePath := "C:\\Windows\\System32\\drivers\\etc\\hosts"
	err1 := ioutil.WriteFile("./hosts", newHostBytes, os.ModeDevice)
	if err1 != nil {
		fmt.Println(err1)
	}
}

func HasHostCheck(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}
