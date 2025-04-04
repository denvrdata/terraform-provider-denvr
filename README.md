# terraform-provider-denvr
A simple terraform provider for managing Denvr resources

[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/denvrdata/terraform-provider-denvr/CI.yml)](https://github.com/denvrdata/terraform-provider-denvr/actions/workflows/CI.yml)
[![Coveralls](https://img.shields.io/coverallsCoverage/github/denvrdata/terraform-provider-denvr)](https://coveralls.io/github/denvrdata/terraform-provider-denvr?branch=main)
[![Denvr Dataworks Docs](https://img.shields.io/badge/denvr_cloud-docs-%234493c5?style=flat)](https://docs.denvrdata.com/docs)
[![Denvr Dataworks Registration](https://img.shields.io/badge/denvr_cloud-registration-%234493c5?style=flat)](https://console.cloud.denvrdata.com/account/register-tenant)

Given a simple terraform configuration, this provider can be used to create and manage Denvr resources.

## Example

```hcl
terraform {
  required_providers {
    denvr = {
      source = "hashicorp.com/denvrdata/denvr"
    }
  }
}

provider "denvr" {}

resource "denvr_vm" "my_denvr_vm" {
    name="terraform-vm"
    rpool="on-demand"
    vpc="denvr-vpc"
    configuration="A100_40GB_PCIe_1x"
    cluster="Msc1"
    ssh_keys=["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC58gLDUttXU67nSxsGxHJtjscN4QT8iyjQFYk9++MFVTaQUnD3D+WUR9eNS/Aj85+swY5wcRyzYyhYb/o+gfy5WyZKC/kpoY+C8EDmcUyt3GeIYjczxP6JY04hEjgseIiZ0wHqr+GMRtGnLIzlX00FdTr5JYbaAWT9qzUVZTeb3U5gyaNBHo8BZDpB1qKThN/4ubWoWwd2Gx010QKX6spsrVdMtrceSglacvzYXogGJblIgJjwjTW0t/kZtmw4ThETLBu7ygG0T0PJFSr8+KD3iFbP9iKmz0v1WgOFZkiNUIuQwhdPBs2kiyKqr3VWE9uPQzss+LGZOgzviMn6E9RQgyMOfPc5sXR636zWUrnnImoPuZo/39gnMoGrAnD/GYbRd/RBG9dI4hUtV3elCKQ3nSybDxREpxLykHQdE5h6L7sMtBzM7SUBklVdYAQx2xfNheR1xWZRMft1r8/jwnvpFYdL6z5TLAXq8Hs8sEw46J2dUvVHo59aGynG13vNbkY14PvQHs1F/obz0oE4aU0s0xSBec8ca+7nYcrTNtlo29nQ7PVtaWF0NovePwsW6fzKqzGCPG6i6gG0IDx86ZBscyLyXhIixpeCTOg4llUr0P3b9vBoNE9X4N8gukpNiMQjowx1Jp3YhJ4v4lqYfCbqUnxJ9VtdUuS49G5pKb3Oxw=="]
    operating_system_image="Ubuntu 22.04.4 LTS"
    personal_storage_mount_path="/home/ubuntu/personal"
    tenant_shared_additional_storage="/home/ubuntu/tenant-shared"
    persist_storage=false
    direct_storage_mount_path="/home/ubuntu/direct-attached"
    root_disk_size=500
}
```
### Creation
```sh
> terraform apply
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/denvrdata/denvr in /Users/rory/go/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  + create

Terraform will perform the following actions:

  # denvr_vm.my_denvr_vm will be created
  + resource "denvr_vm" "my_denvr_vm" {
      + cluster                           = "Msc1"
      + configuration                     = "A100_40GB_PCIe_1x"
      + direct_attached_storage_persisted = false
      + direct_storage_mount_path         = "/home/ubuntu/direct-attached"
      + gpu_type                          = (known after apply)
      + gpus                              = (known after apply)
      + id                                = (known after apply)
      + image                             = (known after apply)
      + ip                                = (known after apply)
      + memory                            = (known after apply)
      + name                              = "terraform-vm"
      + namespace                         = (known after apply)
      + operating_system_image            = "Ubuntu 22.04.4 LTS"
      + persist_storage                   = false
      + personal_storage_mount_path       = "/home/ubuntu/personal"
      + private_ip                        = (known after apply)
      + root_disk_size                    = 500
      + rpool                             = "on-demand"
      + ssh_keys                          = [
          + "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC58gLDUttXU67nSxsGxHJtjscN4QT8iyjQFYk9++MFVTaQUnD3D+WUR9eNS/Aj85+swY5wcRyzYyhYb/o+gfy5WyZKC/kpoY+C8EDmcUyt3GeIYjczxP6JY04hEjgseIiZ0wHqr+GMRtGnLIzlX00FdTr5JYbaAWT9qzUVZTeb3U5gyaNBHo8BZDpB1qKThN/4ubWoWwd2Gx010QKX6spsrVdMtrceSglacvzYXogGJblIgJjwjTW0t/kZtmw4ThETLBu7ygG0T0PJFSr8+KD3iFbP9iKmz0v1WgOFZkiNUIuQwhdPBs2kiyKqr3VWE9uPQzss+LGZOgzviMn6E9RQgyMOfPc5sXR636zWUrnnImoPuZo/39gnMoGrAnD/GYbRd/RBG9dI4hUtV3elCKQ3nSybDxREpxLykHQdE5h6L7sMtBzM7SUBklVdYAQx2xfNheR1xWZRMft1r8/jwnvpFYdL6z5TLAXq8Hs8sEw46J2dUvVHo59aGynG13vNbkY14PvQHs1F/obz0oE4aU0s0xSBec8ca+7nYcrTNtlo29nQ7PVtaWF0NovePwsW6fzKqzGCPG6i6gG0IDx86ZBscyLyXhIixpeCTOg4llUr0P3b9vBoNE9X4N8gukpNiMQjowx1Jp3YhJ4v4lqYfCbqUnxJ9VtdUuS49G5pKb3Oxw==",
        ]
      + status                            = (known after apply)
      + storage                           = (known after apply)
      + storage_type                      = (known after apply)
      + tenancy_name                      = (known after apply)
      + tenant_shared_additional_storage  = "/home/ubuntu/tenant-shared"
      + username                          = (known after apply)
      + vcpus                             = (known after apply)
      + vpc                               = "denvr-vpc"
    }

Plan: 1 to add, 0 to change, 0 to destroy.

Do you want to perform these actions?
  Terraform will perform the actions described above.
  Only 'yes' will be accepted to approve.

  Enter a value: yes

denvr_vm.my_denvr_vm: Creating...
denvr_vm.my_denvr_vm: Creation complete after 8s [id=terraform-vm]

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.
```

### Deletion

```sh
> terraform destroy
╷
│ Warning: Provider development overrides are in effect
│
│ The following provider development overrides are set in the CLI configuration:
│  - hashicorp.com/denvrdata/denvr in /Users/rory/go/bin
│
│ The behavior may therefore not match any released version of the provider and applying changes may cause the state to become incompatible with published releases.
╵
denvr_vm.my_denvr_vm: Refreshing state... [id=terraform-vm]

Terraform used the selected providers to generate the following execution plan. Resource actions are indicated with the following symbols:
  - destroy

Terraform will perform the following actions:

  # denvr_vm.my_denvr_vm will be destroyed
  - resource "denvr_vm" "my_denvr_vm" {
      - cluster                           = "Msc1" -> null
      - configuration                     = "A100_40GB_PCIe_1x" -> null
      - direct_attached_storage_persisted = false -> null
      - direct_storage_mount_path         = "/home/ubuntu/direct-attached" -> null
      - gpu_type                          = "nvidia.com/A100PCIE40GB" -> null
      - gpus                              = 1 -> null
      - id                                = "terraform-vm" -> null
      - image                             = "Ubuntu_22.04.4_LTS" -> null
      - memory                            = 115 -> null
      - name                              = "terraform-vm" -> null
      - namespace                         = "denvr" -> null
      - operating_system_image            = "Ubuntu 22.04.4 LTS" -> null
      - persist_storage                   = false -> null
      - personal_storage_mount_path       = "/home/ubuntu/personal" -> null
      - private_ip                        = "172.16.0.36" -> null
      - root_disk_size                    = 500 -> null
      - rpool                             = "on-demand" -> null
      - ssh_keys                          = [
          - "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC58gLDUttXU67nSxsGxHJtjscN4QT8iyjQFYk9++MFVTaQUnD3D+WUR9eNS/Aj85+swY5wcRyzYyhYb/o+gfy5WyZKC/kpoY+C8EDmcUyt3GeIYjczxP6JY04hEjgseIiZ0wHqr+GMRtGnLIzlX00FdTr5JYbaAWT9qzUVZTeb3U5gyaNBHo8BZDpB1qKThN/4ubWoWwd2Gx010QKX6spsrVdMtrceSglacvzYXogGJblIgJjwjTW0t/kZtmw4ThETLBu7ygG0T0PJFSr8+KD3iFbP9iKmz0v1WgOFZkiNUIuQwhdPBs2kiyKqr3VWE9uPQzss+LGZOgzviMn6E9RQgyMOfPc5sXR636zWUrnnImoPuZo/39gnMoGrAnD/GYbRd/RBG9dI4hUtV3elCKQ3nSybDxREpxLykHQdE5h6L7sMtBzM7SUBklVdYAQx2xfNheR1xWZRMft1r8/jwnvpFYdL6z5TLAXq8Hs8sEw46J2dUvVHo59aGynG13vNbkY14PvQHs1F/obz0oE4aU0s0xSBec8ca+7nYcrTNtlo29nQ7PVtaWF0NovePwsW6fzKqzGCPG6i6gG0IDx86ZBscyLyXhIixpeCTOg4llUr0P3b9vBoNE9X4N8gukpNiMQjowx1Jp3YhJ4v4lqYfCbqUnxJ9VtdUuS49G5pKb3Oxw==",
        ] -> null
      - status                            = "na" -> null
      - storage                           = 1700 -> null
      - storage_type                      = "na" -> null
      - tenancy_name                      = "denvr" -> null
      - tenant_shared_additional_storage  = "/home/ubuntu/tenant-shared" -> null
      - username                          = "rory.finnegan@denvrdata.com" -> null
      - vcpus                             = 10 -> null
      - vpc                               = "denvr-vpc" -> null
        # (1 unchanged attribute hidden)
    }

Plan: 0 to add, 0 to change, 1 to destroy.

Do you really want to destroy all resources?
  Terraform will destroy all your managed infrastructure, as shown above.
  There is no undo. Only 'yes' will be accepted to confirm.

  Enter a value: yes

denvr_vm.my_denvr_vm: Destroying... [id=terraform-vm]
denvr_vm.my_denvr_vm: Destruction complete after 1s

Destroy complete! Resources: 1 destroyed.
```
