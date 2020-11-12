# az-cogsearch-test
The Terraform and script assets here are used to standup a rudimentary demonstration of Azure Cognitive Search managed by IaC principles. It grabs some CSV data from https://registry.opendata.aws/humor-detection/, loads that into the constructed Storage Container, then creates the index and loads the data to that index.

While this exmple uses Terraform and shell scripting, it could well be that a better alternative is to work through [Azure Cloud Shell](https://docs.microsoft.com/en-us/azure/developer/terraform/get-started-cloud-shell), or via custom code in Python. This [article by Ben Keen](https://benalexkeen.com/searching-document-text-at-scale-using-azure-cognitive-search/) shows working through a similar configuration using Python in a Juypter notebook, although some of the APIs have been updated since this example.

It is somewhat disappointing that Terraform does not cover all parts of the IaC puzzle with respect to Azure, and the hybrid approach here is not an optimal solution - shell scripting does not allow for sensitivity with respect to potentially fallible API calls.

## Pre-requisites
It is assumed that you are operating these scripts on a Unix environment - they were developed under MacOS. The following are expected to be available in your operating path in order to execute the scripts:

 - [Azure CLI tool](https://docs.microsoft.com/en-us/cli/azure/install-azure-cli) `az` version 2.14.1 or better
 - [jq](https://stedolan.github.io/jq/) version 1.6 or later
 - [Terraform](https://www.terraform.io/downloads.html) version 0.13.4 or better

You may also need [OpenSSL](https://www.openssl.org) (1.1.0h or later)

It is also assumed that you have a Service Principal available to execute against your target Azure Tenant as a `Contributor`. If you have your own access to the tenant as a Contributor, then you can do the following to create a suitable Service Principal, specifying the name of the principal to create, and using the correct subscription ID

  1. Login as a privileged user with `az login`
  2. `az ad sp create-for-rbac --name NAME --create-cert --scopes="/subscriptions/SUBSCRIPTION-ID"`
  3. move the generated PEM to a suitable secure location (e.g. `~/.ssh`) and set the desired permissions (0400) on it

Using the `create-for-rbac` call will report similar to the following. Note the created name is a pseudo-url, and that the PEM is generated in the root of your user account, not the current directory.

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

The created PEM can be used for logging in as the service principal:

```
$ az login --service-principal -u http://terraformaz -p ~/.ssh/terraformaz.pem --tenant azureleapbeyond.onmicrosoft.com
```

but it cannot be used for authentication by Terraform without [converting it to a PFX](https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/service_principal_client_certificate):

```
$ openssl pkcs12 -inkey terraformaz.pem -in terraformaz.cert -export -out terraformaz.pfx
```

## Usage
After checking this out, the first thing you need to do is create a configuration file in the working directory called `env.rc`, using this as an example:

```
export CLIENT_CERT=~/.ssh/terraformaz.pfx
export SPNAME=http://terraformaz
export SPKEY=~/.ssh/terraformaz.pem
export TENANT=azureleapbeyond
export NAME=rahexample
export LOCATION=uksouth
export TMPLOC=${TMPDIR-/tmp}
```

This file is used by both the scripts and Terraform to access and configure the environment.

| field | comment |
|------ | ------- |
| CLIENT_CERT | path to the PFX file used to allow Terraform to connect as the Service Principal |
| SPNAME | name of the Service Principal |
| SPKEY | path to the PEM used by scripts to connect as the Service Principal |
| TENANT | Name of the target tenant |
| NAME | The prefix used for naming all generated assets - this needs to be alphanumeric, with no spaces. keep it short |
| LOCATION | The Azure region to build the assets in |
| TMPLOC | a directory on the local filesystem used for storing temporary files |

There are three scripts to use. `setup.sh` sets up the assets, `teardown.sh` tears them down, and `search.sh` executes a sample query. These scripts take no command line options:

```
$ ./setup.sh

========================================================================================================
| Fetching data
========================================================================================================
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 1792k  100 1792k    0     0   716k      0  0:00:02  0:00:02 --:--:--  716k
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 1798k  100 1798k    0     0  1093k      0  0:00:01  0:00:01 --:--:-- 1093k
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 1942k  100 1942k    0     0  1163k      0  0:00:01  0:00:01 --:--:-- 1163k
total 13304
-rw-r--r--  1 robert  staff  1835543 12 Nov 16:50 Humorous.csv
-rw-r--r--  1 robert  staff  1841256 12 Nov 16:50 Non-humorous-unbiased.csv
-rw-r--r--  1 robert  staff  1989566 12 Nov 16:50 Non-humours-biased.csv

========================================================================================================
| Logging in
========================================================================================================
[
  {
    "cloudName": "AzureCloud",
    "homeTenantId": "b9f789f9-9772-46b0-9b68-ae52a4b6cfec",
    "id": "93b4e6fc-acb0-44ce-bc65-bcfc9b626edc",
    "isDefault": true,
    "managedByTenants": [],
    "name": "Leap Beyond",
    "state": "Enabled",
    "tenantId": "b9f789f9-9772-46b0-9b68-ae52a4b6cfec",
    "user": {
      "name": "http://terraformaz",
      "type": "servicePrincipal"
    }
  }
]

========================================================================================================
| Executing Terraform
========================================================================================================
.
.
.
Apply complete! Resources: 7 added, 0 changed, 0 destroyed.

========================================================================================================
| Fetch connection string and api key
========================================================================================================
Command group 'search' is in preview. It may be changed/removed in a future release.

========================================================================================================
| Set up datasource
========================================================================================================
Response: 204

========================================================================================================
| Set up search index
========================================================================================================
Response: 204

========================================================================================================
| Create indexer
========================================================================================================
Response: 204

========================================================================================================
| sleeping 60 seconds to allow indexer to run
========================================================================================================

========================================================================================================
| check indexer status
========================================================================================================
{
  "status": "inProgress",
  "errorMessage": null,
  "startTime": "2020-11-12T16:51:11.26Z",
  "endTime": null,
  "itemsProcessed": 10000,
  "itemsFailed": 0,
  "initialTrackingState": null,
  "finalTrackingState": null,
  "errors": [],
  "warnings": [],
  "metrics": null
}

========================================================================================================
| Logging Out
========================================================================================================
```

Although `setup.sh` pauses at the end to check if the indexer has finished indexing the input data, 1 minute may not be enough for it to finish, so it is recommend that you check the status of the indexer via the Azure console as well. The search will not succeed until the index has been populated by the indexer.

```
./search.sh

========================================================================================================
| Logging in
========================================================================================================
[
  {
    "cloudName": "AzureCloud",
.
.
.
  }
]

========================================================================================================
| fetch query key
========================================================================================================
Command group 'search' is in preview. It may be changed/removed in a future release.

========================================================================================================
| run a query
========================================================================================================

Response: 200
{
  "@odata.context": "https://rahexample.search.windows.net/indexes('rahexample')/$metadata#docs(*)",
  "value": [
    {
      "@search.score": 16.585403,
      "question": "Does Donald Trump think BB-8 is a Loser or Hater?",
      "product_description": "Original BB-8 by Sphero (No Droid Trainer)"
    },
    {
      "@search.score": 15.260237,
      "question": "Will I be perceived as a sore loser if I'm seen wearing this shirt in school?",
      "product_description": "Donald Trump for President Make America Great Again T Shirt"
    }
  ]
}

========================================================================================================
| Logging Out
========================================================================================================
```

When finished, the assets should be cleaned up:

```
./teardown.sh

========================================================================================================
| Logging in
========================================================================================================
[
  {
    "cloudName": "AzureCloud",
.
.
.
  }
]

========================================================================================================
| Fetch api key
========================================================================================================
Command group 'search' is in preview. It may be changed/removed in a future release.

========================================================================================================
| Removing indexer
========================================================================================================
Response: 204

========================================================================================================
| Removing index
========================================================================================================
Response: 204

========================================================================================================
| Removing datasource
========================================================================================================
Response: 204

========================================================================================================
| Executing Terraform
========================================================================================================
Initializing the backend...
.
.
.
Destroy complete! Resources: 7 destroyed.

========================================================================================================
| Logging Out
========================================================================================================
```

## Dataset
The dataset contains several CSV files where each row corresponds to a product question related to a product from Amazon. All files have the same structure.

The columns are:
 - question: the question text
 - product_description: short description of the product.
 - image_url: url for the prodcut image.
 - label: 1 if the product question is humorous, 0 otherwise.

## References
https://docs.microsoft.com/en-us/azure/developer/terraform/overview
https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs
https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/service_principal_client_certificate
https://docs.microsoft.com/en-us/azure/search/query-simple-syntax
https://docs.microsoft.com/en-us/azure/search/query-odata-filter-orderby-syntax

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
