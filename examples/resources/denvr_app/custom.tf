terraform {
  required_providers {
    denvr = {
      source = "hashicorp.com/denvrdata/denvr"
    }
  }
}

provider "denvr" {}

resource "denvr_app" "terraform_custom_app" {
  name                         = "terraform-custom-app"
  cluster                      = "Msc1"
  hardware_package_name        = "g-nvidia-1xa100-40gb-pcie-14vcpu-112gb"
  image_cmd_override           = ["nginx"]
  image_repository_hostname    = "https://index.docker.io/v1/"
  image_url                    = "karthequian/helloworld:latest"
  proxy_port                   = 80
  resource_pool                = "reserved-denvr"
  security_context_run_as_root = false
  wait                         = true
}
