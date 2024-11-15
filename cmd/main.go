package main

import (
	"fmt"
	conf "github.com/kasiforce/trade/config"
	"github.com/kasiforce/trade/pkg/util"
	"github.com/kasiforce/trade/repository/db/dao"
	"github.com/kasiforce/trade/routes"
)

func main() {
	loading()
	r := routes.NewRouter()
	fmt.Println("启动成功...")
	_ = r.Run(conf.Config.System.HttpPort)
}

func loading() {
	conf.InitConfig() //配置文件初始化
	util.InitLog()    //日志文件初始化
	dao.InitMySQL()   //数据库初始化
}
