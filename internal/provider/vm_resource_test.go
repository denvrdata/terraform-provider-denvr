package provider

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/denvrdata/go-denvr/result"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// TestVM is a test VM resource model.
// All values prior to GpuType are included in the configuration while
// everything after Vpc is part of the computed response.
var TestVM = vmResourceModel{
	Cluster:                       types.StringValue("Msc1"),
	Configuration:                 types.StringValue("A100_40GB_PCIe_1x"),
	DirectStorageMountPath:        types.StringValue("/home/ubuntu/direct-attached"),
	Name:                          types.StringValue("terraform-vm"),
	OperatingSystemImage:          types.StringValue("Ubuntu 22.04.4 LTS"),
	PersistStorage:                types.BoolValue(false),
	PersonalStorageMountPath:      types.StringValue("/home/ubuntu/personal"),
	RootDiskSize:                  types.Int32Value(500),
	Rpool:                         types.StringValue("on-demand"),
	SshKeys:                       types.ListValueMust(types.StringType, []attr.Value{types.StringValue("ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC58gLbqUnxJ9VtdUuS49G5pKb3Oxw==TEST")}),
	TenantSharedAdditionalStorage: types.StringValue("/home/ubuntu/tenant-shared"),
	Vpc:                           types.StringValue("denvr-vpc"),
	GpuType:                       types.StringValue("nvidia.com/A100PCIE40GB"),
	Gpus:                          types.Int32Value(1),
	Id:                            types.StringValue("terraform-vm"),
	Image:                         types.StringValue("ubuntu-22.04_LTS"),
	Ip:                            types.StringValue("198.16.0.37"),
	Memory:                        types.Int64Value(115),
	Namespace:                     types.StringValue("denvr"),
	PrivateIp:                     types.StringValue("172.16.0.36"),
	Status:                        types.StringValue("na"),
	Storage:                       types.Int64Value(1700),
	StorageType:                   types.StringValue("na"),
	TenancyName:                   types.StringValue("denvr"),
	Username:                      types.StringValue("test@foobar.com"),
	Vcpus:                         types.Int32Value(10),
}

var TestAuthResult = `
{
	"result": {
		"accessToken": "access1",
		"refreshToken": "refresh",
		"expireInSeconds": 600,
		"refreshTokenExpireInSeconds": 3600
	}
}
`

var TestVirtualServerResult = fmt.Sprintf(`
{
	"result": {
		"gpu_type": "%s",
		"gpus": %d,
		"id": "%s",
		"image": "%s",
		"ip": "%s",
		"memory": %d,
		"namespace": "%s",
		"privateIp": "%s",
		"status": "%s",
		"storage": %d,
		"storageType": "%s",
		"tenancy_name": "%s",
		"username": "%s",
		"vcpus": %d
	}
}
`,
	TestVM.GpuType.ValueString(),
	TestVM.Gpus.ValueInt32(),
	TestVM.Id.ValueString(),
	TestVM.Image.ValueString(),
	TestVM.Ip.ValueString(),
	TestVM.Memory.ValueInt64(),
	TestVM.Namespace.ValueString(),
	TestVM.PrivateIp.ValueString(),
	TestVM.Status.ValueString(),
	TestVM.Storage.ValueInt64(),
	TestVM.StorageType.ValueString(),
	TestVM.TenancyName.ValueString(),
	TestVM.Username.ValueString(),
	TestVM.Vcpus.ValueInt32(),
)

var resourceConfig = fmt.Sprintf(`
resource "denvr_vm" "test" {
	name = "%s"
	rpool = "%s"
	vpc = "%s"
	configuration = "%s"
	cluster = "%s"
	ssh_keys = ["%s"]
	operating_system_image = "%s"
	personal_storage_mount_path = "%s"
	tenant_shared_additional_storage = "%s"
	persist_storage = %t
	direct_storage_mount_path = "%s"
	root_disk_size = %d
}
`,
	TestVM.Name.ValueString(),
	TestVM.Rpool.ValueString(),
	TestVM.Vpc.ValueString(),
	TestVM.Configuration.ValueString(),
	TestVM.Cluster.ValueString(),
	TestVM.SshKeys.Elements()[0].(types.String).ValueString(),
	TestVM.OperatingSystemImage.ValueString(),
	TestVM.PersonalStorageMountPath.ValueString(),
	TestVM.TenantSharedAdditionalStorage.ValueString(),
	TestVM.PersistStorage.ValueBool(),
	TestVM.DirectStorageMountPath.ValueString(),
	TestVM.RootDiskSize.ValueInt32(),
)

func TestAccVMResource(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc(
		"/",
		func(resp http.ResponseWriter, req *http.Request) {
			fmt.Printf("Request received: %s %s\n", req.Method, req.URL.Path)
			resp.WriteHeader(http.StatusNotFound)
			resp.Write([]byte(`{"error": {"message": "Path not found"}}`))
		},
	)
	mux.HandleFunc(
		"/api/TokenAuth/Authenticate",
		func(resp http.ResponseWriter, req *http.Request) {
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte(TestAuthResult))
		},
	)
	mux.HandleFunc(
		"/api/v1/servers/virtual/CreateServer",
		func(resp http.ResponseWriter, req *http.Request) {
			//fmt.Println(TestVirtualServerResult)
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte(TestVirtualServerResult))
		},
	)
	mux.HandleFunc(
		"/api/v1/servers/virtual/DestroyServer",
		func(resp http.ResponseWriter, req *http.Request) {
			//fmt.Println(TestVirtualServerResult)
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte(TestVirtualServerResult))
		},
	)
	server := httptest.NewServer(mux)
	defer server.Close()

	content := fmt.Sprintf(
		`[defaults]
        server = "%s"
        api = "v2"
        cluster = "Hou1"
        tenant = "denvr"
        vpcid = "denvr"
        rpool = "reserved-denvr"
        retries = 5

        [credentials]
        username = "test@foobar.com"
        password = "test.foo.bar.baz"`,
		server.URL,
	)

	f := result.Wrap(os.CreateTemp("", "test-newconfig-tmpfile-")).Unwrap()
	defer f.Close()
	defer os.Remove(f.Name())
	result.Wrap(f.Write([]byte(content))).Unwrap()

	// Use the DENVR_CONFIG environment variable for our tests
	os.Setenv("DENVR_CONFIG", f.Name())

	resource.Test(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: providerConfig + resourceConfig,
				Check: resource.ComposeTestCheckFunc(
					// Verify provided values
					resource.TestCheckResourceAttr("denvr_vm.test", "name", "terraform-vm"),
					resource.TestCheckResourceAttr("denvr_vm.test", "rpool", "on-demand"),
					resource.TestCheckResourceAttr("denvr_vm.test", "vpc", "denvr-vpc"),
					resource.TestCheckResourceAttr("denvr_vm.test", "configuration", "A100_40GB_PCIe_1x"),
					resource.TestCheckResourceAttr("denvr_vm.test", "cluster", "Msc1"),
					resource.TestCheckResourceAttr("denvr_vm.test", "ssh_keys.0", "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACAQC58gLbqUnxJ9VtdUuS49G5pKb3Oxw==TEST"),
					resource.TestCheckResourceAttr("denvr_vm.test", "operating_system_image", "Ubuntu 22.04.4 LTS"),
					resource.TestCheckResourceAttr("denvr_vm.test", "personal_storage_mount_path", "/home/ubuntu/personal"),
					resource.TestCheckResourceAttr("denvr_vm.test", "tenant_shared_additional_storage", "/home/ubuntu/tenant-shared"),
					resource.TestCheckResourceAttr("denvr_vm.test", "persist_storage", "false"),
					resource.TestCheckResourceAttr("denvr_vm.test", "direct_storage_mount_path", "/home/ubuntu/direct-attached"),
					resource.TestCheckResourceAttr("denvr_vm.test", "root_disk_size", "500"),
					// TODO: Verify computed values?
				),
			},
			// // ImportState testing
			// {
			// 	ResourceName:      "denvr_vm.test",
			// 	ImportState:       true,
			// 	ImportStateVerify: true,
			// },
		},
	})
}
