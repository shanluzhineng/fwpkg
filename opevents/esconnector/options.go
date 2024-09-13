package esconnector

import "github.com/elastic/go-elasticsearch/v8"

var (
	// Elasticsearch client
	elasticsearchClient      *elasticsearch.Client
	elasticsearchTypedClient *elasticsearch.TypedClient

	// es索引名
	esIndexNames struct {
		OpEventLogIndex string
	} = struct{ OpEventLogIndex string }{
		OpEventLogIndex: "opevent_log",
	}
)
