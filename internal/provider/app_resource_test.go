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

// TestCatalogApp is a test app resource model for catalog based applications.
var TestCatalogApp = appResourceModel{
	ApplicationCatalogItemName:    types.StringValue("jupyter-notebook"),
	ApplicationCatalogItemVersion: types.StringValue("python-3.11.9"),
	Cluster:                       types.StringValue("Msc1"),
	Dns:                           types.StringValue(""),
	EnvironmentVariables:          types.MapNull(types.StringType),
	HardwarePackageName:           types.StringValue("g-nvidia-1xa100-40gb-pcie-14vcpu-112gb"),
	ImageCmdOverride:              types.ListNull(types.StringType),
	ImageRepositoryHostname:       types.StringValue(""),
	ImageRepositoryPassword:       types.StringValue(""),
	ImageRepositoryUsername:       types.StringValue(""),
	ImageUrl:                      types.StringValue(""),
	JupyterToken:                  types.StringValue("abc123"),
	Id:                            types.StringValue("terraform-app"),
	Ip:                            types.StringValue(""),
	Name:                          types.StringValue("terraform-app"),
	PersistDirectAttachedStorage:  types.BoolValue(false),
	PersonalSharedStorage:         types.BoolValue(false),
	PrivateIp:                     types.StringValue("172.16.0.96"),
	ProxyPort:                     types.Int32Value(-1),
	ResourcePool:                  types.StringValue("on-demand"),
	SecurityContextContainerGid:   types.Int32Value(-1),
	SecurityContextContainerUid:   types.Int32Value(-1),
	SecurityContextRunAsRoot:      types.BoolValue(false),
	SshKeys:                       types.ListNull(types.StringType),
	Status:                        types.StringValue("UNKNOWN"),
	Tenant:                        types.StringValue("denvr"),
	TenantSharedStorage:           types.BoolValue(false),
	Username:                      types.StringValue("test@foobar.com"),
	Wait:                          types.BoolValue(false),
	Interval:                      types.Int64Value(30),
	Timeout:                       types.Int64Value(600),
}

// TestCustomApp is a test app resource model for custom container images.
var TestCustomApp = appResourceModel{
	ApplicationCatalogItemName:    types.StringValue(""),
	ApplicationCatalogItemVersion: types.StringValue(""),
	Cluster:                       types.StringValue("Msc1"),
	Dns:                           types.StringValue(""),
	EnvironmentVariables:          types.MapNull(types.StringType),
	HardwarePackageName:           types.StringValue("g-nvidia-1xa100-40gb-pcie-14vcpu-112gb"),
	ImageCmdOverride:              types.ListValueMust(types.StringType, []attr.Value{types.StringValue("nginx")}),
	ImageRepositoryHostname:       types.StringValue("https://index.docker.io/v1/"),
	ImageRepositoryPassword:       types.StringValue(""),
	ImageRepositoryUsername:       types.StringValue(""),
	ImageUrl:                      types.StringValue("karthequian/helloworld:latest"),
	JupyterToken:                  types.StringValue(""),
	Id:                            types.StringValue("terraform-custom-app"),
	Ip:                            types.StringValue(""),
	Name:                          types.StringValue("terraform-custom-app"),
	PersistDirectAttachedStorage:  types.BoolValue(false),
	PersonalSharedStorage:         types.BoolValue(false),
	PrivateIp:                     types.StringValue("172.16.0.97"),
	ProxyPort:                     types.Int32Value(80),
	ResourcePool:                  types.StringValue("reserved-denvr"),
	SecurityContextContainerUid:   types.Int32Value(-1),
	SecurityContextContainerGid:   types.Int32Value(-1),
	SecurityContextRunAsRoot:      types.BoolValue(false),
	SshKeys:                       types.ListNull(types.StringType),
	Status:                        types.StringValue("UNKNOWN"),
	Tenant:                        types.StringValue("denvr"),
	TenantSharedStorage:           types.BoolValue(false),
	Username:                      types.StringValue("test@foobar.com"),
	Wait:                          types.BoolValue(false),
	Interval:                      types.Int64Value(30),
	Timeout:                       types.Int64Value(600),
}

var AppTestAuthResult = `
{
	"result": {
		"accessToken": "access1",
		"refreshToken": "refresh",
		"expireInSeconds": 600,
		"refreshTokenExpireInSeconds": 3600
	}
}
`

func makeApplicationsApiOverviewResponse(app appResourceModel, desiredStatus string) string {
	return fmt.Sprintf(
		`
		{
			"result": {
				"applicationCatalogItemName": "%s",
				"applicationCatalogItemVersion": "%s",
				"cluster": "%s",
				"createdBy": "%s",
				"dns": "%s",
				"hardwarePackage": "%s",
				"id": "%s",
				"privateIP": "%s",
				"publicIP": "%s",
				"resourcePool": "%s",
				"sshUsername": "%s",
				"status": "%s",
				"tenant": "%s"
			}
		}
		`,
		app.ApplicationCatalogItemName.ValueString(),
		app.ApplicationCatalogItemVersion.ValueString(),
		app.Cluster.ValueString(),
		app.Username.ValueString(),
		app.Dns.ValueString(),
		app.HardwarePackageName.ValueString(),
		app.Id.ValueString(),
		app.PrivateIp.ValueString(),
		app.Ip.ValueString(),
		app.ResourcePool.ValueString(),
		app.Username.ValueString(),
		desiredStatus,
		app.Tenant.ValueString(),
	)
}

func makeApplicationsApiDetailsResponse(app appResourceModel, desiredStatus string) string {
	// NOTE: We only care about the InstanceDetails in our codebase
	return fmt.Sprintf(
		`
		{
			"result": {
				"InstanceDetails": {
					"cluster": "%s",
					"createdBy": "%s",
					"dns": "%s",
					"hardwarePackage": "%s",
					"id": "%s",
					"privateIP": "%s",
					"publicIP": "%s",
					"resourcePool": "%s",
					"sshUsername": "%s",
					"status": "%s",
					"tenant": "%s"
				}
			}
		}
		`,
		app.Cluster.ValueString(),
		app.Username.ValueString(),
		app.Dns.ValueString(),
		app.HardwarePackageName.ValueString(),
		app.Id.ValueString(),
		app.PrivateIp.ValueString(),
		app.Ip.ValueString(),
		app.ResourcePool.ValueString(),
		app.Username.ValueString(),
		desiredStatus,
		app.Tenant.ValueString(),
	)
}

func makeApplicationsApiCommandResponse(app appResourceModel) string {
	return fmt.Sprintf(
		`{ "result": { "cluster": "%s", "id": "%s" } }`,
		app.Cluster.ValueString(),
		app.Id.ValueString(),
	)
}

var catalogAppResourceConfig = fmt.Sprintf(`
resource "denvr_app" "test_catalog" {
 name = "%s"
 cluster = "%s"
 hardware_package_name = "%s"
 application_catalog_item_name = "%s"
 application_catalog_item_version = "%s"
 resource_pool = "%s"
 jupyter_token = "%s"
 wait = %t
 interval = %d
 timeout = %d
}
`,
	TestCatalogApp.Name.ValueString(),
	TestCatalogApp.Cluster.ValueString(),
	TestCatalogApp.HardwarePackageName.ValueString(),
	TestCatalogApp.ApplicationCatalogItemName.ValueString(),
	TestCatalogApp.ApplicationCatalogItemVersion.ValueString(),
	TestCatalogApp.ResourcePool.ValueString(),
	TestCatalogApp.JupyterToken.ValueString(),
	TestCatalogApp.Wait.ValueBool(),
	TestCatalogApp.Interval.ValueInt64(),
	TestCatalogApp.Timeout.ValueInt64(),
)

var customAppResourceConfig = fmt.Sprintf(`
resource "denvr_app" "test_custom" {
 name = "%s"
 cluster = "%s"
 hardware_package_name = "%s"
 image_cmd_override = ["%s"]
 image_repository_hostname = "%s"
 image_url = "%s"
 proxy_port = %d
 resource_pool = "%s"
 security_context_run_as_root = %t
 wait = %t
 interval = %d
 timeout = %d
}
`,
	TestCustomApp.Name.ValueString(),
	TestCustomApp.Cluster.ValueString(),
	TestCustomApp.HardwarePackageName.ValueString(),
	TestCustomApp.ImageCmdOverride.Elements()[0].(types.String).ValueString(),
	TestCustomApp.ImageRepositoryHostname.ValueString(),
	TestCustomApp.ImageUrl.ValueString(),
	TestCustomApp.ProxyPort.ValueInt32(),
	TestCustomApp.ResourcePool.ValueString(),
	TestCustomApp.SecurityContextRunAsRoot.ValueBool(),
	TestCustomApp.Wait.ValueBool(),
	TestCustomApp.Interval.ValueInt64(),
	TestCustomApp.Timeout.ValueInt64(),
)

func TestAccAppResource_basic(t *testing.T) {
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
			resp.Write([]byte(AppTestAuthResult))
		},
	)
	catalogStatus := "UNKNOWN"
	mux.HandleFunc(
		"/api/v1/servers/applications/CreateCatalogApplication",
		func(resp http.ResponseWriter, req *http.Request) {
			catalogStatus = "UNKNOWN"
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte(makeApplicationsApiOverviewResponse(TestCatalogApp, catalogStatus)))
		},
	)
	mux.HandleFunc(
		"/api/v1/servers/applications/CreateCustomApplication",
		func(resp http.ResponseWriter, req *http.Request) {
			catalogStatus = "UNKNOWN"
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte(makeApplicationsApiOverviewResponse(TestCatalogApp, catalogStatus)))
		},
	)
	mux.HandleFunc(
		"/api/v1/servers/applications/GetApplicationDetails",
		func(resp http.ResponseWriter, req *http.Request) {
			switch catalogStatus {
			case "UNKNOWN":
				catalogStatus = "PENDING"
			case "PENDING":
				catalogStatus = "INITIALIZING"
			case "INITIALIZING":
				catalogStatus = "RUNNING"
			}
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte(makeApplicationsApiDetailsResponse(TestCatalogApp, catalogStatus)))
		},
	)
	mux.HandleFunc(
		"/api/v1/servers/applications/DestroyApplication",
		func(resp http.ResponseWriter, req *http.Request) {
			resp.WriteHeader(http.StatusOK)
			resp.Write([]byte(makeApplicationsApiCommandResponse(TestCatalogApp)))
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

	resource.Test(
		t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + catalogAppResourceConfig,
					Check: resource.ComposeTestCheckFunc(
						// Verify provided values
						resource.TestCheckResourceAttr("denvr_app.test_catalog", "name", "terraform-app"),
						resource.TestCheckResourceAttr("denvr_app.test_catalog", "cluster", "Msc1"),
						resource.TestCheckResourceAttr("denvr_app.test_catalog", "hardware_package_name", "g-nvidia-1xa100-40gb-pcie-14vcpu-112gb"),
						resource.TestCheckResourceAttr("denvr_app.test_catalog", "application_catalog_item_name", "jupyter-notebook"),
						resource.TestCheckResourceAttr("denvr_app.test_catalog", "application_catalog_item_version", "python-3.11.9"),
						resource.TestCheckResourceAttr("denvr_app.test_catalog", "resource_pool", "on-demand"),
						resource.TestCheckResourceAttr("denvr_app.test_catalog", "jupyter_token", "abc123"),
					),
				},
			},
		})

	resource.Test(
		t, resource.TestCase{
			ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
			Steps: []resource.TestStep{
				{
					Config: providerConfig + customAppResourceConfig,
					Check: resource.ComposeTestCheckFunc(
						// Verify provided values
						resource.TestCheckResourceAttr("denvr_app.test_custom", "name", "terraform-custom-app"),
						resource.TestCheckResourceAttr("denvr_app.test_custom", "cluster", "Msc1"),
						resource.TestCheckResourceAttr("denvr_app.test_custom", "hardware_package_name", "g-nvidia-1xa100-40gb-pcie-14vcpu-112gb"),
						resource.TestCheckResourceAttr("denvr_app.test_custom", "image_cmd_override.0", "nginx"),
						resource.TestCheckResourceAttr("denvr_app.test_custom", "image_repository_hostname", "https://index.docker.io/v1/"),
						resource.TestCheckResourceAttr("denvr_app.test_custom", "image_url", "karthequian/helloworld:latest"),
						resource.TestCheckResourceAttr("denvr_app.test_custom", "proxy_port", "80"),
						resource.TestCheckResourceAttr("denvr_app.test_custom", "resource_pool", "reserved-denvr"),
						resource.TestCheckResourceAttr("denvr_app.test_custom", "security_context_run_as_root", "false"),
					),
				},
			},
		})
}
