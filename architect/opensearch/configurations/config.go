package configurations

type config struct {
	ElasticSearchConfig ElasticSearchConfig `json:"elastic_search_config"`
	LoggerConfig        LoggerConfig        `json:"logger_config"`
	BulkIndexerConfig   BulkIndexerConfig   `json:"bulk_indexer_config"`
}

type ElasticSearchConfig struct {
	HostName string `json:"elastic_search_host"`
	Port     string `json:"elastic_search_port"`
}

type LoggerConfig struct {
	LogLevel string `json:"log_level"`
}
type BulkIndexerConfig struct {
	BatchSize int `json:"batch_size"`
}
