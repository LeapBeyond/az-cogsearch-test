# azure terraform search.

https://docs.microsoft.com/en-us/azure/developer/terraform/overview

Need:   
* resource group
* storage account
* blob container
* blobs
* search service


service principal https://docs.microsoft.com/en-us/cli/azure/create-an-azure-service-principal-azure-cli

1. logged in as root account with `az login`
2. `az ad sp create-for-rbac --name terraformaz --create-cert --scopes="/subscriptions/93b4e6fc-acb0-44ce-bc65-bcfc9b626edc"`
3. moved the generated PEM to ~/.ssh/terraformaz.pem

```
{
  "appId": "e0d8c6a7-3e6b-4827-a0af-a191818e0ab7",
  "displayName": "terraformaz",
  "fileWithCertAndPrivateKey": "/Users/robert/tmpi8_s3cqa.pem",
  "name": "http://terraformaz",
  "password": null,
  "tenant": "b9f789f9-9772-46b0-9b68-ae52a4b6cfec"
}
```

az login --service-principal -u http://terraformaz -p ~/.ssh/terraformaz.pem --tenant azureleapbeyond.onmicrosoft.com

https://nedinthecloud.com/2019/07/16/demystifying-azure-ad-service-principals/

https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs


can't use a PEM for the service principal,. needs to be a pfx = https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/service_principal_client_certificate

openssl pkcs12 -inkey terraformaz.pem -in terraformaz.cert -export -out terraformaz.pfx

buggering bollocks -> index cannot be added via terraform - see https://github.com/rozele/azure-cognitive-search-skills-terraform/blob/main/deploy/terraform/search.tf for an example, have to shove JSON in via curl


note that an alternative may be to work through azure cloud shell -> https://docs.microsoft.com/en-us/azure/developer/terraform/get-started-cloud-shell


test data from https://registry.opendata.aws/humor-detection/

Humor Detection in Product Question Answering Systems Dataset

Questions extracted from the Amazon website.

Please cite:
Humor Detection in Product Question Answering Systems. Yftah Ziser, Elad Kravi & David Carmel. SIGIR 2020.

Authors:
- Yftah Ziser (yftahz@amazon.com)
- Elad Kravi (ekravi@amazon.com)
- David Carmel (dacarmel@amazon.com)

Dataset Structure:
The dataset contains 3 csv file where each row corresponds to a product question.
Humorous.csv - contains the humorous product questions, link for direct download - https://humor-detection-pds.s3-us-west-2.amazonaws.com/Humorous.csv
Non-humorous-unbiased.csv - contains the non-humorous prodcut questions, from the same products as the humorous one, link for direct download -  https://humor-detection-pds.s3-us-west-2.amazonaws.com/Non-humorous-unbiased.csv
Non-humorous-biased.csv - contains the non-humorous prodcut questions, from randomly selected products, link for direct download -  https://humor-detection-pds.s3-us-west-2.amazonaws.com/Non-humours-biased.csv
All files have the same structure
The columns are:
- question: the question text
- product_description: short description of the product.
- image_url: url for the prodcut image.
- label: 1 if the product question is humorous, 0 otherwise.

================
data source definition
```
{
  "@odata.context": "https://cogservtesttwo.search.windows.net/$metadata#datasources/$entity",
  "@odata.etag": "\"0x8D884C431F1A46B\"",
  "name": "cogservtesttwo",
  "description": null,
  "type": "azureblob",
  "subtype": null,
  "credentials": {
    "connectionString": "DefaultEndpointsProtocol=https;AccountName=cogservtesttwo;AccountKey=..."
  },
  "container": {
    "name": "cogservtesttwo",
    "query": null
  },
  "dataChangeDetectionPolicy": null,
  "dataDeletionDetectionPolicy": {
    "@odata.type": "#Microsoft.Azure.Search.SoftDeleteColumnDeletionDetectionPolicy",
    "softDeleteColumnName": "question",
    "softDeleteMarkerValue": "question"
  },
  "encryptionKey": null
}
```

index definition

```
{
  "name": "azureblob-index",
  "fields": [
    {
      "name": "question",
      "type": "Edm.String",
      "facetable": false,
      "filterable": false,
      "key": false,
      "retrievable": true,
      "searchable": true,
      "sortable": false,
      "analyzer": "standard.lucene",
      "indexAnalyzer": null,
      "searchAnalyzer": null,
      "synonymMaps": [],
      "fields": []
    },
    {
      "name": "product_description",
      "type": "Edm.String",
      "facetable": false,
      "filterable": false,
      "key": false,
      "retrievable": true,
      "searchable": true,
      "sortable": false,
      "analyzer": "standard.lucene",
      "indexAnalyzer": null,
      "searchAnalyzer": null,
      "synonymMaps": [],
      "fields": []
    },
    {
      "name": "image_url",
      "type": "Edm.String",
      "facetable": false,
      "filterable": false,
      "key": false,
      "retrievable": true,
      "searchable": false,
      "sortable": false,
      "analyzer": null,
      "indexAnalyzer": null,
      "searchAnalyzer": null,
      "synonymMaps": [],
      "fields": []
    },
    {
      "name": "label",
      "type": "Edm.String",
      "facetable": true,
      "filterable": true,
      "key": false,
      "retrievable": true,
      "searchable": true,
      "sortable": true,
      "analyzer": "standard.lucene",
      "indexAnalyzer": null,
      "searchAnalyzer": null,
      "synonymMaps": [],
      "fields": []
    },
    {
      "name": "AzureSearch_DocumentKey",
      "type": "Edm.String",
      "facetable": false,
      "filterable": false,
      "key": true,
      "retrievable": true,
      "searchable": false,
      "sortable": false,
      "analyzer": null,
      "indexAnalyzer": null,
      "searchAnalyzer": null,
      "synonymMaps": [],
      "fields": []
    },
    {
      "name": "metadata_storage_content_type",
      "type": "Edm.String",
      "facetable": false,
      "filterable": true,
      "key": false,
      "retrievable": true,
      "searchable": false,
      "sortable": true,
      "analyzer": null,
      "indexAnalyzer": null,
      "searchAnalyzer": null,
      "synonymMaps": [],
      "fields": []
    },
    {
      "name": "metadata_storage_size",
      "type": "Edm.Int64",
      "facetable": false,
      "filterable": true,
      "retrievable": true,
      "sortable": true,
      "analyzer": null,
      "indexAnalyzer": null,
      "searchAnalyzer": null,
      "synonymMaps": [],
      "fields": []
    },
    {
      "name": "metadata_storage_last_modified",
      "type": "Edm.DateTimeOffset",
      "facetable": false,
      "filterable": false,
      "retrievable": true,
      "sortable": false,
      "analyzer": null,
      "indexAnalyzer": null,
      "searchAnalyzer": null,
      "synonymMaps": [],
      "fields": []
    },
    {
      "name": "metadata_storage_name",
      "type": "Edm.String",
      "facetable": false,
      "filterable": false,
      "key": false,
      "retrievable": true,
      "searchable": false,
      "sortable": false,
      "analyzer": null,
      "indexAnalyzer": null,
      "searchAnalyzer": null,
      "synonymMaps": [],
      "fields": []
    },
    {
      "name": "metadata_storage_path",
      "type": "Edm.String",
      "facetable": false,
      "filterable": false,
      "key": false,
      "retrievable": true,
      "searchable": false,
      "sortable": false,
      "analyzer": null,
      "indexAnalyzer": null,
      "searchAnalyzer": null,
      "synonymMaps": [],
      "fields": []
    },
    {
      "name": "metadata_storage_file_extension",
      "type": "Edm.String",
      "facetable": false,
      "filterable": false,
      "key": false,
      "retrievable": true,
      "searchable": false,
      "sortable": false,
      "analyzer": null,
      "indexAnalyzer": null,
      "searchAnalyzer": null,
      "synonymMaps": [],
      "fields": []
    }
  ],
  "suggesters": [],
  "scoringProfiles": [],
  "defaultScoringProfile": "",
  "corsOptions": null,
  "analyzers": [],
  "charFilters": [],
  "tokenFilters": [],
  "tokenizers": [],
  "@odata.etag": "\"0x8D884C432FAA26E\""
}
```


indexer definition
```
{
  "@odata.context": "https://cogservtesttwo.search.windows.net/$metadata#indexers/$entity",
  "@odata.etag": "\"0x8D884C433714ECB\"",
  "name": "azureblob-indexer",
  "description": "",
  "dataSourceName": "cogservtesttwo",
  "skillsetName": null,
  "targetIndexName": "azureblob-index",
  "disabled": null,
  "schedule": {
    "interval": "PT1H",
    "startTime": "2020-11-09T15:29:07.845Z"
  },
  "parameters": {
    "batchSize": null,
    "maxFailedItems": 0,
    "maxFailedItemsPerBatch": 0,
    "base64EncodeKeys": null,
    "configuration": {
      "dataToExtract": "contentAndMetadata",
      "parsingMode": "delimitedText",
      "firstLineContainsHeaders": true,
      "delimitedTextDelimiter": ",",
      "delimitedTextHeaders": ""
    }
  },
  "fieldMappings": [
    {
      "sourceFieldName": "AzureSearch_DocumentKey",
      "targetFieldName": "AzureSearch_DocumentKey",
      "mappingFunction": {
        "name": "base64Encode"
      }
    }
  ],
  "outputFieldMappings": [],
  "cache": null,
  "encryptionKey": null
}
```


https://docs.microsoft.com/en-us/azure/search/query-simple-syntax
https://docs.microsoft.com/en-us/azure/search/query-odata-filter-orderby-syntax
search=donald trump+loser&$select=question,product_description&orderby=@search.score
