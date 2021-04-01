package main

import(
	"os"
	"runtime"
	"syscall"
	"time"
	"fmt"
)

func main()  {
	path:=os.Args[1]
	//timeLayout := "0000-00-00 00:00:00" 
    //loc, _ := time.LoadLocation("Asia/Shanghai")
	timestamp := int64(GetFileCreateTime(path))
	//时间戳转化为日期
	println(timestamp)  
    datetime := time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
    fmt.Println(datetime) 

}


func GetFileCreateTime(path string) int64{
    osType := runtime.GOOS
    fileInfo, _ := os.Stat(path)
    if osType == "linux" {
        stat_t := fileInfo.Sys().(*syscall.Stat_t)
        tCreate := int64(stat_t.Ctim.Sec)
        return tCreate
    }
    return time.Now().Unix()
}