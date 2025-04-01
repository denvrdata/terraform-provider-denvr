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
