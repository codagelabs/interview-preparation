package es_query_builder

import (
	"bytes"
	"context"
	"encoding/json"
)

type ElasticSearchQueryBuilder interface {
	WithSize(size int) ElasticSearchQueryBuilder
	WithMatch(matchString string) ElasticSearchQueryBuilder
	BindQuery(buf *bytes.Buffer) ElasticSearchQueryBuilder
	WithMatchAll() ElasticSearchQueryBuilder
	WithMultiMatch(matchString string, fields []string) ElasticSearchQueryBuilder
	WithSuggester(matchString string, fields []string) ElasticSearchQueryBuilder
	Error() error
}

type elasticType map[string]interface{}

func NewElasticSearchQueryBuilder() ElasticSearchQueryBuilder {
	return newElasticQueryBuilder()
}

type elasticQueryBuilder struct {
	esQuery elasticType
	err     error
	ctx     context.Context
}

func (elasticQueryBuilder elasticQueryBuilder) WithMultiMatch(matchString string, fields []string) ElasticSearchQueryBuilder {
	elasticQueryBuilder.esQuery["query"].(elasticType)["multi_match"] = elasticType{
		"query":  matchString,
		"fields": fields,
	}
	return elasticQueryBuilder
}

func (elasticQueryBuilder elasticQueryBuilder) WithMatchAll() ElasticSearchQueryBuilder {

	elasticQueryBuilder.esQuery["query"].(elasticType)["match_all"] = elasticType{}
	return elasticQueryBuilder
}

func (elasticQueryBuilder elasticQueryBuilder) Error() error {
	return elasticQueryBuilder.err
}

func (elasticQueryBuilder elasticQueryBuilder) BindQuery(buf *bytes.Buffer) ElasticSearchQueryBuilder {
	elasticQueryBuilder.err = json.NewEncoder(buf).Encode(elasticQueryBuilder.esQuery)
	return elasticQueryBuilder
}

func (elasticQueryBuilder elasticQueryBuilder) WithMatch(matchString string) ElasticSearchQueryBuilder {
	elasticQueryBuilder.esQuery["query"].(elasticType)["match"] = elasticType{
		"product_name": matchString,
	}
	return elasticQueryBuilder
}

func (elasticQueryBuilder elasticQueryBuilder) WithSize(size int) ElasticSearchQueryBuilder {
	elasticQueryBuilder.esQuery["size"] = size
	return elasticQueryBuilder
}

func (elasticQueryBuilder elasticQueryBuilder) WithSuggester(matchString string, fields []string) ElasticSearchQueryBuilder {
	elasticQueryBuilder.esQuery["query"].(elasticType)["multi_match"] = elasticType{
		"query":  matchString,
		"type":   "bool_prefix",
		"fields": fields,
	}
	return elasticQueryBuilder
}

func newElasticQueryBuilder() *elasticQueryBuilder {
	return &elasticQueryBuilder{esQuery: elasticType{
		"query": elasticType{},
	},
	}
}
