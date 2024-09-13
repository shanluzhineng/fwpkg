package esconnector

import (
	"fmt"

	"github.com/shanluzhineng/fwpkg/system/log"
	jsonUtil "github.com/shanluzhineng/fwpkg/utils/json"

	"github.com/elastic/go-elasticsearch/v8/esutil"
	"go.uber.org/zap/zapcore"

	opevent "github.com/shanluzhineng/fwpkg/opevents/pkg"
)

// 处理push_task的事件
type OpEventLogESConnectService struct {
}

var _ opevent.IOpEventLogStoreService = (*OpEventLogESConnectService)(nil)

// new一个实例
func NewOpEventLogESConnectService() *OpEventLogESConnectService {
	return &OpEventLogESConnectService{}
}

// 分页搜索事件日志结果
// func (service *opEventLogService) GetPageList(input pkg.OpEventLogSearch) (list []pkg.OpEventLog, total int64, err error) {
// 	boolQuery := elastic.NewBoolQuery()
// 	if len(input.AndroidId) > 0 {
// 		boolQuery = boolQuery.Must(elastic.NewTermsQuery("androidId", input.AndroidId))
// 	}
// 	if input.AccountId > 0 {
// 		boolQuery = boolQuery.Must(elastic.NewTermsQuery("accountId", input.AccountId))
// 	}
// 	if len(input.PushTaskId) > 0 {
// 		boolQuery = boolQuery.Must(elastic.NewTermsQuery("pushTaskId", input.PushTaskId))
// 	}
// 	if len(input.ScriptTaskId) > 0 {
// 		boolQuery = boolQuery.Must(elastic.NewTermsQuery("scriptTaskId", input.ScriptTaskId))
// 	}
// 	if len(input.ScriptTaskInfo.TaskItemId) > 0 {
// 		boolQuery = boolQuery.Must(elastic.NewTermsQuery("scriptTaskInfo.taskItemId", input.ScriptTaskInfo.TaskItemId))
// 	}
// 	body := elastic.NewSearchBodyBuilder().Query(boolQuery)
// 	if len(input.OrderBy) <= 0 {
// 		//构建默认的排序
// 		body = body.SortBy(elastic.SortInfo{Field: "onTimestamp", Ascending: false})
// 	} else {
// 		var orderFieldList []request.Order
// 		if jsonUtil.JsonStringToObject(input.OrderBy, &orderFieldList) != nil && len(orderFieldList) > 0 {
// 			//使用传入的参数来构建索引
// 			for _, eachOrderValue := range orderFieldList {
// 				body = body.SortBy(elastic.SortInfo{Field: eachOrderValue.Field, Ascending: eachOrderValue.Ascending})
// 			}
// 		}
// 	}
// 	bodyString, err := body.BuildBody()
// 	if err != nil {
// 		err = fmt.Errorf("执行es查询时出现异常,%s", err.Error())
// 		log.Logger.Error(err.Error())
// 		return nil, total, err
// 	}

// 	log.Logger.Sugar().Debugf("执行es查询: %s", body)
// 	offset := input.PageSize * (input.Page - 1)
// 	res, err := options.ElasticsearchClient.Search(
// 		options.ElasticsearchClient.Search.WithIndex(options.ESIndexNames.OpEventLogIndex),
// 		options.ElasticsearchClient.Search.WithBody(strings.NewReader(bodyString)),
// 		options.ElasticsearchClient.Search.WithSize(input.PageSize),
// 		options.ElasticsearchClient.Search.WithFrom(offset),
// 	)
// 	if err != nil {
// 		err = fmt.Errorf("执行es查询时出现异常,%s", err.Error())
// 		log.Logger.Error(err.Error())
// 		return nil, total, err
// 	}
// 	defer res.Body.Close()

// 	if res.IsError() {
// 		err = fmt.Errorf("执行es查询时,es返回错误,%s", res.String())
// 		log.Logger.Error(err.Error())
// 		return nil, total, err
// 	}

// 	//解析查询结果
// 	var result elastic.SearchResult
// 	if err = json.NewDecoder(res.Body).Decode(&result); err != nil {
// 		err = fmt.Errorf("解析 Elasticsearch 结果发生错误 %s", err.Error())
// 		log.Logger.Error(err.Error())
// 		return nil, total, err
// 	}
// 	if result.Hits.TotalHits != nil && len(result.Hits.Hits) > 0 {
// 		total = result.Hits.TotalHits.Value
// 	}
// 	if total <= 0 {
// 		return list, total, nil
// 	}

// 	list = make([]pkg.OpEventLog, 0)
// 	for _, eachItem := range result.Hits.Hits {
// 		currentTaskEvent := pkg.OpEventLog{}
// 		err := json.Unmarshal(eachItem.Source, &currentTaskEvent)
// 		if err != nil {
// 			log.Logger.Error(fmt.Sprintf("将es中读取到的日志反序列化成ScriptTaskEvent事出现异常,id:%s,异常信息:%s", eachItem.Id, err.Error()))
// 			continue
// 		}
// 		list = append(list, currentTaskEvent)
// 	}
// 	return
// }

// 插入一条记录
func (service *OpEventLogESConnectService) Save(item *opevent.OpEventLog) error {
	item.BeforeCreate()

	res, err := elasticsearchClient.Index(esIndexNames.OpEventLogIndex,
		esutil.NewJSONReader(&item),
		elasticsearchClient.Index.WithRefresh("true"))
	if err != nil {
		err = fmt.Errorf("在插入push_task的任务事件到es中时出现异常,事件数据:%s,异常信息:%s", jsonUtil.ObjectToJson(item), err.Error())
		log.Logger.Error(err.Error())
		return err
	}
	defer res.Body.Close()
	if res.IsError() {
		err = fmt.Errorf("插入push_task的任务事件到es中时失败,事件数据:%s,错误信息:%s", jsonUtil.ObjectToJson(item), res.String())
		log.Logger.Error(err.Error())
		return err
	}
	return nil

}

// 保存一条调试级别的任务日志
func (service *OpEventLogESConnectService) SaveDebugOpEventLog(message string, opts ...opevent.EventLogOption) (err error) {
	return service.SaveOpEventLog(message, zapcore.DebugLevel, opts...)
}

// 保存一条警告级别的任务日志
func (service *OpEventLogESConnectService) SaveWarnOpEventLog(message string, opts ...opevent.EventLogOption) (err error) {
	return service.SaveOpEventLog(message, zapcore.WarnLevel, opts...)
}

// 保存一条错误级别的任务日志
func (service *OpEventLogESConnectService) SaveErrorOpEventLog(message string, opts ...opevent.EventLogOption) (err error) {
	return service.SaveOpEventLog(message, zapcore.ErrorLevel, opts...)
}

// 保存一条事件日志
func (service *OpEventLogESConnectService) SaveOpEventLog(message string, logLevel zapcore.Level, opts ...opevent.EventLogOption) (err error) {

	newTaskEvent := opevent.NewDefaultOpEventLog(logLevel, message)
	for _, eachOpt := range opts {
		//调用回调
		eachOpt(newTaskEvent)
	}
	return service.Save(newTaskEvent)
}
