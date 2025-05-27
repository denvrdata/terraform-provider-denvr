terraform {
  required_providers {
    denvr = {
      source = "hashicorp.com/denvrdata/denvr"
    }
  }
}

provider "denvr" {}

resource "denvr_app" "terraform_app" {
  name                             = "terraform-app"
  cluster                          = "Msc1"
  hardware_package_name            = "g-nvidia-1xa100-40gb-pcie-14vcpu-112gb"
  application_catalog_item_name    = "jupyter-notebook"
  application_catalog_item_version = "python-3.11.9"
  resource_pool                    = "on-demand"
  jupyter_token                    = "abc123"
}
