package main

import (
  "fmt"
  "flag"
  "os/exec"
  "regexp"
)

func check(e error) {
  if e != nil { panic(e) }
}

func devicenames() (string,string) {
  out, e := exec.Command("ip", "a").Output()
  check(e)

  wifi := regexp.MustCompile(`2: (.*?):`)
  wired := regexp.MustCompile(`3: (.*?):`)

  return wired.FindStringSubmatch(string(out))[1], wifi.FindStringSubmatch(string(out))[1]
}

func check_wpasupplicant() {

}

func main() {
  dt := flag.String("type", "wifi", "device type: wifi or wired")
  gw := flag.String("gateway", "", "gateway")
  nm := flag.String("netmask", "", "netmask")
  essid := flag.String("essid", "", "Your Wi-Fi's name")
  password := flag.String("password", "", "Your Wi-Fi's password")

  flag.Parse()

  if *gw == "" || *nm == "" {
    panic(`Because there's no dhclient in rescue mode, rescue-network
will try to set up a static network for you. You must give the network
gateway and netmask. You can connect the network with a Phone and see
them in Wi-Fi settings.`)
  }

  if *dt == "wifi" && *essid == "" {
    panic(`You chose to connect to a Wi-Fi network. Please give its
ESSID with '--essid='.`)
  }

  if *password == "" {
    fmt.Println("Warning: You didn't give your Wi-Fi password. Open Network?")
  }

  wired, wifi := devicenames()
  fmt.Println(wired)
  fmt.Println(wifi)
}
