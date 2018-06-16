package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func errChk(e error) {
	if e != nil {
		panic(e)
	}
}

func permissionChk() {
	if os.Getuid() != 0 {
		panic("Must be root to exectuate this program")
	}
}

func getDeviceNames() (string, string) {
	out, e := exec.Command("ip", "a").Output()
	errChk(e)

	wired := regexp.MustCompile(`2: (.*?):`)
	wifi := regexp.MustCompile(`3: (.*?):`)

	return wired.FindStringSubmatch(string(out))[1], wifi.FindStringSubmatch(string(out))[1]
}

func killWpaSupplicant() {
	out, e := exec.Command("ps", "-A").Output()
	errChk(e)
	// 1257 ?        00:00:00 wpa_supplicant
	re := regexp.MustCompile(`(?m)^[^\d]+(\d+).*wpa_supplicant$`)

	if re.MatchString(string(out)) {
		for _, r := range re.FindAllStringSubmatch(string(out), -1) {
			pid := r[1]
			fmt.Println(pid)
			_, e = exec.Command("/usr/bin/kill", "-9", pid).Output()
			errChk(e)
		}
	}
}

func getCurrentDir() string {
	out, err := os.Getwd()
	errChk(err)
	return out
}

func createWpaSupplicantConfig(essid, passwd string) string {
	dir := getCurrentDir()
	file := filepath.Join(dir, "wpa_supplicant.conf")

	if _, err := os.Stat(file); err == nil {
		err = os.Remove(file)
		errChk(err)
	}

	f, err := os.Create(file)
	errChk(err)
	defer f.Close()

	conf, err := exec.Command("/usr/sbin/wpa_passphrase", essid, passwd).Output()
	errChk(err)
	_, err = f.Write(conf)
	errChk(err)
	f.Sync()

	return file
}

func runWpaSupplicantConfig(dev, conf string) {
	_, err := exec.Command("/usr/sbin/wpa_supplicant", "-B", "-Dwext", "-i"+dev, "-c"+conf).Output()
	errChk(err)
}

func randomBetween(min, max int) int {
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}

func formIPAddr(gw string) string {
	// 192.168.1.0
	re := regexp.MustCompile(`(\d+\.\d+\.\d+\.)\d+`)
	// a higher number means less conflict ip
	return re.FindStringSubmatch(gw)[1] + strconv.Itoa(randomBetween(190, 254))
}

func formNetMask(num string) string {
	var nm string
	n, err := strconv.Atoi(num)
	errChk(err)
	for i := 0; i < n/8; i++ {
		nm += "255."
	}
	rest := 32 - n
	for j := 0; j < rest/8; j++ {
		nm += "0."
	}
	return nm[0 : len(nm)-1]
}

func removeExIPAddr(dev string) {
	out, err := exec.Command("/sbin/ip", "a").Output()
	errChk(err)

	re := regexp.MustCompile(`inet (\d+[^/]+\d+)/(\d+) scope global.*` + regexp.QuoteMeta(dev))
	if re.MatchString(string(out)) {
		m := re.FindStringSubmatch(string(out))
		addr := m[1]
		netmask := formNetMask(m[2])
		_, err := exec.Command("/sbin/ip", "a", "delete", addr+"/"+netmask, "dev", dev).Output()
		errChk(err)
	}
}

func setIPAddr(gateway, netmask, dev string) string {
	removeExIPAddr(dev)
	ip := formIPAddr(gateway)
	fmt.Println("using static IP:" + ip)
	_, err := exec.Command("/sbin/ip", "a", "add", ip+"/"+netmask, "dev", dev).Output()
	errChk(err)
	return ip
}

func removeRoute(dev string) {
	//default via 192.168.31.1 dev wlp0s20f0u9 proto dhcp metric 600
	//192.168.31.0/24 dev wlp0s20f0u9 proto kernel scope link src 192.168.31.221 metric 600
	out, err := exec.Command("/sbin/ip", "route", "list").Output()
	errChk(err)
	re := regexp.MustCompile(`dev ` + regexp.QuoteMeta(dev))
	routes := strings.Split(string(out), "\n")
	if len(routes) > 0 {
		for _, r := range routes {
			if re.MatchString(r) {
				// remove linkdown
				cmd := "route del " + strings.Replace(r, " linkdown ", "", -1)
				_, err = exec.Command("/sbin/ip", strings.Split(cmd, " ")[0:]...).Output()
				errChk(err)
			}
		}
	}
}

func countNetMask(nm string) string {
	re := regexp.MustCompile("255")
	num := len(re.FindAllStringSubmatch(nm, -1))
	return strconv.Itoa(num * 8)
}

func findIPRange(gw string) string {
	re := regexp.MustCompile(`(\d+\.\d+\.\d+\.).*`)
	return re.FindStringSubmatch(gw)[1] + "0"
}

func setRoute(ip, gateway, netmask, dev string) {
	removeRoute(dev)
	nm := countNetMask(netmask)
	ran := findIPRange(gateway)
	_, err := exec.Command("/sbin/ip", "route", "add", ran+"/"+nm, "dev", dev, "proto", "kernel", "scope", "link", "src", ip, "metric", "600").Output()
	errChk(err)
	_, err = exec.Command("/sbin/ip", "route", "add", "default", "via", gateway, "dev", dev, "proto", "dhcp", "metric", "600").Output()
	errChk(err)
}

func setDNS() {
	out, err := ioutil.ReadFile("/etc/resolv.conf")
	errChk(err)
	f, err := os.OpenFile("/etc/resolv.conf", os.O_APPEND|os.O_WRONLY, 0644)
	errChk(err)
	defer f.Close()

	re := regexp.MustCompile(`(?m)^nameserver.*?$`)
	if !re.MatchString(string(out)) {
		f.WriteString("nameserver 8.8.8.8\n")
		f.WriteString("nameserver 8.8.4.4\n")
		f.Sync()
	}
}

func main() {
	var dt string
	var gw string
	var nm string
	var essid string
	var passwd string

	flag.StringVar(&dt, "device", "wifi", "device type: wifi or wired")
	flag.StringVar(&gw, "gateway", "", "gateway")
	flag.StringVar(&nm, "netmask", "", "netmask")
	flag.StringVar(&essid, "essid", "", "Your Wi-Fi's name")
	flag.StringVar(&passwd, "password", "", "Your Wi-Fi's password")

	flag.Parse()

	if len(gw) == 0 || len(nm) == 0 {
		panic(`Because there's no dhclient in rescue mode, rescue-network
will try to set up a static network for you. You must give the network
gateway and netmask. You can connect the network with a Phone and see
them in Wi-Fi settings.`)
	}

	if dt == "wifi" && len(essid) == 0 {
		panic(`You chose to connect to a Wi-Fi network. Please give its
ESSID with '--essid='.`)
	}

	if len(passwd) == 0 {
		fmt.Println("Warning: You didn't give your Wi-Fi password. Open Network?")
	}

	permissionChk()

	wired, wifi := getDeviceNames()
	fmt.Println("ethernet:" + wired)
	fmt.Println("wireless:" + wifi)
	dev := wired

	if dt == "wifi" {
		dev = wifi
		killWpaSupplicant()
		conf := createWpaSupplicantConfig(essid, passwd)
		fmt.Println("new wpa_supplicant configuration created at:" + conf)
		runWpaSupplicantConfig(dev, conf)
	}
	ip := setIPAddr(gw, nm, wifi)
	fmt.Println(ip)
	setRoute(ip, gw, nm, wifi)
	setDNS()
	fmt.Println("successfully rescued network.")
}
