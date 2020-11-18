#!/usr/bin/env bash

# make sure we are in the same directory as the script
cd $(dirname $0)
[[ -s ./env.rc ]] || exit 1
source ./env.rc

# terminate script on error
set -e

export SUB_ID=$(jq -r ".subscriptionId.value" < parameters.json)
export NAME=$(jq -r ".baseName.value" < parameters.json)
export LOCATION=$(jq -r ".targetLocation.value" < parameters.json)
export TEMPLATE=example.json

function banner {
  printf "\n========================================================================================================\n"
  printf "| $1\n"
  printf "========================================================================================================\n"
}

# logout on any exit
trap 'banner "Logging Out"; az logout' EXIT

# fetch our CSV data
banner "Fetching data"
mkdir -p ./data 2>/dev/null
for FIL in Humorous.csv Non-humorous-unbiased.csv Non-humours-biased.csv
do
  curl -o ./data/$FIL https://humor-detection-pds.s3-us-west-2.amazonaws.com/$FIL
done
ls -l ./data

banner "Logging in"
az login --service-principal -u $SPNAME -p $SPKEY --tenant $TENANT.onmicrosoft.com
az account set  --subscription $SUB_ID

banner "Create resource group"
az deployment sub create \
  --name $NAME \
  --location $LOCATION \
  --template-file resource_group.json \
  --parameters @parameters.json

banner "Deploy template"
az deployment group create \
  --resource-group $NAME \
  --template-file $TEMPLATE \
  --parameters @parameters.json

banner "Waiting for creation"
az deployment group wait \
  --name ${TEMPLATE%".json"} \
  --resource-group $NAME \
  --timeout 120 \
  --created

banner "Fetch connection string and api key"
CONN_STRING=$(az storage account show-connection-string --name $NAME | jq -r ".connectionString")
API_KEY=$(az search admin-key show --resource-group $NAME --service-name ${NAME}-search | jq -r ".primaryKey")

banner "Load files to blob storage"
for FIL in $(ls ./data)
do
  az storage blob upload --container-name $NAME --file data/$FIL --name $FIL --connection-string $CONN_STRING
done

# see https://docs.microsoft.com/en-us/rest/api/searchservice/create-data-source
banner "Set up datasource"
JSON=$TMPLOC/$$.json
cat <<EOF > $JSON
{
  "name": "$NAME",
  "type": "azureblob",
  "credentials": {"connectionString": "$CONN_STRING"},
  "container": {"name": "$NAME"}
}
EOF

curl -H "Content-Type: application/json" -H "api-key: $API_KEY"  -H "Prefer: return=minimal" \
  -w"\nResponse: %{http_code}\n" \
  -X POST \
  -d @$JSON \
  https://${NAME}-search.search.windows.net/datasources?api-version=2020-06-30

# see https://docs.microsoft.com/en-us/rest/api/searchservice/create-index
banner "Set up search index"
cat <<EOF > $JSON
{
  "name": "$NAME",
  "fields": [
    {
      "name": "question", "type": "Edm.String",
      "facetable": false, "filterable": false, "key": false, "retrievable": true, "searchable": true, "sortable": false,
      "analyzer": "standard.lucene", "indexAnalyzer": null, "searchAnalyzer": null,
      "synonymMaps": [], "fields": []
    },
    {
      "name": "product_description", "type": "Edm.String",
      "facetable": false, "filterable": false, "key": false, "retrievable": true, "searchable": true, "sortable": false,
      "analyzer": "standard.lucene", "indexAnalyzer": null, "searchAnalyzer": null,
      "synonymMaps": [], "fields": []
    },
    {
      "name": "image_url", "type": "Edm.String",
      "facetable": false, "filterable": false, "key": false, "retrievable": true, "searchable": false, "sortable": false,
      "analyzer": null, "indexAnalyzer": null, "searchAnalyzer": null,
      "synonymMaps": [], "fields": []
    },
    {
      "name": "label", "type": "Edm.String",
      "facetable": true, "filterable": true, "key": false, "retrievable": true, "searchable": true, "sortable": true,
      "analyzer": "standard.lucene", "indexAnalyzer": null, "searchAnalyzer": null,
      "synonymMaps": [], "fields": []
    },
    {
      "name": "AzureSearch_DocumentKey", "type": "Edm.String",
      "facetable": false, "filterable": false, "key": true, "retrievable": true, "searchable": false, "sortable": false,
      "analyzer": null, "indexAnalyzer": null, "searchAnalyzer": null,
      "synonymMaps": [], "fields": []
    },
    {
      "name": "metadata_storage_content_type", "type": "Edm.String",
      "facetable": false, "filterable": true, "key": false, "retrievable": true, "searchable": false, "sortable": true,
      "analyzer": null, "indexAnalyzer": null, "searchAnalyzer": null,
      "synonymMaps": [], "fields": []
    },
    {
      "name": "metadata_storage_size", "type": "Edm.Int64",
      "facetable": false, "filterable": true, "retrievable": true, "sortable": true, "analyzer": null,
      "indexAnalyzer": null, "searchAnalyzer": null,
      "synonymMaps": [], "fields": []
    },
    {
      "name": "metadata_storage_last_modified", "type": "Edm.DateTimeOffset",
      "facetable": false, "filterable": false, "retrievable": true, "sortable": false, "analyzer": null,
      "indexAnalyzer": null, "searchAnalyzer": null,
      "synonymMaps": [], "fields": []
    },
    {
      "name": "metadata_storage_name", "type": "Edm.String",
      "facetable": false, "filterable": false, "key": false, "retrievable": true, "searchable": false, "sortable": false,
      "analyzer": null, "indexAnalyzer": null, "searchAnalyzer": null,
      "synonymMaps": [], "fields": []
    },
    {
      "name": "metadata_storage_path", "type": "Edm.String",
      "facetable": false, "filterable": false, "key": false, "retrievable": true, "searchable": false, "sortable": false,
      "analyzer": null, "indexAnalyzer": null, "searchAnalyzer": null,
      "synonymMaps": [], "fields": []
    },
    {
      "name": "metadata_storage_file_extension", "type": "Edm.String",
      "facetable": false, "filterable": false, "key": false, "retrievable": true, "searchable": false, "sortable": false,
      "analyzer": null, "indexAnalyzer": null, "searchAnalyzer": null,
      "synonymMaps": [], "fields": []
    }
  ]
}
EOF

curl -H "Content-Type: application/json" -H "api-key: $API_KEY" -H "Prefer: return=minimal" \
  -w"\nResponse: %{http_code}\n" \
  -X PUT \
  -d @$JSON \
  https://${NAME}-search.search.windows.net/indexes/${NAME}?api-version=2020-06-30

NOW=$(date +%Y-%m-%dT%H:%M:%SZ)
banner "Create indexer"
cat <<EOF > $JSON
{
  "name": "$NAME",
  "description": "",
  "dataSourceName": "$NAME",
  "skillsetName": null,
  "targetIndexName": "$NAME",
  "disabled": null,
  "schedule": {
    "interval": "PT1H",
    "startTime": "$NOW"
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
  "outputFieldMappings": []
}
EOF

curl -H "Content-Type: application/json" -H "api-key: $API_KEY" -H "Prefer: return=minimal" \
  -w"\nResponse: %{http_code}\n" \
  -X PUT \
  -d @$JSON \
  https://${NAME}-search.search.windows.net/indexers/${NAME}?api-version=2020-06-30

banner "sleeping 60 seconds to allow indexer to run"
sleep 60

banner "check indexer status"
curl -s -o $TMPLOC/$$.out -H "Content-Type: application/json" -H "api-key: $API_KEY" \
  https://$NAME.search.windows.net/indexers/${NAME}/status?api-version=2020-06-30
jq -r ".lastResult" $TMPLOC/$$.out
rm $TMPLOC/$$.out

rm $JSON
