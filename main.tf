# ----------------------------------------------------------------------------------------------------------------
# resource group to contain everything
# ----------------------------------------------------------------------------------------------------------------
resource azurerm_resource_group rg {
  name     = var.base_name
  location = var.rg_location
  tags     = merge({ "Name" = "CogServTestTwo" }, var.tags)
}

# ----------------------------------------------------------------------------------------------------------------
# storage account for holding the blobs
# ----------------------------------------------------------------------------------------------------------------
resource azurerm_storage_account cogserv {
  name                      = var.base_name
  resource_group_name       = azurerm_resource_group.rg.name
  location                  = azurerm_resource_group.rg.location
  account_tier              = "Standard"
  account_replication_type  = "GRS"
  access_tier               = "Hot"
  enable_https_traffic_only = true
  min_tls_version           = "TLS1_2"
  allow_blob_public_access  = true
  tags                      = merge({ "Name" = "CogServTestTwo" }, var.tags)
}

resource azurerm_storage_container cogserv {
  name                  = var.base_name
  storage_account_name  = azurerm_storage_account.cogserv.name
  container_access_type = "private"
}

resource azurerm_storage_blob files {
  for_each = fileset("${path.root}/data", "*.csv")

  name                   = each.key
  source                 = "${path.root}/data/${each.key}"
  type                   = "Block"
  storage_account_name   = azurerm_storage_account.cogserv.name
  storage_container_name = azurerm_storage_container.cogserv.name
}


# ----------------------------------------------------------------------------------------------------------------
# search service for the index
# ----------------------------------------------------------------------------------------------------------------
resource azurerm_search_service cogserv {
  name                          = var.base_name
  resource_group_name           = azurerm_resource_group.rg.name
  location                      = azurerm_resource_group.rg.location
  public_network_access_enabled = true
  sku                           = "basic"
  tags                          = merge({ "Name" = "CogServTestTwo" }, var.tags)
}
