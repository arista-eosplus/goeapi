package main

import (
    "github.com/fredhsu/go-eapi"
    "fmt"
    "github.com/mitchellh/mapstructure"
)

func main() {
	cmds := []string{"show version", "show interfaces"}
	url := "https://admin:admin@dbrl3-leaf1/command-api/"
	jr := eapi.Call(url, cmds)
    var sv eapi.ShowVersion
    err := mapstructure.Decode(jr.Result[0], &sv)
    if err != nil {
        panic(err)
    }
	fmt.Println("\nVersion: ", sv.Version)
    //configCmds := []string{"enable", "configure", "interface ethernet 1", "descr go"}
    configCmds := []string{"enable", "configure", "aaa root secret arista"}
    jr = eapi.Call(url, configCmds)
	fmt.Println("result: ", jr.Result)
	fmt.Println("error: ", jr.Error)
    cmds = []string{"show interfaces ethernet 1"}
    jr = eapi.Call(url, cmds)
    var si eapi.ShowInterfaces
    err = mapstructure.Decode(jr.Result[0], &si)
    if err != nil {
        panic(err)
    }
	fmt.Println("result: ", si) 
	fmt.Println("result: ", si.Interfaces["Ethernet1"].Description)
	fmt.Println("result: ", si.Interfaces["Ethernet1"].InterfaceStatistics)
	fmt.Println("result: ", si.Interfaces["Ethernet1"].Mtu)
	fmt.Println("result: ", si.Interfaces["Ethernet1"].LineProtocolStatus)
}
