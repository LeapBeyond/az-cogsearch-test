package main

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/azure-sdk-for-go/profiles/latest/storage/mgmt/storage"
	"github.com/Azure/azure-storage-blob-go/azblob"
	"leapbeyond.ai/http"
	"leapbeyond.ai/models"
	"log"
	"net/url"
	"os"
	"path/filepath"
)

// createStorageAccount tries to setup a storage account
func createStorageAccount(session *models.AzureSession) error {
	name := makeStorageAccountName(session)
	log.Printf("Begin creating storage account %s", name)

	client := storage.NewAccountsClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer

	allowBlobPublicAccess := true
	hnsEnabled := false
	httpsOnly := true
	isTrue := true

	params := storage.AccountCreateParameters{
		Sku:      &storage.Sku{Name: "Standard_GRS", Tier: "Standard"},
		Kind:     "StorageV2",
		Location: &session.TargetLocation,
		Tags: map[string]*string{
			"Name":    &name,
			"Client":  &clientTag,
			"Owner":   &ownerTag,
			"Project": &projectTag,
		},
		AccountPropertiesCreateParameters: &storage.AccountPropertiesCreateParameters{
			MinimumTLSVersion:     "TLS1_2",
			AllowBlobPublicAccess: &allowBlobPublicAccess,
			IsHnsEnabled:          &hnsEnabled,
			NetworkRuleSet: &storage.NetworkRuleSet{
				Bypass:        "AzureServices",
				DefaultAction: "Allow",
			},
			EnableHTTPSTrafficOnly: &httpsOnly,
			Encryption: &storage.Encryption{
				KeySource: "Microsoft.storage",
				Services: &storage.EncryptionServices{
					Blob: &storage.EncryptionService{
						KeyType: "Account",
						Enabled: &isTrue,
					},
					File: &storage.EncryptionService{
						KeyType: "Account",
						Enabled: &isTrue,
					},
				},
			},
			AccessTier: "Hot",
		},
	}

	future, err := client.Create(ctx, session.ResourceGroupName, name, params)
	if err != nil {
		return fmt.Errorf("failed to create stoarage account %s: %v", name, err)
	}
	err = future.WaitForCompletionRef(ctx, client.Client)

	if err != nil {
		return fmt.Errorf("failed on waiting to create storage account %s: %v", name, err)
	}

	session.StorageAccountName = name
	log.Printf("storage account: %s (%s)", name, future.Status())

	return nil
}

// createBlobStorage tries to create a Blob container in the nominated storage account.
// Note that ARM requires creating a blobService into which the container is built, but for some
// bizarre reason the API adds the container directly into the storage accoung.
func createBlobStorage(session *models.AzureSession) error {
	log.Printf("Begin creating blob container in  %s", session.StorageAccountName)

	client := storage.NewBlobContainersClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer

	containerName := makeContainerName(session)
	resourceType := "Microsoft.Storage/storageAccounts/blobServices/containers"
	denyEncryptionScopeOverride := false
	defaultEncryptionScope := "$account-encryption-key"
	blobContainer := storage.BlobContainer{
		Type: &resourceType,
		Name: &containerName,
		ContainerProperties: &storage.ContainerProperties{
			DefaultEncryptionScope:      &defaultEncryptionScope,
			DenyEncryptionScopeOverride: &denyEncryptionScopeOverride,
			PublicAccess:                "none",
		},
	}

	container, err := client.Create(ctx, session.ResourceGroupName, session.StorageAccountName, containerName, blobContainer)
	if err != nil {
		return fmt.Errorf("failed to create blob container %s: %v", containerName, err)
	}

	log.Printf("blob storage container %s created", *container.ID)
	session.ContainerName = containerName
	return nil
}

// getConnectionString tries to assemble the connection string for the blob storage
func getConnectionString(session *models.AzureSession) error {
	log.Printf("Begining to search for connection string for storage account %s", session.StorageAccountName)

	key, err := getAccountKey(session)
	if err != nil {
		return err
	}

	log.Printf("Connection String fetched")
	session.ConnectionString = fmt.Sprintf("DefaultEndpointsProtocol=https;EndpointSuffix=core.windows.net;AccountName=%s;AccountKey=%s", session.StorageAccountName, key)
	return nil
}

// getAccountKey fetches the storage account key
func getAccountKey(session *models.AzureSession) (string, error) {
	client := storage.NewAccountsClient(session.SubscriptionID)
	client.Authorizer = session.Authorizer

	result, err := client.ListKeys(ctx, session.ResourceGroupName, session.StorageAccountName, "kerb")
	if err != nil {
		return "", fmt.Errorf("Failed to list storage account keys for %s: %v", session.StorageAccountName, err)
	}

	if len(*result.Keys) < 1 {
		return "", fmt.Errorf("Expected at least one key for %s", session.StorageAccountName)
	}

	return *(*result.Keys)[0].Value, nil
}

// storeBlobs stores a set of files in the specified container.
// Note that the blob name is the base file name of the provided path
func storeBlobs(session *models.AzureSession) error {
	accountKey, err := getAccountKey(session)
	if err != nil {
		return err
	}

	credential, err := azblob.NewSharedKeyCredential(session.StorageAccountName, accountKey)
	if err != nil {
		return err
	}
	pipeLine := azblob.NewPipeline(credential, azblob.PipelineOptions{})
	URL, _ := url.Parse(fmt.Sprintf("https://%s.blob.core.windows.net/%s", session.StorageAccountName, session.ContainerName))
	containerURL := azblob.NewContainerURL(*URL, pipeLine)
	for _, path := range session.DataFiles {
		fileName := filepath.Base(path)
		blobURL := containerURL.NewBlockBlobURL(fileName)
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		_, err = azblob.UploadFileToBlockBlob(ctx, file, blobURL, azblob.UploadToBlockBlobOptions{
			BlockSize:   4 * 1024 * 1024,
			Parallelism: 16})
	}

	err = listBlobs(containerURL)
	if err != nil {
		return err
	}

	return nil
}

// listBlobs reads the list of blobs in the container
func listBlobs(containerURL azblob.ContainerURL) error {
	for marker := (azblob.Marker{}); marker.NotDone(); {
		listBlob, err := containerURL.ListBlobsFlatSegment(ctx, marker, azblob.ListBlobsSegmentOptions{})
		if err != nil {
			return err
		}
		marker = listBlob.NextMarker

		for _, blobInfo := range listBlob.Segment.BlobItems {
			log.Printf("Blob name: " + blobInfo.Name)
		}
	}
	return nil
}

// createDataSource creates the data source to be used by the search indexer
func createDataSource(session *models.AzureSession) error {
	log.Printf("Start creating data source")

	dsName := makeDataSourceName(session)

	// if datasource name already exists, don't try to recreate it
	names, err := listDataSources(session.SearchServiceName, session.AccountKeys.Primary)
	if err == nil {
		for _, name := range names {
			if name == dsName {
				session.DataSourceName=dsName
				return nil
			}
		}

	}

	searchUrl := fmt.Sprintf("https://%s.search.windows.net/datasources?api-version=2020-06-30", session.SearchServiceName)

	dataSource := models.DataSource{
		Name: dsName,
		Type: "azureblob",
		Credentials: models.DSCredentials{
			ConnectionString: session.ConnectionString,
		},
		Container: models.DSContainer{
			Name: session.ContainerName,
		},
	}

	_, err = http.Post(searchUrl, dataSource, map[string]string{"api-key": session.AccountKeys.Primary})
	if err != nil {
		return fmt.Errorf("Failed to post the request: %v", err)
	}

	session.DataSourceName=dsName
	log.Printf("Finish creating data source")
	return nil
}

func listDataSources(searchServiceName, apiKey string) ([]string, error) {
	log.Printf("Start listing data sources")
	names := []string{}

	searchUrl := fmt.Sprintf("https://%s.search.windows.net/datasources?api-version=2020-06-30&$select=name", searchServiceName)
	body, err := http.Get(searchUrl, map[string]string{"api-key": apiKey})
	if err != nil {
		return names, fmt.Errorf("GET on datasources failed: %v", err)
	}

	dsList := models.DataSourceResponse{}
	err = json.Unmarshal(body, &dsList)
	if err != nil {
		return names, fmt.Errorf("JSON unmarshal on datasources failed: %v", err)
	}

	for _, value := range dsList.Value {
		names = append(names, value.Name)
	}

	log.Printf("finish listing data sources")
	return names, nil
}
