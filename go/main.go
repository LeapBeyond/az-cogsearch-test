// Package main provides an example of scripting against Azure.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"leapbeyond.ai/models"
	"log"
	"os"
)

var cfgPath string
var deleteFlag bool
var searchFlag bool
var clientTag = "Leap Beyond"
var ownerTag = "Robert"
var projectTag = "IU-UK"
var ctx = context.Background()

const (
	dataDir = "data"
)

func init() {
	flag.StringVar(&cfgPath, "n", "parameters.json", "config file")
	flag.BoolVar(&deleteFlag, "d", false, "delete example")
	flag.BoolVar(&searchFlag, "s", false, "perform a search")
	flag.Parse()
}

// doDelete performs deletion and clean up operations using the provided session
func doDelete(sess *models.AzureSession) error {
	log.Println("Cleaning up...")
	sess.ResourceGroupName = makeResourceGroupName(sess)

	// remove the downloaded data files
	err := os.RemoveAll(dataDir)
	if err != nil {
		return fmt.Errorf("Failed to remove %s: %v", dataDir, err)
	}

	// destroy the resource group
	err = destroyResourceGroup(sess)
	if err != nil {
		return fmt.Errorf("Failed to remove resource group %s?: %v", sess.ResourceGroupName, err)
	}

	return nil
}

// doCreate tries to setup all the assets using the provided session
func doCreate(session *models.AzureSession) error {
	log.Println("Setting up...")

	// Create a resource group, using the base name from the configuration
	err := createResourceGroup(session)
	if err != nil {
		return err
	}

	// Create the search service
	err = createSearchService(session)
	if err != nil {
		return err
	}

	// fetch the search service keys
	err = getAccountKeys(session)
	if err != nil {
		return err
	}

	// Create the storage account
	err = createStorageAccount(session)
	if err != nil {
		return err
	}

	// add a blob container to the storage account
	err = createBlobStorage(session)
	if err != nil {
		return err
	}

	// fetch the connection string for the storage account
	err = getConnectionString(session)
	if err != nil {
		return err
	}

	// fetch the data files that we are going to load to the index
	files, err := fetchDataFiles("https://humor-detection-pds.s3-us-west-2.amazonaws.com/",
		[]string{
			"Humorous.csv",
			"Non-humorous-unbiased.csv",
			"Non-humours-biased.csv",
		},
		dataDir)
	if err != nil {
		return err
	}
	session.DataFiles = files

	// push the test data files into our storage container
	err = storeBlobs(session)
	if err != nil {
		return err
	}

	// create a data source to be used by the indexer
	err = createDataSource(session)
	if err != nil {
		return err
	}

	// create the index in our search service
	err = createSearchIndex(session)
	if err != nil {
		return err
	}

	// create the indexer
	err = createIndexer(session)
	if err != nil {
		return err
	}
	return nil
}

// doSearch tries to perform a (fixed) search across the created index
func doSearch(sess *models.AzureSession) error {
	sess.ResourceGroupName = makeResourceGroupName(sess)
	sess.SearchServiceName = makeServiceName(sess)
	sess.SearchIndexName = makeSearchIndexName(sess)
	// fetch the search service keys
	err := getQueryKey(sess)
	if err != nil {
		return err
	}

	result, err := executeSearch(sess, "search=donald%20trump%2Bloser&%24select=question%2Cproduct_description&orderby=%40search.score")
	if err != nil {
		return err
	}

	searchResult := models.SearchResult{}
	err = json.Unmarshal(result, &searchResult)
	if err != nil {
		return fmt.Errorf("Failed to read JSON result: %v", err)
	}

	for idx, value := range searchResult.Value {
		log.Printf("  %d   Product : %s", idx, value.ProductDescription)
		log.Printf("  %d   Question: %s", idx, value.Question)
		log.Printf("  %d   Score   : %f", idx, value.Score)
	}

	return nil
}

func main() {
	log.Println("Starting...")

	// Get the configuration from the supplied configuration file
	cfg, err := getConfig(&cfgPath)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Subscription ID: %s\n", cfg.SubscriptionId)

	// Setup a a session object we can use to carry around useful information and resources
	sess, err := makeSession(cfg)
	if err != nil {
		log.Fatal(err)
	}

	if searchFlag {
		err = doSearch(sess)
	} else if deleteFlag {
		err = doDelete(sess)
	} else {
		err = doCreate(sess)
	}

	if err != nil {
		log.Fatal(err)
	}
	log.Println("Finishing...")
}
