package main

import (
	"encoding/json"
	"fmt"
	"github.com/Azure/go-autorest/autorest/azure/auth"
	"io"
	"io/ioutil"
	"leapbeyond.ai/models"
	"log"
	"net/http"
	"os"
	"time"
)

// getConfig tries to load the configuration file that it is pointed at.
// It returns the configuration if possible, otherwise returns an error.
func getConfig(path *string) (*models.Configuration, error) {
	cfg := models.Configuration{}

	file, err := ioutil.ReadFile(*path)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration: %v", err)
	}

	err = json.Unmarshal([]byte(file), &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON from %s: %v", *path, err)
	}

	return &cfg, nil
}

// downloadFile tries to GET from the provided url and write the content to the specified filename.
func downloadFile(url string, fileName string) error {
	response, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("http get on %s failed: %v", url, err)
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("non-200 response from %s", url)
	}

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, response.Body)
	if err != nil {
		return fmt.Errorf("failed to write to %s: %v", fileName, err)
	}

	return nil
}

// fetchDataFiles tries to retrieve the set of specified remote files and write them into the named directory
// it returns a slice of pathnames for the downloaded files
func fetchDataFiles(baseUrl string, files []string, dirName string) ([]string, error) {
	log.Printf("Begin fetching data files")
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		return []string{}, fmt.Errorf("could not make directory %s: %v", dirName, err)
	}

	paths := make([]string, len(files))
	for index, file := range files {
		url := baseUrl + file
		fileName := dirName + "/" + file

		err := downloadFile(url, fileName)
		if err != nil {
			return []string{}, err
		}

		paths[index] = fileName
	}
	log.Printf("Finish fetching data files")
	return paths, nil
}

// makeSession attempts to assemble a "session" object
func makeSession(cfg *models.Configuration) (*models.AzureSession, error) {
	certificateAuthorizer := auth.NewClientCertificateConfig(cfg.ServicePrincipalKey, "", cfg.ClientId, cfg.Tenant)
	authorizer, err := certificateAuthorizer.Authorizer()
	if err != nil {
		return nil, fmt.Errorf("failed to construct an authoriser: %v", err)
	}

	sess := models.AzureSession{
		SubscriptionID: cfg.SubscriptionId,
		Authorizer:     authorizer,
		Configuration:  cfg,
	}

	return &sess, nil
}

// TimeStampNow returns a format timestamp for the current second
func TimeStampNow() string {
	return timeStamp(time.Now())
}

// timeStamp formats the provided time as YYYY-mm-ddTHH:MM:SSZ
func timeStamp(t time.Time) string {
	return t.Format("2006-01-02T03:04:05Z")
}

// create the Search Service Name based on current session details
func makeServiceName(session *models.AzureSession) string {
	return session.ServiceName
}

func makeSearchIndexName(session *models.AzureSession) string {
	return session.BaseName
}

func makeIndexerName(session *models.AzureSession) string {
	return session.BaseName
}

func makeResourceGroupName(session *models.AzureSession) string {
	return session.BaseName
}

func makeStorageAccountName(session *models.AzureSession) string {
	return session.BaseName
}

func makeContainerName(session *models.AzureSession) string {
	return session.BaseName
}

func makeDataSourceName(session *models.AzureSession) string {
	if session.ContainerName == "" {
		panic("Container name needs to be specified before data source")
	}
	return session.ContainerName
}