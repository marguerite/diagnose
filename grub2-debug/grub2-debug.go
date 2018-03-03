package main

import (
  "fmt"
  "encoding/hex"
  "io/ioutil"
  "os"
  "path/filepath"
  "strconv"
  "strings"
)

func check(e error) {
  if e != nil { panic(e) }
}

func detect_uefi() bool {
  if _, e := os.Stat("/sys/firmware/efi"); e != nil {
    return false
  }
  return true
}

func detect_secureboot() bool {
  if detect_uefi() {
    if path, e := filepath.Glob("/sys/firmware/efi/efivars/SecureBoot*"); e == nil {
      bytes, e := ioutil.ReadFile(path[0])
      check(e)
      // 00000000  06 00 00 00 00                                    |.....| 
      str := strings.TrimSpace(strings.Split(hex.Dump(bytes), "|")[0])
      b, e := strconv.ParseBool(string(str[len(str) - 1]))
      check(e)
      return b
    } else {
      return false
    }
  }
  return false
}

func main() {
  fmt.Printf("UEFI Flag: %t\n", detect_uefi())
  fmt.Printf("SecureBoot Flag: %t\n", detect_secureboot())
}
