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

# get the access token, and derive the tenant and subscription id
ACCCESS=$(az account get-access-token)
TOKEN=$(jq -r ".accessToken" <<<$ACCCESS)
TENANT_ID=$(jq -r ".tenant" <<<$ACCCESS)
SUB_ID=$(jq -r ".subscription" <<<$ACCCESS)

APP_ID=$(az ad sp show --id $SPNAME | jq -r ".appId")

# setup the terraform.tfvars
cat <<EOF > ./terraform.tfvars
client_certificate_path = "$CLIENT_CERT"
tenant_id               = "$TENANT_ID"
subscription_id         = "$SUB_ID"
client_id               = "$APP_ID"
base_name               = "$NAME"
rg_location             = "$LOCATION"
EOF

# need connection string
banner "Fetch api key"
API_KEY=$(az search admin-key show --resource-group $NAME --service-name $NAME | jq -r ".primaryKey")

banner "Removing indexer"
curl -H "Content-Type: application/json" -H "api-key: $API_KEY" \
  -w"\nResponse: %{http_code}\n" \
  -X DELETE \
  https://$NAME.search.windows.net/indexers/${NAME}?api-version=2020-06-30

banner "Removing index"
curl -H "Content-Type: application/json" -H "api-key: $API_KEY" \
  -w"\nResponse: %{http_code}\n" \
  -X DELETE \
  https://$NAME.search.windows.net/indexes/${NAME}?api-version=2020-06-30

banner "Removing datasource"
curl -H "Content-Type: application/json" -H "api-key: $API_KEY" \
  -w"\nResponse: %{http_code}\n" \
  -X DELETE \
  https://$NAME.search.windows.net/datasources/${NAME}?api-version=2020-06-30


banner "Executing Terraform"
terraform init
terraform destroy --auto-approve
rm ./terraform.tfvars
