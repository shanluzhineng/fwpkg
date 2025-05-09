package mongodb

import (
	"fmt"
	"time"

	"github.com/shanluzhineng/configurationx"
	"github.com/shanluzhineng/configurationx/options/mongodb"
	"github.com/shanluzhineng/fwpkg/app"
	"github.com/shanluzhineng/fwpkg/app/web"
	"github.com/shanluzhineng/fwpkg/mongodbr"
	"github.com/shanluzhineng/fwpkg/system/log"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func initMongodbConfigurator(wa web.WebApplication) {
	if app.HostApplication.SystemConfig().App.IsRunInCli {
		return
	}
	initMongodb()
}

// 初始化mongodb
func initMongodb() {
	mongodbOptions := configurationx.GetInstance().Mongodb

	if len(mongodbOptions.MongodbList) <= 0 {
		err := fmt.Errorf("无法初始化mongodb,没有配置好Mongodb的字符串连接")
		log.Logger.Error(err.Error())
		panic(err)
	}
	for eachKey, eachOption := range mongodbOptions.MongodbList {
		var client *mongo.Client
		var err error
		opts := make([]func(*options.ClientOptions), 0)
		if eachOption.EnableCommandMonitor != nil && *eachOption.EnableCommandMonitor {
			// enable command monitor
			opts = append(opts, mongodbr.EnableMongodbMonitor())
		}

		if eachKey == mongodb.AliasName_Default {
			if mongodbr.DefaultClient == nil {
				client, err = mongodbr.SetupDefaultClient(eachOption.Uri, opts...)
				if err != nil {
					log.Logger.Error(err.Error())
					panic(err)
				}
			}
		} else {
			client, err = mongodbr.RegistClient(eachKey, eachOption.Uri, func(co *options.ClientOptions) {})
			if err != nil {
				log.Logger.Error(err.Error())
				panic(err)
			}
		}
		//测试ping
		for {
			err = mongodbr.Ping(client)
			if err == nil {
				break
			}
			log.Logger.Warn(err.Error())
			log.Logger.Warn(fmt.Sprintf("2s后重新测试...，uri: %s", eachOption.Uri))
			time.Sleep(2 * time.Second)
		}
		log.Logger.Info(fmt.Sprintf(">>> mongo init DONE，uri: %s", eachOption.Uri))
	}
}
