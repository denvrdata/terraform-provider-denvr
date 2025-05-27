# Denvr Dataworks Provider

A simple terraform provider for managing Denvr resources

[![GitHub Actions Workflow Status](https://img.shields.io/github/actions/workflow/status/denvrdata/terraform-provider-denvr/CI.yml)](https://github.com/denvrdata/terraform-provider-denvr/actions/workflows/CI.yml)
[![Coveralls](https://img.shields.io/coverallsCoverage/github/denvrdata/terraform-provider-denvr)](https://coveralls.io/github/denvrdata/terraform-provider-denvr?branch=main)
[![Terraform](https://img.shields.io/badge/terraform-latest-blue?logo=terraform&logoColor=white)](https://registry.terraform.io/providers/denvrdata/denvr/latest/docs)
[![Denvr Dataworks Docs](https://img.shields.io/badge/denvr_cloud-docs-%234493c5?style=flat)](https://docs.denvrdata.com/docs)
[![Denvr Dataworks Registration](https://img.shields.io/badge/denvr_cloud-registration-%234493c5?style=flat)](https://console.cloud.denvrdata.com/account/register-tenant)

Given a simple terraform configuration, this provider can be used to create and manage Denvr resources.

## Example

```terraform
terraform {
  required_providers {
    denvr = {
      source = "hashicorp.com/denvrdata/denvr"
    }
  }
}

provider "denvr" {}

resource "denvr_vm" "terraform_vm" {
  name                             = "terraform-vm"
  rpool                            = "reserved-denvr"
  vpc                              = "denvr-vpc"
  configuration                    = "A100_40GB_SXM_1x"
  cluster                          = "Msc1"
  ssh_keys                         = ["ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC58gLDUttXU67nSxsGxHJtjscN4QT8iyjQFYk9++MFVTaQUnD3D+WUR9eNS/Aj85+swY5wcRyzYyhYb/o+gfy5WyZKC/kpoY+C8EDmcUyt3GeIYjczxP6JY04hEjgseIiZ0wHqr+GMRtGnLIzlX00FdTr5JYbaAWT9qzUVZTeb3U5gyaNBHo8BZDpB1qKThN/4ubWoWwd2Gx010QKX6spsrVdMtrceSglacvzYXogGJblIgJjwjTW0t/kZtmw4ThETLBu7ygG0T0PJFSr8+KD3iFbP9iKmz0v1WgOFZkiNUIuQwhdPBs2kiyKqr3VWE9uPQzss+LGZOgzviMn6E9RQgyMOfPc5sXR636zWUrnnImoPuZo/39gnMoGrAnD/GYbRd/RBG9dI4hUtV3elCKQ3nSybDxREpxLykHQdE5h6L7sMtBzM7SUBklVdYAQx2xfNheR1xWZRMft1r8/jwnvpFYdL6z5TLAXq8Hs8sEw46J2dUvVHo59aGynG13vNbkY14PvQHs1F/obz0oE4aU0s0xSBec8ca+7nYcrTNtlo29nQ7PVtaWF0NovePwsW6fzKqzGCPG6i6gG0IDx86ZBscyLyXhIixpeCTOg4llUr0P3b9vBoNE9X4N8gukpNiMQjowx1Jp3YhJ4v4lqYfCbqUnxJ9VtdUuS49G5pKb3Oxw=="]
  operating_system_image           = "Ubuntu 22.04.4 LTS"
  personal_storage_mount_path      = "/home/ubuntu/personal"
  tenant_shared_additional_storage = "/home/ubuntu/tenant-shared"
  persist_storage                  = false
  direct_storage_mount_path        = "/home/ubuntu/direct-attached"
  root_disk_size                   = 500
  wait                             = true
}
```

### Contributing

### Issues

The easiest way to help the terraform-provider-denvr project is to open [issues](https://github.com/denvrdata/terraform-provider-denvr/issues/new/choose).
That being said, here is a list of common details we recommend including:

- Objective
- Code snippets
- Logs
- Stacktrace
- Settings
- Environments (e.g., docker container)

## Pull Requests (PRs)

Our primary workflow for making changes to `terraform-provider-denvr` is with pull requests (PRs).
This process involves:

1. Forking the original repo (from the browser)
2. Clone your fork with `git clone`
3. Make your changes locally
4. Add your modified files with `git add`
5. Commit your changes with a descriptive message `git commit`
6. Push your change back up to your fork `git push`
7. Open a pull request on the original repository (from the browser)

Again, the [GitHub docs](https://docs.github.com/en/pull-requests/collaborating-with-pull-requests/proposing-changes-to-your-work-with-pull-requests/creating-a-pull-request-from-a-fork) already have a detailed summary of this workflow.

When opening a pull request it's best to:

1. Summarize what you're changing and why
2. Keep you changes minimial
2. Include any needed tests and docs
3. Ensure all tests, linting and coverage checks pass

## Checks

We aren't golang or terraform experts, but we do ask that you run a few checks as part of any contributions.

- `gofmt -w .` helps ensure that the code in this repo is consistent and easy to read (if/when things break)
- `go mod tidy` helps us keep our `go.mod` file in sync and removes any unnecessary bloat
- `go test ./...` just checks that our tests are still passing
- `golangci-lint run -D errcheck -D staticcheck` just runs some basic linting (Optional)
  - We typically ignore errcheck and staticcheck as they're a bit noisy while this repo is still in development.

The CI process checks these anyway, but running them locally can save a bit of time :)

### Tips

#### Debug Logs

To get more detail output from either `go test` or `terraform apply`, I usually prefix my command with `TF_LOG=debug`.
```
TF_LOG=DEBUG go test --race --covermode atomic --coverprofile=full-report.cov ./...
```

#### Local go-denvr

If I'm debugging something with the go-denvr SDK I'll typically add something like:
```
replace github.com/denvrdata/go-denvr => /Users/rory/repos/denvrdata/go-denvr
```
to my `go.mod` file.

#### Don't forget to rebuild before running `apply` :)

```
go build -o ~/go/bin/terraform-provider-denvr
```
