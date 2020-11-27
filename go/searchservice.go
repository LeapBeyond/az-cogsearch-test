package main

import (
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/search/mgmt/search"
	"leapbeyond.ai/http"
	"leapbeyond.ai/models"
	"log"
)

// createSearchService tries to setup the search service, returning an error if something went wrong
func createSearchService(session *models.AzureSession) error {
	name := makeServiceName(session)

	log.Printf("Begin creating search service %s", name)

	client := search.NewServicesClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer

	resourceType := "Microsoft.Search/searchServices"

	params := search.Service{
		Name:     &name,
		Type:     &resourceType,
		Location: &session.TargetLocation,
		Sku:      &search.Sku{Name: "basic"},
		Tags: map[string]*string{
			"Name":    &name,
			"Client":  &clientTag,
			"Owner":   &ownerTag,
			"Project": &projectTag,
		},
	}

	future, err := client.CreateOrUpdate(ctx, session.ResourceGroupName, name, params, nil)
	if err != nil {
		return fmt.Errorf("failed to create search service %s: %v", name, err)
	}
	err = future.WaitForCompletionRef(ctx, client.Client)
	if err != nil {
		return fmt.Errorf("failed on waiting to create search service %s: %v", name, err)
	}
	log.Printf("search service: %s (%s)", name, future.Status())

	session.SearchServiceName = name
	return nil
}

// getAccountKeys tries to fetch the api keys for the search service
func getAccountKeys(session *models.AzureSession) error {
	log.Printf("Fetching search service keys")
	client := search.NewAdminKeysClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.Get(ctx, session.ResourceGroupName, session.SearchServiceName, nil)
	if err != nil {
		return fmt.Errorf("failed to retrieve access keys for search Service %s: %v", session.SearchServiceName, err)
	}

	keys := models.SearchKey{
		Primary:   *result.PrimaryKey,
		Secondary: *result.SecondaryKey,
	}
	log.Printf("Search service keys fetched")
	session.AccountKeys = keys
	return nil
}

// getQueryKey tries to get the first available query key.
func getQueryKey(session *models.AzureSession) error {
	log.Printf("Fetching query key")
	client:=search.NewQueryKeysClient(session.SubscriptionID)
	client.Authorizer=session.Authorizer

	result, err := client.ListBySearchService(ctx, session.ResourceGroupName, session.SearchServiceName, nil)
	if err!=nil {
		return fmt.Errorf("Failed to list search query keys: %v", err)
	}

	keys := result.Values()
	if len(keys)<1 {
		return fmt.Errorf("no keys found")
	}
	session.QueryKey = *keys[0].Key
	log.Printf("Finishd fetching query key")

	return nil
}

// createSearchIndex tries to create a search index.
func createSearchIndex(session *models.AzureSession) error {
	log.Printf("Start creating search Index")

	analyzer := "standard.lucene"

	name := makeSearchIndexName(session)
	url := fmt.Sprintf("https://%s.search.windows.net/indexes/%s?api-version=2020-06-30", session.SearchServiceName, name)

	searchIndex := models.SearchIndex{
		Name: name,
		Fields: []models.SearchIndexField{
			{Name: "question", Type: "Edm.String", Retrievable: true, Searchable: true, Analyzer: &analyzer},
			{Name: "product_description", Type: "Edm.String", Retrievable: true, Searchable: true, Analyzer: &analyzer},
			{Name: "image_url", Type: "Edm.String", Retrievable: true},
			{Name: "label", Type: "Edm.String", Facetable: true, Filterable: true, Retrievable: true, Searchable: true, Sortable: true, Analyzer: &analyzer},
			{Name: "AzureSearch_DocumentKey", Type: "Edm.String", Key: true, Retrievable: true},
			{Name: "metadata_storage_content_type", Type: "Edm.String", Filterable: true, Retrievable: true, Sortable: true},
			{Name: "metadata_storage_size", Type: "Edm.Int64", Filterable: true, Retrievable: true, Sortable: true},
			{Name: "metadata_storage_last_modified", Type: "Edm.DateTimeOffset", Retrievable: true},
			{Name: "metadata_storage_name", Type: "Edm.String", Retrievable: true},
			{Name: "metadata_storage_path", Type: "Edm.String", Retrievable: true},
			{Name: "metadata_storage_file_extension", Type: "Edm.String", Retrievable: true},
		},
	}

	_, err := http.Put(url, searchIndex, map[string]string{"api-key": session.AccountKeys.Primary, "Prefer": "return=minimal"})
	if err != nil {
		return fmt.Errorf("Failed to post the request: %v", err)
	}

	session.SearchIndexName = name
	log.Printf("Finish creating search Index")

	return nil
}

// createIndexer builds the search indexer
func createIndexer(session *models.AzureSession) error {
	log.Printf("Start creating search indexer")

	name := makeIndexerName(session)
	url := fmt.Sprintf("https://%s.search.windows.net/indexers/%s?api-version=2020-06-30", session.SearchServiceName, name)

	indexer := models.Indexer{
		Name:            name,
		DataSourceName:  session.DataSourceName,
		TargetIndexName: session.SearchIndexName,
		Schedule: &models.IndexerSchedule{
			Interval:  "PT1H",
			StartTime: TimeStampNow(),
		},
		Parameters: &models.IndexerParameters{
			Configuration: &models.IndexerBlobParameters{
				ParsingMode:              "delimitedText",
				DelimitedTextHeaders:     "",
				DelimitedTextDelimiter:   ",",
				FirstLineContainsHeaders: true,
				DataToExtract:            "contentAndMetadata",
			},
		},
		FieldMappings: []models.IndexerFieldMappings{
			{
				SourceFieldName: "AzureSearch_DocumentKey",
				TargetFieldName: "AzureSearch_DocumentKey", MappingFunction: &models.IndexerMappingFunction{Name: "base64Encode"}},
		},
		OutputFieldMappings: []models.IndexerFieldMappings{},
	}

	_, err := http.Put(url, indexer, map[string]string{"api-key": session.AccountKeys.Primary, "Prefer": "return=minimal"})
	if err != nil {
		return fmt.Errorf("Failed to post the request: %v", err)
	}

	log.Printf("Finish creating search indexer")
	return nil
}

//executeSearch attempts to execute the provided query
func executeSearch(session *models.AzureSession, query string) ([]byte, error) {
	log.Printf("Begin query execution")

	url:=fmt.Sprintf("https://%s.search.windows.net/indexes/%s/docs?api-version=2020-06-30&%s", session.SearchServiceName, session.SearchIndexName, query)

	response,err := http.Get(url, map[string]string{"api-key": session.QueryKey})
	if err!=nil {
		return []byte{}, fmt.Errorf("Failed to get a query result: %v", err)
	}

	log.Printf("Finish query execution")
	return response, nil
}