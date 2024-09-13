package esconnector

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/elastic/elastic-transport-go/v8/elastictransport"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/shanluzhineng/configurationx"
	esConfiguration "github.com/shanluzhineng/configurationx/options/elasticsearch"
	"github.com/shanluzhineng/fwpkg/app"
	"github.com/shanluzhineng/fwpkg/system/log"
	jsonUtil "github.com/shanluzhineng/fwpkg/utils/json"
	"go.uber.org/zap"
)

// 初始化es
func InitElasticsearch() {
	defaultOptions := configurationx.GetInstance().Elasticsearch.GetDefaultOptions()
	if defaultOptions == nil {
		log.Logger.Info("没有配置好elasticsearch,无法初始化opevents模块")
		return
	}
	esConfig, currentClient := setupElasticsearchClient(defaultOptions)
	if currentClient == nil || esConfig == nil {
		return
	}
	elasticsearchTypedClient = setupElasticsearchTypedClient(esConfig)
	elasticsearchClient = currentClient
	//注册到ioc中
	app.Context.RegistInstance(currentClient)
	app.Context.RegistInstance(elasticsearchTypedClient)

	log.Logger.Info(fmt.Sprintf("elasticsearch初始化成功,url:%s", strings.Join(defaultOptions.Addresses, ",")))
	//初始化各个索引
	initAllIndex()
}

func setupElasticsearchClient(elasticsearchOptions *esConfiguration.ElasticsearchOptions) (*elasticsearch.Config, *elasticsearch.Client) {
	log.Logger.Info("准备初始化elasticsearch...")

	config := setupElasticsearchConfig(elasticsearchOptions)
	esClient, err := elasticsearch.NewClient(*config)
	if err != nil {
		log.Logger.Error("初始化es时出现异常", zap.Error(err))
		panic(err)
	}
	pingResponse, err := esClient.Ping()
	if err != nil {
		err = fmt.Errorf("初始化es时出现异常,es配置:%s,异常信息:%s", jsonUtil.ObjectToJson(elasticsearchOptions), err.Error())
		log.Logger.Error(err.Error())
		panic(err)
	}
	defer pingResponse.Body.Close()
	if pingResponse.IsError() {
		err = fmt.Errorf("初始化es时出现错误,es配置:%s,错误信息:%s", jsonUtil.ObjectToJson(elasticsearchOptions), pingResponse.String())
		log.Logger.Error(err.Error())
		panic(err)
	}
	return config, esClient
}

func setupElasticsearchConfig(elasticsearchOptions *esConfiguration.ElasticsearchOptions) *elasticsearch.Config {
	config := &elasticsearch.Config{
		Addresses: elasticsearchOptions.Addresses,
		Username:  elasticsearchOptions.Username,
		Password:  elasticsearchOptions.Password,
	}
	//输出debug日志
	if elasticsearchOptions.Debuglog {
		config.Logger = &elastictransport.ColorLogger{
			Output:            os.Stdout,
			EnableRequestBody: true,
		}
	}
	//set transport timeout fields
	if elasticsearchOptions.RequestTimout != nil && *elasticsearchOptions.RequestTimout > 0 {
		transport := http.DefaultTransport.(*http.Transport)
		transport.ResponseHeaderTimeout = (*elasticsearchOptions.RequestTimout) * time.Millisecond
		transport.DialContext = (&net.Dialer{
			Timeout:   (*elasticsearchOptions.RequestTimout) * time.Millisecond,
			KeepAlive: 30 * time.Second,
		}).DialContext
		config.Transport = transport
	}
	return config
}

func setupElasticsearchTypedClient(esConfig *elasticsearch.Config) *elasticsearch.TypedClient {
	esClient, err := elasticsearch.NewTypedClient(*esConfig)
	if err != nil {
		log.Logger.Error("创建elasticsearch.TypedClient时出现异常", zap.Error(err))
		panic(err)
	}
	return esClient
}

func initAllIndex() {
	initOpEventLogIndex()
}
