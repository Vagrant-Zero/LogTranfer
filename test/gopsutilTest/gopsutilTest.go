package main

import (
	client "github.com/influxdata/influxdb1-client/v2"
	"github.com/shirou/gopsutil/v3/cpu"
	"log"
	"time"
)

func main() {
	httpClient := connInflux()
	cpuPercent := getCpuInfo()
	ticker := time.Tick(time.Second)
	for {
		select {
		case <-ticker:
			writesPoints(httpClient, cpuPercent)
			log.Println("发送了一条消息!")
		}
	}

}

func getCpuInfo() float64 {
	// 获取cpu使用率
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

func connInflux() client.Client {
	cli, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     "http://127.0.0.1:8086",
		Username: "admin",
		Password: "",
	})
	if err != nil {
		log.Fatal(err)
	}
	return cli
}

func writesPoints(cli client.Client, percent float64) {
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Database:  "monitor",
		Precision: "s", //精度，默认ns
	})
	if err != nil {
		log.Fatal(err)
	}
	tags := map[string]string{"cpu": "cpu0"}
	fields := map[string]interface{}{
		"cpu_percent": percent,
	}

	pt, err := client.NewPoint("cpu_percent", tags, fields, time.Now())
	if err != nil {
		log.Fatal(err)
	}
	bp.AddPoint(pt)
	err = cli.Write(bp)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("insert success")
}
