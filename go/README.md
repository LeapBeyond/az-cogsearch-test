# go

This example uses Go to demonstrate using a 3GL to create and use the search example. It endeavours to use the [Azure SDK for Go](https://docs.microsoft.com/en-us/azure/developer/go/), but as with ARM and Terraform, this lags well behind the Azure API, so the parts of the example for creating the search index, and using the search, still falls back on using direct HTTP calls to the API.

## Usage

After checking this out you need to create or update a `parameters.json` that is used by the example. The application can use a different specified parameter file for input, but defaults to `parameters.json`. It should resemble:

```
{
  "baseName": "cogsearch",
  "serviceName": "cogsearch-search",
  "targetLocation" : "uksouth",
  "subscriptionId" :  "93b4e6fc-acb0-44ce-bc65-bcfc9b626edc",
  "clientId": "e0d8c6a7-3e6b-4827-a0af-a191818e0ab7",
  "servicePrincipal": "http://terraformaz",
  "servicePrincipalKey": "/Users/robert/.ssh/terraformaz.pfx",
  "tenant": "azureleapbeyond.onmicrosoft.com"
}
```


| field | comment |
|------ | ------- |
| baseName | the base of most names used in the configuration |
| serviceName | the name of the search service to create |
| subscriptionId | the subscription ID |
| clientId | the client or app ID for the service principal (future fixes may be able to find this from the principal) |
| servicePrincipal | name of the Service Principal |
| servicePrincipalKey| path to the PFX used to connect as the Service Principal |
| tenant | Name of the target tenant |

To build the application, first [download and install Go](https://golang.org/doc/install). Next, compile the code:

```
$ go build
$ ls azcogsearch
azcogsearch
```

The tool has command line help:

```
% ./azcogsearch -help
Usage of ./azcogsearch:
  -d	delete example
  -n string
    	config file (default "parameters.json")
  -s	perform a search
```

The default option is to build the assets, which takes about a minute:

```
$ ./azcogsearch

2020/11/27 17:15:36 Starting...
2020/11/27 17:15:36 Subscription ID: 93b4e6fc-acb0-44ce-bc65-bcfc9b626edc
2020/11/27 17:15:36 Setting up...
2020/11/27 17:15:36 Begin creating resource group cogsearch
2020/11/27 17:15:37 resource group created: /subscriptions/93b4e6fc-acb0-44ce-bc65-bcfc9b626edc/resourceGroups/cogsearch
2020/11/27 17:15:37 Begin creating search service cogsearch-search
2020/11/27 17:15:41 search service: cogsearch-search (succeeded)
2020/11/27 17:15:41 Fetching search service keys
2020/11/27 17:15:41 Search service keys fetched
2020/11/27 17:15:41 Begin creating storage account cogsearch
2020/11/27 17:16:05 storage account: cogsearch (Succeeded)
2020/11/27 17:16:05 Begin creating blob container in  cogsearch
2020/11/27 17:16:05 blob storage container /subscriptions/93b4e6fc-acb0-44ce-bc65-bcfc9b626edc/resourceGroups/cogsearch/providers/Microsoft.Storage/storageAccounts/cogsearch/blobServices/default/containers/cogsearch created
2020/11/27 17:16:05 Begining to search for connection string for storage account cogsearch
2020/11/27 17:16:05 Connection String fetched
2020/11/27 17:16:05 Begin fetching data files
2020/11/27 17:16:10 Finish fetching data files
2020/11/27 17:16:11 Blob name: Humorous.csv
2020/11/27 17:16:11 Blob name: Non-humorous-unbiased.csv
2020/11/27 17:16:11 Blob name: Non-humours-biased.csv
2020/11/27 17:16:11 Start creating data source
2020/11/27 17:16:11 Start listing data sources
2020/11/27 17:16:12 finish listing data sources
2020/11/27 17:16:15 Finish creating data source
2020/11/27 17:16:15 Start creating search Index
2020/11/27 17:16:19 Finish creating search Index
2020/11/27 17:16:19 Start creating search indexer
2020/11/27 17:16:21 Finish creating search indexer
2020/11/27 17:16:21 Finishing...
```

searching uses the same tool:

```
$ ./azcogsearch -s true

2020/11/27 17:17:47 Starting...
2020/11/27 17:17:47 Subscription ID: 93b4e6fc-acb0-44ce-bc65-bcfc9b626edc
2020/11/27 17:17:47 Fetching query key
2020/11/27 17:17:48 Finishd fetching query key
2020/11/27 17:17:48 Begin query execution
2020/11/27 17:17:49 Finish query execution
2020/11/27 17:17:49   0   Product : Original BB-8 by Sphero (No Droid Trainer)
2020/11/27 17:17:49   0   Question: Does Donald Trump think BB-8 is a Loser or Hater?
2020/11/27 17:17:49   0   Score   : 16.729845
2020/11/27 17:17:49   1   Product : Donald Trump for President Make America Great Again T Shirt
2020/11/27 17:17:49   1   Question: Will I be perceived as a sore loser if I'm seen wearing this shirt in school?
2020/11/27 17:17:49   1   Score   : 14.836155
2020/11/27 17:17:49 Finishing...
```

as does tearing the resources down:

```
./azcogsearch -d true

2020/11/27 17:18:52 Starting...
2020/11/27 17:18:52 Subscription ID: 93b4e6fc-acb0-44ce-bc65-bcfc9b626edc
2020/11/27 17:18:52 Cleaning up...
2020/11/27 17:18:52 Waiting to delete resource group cogsearch
2020/11/27 17:20:37 Finishing...
```

## Dataset
The dataset contains several CSV files where each row corresponds to a product question related to a product from Amazon. All files have the same structure.

The columns are:
 - question: the question text
 - product_description: short description of the product.
 - image_url: url for the prodcut image.
 - label: 1 if the product question is humorous, 0 otherwise.

## References

 - https://docs.microsoft.com/en-us/rest/api/azure/
 - https://github.com/Azure/azure-sdk-for-go
 - https://github.com/Azure-Samples/azure-sdk-for-go-samples
 - https://docs.microsoft.com/en-us/azure/developer/go/


## License
Copyright 2020 Leap Beyond Emerging Technologies B.V.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
