#!/usr/bin/env bash

# make sure we are in the same directory as the script
cd $(dirname $0)
[[ -s ./env.rc ]] || exit 1
source ./env.rc

# terminate script on error
set -e

function banner {
  printf "\n========================================================================================================\n"
  printf "| $1\n"
  printf "========================================================================================================\n"
}

# logout on any exit
trap 'banner "Logging Out"; az logout' EXIT

banner "Logging in"
az login --only-show-errors --service-principal -u $SPNAME -p $SPKEY  --tenant $TENANT.onmicrosoft.com

banner "fetch query key"
API_KEY=$(az search query-key list --resource-group $NAME --service-name $NAME | jq -r ".[0].key")

# search=donald trump+loser&$select=question,product_description&orderby=@search.score
banner "run a query"

QRY='search=donald%20trump%2Bloser&%24select=question%2Cproduct_description&orderby=%40search.score'

curl -H "Content-Type: application/json" -H "api-key: $API_KEY" \
  -s -o $TMPLOC/$$.json \
  -w"\nResponse: %{http_code}\n" \
  "https://$NAME.search.windows.net/indexes/$NAME/docs?api-version=2020-06-30&$QRY"

jq < $TMPLOC/$$.json

rm $TMPLOC/$$.json
