# arm

This example uses a combination of ARM Templates and API Calls

## Usage

After checking this out, the first thing you need to do is create a configuration file in the working directory called `env.rc`, using this as an example:

```
export CLIENT_CERT=~/.ssh/terraformaz.pfx
export SPNAME=http://terraformaz
export SPKEY=~/.ssh/terraformaz.pem
export TENANT=azureleapbeyond
export TMPLOC=${TMPDIR-/tmp}

```

This file is used by both the scripts to access and configure the environment.

| field | comment |
|------ | ------- |
| CLIENT_CERT | path to the PFX file used to allow Terraform to connect as the Service Principal |
| SPNAME | name of the Service Principal |
| SPKEY | path to the PEM used by scripts to connect as the Service Principal |
| TENANT | Name of the target tenant |
| TMPLOC | a directory on the local filesystem used for storing temporary files |

You may also need to update the default parameters used by ARM:

```
{
  "baseName": {
    "value": "cogsearch"
  },
  "targetLocation" : {
    "value": "uksouth"
  },
  "subscriptionId" : {
    "value": "93b4e6fc-acb0-44ce-bc65-bcfc9b626edc"
  }
}
```

There are three scripts to use. `setup.sh` sets up the assets, `teardown.sh` tears them down, and `search.sh` executes a sample query. These scripts take no command line options:

```
$ ./setup.sh

========================================================================================================
| Fetching data
========================================================================================================
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 1792k  100 1792k    0     0   933k      0  0:00:01  0:00:01 --:--:--  933k
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 1798k  100 1798k    0     0   982k      0  0:00:01  0:00:01 --:--:--  982k
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 1942k  100 1942k    0     0  1035k      0  0:00:01  0:00:01 --:--:-- 1035k
total 12672
-rw-r--r--  1 robert  staff  1835543 18 Nov 13:19 Humorous.csv
-rw-r--r--  1 robert  staff  1841256 18 Nov 13:19 Non-humorous-unbiased.csv
-rw-r--r--  1 robert  staff  1989566 18 Nov 13:19 Non-humours-biased.csv

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
| Create resource group
========================================================================================================
{
   ...
}

========================================================================================================
| Deploy template
========================================================================================================
{
   ...
}

========================================================================================================
| Waiting for creation
========================================================================================================

========================================================================================================
| Fetch connection string and api key
========================================================================================================
Command group 'search' is in preview. It may be changed/removed in a future release.

========================================================================================================
| Load files to blob storage
========================================================================================================
Finished[#############################################################]  100.0000%
{
  "etag": "\"0x8D88BC4CF78AB99\"",
  "lastModified": "2020-11-18T13:21:08+00:00"
}
Finished[#############################################################]  100.0000%
{
  "etag": "\"0x8D88BC4D01D6BF4\"",
  "lastModified": "2020-11-18T13:21:09+00:00"
}
Finished[#############################################################]  100.0000%
{
  "etag": "\"0x8D88BC4D0CA43ED\"",
  "lastModified": "2020-11-18T13:21:10+00:00"
}

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

========================================================================================================
| Logging Out
========================================================================================================
```

Although `setup.sh` pauses at the end to check if the indexer has finished indexing the input data, 1 minute may not be enough for it to finish, so it is recommend that you check the status of the indexer via the Azure console as well. The search will not succeed until the index has been populated by the indexer.

```
$ ./search.sh

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
$ ./teardown.sh

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
| Un-deploy template
========================================================================================================

========================================================================================================
| Waiting for deletion
========================================================================================================

========================================================================================================
| Remove resource group
========================================================================================================

========================================================================================================
| Logging out
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

 - https://docs.microsoft.com/en-us/azure/developer/terraform/overview
 - https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs
 - https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/service_principal_client_certificate
 - https://docs.microsoft.com/en-us/azure/search/query-simple-syntax
 - https://docs.microsoft.com/en-us/azure/search/query-odata-filter-orderby-syntax

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
