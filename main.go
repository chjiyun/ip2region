package main

import (
	"flag"
	jsoniter "github.com/json-iterator/go"
	"github.com/lionsoul2014/ip2region/v1.0/binding/golang/ip2region"
	"log"
	"net/http"
	"os"
)

var port = ""
var region *ip2region.Ip2Region

const dbPath = "./ip2region.db"

type Result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type IpInfo struct {
	Ip       string `json:"ip"`
	Country  string `json:"country"`  // 国家
	Province string `json:"province"` // 省
	City     string `json:"city"`     // 市
	County   string `json:"county"`   // 县、区
	Region   string `json:"region"`   // 区域位置
	ISP      string `json:"isp"`      // 互联网服务提供商
}

// 先于main执行
func init() {
	p := flag.String("p", "7777", "port")
	flag.Parse()
	port = *p

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		panic(err)
	}
	var err error
	region, err = ip2region.New(dbPath)
	if err != nil {
		panic(err)
	}
	defer region.Close()
}

func main() {
	addr := ":" + port
	log.Println("server is running on", "http://127.0.0.1"+addr)
	http.HandleFunc("/", queryIpInfo)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func queryIpInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("content-type", "application/json")
	defer func() {
		if err := recover(); err != nil {
			log.Println("查询ip发生错误", err)
		}
	}()

	if r.URL.Path != "/" {
		w.WriteHeader(404)
		res, _ := jsoniter.Marshal(&Result{Code: 404, Msg: "not found"})
		w.Write(res)
		return
	}

	r.ParseForm()
	ip := r.FormValue("ip")

	if ip == "" {
		res, _ := jsoniter.Marshal(Result{Code: 1, Msg: "ip is required"})
		w.Write(res)
		return
	}

	info, err := region.MemorySearch(ip)
	if err != nil {
		res, _ := jsoniter.Marshal(Result{Code: 1, Msg: err.Error()})
		w.Write(res)
		return
	}

	ipInfo := &IpInfo{
		Ip:       ip,
		ISP:      info.ISP,
		Country:  info.Country,
		Province: info.Province,
		City:     info.City,
		County:   "",
		Region:   info.Region,
	}
	res, _ := jsoniter.Marshal(Result{Code: 0, Data: ipInfo})
	w.Write(res)
}
