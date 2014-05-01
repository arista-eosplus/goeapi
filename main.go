package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
    "github.com/mitchellh/mapstructure"
)

type Parameters struct {
	Version int      `json:"version"`
	Cmds    []string `json:"cmds"`
	Format  string   `json:"format"`
}

type Request struct {
	Jsonrpc string     `json:"jsonrpc"`
	Method  string     `json:"method"`
	Params  Parameters `json:"params"`
	Id      string     `json:"id"`
}

type JsonRpcResponse struct {
	Jsonrpc string        `json:"jsonrpc"`
	Result  []map[string]interface{} `json:"result"`
    Error   map[string]interface{} `json:"error"`
	Id      string        `json:"id"`
}

type ShowVersion struct {
	ModelName        string  `json:"modelName"`
	InternalVersion  string  `json:"internalVersion"`
	SystemMacAddress string  `json:"systemMacAddress"`
	SerialNumber     string  `json:"serialNumber"`
	MemTotal         float64 `json:"memTotal"`
	BootupTimestap   float64 `json:"bootupTimestamp"`
	MemFree          float64 `json:"memFree"`
	Version          string  `json:"version"`
	Architecture     string  `json:"architecture"`
	InternalBuildId  string  `json:"internalBuildId"`
	HardwareRevision string  `json:"hardwareRevision"`
}

type ShowInterfaces struct {
    Interfaces  map[string]Interface  `json:"interfaces"`
}

type Interface struct {
    Description string  `json:"description"`
}

func eapiCall(url string, cmds []string) JsonRpcResponse {
	p := Parameters{1, cmds, "json"}
	req := Request{"2.0", "runCmds", p, "1"}
	buf, err := json.Marshal(req)
	if err != nil {
		panic(err)
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(buf))
	defer resp.Body.Close()

	if err != nil {
		panic(err)
	}
	return decodeEapiResponse(resp)
}

func decodeEapiResponse(resp *http.Response) JsonRpcResponse {
	dec := json.NewDecoder(resp.Body)
	var v JsonRpcResponse
	if err := dec.Decode(&v); err != nil {
		log.Println(err)
	}
	return v
}

func main() {
	cmds := []string{"show version", "show interfaces"}
	url := "http://admin:admin@192.168.56.101/command-api/"
	jr := eapiCall(url, cmds)
    var sv ShowVersion
    err := mapstructure.Decode(jr.Result[0], &sv)
    if err != nil {
        panic(err)
    }
	fmt.Println("\nVersion: ", sv.Version)
    configCmds := []string{"enable", "configure", "interface ethernet 1", "descr go"}
    jr = eapiCall(url, configCmds)
	fmt.Println("result: ", jr.Result)
	fmt.Println("error: ", jr.Error)
    cmds = []string{"show interfaces ethernet 1"}
    jr = eapiCall(url, cmds)
    var si ShowInterfaces
    err = mapstructure.Decode(jr.Result[0], &si)
    if err != nil {
        panic(err)
    }
	fmt.Println("result: ", si) 
	fmt.Println("result: ", si.Interfaces["Ethernet1"].Description)
}
