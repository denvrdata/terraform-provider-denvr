package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/denvrdata/go-denvr/api/v1/servers/applications"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/mapdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource = &appResource{}
)

type appResource struct{}

type appResourceModel struct {
	ApplicationCatalogItemName    types.String `tfsdk:"application_catalog_item_name"`
	ApplicationCatalogItemVersion types.String `tfsdk:"application_catalog_item_version"`
	Cluster                       types.String `tfsdk:"cluster"`
	Dns                           types.String `tfsdk:"dns"`
	EnvironmentVariables          types.Map    `tfsdk:"environment_variables"`
	HardwarePackageName           types.String `tfsdk:"hardware_package_name"`
	ImageCmdOverride              types.List   `tfsdk:"image_cmd_override"`
	ImageRepositoryHostname       types.String `tfsdk:"image_repository_hostname"`
	ImageRepositoryPassword       types.String `tfsdk:"image_repository_password"`
	ImageRepositoryUsername       types.String `tfsdk:"image_repository_username"`
	ImageUrl                      types.String `tfsdk:"image_url"`
	JupyterToken                  types.String `tfsdk:"jupyter_token"`
	Id                            types.String `tfsdk:"id"`
	Ip                            types.String `tfsdk:"ip"`
	Name                          types.String `tfsdk:"name"`
	PersistDirectAttachedStorage  types.Bool   `tfsdk:"persist_direct_attached_storage"`
	PersonalSharedStorage         types.Bool   `tfsdk:"personal_shared_storage"`
	PrivateIp                     types.String `tfsdk:"private_ip"`
	ProxyPort                     types.Int32  `tfsdk:"proxy_port"`
	ReadinessWatcherPort          types.Int32  `tfsdk:"readiness_watcher_port"`
	ResourcePool                  types.String `tfsdk:"resource_pool"`
	SecurityContextContainerGid   types.Int32  `tfsdk:"security_context_container_gid"`
	SecurityContextContainerUid   types.Int32  `tfsdk:"security_context_container_uid"`
	SecurityContextRunAsRoot      types.Bool   `tfsdk:"security_context_run_as_root"`
	SshKeys                       types.List   `tfsdk:"ssh_keys"`
	Status                        types.String `tfsdk:"status"`
	Tenant                        types.String `tfsdk:"tenant"`
	TenantSharedStorage           types.Bool   `tfsdk:"tenant_shared_storage"`
	Username                      types.String `tfsdk:"username"`
	Wait                          types.Bool   `tfsdk:"wait"`
	Interval                      types.Int64  `tfsdk:"interval"`
	Timeout                       types.Int64  `tfsdk:"timeout"`
}

func NewAppResource() resource.Resource {
	return &appResource{}
}

func (r *appResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_app"
}

func (r *appResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "App resource schema",
		Description:         "Schema for App resource configuration and management",
		Attributes: map[string]schema.Attribute{
			"application_catalog_item_name": schema.StringAttribute{
				Optional: true,
			},
			"application_catalog_item_version": schema.StringAttribute{
				Optional: true,
			},
			"cluster": schema.StringAttribute{
				Required: true,
			},
			"dns": schema.StringAttribute{
				Computed: true,
			},
			"environment_variables": schema.MapAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     mapdefault.StaticValue(types.MapValueMust(types.StringType, nil)),
			},
			"hardware_package_name": schema.StringAttribute{
				Required: true,
			},
			"image_cmd_override": schema.ListAttribute{
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				Default:     listdefault.StaticValue(types.ListValueMust(types.StringType, nil)),
			},
			"image_repository_hostname": schema.StringAttribute{
				Optional: true,
			},
			"image_repository_password": schema.StringAttribute{
				Optional: true,
			},
			"image_repository_username": schema.StringAttribute{
				Optional: true,
			},
			"image_url": schema.StringAttribute{
				Optional: true,
			},
			"jupyter_token": schema.StringAttribute{
				Optional: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"ip": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"persist_direct_attached_storage": schema.BoolAttribute{
				Optional: true,
			},
			"personal_shared_storage": schema.BoolAttribute{
				Optional: true,
			},
			"private_ip": schema.StringAttribute{
				Computed: true,
			},
			"proxy_port": schema.Int32Attribute{
				Optional: true,
			},
			"readiness_watcher_port": schema.Int32Attribute{
				Optional: true,
			},
			"resource_pool": schema.StringAttribute{
				Required: true,
			},
			"security_context_container_gid": schema.Int32Attribute{
				Optional: true,
			},
			"security_context_container_uid": schema.Int32Attribute{
				Optional: true,
			},
			"security_context_run_as_root": schema.BoolAttribute{
				Optional: true,
			},
			"ssh_keys": schema.ListAttribute{
				Optional:    true,
				ElementType: types.StringType,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"tenant": schema.StringAttribute{
				Computed: true,
			},
			"tenant_shared_storage": schema.BoolAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Computed: true,
			},
			"wait": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"interval": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(30),
			},
			"timeout": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				Default:  int64default.StaticInt64(600),
			},
		},
	}
}

func (r *appResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Reading Terraform plan data into appResourceModel")
	var data appResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Constructing application client")
	client := applications.NewClient()

	var app *applications.ApplicationsApiOverview
	var err error

	// If an image repository hostname is provided then we must be creating a custom application.
	// Otherwise fallback to trying to create a catalog application.
	// If this assumption is wrong we should get an error we can provide to the user.
	if data.ImageRepositoryHostname.ValueString() != "" {
		app, err = createCustomApplication(ctx, client, data)
	} else {
		app, err = createCatalogApplication(ctx, client, data)
	}
	if err != nil {
		resp.Diagnostics.AddError("Error creating application", err.Error())
		return
	} else if app == nil {
		resp.Diagnostics.AddError("Error creating application", "Application was not created and no error was returned.")
		return
	} else if app.Id == nil {
		resp.Diagnostics.AddError("Error creating application", "Returned application Id is nil")
		return
	} else if app.Cluster == nil {
		resp.Diagnostics.AddError("Error creating application", "Returned application Cluster is nil")
		return
	}

	appJson, err := json.MarshalIndent(app, "", "\t")
	if err != nil {
		resp.Diagnostics.AddError("Error marshaling application response", err.Error())
		return
	}
	tflog.Debug(ctx, string(appJson))

	data = updateState(ctx, data, *app)

	if data.Wait.ValueBool() {
		tflog.Debug(ctx, "Waiting for application to be ready")
		getParams := applications.GetApplicationDetailsParams{
			Id:      *app.Id,
			Cluster: *app.Cluster,
		}

		start := time.Now()
		for {
			if time.Since(start) > (time.Duration(data.Timeout.ValueInt64()) * time.Second) {
				resp.Diagnostics.AddError("Timeout Error", "Waiting for application to come \"ONLINE\" timed out")
				return
			}

			details, err := client.GetApplicationDetails(ctx, &getParams)
			if err != nil {
				resp.Diagnostics.AddError("Error checking application status", err.Error())
				return
			} else if details == nil {
				resp.Diagnostics.AddError("Error checking application status", "Returned application details is nil")
				return
			} else if details.InstanceDetails == nil {
				resp.Diagnostics.AddError("Error checking application status", "Returned application instance details is nil")
				return
			} else if details.InstanceDetails.Status == nil {
				resp.Diagnostics.AddError("Error checking application status", "Returned application instance status is nil")
				return
			}

			if *details.InstanceDetails.Status == "ONLINE" {
				data = updateState(ctx, data, *details.InstanceDetails)
				break
			}

			time.Sleep(time.Duration(data.Interval.ValueInt64()) * time.Second)
		}
	}

	// Save data into Terraform state
	tflog.Debug(ctx, "Saving updated application Terraform state ")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *appResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data appResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic
	getParams := applications.GetApplicationDetailsParams{
		Id:      data.Id.ValueString(),
		Cluster: data.Cluster.ValueString(),
	}

	tflog.Debug(ctx, "Constructing application service client")
	client := applications.NewClient()

	tflog.Debug(ctx, "Making applications get request")
	details, err := client.GetApplicationDetails(ctx, &getParams)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("\"%s\" not found", getParams.Id)) {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Error getting application", err.Error())
		}
		return
	}

	tflog.Debug(ctx, "Updating application resource state")
	data = updateState(ctx, data, *details.InstanceDetails)

	// Save data into Terraform state
	// tflog.Debug(ctx, "Saving updated virtual machine Terraform state ")
	// resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *appResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data appResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *appResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Reading Terraform plan data into appResourceModel")
	var data appResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Constructing application deletion request")
	destroyParams := applications.DestroyApplicationParams{
		Id:      data.Id.ValueString(),
		Cluster: data.Cluster.ValueString(),
	}

	tflog.Debug(ctx, "Constructing application service client")
	client := applications.NewClient()

	tflog.Debug(ctx, "Making application deletion request")
	app, err := client.DestroyApplication(ctx, &destroyParams)
	if err != nil {
		if strings.Contains(err.Error(), fmt.Sprintf("\"%s\" not found", destroyParams.Id)) {
			resp.State.RemoveResource(ctx)
		} else {
			resp.Diagnostics.AddError("Error deleting application", err.Error())
		}
		return
	}

	tflog.Debug(ctx, "Updating application resource state")
	appJson, err := json.MarshalIndent(app, "", "\t")
	if err != nil {
		resp.Diagnostics.AddError("Error marshaling application response", err.Error())
		return
	}
	tflog.Debug(ctx, string(appJson))
}

func createCatalogApplication(ctx context.Context, client applications.Client, data appResourceModel) (*applications.ApplicationsApiOverview, error) {
	tflog.Debug(ctx, "Constructing catalog application request")

	// Convert SSH keys
	var sshKeys []string
	if !data.SshKeys.IsNull() {
		if diags := data.SshKeys.ElementsAs(ctx, &sshKeys, false); diags.HasError() {
			return nil, fmt.Errorf("error parsing SSH keys: %v", diags)
		}
	}

	// Construct the request body
	appReq := applications.CreateCatalogApplicationJSONRequestBody{
		ApplicationCatalogItemName:    data.ApplicationCatalogItemName.ValueString(),
		ApplicationCatalogItemVersion: data.ApplicationCatalogItemVersion.ValueString(),
		Cluster:                       data.Cluster.ValueString(),
		HardwarePackageName:           data.HardwarePackageName.ValueString(),
		JupyterToken:                  data.JupyterToken.ValueStringPointer(),
		Name:                          data.Name.ValueString(),
		PersistDirectAttachedStorage:  data.PersistDirectAttachedStorage.ValueBoolPointer(),
		PersonalSharedStorage:         data.PersonalSharedStorage.ValueBoolPointer(),
		ResourcePool:                  data.ResourcePool.ValueStringPointer(),
		SshKeys:                       &sshKeys,
		TenantSharedStorage:           data.TenantSharedStorage.ValueBoolPointer(),
	}

	tflog.Debug(ctx, "Making catalog application creation request")
	app, err := client.CreateCatalogApplication(ctx, appReq)
	if err != nil {
		return nil, fmt.Errorf("create catalog application failed: %w", err)
	}

	return app, nil
}

func createCustomApplication(ctx context.Context, client applications.Client, data appResourceModel) (*applications.ApplicationsApiOverview, error) {
	tflog.Debug(ctx, "Constructing custom application request")

	// Convert environment variables
	var envVars map[string]*string
	if !data.EnvironmentVariables.IsNull() {
		if diags := data.EnvironmentVariables.ElementsAs(ctx, &envVars, false); diags.HasError() {
			return nil, fmt.Errorf("error parsing environment variables: %v", diags)
		}
	}

	// Convert image command override
	var imageCmdOverride []string
	if !data.ImageCmdOverride.IsNull() {
		if diags := data.ImageCmdOverride.ElementsAs(ctx, &imageCmdOverride, false); diags.HasError() {
			return nil, fmt.Errorf("error parsing image command override: %v", diags)
		}
	}

	// Construct the request body
	appReq := applications.CreateCustomApplicationJSONRequestBody{
		Cluster:              data.Cluster.ValueString(),
		EnvironmentVariables: &envVars,
		HardwarePackageName:  data.HardwarePackageName.ValueString(),
		ImageCmdOverride:     &imageCmdOverride,
		ImageRepository: applications.ImageRepositoryDto{
			Hostname: data.ImageRepositoryHostname.ValueString(),
			Username: data.ImageRepositoryUsername.ValueStringPointer(),
			Password: data.ImageRepositoryPassword.ValueStringPointer(),
		},
		ImageUrl:                     data.ImageUrl.ValueString(),
		Name:                         data.Name.ValueString(),
		PersistDirectAttachedStorage: data.PersistDirectAttachedStorage.ValueBoolPointer(),
		PersonalSharedStorage:        data.PersonalSharedStorage.ValueBoolPointer(),
		ProxyPort:                    data.ProxyPort.ValueInt32Pointer(),
		ReadinessWatcherPort:         data.ReadinessWatcherPort.ValueInt32Pointer(),
		ResourcePool:                 data.ResourcePool.ValueStringPointer(),
		TenantSharedStorage:          data.TenantSharedStorage.ValueBoolPointer(),
	}

	tflog.Debug(ctx, "Making custom application creation request")
	app, err := client.CreateCustomApplication(ctx, appReq)
	if err != nil {
		return nil, fmt.Errorf("create custom application failed: %w", err)
	}

	return app, nil
}

func updateState(ctx context.Context, data appResourceModel, info interface{}) appResourceModel {
	tflog.Debug(ctx, "Updating application state")

	// Update data state with application details
	var id string
	var status, publicIp, privateIp, dns, createdBy, tenant *string
	var persistDirectAttachedStorage, personalSharedStorage, tenantSharedStorage *bool

	extractFields := func(
		_id, _status, _publicIp, _privateIp, _dns, _createdBy, _tenant *string,
		_persistDirectAttachedStorage, _personalSharedStorage, _tenantSharedStorage *bool) {
		if _id != nil {
			id = *_id
		}

		status = _status
		publicIp = _publicIp
		privateIp = _privateIp
		dns = _dns
		createdBy = _createdBy
		tenant = _tenant
		persistDirectAttachedStorage = _persistDirectAttachedStorage
		personalSharedStorage = _personalSharedStorage
		tenantSharedStorage = _tenantSharedStorage
	}

	// Check which type of info we received and extract fields accordingly
	switch v := info.(type) {
	case applications.ApplicationsApiOverview:
		extractFields(
			v.Id, v.Status, v.PublicIp, v.PrivateIp, v.Dns, v.CreatedBy, v.Tenant,
			v.PersistedDirectAttachedStorage, v.PersonalSharedStorage, v.TenantSharedStorage,
		)
	case applications.InstanceDetails:
		extractFields(
			v.Id, v.Status, v.PublicIp, v.PrivateIp, v.Dns, v.CreatedBy, v.Tenant,
			v.PersistedDirectAttachedStorage, v.PersonalSharedStorage, v.TenantSharedStorage,
		)
	default:
		tflog.Error(ctx, fmt.Sprintf("Unexpected info type: %T", info))
		return data
	}

	data.Id = types.StringValue(id)

	// Map other fields if they exist
	if status != nil {
		data.Status = types.StringValue(*status)
		tflog.Debug(ctx, "Application status: "+*status)
	} else if data.Status.IsUnknown() {
		tflog.Debug(ctx, "Application status is UNKNOWN")
		data.Status = types.StringValue("UNKNOWN")
	} else if data.Status.IsNull() {
		tflog.Debug(ctx, "Application status is NULL")
		data.Status = types.StringValue("NULL")
	} else if data.Status.Equal(types.StringValue("")) {
		tflog.Debug(ctx, "Application status is default")
	} else {
		tflog.Debug(ctx, "Application status: \"\"")
		data.Status = types.StringValue("")
	}

	if publicIp != nil {
		data.Ip = types.StringValue(*publicIp)
	} else {
		data.Ip = types.StringValue("NA")
	}
	if privateIp != nil {
		data.PrivateIp = types.StringValue(*privateIp)
	} else {
		data.PrivateIp = types.StringValue("NA")
	}
	if dns != nil {
		data.Dns = types.StringValue(*dns)
	} else {
		data.Dns = types.StringValue("NA")
	}
	if createdBy != nil {
		data.Username = types.StringValue(*createdBy)
	}
	if tenant != nil {
		data.Tenant = types.StringValue(*tenant)
	}
	if persistDirectAttachedStorage != nil {
		data.PersistDirectAttachedStorage = types.BoolValue(*persistDirectAttachedStorage)
	}
	if personalSharedStorage != nil {
		data.PersonalSharedStorage = types.BoolValue(*personalSharedStorage)
	}
	if tenantSharedStorage != nil {
		data.TenantSharedStorage = types.BoolValue(*tenantSharedStorage)
	}

	return data
}
