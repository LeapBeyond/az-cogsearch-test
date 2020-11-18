# az-cogsearch-test
The Terraform and script assets here are used to standup a rudimentary demonstration of Azure Cognitive Search managed by IaC principles. It grabs some CSV data from https://registry.opendata.aws/humor-detection/, loads that into the constructed Storage Container, then creates the index and loads the data to that index.

While this exmple uses Terraform/Azure ARM and shell scripting, it could well be that a better alternative is to work through [Azure Cloud Shell](https://docs.microsoft.com/en-us/azure/developer/terraform/get-started-cloud-shell), or via custom code in Python. This [article by Ben Keen](https://benalexkeen.com/searching-document-text-at-scale-using-azure-cognitive-search/) shows working through a similar configuration using Python in a Juypter notebook, although some of the APIs have been updated since this example.

It is somewhat disappointing neither Terraform nor Azure ARM covers all parts of the IaC puzzle with respect to Azure, and the hybrid approaches here are not an optimal solution - shell scripting does not allow for sensitivity with respect to potentially fallible API calls.

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
$ openssl pkcs12 -inkey terraformaz.pem -in terraformaz.pem -export -out terraformaz.pfx
```

## Usage

The subdirectories under here provide different ways of implementing the same search demonstration:

- `arm` uses a combination of ARM templates and API calls
- `hybrid` uses a combination of Terraform and API calls

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
