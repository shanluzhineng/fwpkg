package esconnector

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/shanluzhineng/fwpkg/system/log"
)

var opEventLogIndexMapping = `
{
	"settings":{
		"number_of_shards": 1,
		"number_of_replicas": 0
	},
	"mappings":{
		"properties":{
			"id":{
				"type":"keyword"
			},
			"tenantId":{
				"type":"keyword"
			},
			"reportTimestamp":{
				"type":"long"
			},
			"reportTime":{
				"type":"date"
			},
			"creationTime":{
				"type":"date"
			},
			"onTimestamp":{
				"type":"long"
			},
			"creatorId":{
				"type":"keyword"
			},
			"accountId":{
				"type":"long"
			},
			"logLevel":{
				"type":"keyword"
			},
			"ipAddress":{
				"type":"keyword"
			},
			"appName":{
				"type":"keyword"
			},
			"androidId":{
				"type":"keyword"
			},
			"deviceMobileNo":{
				"type":"keyword"
			},
			"dataFolderName":{
				"type":"keyword"
			},
			"eventMessage":{
				"type":"text"
			},
			"source":{
				"type":"keyword"
			},
			"opAction":{
				"type":"keyword"
			},
			"correlationId":{
				"type":"keyword"
			},

			"taskInfo":{
				"type":"object"
			},
			"sceneTaskInfo":{
				"type":"object"
			}
		}
	}
}
`

// 检测opevent_log索引是否存在
func initOpEventLogIndex() {
	res, err := elasticsearchClient.Indices.Exists([]string{esIndexNames.OpEventLogIndex})
	if err != nil {
		err = fmt.Errorf("向es检测索引是否存在时出现异常: %s", err.Error())
		log.Logger.Error(err.Error())
	}
	if res.StatusCode == http.StatusNotFound {
		//不存在，则创建索引
		createOpEventLogIndex()
	}
	defer res.Body.Close()
}

// 创建索引
func createOpEventLogIndex() {
	res, err := elasticsearchClient.Indices.Create(esIndexNames.OpEventLogIndex,
		elasticsearchClient.Indices.Create.WithBody(strings.NewReader(opEventLogIndexMapping)),
	)
	if err != nil {
		err = fmt.Errorf("无法创建 %s 索引,错误信息:%s", esIndexNames.OpEventLogIndex, err.Error())
		log.Logger.Error(err.Error())
		panic(err)
	}
	defer res.Body.Close()

	if res.IsError() {
		err = fmt.Errorf("创建 %s 索引时,es返回错误,%s", esIndexNames.OpEventLogIndex, res.String())
		log.Logger.Error(err.Error())
		panic(err)
	}
	log.Logger.Info(fmt.Sprintf("创建 %s 索引成功", esIndexNames.OpEventLogIndex))
}
