variable "tags" {
  type = map(string)

  default = {
    "Owner"   = "Robert"
    "Client"  = "Leap Beyond"
    "Project" = "IU-UK"
  }
}

variable base_name {
  type    = string
  default = "cogservtesttwo"
}

variable rg_location {
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
