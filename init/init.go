package init

import (
	"ehelp/x/config"
	"ehelp/x/db/mongodb"
	"ehelp/x/fcm"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"os"
	"path"
)

func init() {
	loadConfig()
	initLog()
	initDB()
	//initCache()
	initFcm()
}

var context *config.Context

func loadConfig() {
	context, _ = config.LoadContext("app.conf", []string{""})
}

func initLog() {
	//config for gin request log
	{
		f, _ := os.Create(path.Join("log", "gin.log"))
		gin.DefaultWriter = io.MultiWriter(f, os.Stdout)
	}
	//config for app log use glog
	{
		logDir, _ := context.String("glog.log_dir")
		logStd, _ := context.String("glog.alsologtostderr")
		flag.Set("alsologtostderr", logStd)
		flag.Set("log_dir", logDir)
		flag.Parse()
	}
}
func initDB() {
	fmt.Println("init db")
	// Read configuration.
	mongodb.MaxPool = context.IntDefault("mongo.maxPool", 0)
	mongodb.PATH, _ = context.String("mongo.path")
	mongodb.DBNAME, _ = context.String("mongo.database")
	mongodb.CheckAndInitServiceConnection()
}

func initFcm() {
	fcm.FCM_SERVER_KEY_CUSTOMER, _ = context.String("fcm.serverkey.customer")
	fcm.FCM_SERVER_KEY_EMPLOYEE, _ = context.String("fcm.serverkey.employee")
	fcm.LINK_AVATAR, _ = context.String("server.avatar")
	fcm.NewFcmApp(fcm.FCM_SERVER_KEY_CUSTOMER, fcm.FCM_SERVER_KEY_EMPLOYEE)
}

// func initCache() {
// 	rest.AssertNil(cache.SetCacheSystem())
// }
