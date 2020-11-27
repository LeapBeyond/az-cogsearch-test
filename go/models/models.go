// Package models contains various structures to encapsulate the requests and responses being processed
package models

import "github.com/Azure/go-autorest/autorest"

// SearchKey contains the primary and secondary search keys for the search service
type SearchKey struct {
	Primary   string
	Secondary string
}

// Configuration contains the configuration loaded from the supplied configuration file
type Configuration struct {
	BaseName            string `json:"baseName"`
	ServiceName         string `json:"serviceName"`
	TargetLocation      string `json:"targetLocation"`
	SubscriptionId      string `json:"subscriptionId"`
	ServicePrincipal    string `json:"servicePrincipal"`
	ServicePrincipalKey string `json:"servicePrincipalKey"`
	Tenant              string `"json:"tenant"`
	ClientId            string `json:"clientId"`
}

// AzureSession is an object representing session for subscription
type AzureSession struct {
	SubscriptionID     string
	Authorizer         autorest.Authorizer
	ResourceGroupName  string
	SearchServiceName  string
	StorageAccountName string
	ContainerName      string
	ConnectionString   string
	DataFiles          []string
	AccountKeys        SearchKey
	DataSourceName     string
	SearchIndexName    string
	*Configuration
	QueryKey string
}

// DataSource holds the JSON definition of a datasource for the search index
type DataSource struct {
	Name        string        `json:"name"`
	Type        string        `json:"type"`
	Credentials DSCredentials `json:"credentials"`
	Container   DSContainer   `json:"container"`
}

type DSCredentials struct {
	ConnectionString string `json:"connectionString"`
}

type DSContainer struct {
	Name string `json:"name"`
}

type DataSourceResponse struct {
	Context string `json:"@odata.context"`
	Value   []struct {
		Name string `json:"name"`
	} `json:"value"`
}

// SearchIndex is used to describe the index being built in search
type SearchIndex struct {
	Name   string             `json:"name"`
	Fields []SearchIndexField `json:"fields"`
}

type SearchIndexField struct {
	Name           string   `json:"name"`
	Type           string   `json:"type"`
	Facetable      bool     `json:"facetable"`
	Filterable     bool     `json:"filterable"`
	Key            bool     `json:"key"`
	Retrievable    bool     `json:"retrievable"`
	Searchable     bool     `json:"searchable"`
	Sortable       bool     `json:"sortable"`
	Analyzer       *string  `json:"analyzer"`
	IndexAnalyzer  *string  `json:"indexAnalyzer"`
	SearchAnalyzer *string  `json:"searchAnalyzer"`
	SynonymMaps    []string `json:"synonymMaps,omitempty"`
	Fields         []string `json:"fields,omitempty"`
}

// Indexer is used to describe the indexer for searching
// Refer to https://docs.microsoft.com/en-us/rest/api/searchservice/create-indexer
type Indexer struct {
	Name                string                 `json:"name"`
	Description         string                 `json:"description,omitempty"`
	DataSourceName      string                 `json:"dataSourceName"`
	SkillSetName        string                 `json:"skillSetName,omitempty"`
	TargetIndexName     string                 `json:"targetIndexName"`
	Disabled            bool                   `json:"disabled"`
	Schedule            *IndexerSchedule       `json:"schedule,omitempty"`
	Parameters          *IndexerParameters     `json:"parameters,omitempty"`
	FieldMappings       []IndexerFieldMappings `json:"fieldMappings"`
	OutputFieldMappings []IndexerFieldMappings `json:"outputFieldMappings"`
}

type IndexerSchedule struct {
	Interval string `json:"interval"`
	// expected to be YYYY-mm-ddTHH:MM:SSZ, eg 2020-11-27T10:42:49Z
	StartTime string `json:"startTime,omitempty"`
}

type IndexerParameters struct {
	BatchSize              int                    `json:"batchSize,omitempty"`
	MaxFailedItems         int                    `json:"maxFailedItems,omitempty"`
	MaxFailedItemsPerBatch int                    `json:"maxFailedItemsPerBatch,omitempty"`
	ExecutionEnvironment   string                 `json:"executionEnvironment,omitempty"`
	Base64EncodeKeys       string                 `json:"base64EncodeKeys,omitempty"`
	Configuration          *IndexerBlobParameters `json:"configuration,omitempty"`
}

type IndexerBlobParameters struct {
	ParsingMode                                   string `json:"parsingMode"`
	ExcludedFileNameExtensions                    string `json:"excludedFileNameExtensions,omitempty"`
	IndexedFileNameExtensions                     string `json:"indexedFileNameExtensions,omitempty"`
	FailOnUnsupportedContentType                  bool   `json:"failOnUnsupportedContentType,omitempty"`
	FailOnUnprocessableDocument                   bool   `json:"failOnUnprocessableDocument,omitempty"`
	IndexStorageMetadataOnlyForOversizedDocuments bool   `json:"indexStorageMetadataOnlyForOversizedDocuments,omitempty"`
	DelimitedTextHeaders                          string `json:"delimitedTextHeaders"`
	DelimitedTextDelimiter                        string `json:"delimitedTextDelimiter"`
	FirstLineContainsHeaders                      bool   `json:"firstLineContainsHeaders,omitempty"`
	DocumentRoot                                  string `json:"documentRoot,omitempty"`
	DataToExtract                                 string `json:"dataToExtract"`
	ImageAction                                   string `json:"imageAction,omitempty"`
	AllowSkillsetToReadFileData                   bool   `json:"allowSkillsetToReadFileData,omitempty"`
	PdfTextRotationAlgorithm                      string `json:"pdfTextRotationAlgorithm,omitempty"`
}

type IndexerFieldMappings struct {
	SourceFieldName string                  `json:"sourceFieldName"`
	TargetFieldName string                  `json:"targetFieldName"`
	MappingFunction *IndexerMappingFunction `json:"mappingFunction,omitempty"`
}

type IndexerMappingFunction struct {
	Name string `json:"name"`
}

type SearchResult struct {
	Context string `json:"@odata.context"`
	Value   []struct {
		Score              float64 `json:"@search.score"`
		ProductDescription string  `json:"product_description"`
		Question           string  `json:"question"`
	} `json:"value"`
}
