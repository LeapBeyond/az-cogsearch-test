variable "tags" {
  type = map(string)

  default = {
    "Owner"   = "Robert"
    "Client"  = "Leap Beyond"
    "Project" = "IU-UK"
  }
}

variable base_name {
  description = "base of all names"
  type    = string
}

variable rg_location {
  description = "location to install into"
  type    = string
  default = "uksouth"
}

/* variables to inject via terraform.tfvars */

variable client_certificate_path {
  type = string
}

variable tenant_id {
  type = string
}

variable subscription_id {
  type = string
}

variable client_id {
  type = string
}
