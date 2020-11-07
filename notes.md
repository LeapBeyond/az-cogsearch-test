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
