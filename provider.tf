provider azurerm {
  version                 = "~>2.0"
  client_certificate_path = var.client_certificate_path
  client_id               = var.client_id
  subscription_id         = var.subscription_id
  tenant_id               = var.tenant_id
  features {}
}
