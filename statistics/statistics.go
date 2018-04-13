package statistics

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"strings"
)

var Record_Name = "stat.json"
var TrafficCh chan string = make(chan string, 2048)

func init() {
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SetConfigName("stat")
	if exist, _ := PathExists(Record_Name); !exist {
		viper.Set("0", 0)
		viper.WriteConfigAs(Record_Name)
	}
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	go DoTrafficStat()
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

func DoTrafficStat() {
	var acc int64 = 0
	for {
		tf := <-TrafficCh
		acc++
		v := strings.Split(tf, ":")
		if len(v) != 2 {
			continue
		}
		port := v[0]
		byteCnt, _ := strconv.ParseInt(v[1], 10, 64)
		lastByteCnt := viper.GetInt64(port) + byteCnt
		viper.Set(port, lastByteCnt)

		if acc%500 == 0 {
			viper.WriteConfigAs(Record_Name)
		}
	}

}

func SendTrafficCnt(accountName string, byteCnt int64) {

	TrafficCh <- fmt.Sprintf("%v:%v", accountName, byteCnt)
}
