#!/usr/bin/env bash

# make sure we are in the same directory as the script
cd $(dirname $0)
[[ -s ./env.rc ]] || exit 1
source ./env.rc

export SUB_ID=$(jq -r ".subscriptionId.value" < parameters.json)
export NAME=$(jq -r ".baseName.value" < parameters.json)
export LOCATION=$(jq -r ".targetLocation.value" < parameters.json)
export TEMPLATE=example.json

function banner {
  printf "\n========================================================================================================\n"
  printf "| $1\n"
  printf "========================================================================================================\n"
}

banner "Logging in"
az login --service-principal -u $SPNAME -p $SPKEY --tenant $TENANT.onmicrosoft.com
az account set --subscription $SUB_ID

# need connection string
banner "Fetch api key"
API_KEY=$(az search admin-key show --resource-group $NAME --service-name ${NAME}-search | jq -r ".primaryKey")

banner "Removing indexer"
curl -H "Content-Type: application/json" -H "api-key: $API_KEY" \
  -w"\nResponse: %{http_code}\n" \
  -X DELETE \
  https://$NAME-search.search.windows.net/indexers/${NAME}?api-version=2020-06-30

banner "Removing index"
curl -H "Content-Type: application/json" -H "api-key: $API_KEY" \
  -w"\nResponse: %{http_code}\n" \
  -X DELETE \
  https://$NAME-search.search.windows.net/indexes/${NAME}?api-version=2020-06-30

banner "Removing datasource"
curl -H "Content-Type: application/json" -H "api-key: $API_KEY" \
  -w"\nResponse: %{http_code}\n" \
  -X DELETE \
  https://$NAME-search.search.windows.net/datasources/${NAME}?api-version=2020-06-30


banner "Un-deploy template"
az deployment group delete \
  --name $TEMPLATE  \
  --resource-group $NAME \
  --subscription $SUB_ID

banner "Waiting for deletion"
az deployment group wait \
  --name $TEMPLATE \
  --resource-group $NAME \
  --timeout 120 \
  --deleted

banner "Remove resource group"
az deployment sub delete \
   --name $NAME

az group delete \
  --name $NAME \
  --yes

banner "Logging out"
az logout
