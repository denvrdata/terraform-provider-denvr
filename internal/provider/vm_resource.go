package provider

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httputil"

	"github.com/denvrdata/go-denvr/api/v1/servers/virtual"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource = &vmResource{}
)

type vmResource struct{}

type vmResourceModel struct {
	Cluster                        types.String `tfsdk:"cluster"`
	Configuration                  types.String `tfsdk:"configuration"`
	DirectAttachedStoragePersisted types.Bool   `tfsdk:"direct_attached_storage_persisted"`
	DirectStorageMountPath         types.String `tfsdk:"direct_storage_mount_path"`
	GpuType                        types.String `tfsdk:"gpu_type"`
	Gpus                           types.Int32  `tfsdk:"gpus"`
	Id                             types.String `tfsdk:"id"`
	Image                          types.String `tfsdk:"image"`
	Ip                             types.String `tfsdk:"ip"`
	Memory                         types.Int64  `tfsdk:"memory"`
	Name                           types.String `tfsdk:"name"`
	Namespace                      types.String `tfsdk:"namespace"`
	OperatingSystemImage           types.String `tfsdk:"operating_system_image"`
	PersistStorage                 types.Bool   `tfsdk:"persist_storage"`
	PersonalStorageMountPath       types.String `tfsdk:"personal_storage_mount_path"`
	PrivateIp                      types.String `tfsdk:"private_ip"`
	RootDiskSize                   types.Int32  `tfsdk:"root_disk_size"`
	Rpool                          types.String `tfsdk:"rpool"`
	SshKeys                        types.List   `tfsdk:"ssh_keys"`
	Status                         types.String `tfsdk:"status"`
	Storage                        types.Int64  `tfsdk:"storage"`
	StorageType                    types.String `tfsdk:"storage_type"`
	TenancyName                    types.String `tfsdk:"tenancy_name"`
	TenantSharedAdditionalStorage  types.String `tfsdk:"tenant_shared_additional_storage"`
	Username                       types.String `tfsdk:"username"`
	Vcpus                          types.Int32  `tfsdk:"vcpus"`
	Vpc                            types.String `tfsdk:"vpc"`
}

func NewVmResource() resource.Resource {
	return &vmResource{}
}

func (r *vmResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vm"
}

func (r *vmResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"cluster": schema.StringAttribute{
				Required: true,
			},
			"configuration": schema.StringAttribute{
				Required: true,
			},
			"direct_attached_storage_persisted": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"direct_storage_mount_path": schema.StringAttribute{
				Optional: true,
				Computed: true,
			},
			"gpu_type": schema.StringAttribute{
				Computed: true,
			},
			"gpus": schema.Int32Attribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
			"image": schema.StringAttribute{
				Computed: true,
			},
			"ip": schema.StringAttribute{
				Computed: true,
			},
			"memory": schema.Int64Attribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
			},
			"namespace": schema.StringAttribute{
				Computed: true,
			},
			"operating_system_image": schema.StringAttribute{
				Required: true,
			},
			"persist_storage": schema.BoolAttribute{
				Optional: true,
				Computed: true,
				Default:  booldefault.StaticBool(false),
			},
			"personal_storage_mount_path": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("/home/ubuntu/personal"),
			},
			"private_ip": schema.StringAttribute{
				Computed: true,
			},
			"root_disk_size": schema.Int32Attribute{
				Required: true,
			},
			"rpool": schema.StringAttribute{
				Required: true,
			},
			"ssh_keys": schema.ListAttribute{
				Required:    true,
				ElementType: types.StringType,
			},
			"status": schema.StringAttribute{
				Computed: true,
			},
			"storage": schema.Int64Attribute{
				Computed: true,
			},
			"storage_type": schema.StringAttribute{
				Computed: true,
			},
			"tenancy_name": schema.StringAttribute{
				Computed: true,
			},
			"tenant_shared_additional_storage": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("/home/ubuntu/tenant-shared"),
			},
			"username": schema.StringAttribute{
				Computed: true,
			},
			"vcpus": schema.Int32Attribute{
				Computed: true,
			},
			"vpc": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *vmResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Reading Terraform plan data into vmResourceModel")
	var data vmResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Constructing virtual server request")
	serverReq := virtual.CreateServerJSONRequestBody{
		Cluster:                       data.Cluster.ValueString(),
		Configuration:                 data.Configuration.ValueString(),
		DirectStorageMountPath:        data.DirectStorageMountPath.ValueStringPointer(),
		Name:                          data.Name.ValueStringPointer(),
		OperatingSystemImage:          data.OperatingSystemImage.ValueStringPointer(),
		PersistStorage:                data.PersistStorage.ValueBoolPointer(),
		PersonalStorageMountPath:      data.PersonalStorageMountPath.ValueStringPointer(),
		RootDiskSize:                  data.RootDiskSize.ValueInt32Pointer(),
		Rpool:                         data.Rpool.ValueStringPointer(),
		SshKeys:                       []string{},
		TenantSharedAdditionalStorage: data.TenantSharedAdditionalStorage.ValueStringPointer(),
		Vpc:                           data.Vpc.ValueString(),
	}

	// Ugly hack cause data.SshKeys.Elements() seems to complain?
	var sshkeys []string
	if diags := data.SshKeys.ElementsAs(ctx, &sshkeys, false); diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}
	for _, key := range sshkeys {
		serverReq.SshKeys = append(serverReq.SshKeys, key)
	}

	tflog.Debug(ctx, "Constructing virtual machine service client")
	client := virtual.NewClient()

	tflog.Debug(ctx, "Making virtual machine creation request")
	rsp, err := client.CreateServerRaw(ctx, serverReq)
	if err != nil {
		resp.Diagnostics.AddError("Error creating server", err.Error())
		return
	}
	rspDump, err := httputil.DumpResponse(rsp, true)
	if err != nil {
		resp.Diagnostics.AddError("Error dumping response", err.Error())
		return
	}
	tflog.Debug(ctx, string(rspDump))

	server, err := ParseVirtualServerDetailsItem(rsp)
	if err != nil {
		resp.Diagnostics.AddError("Error parsing server response", err.Error())
		return
	}

	tflog.Debug(ctx, "Updating virtual machine resource state")
	serverJson, err := json.MarshalIndent(server, "", "\t")
	if err != nil {
		resp.Diagnostics.AddError("Error marshaling server response", err.Error())
		return
	}
	tflog.Debug(ctx, string(serverJson))
	data.GpuType = types.StringValue(*server.GpuType)
	data.Gpus = types.Int32Value(*server.Gpus)
	data.Id = types.StringValue(*server.Id)
	data.Image = types.StringValue(*server.Image)
	data.Ip = types.StringValue(*server.Ip)
	data.Memory = types.Int64Value(*server.Memory)
	data.Namespace = types.StringValue(*server.Namespace)
	data.PrivateIp = types.StringValue(*server.PrivateIp)
	data.Status = types.StringValue(*server.Status)
	data.Storage = types.Int64Value(*server.Storage)
	data.StorageType = types.StringValue(*server.StorageType)
	data.TenancyName = types.StringValue(*server.TenancyName)
	data.Username = types.StringValue(*server.Username)
	data.Vcpus = types.Int32Value(*server.Vcpus)

	tflog.Debug(ctx, "Handling optional Direct Attached Storage Persisted")
	if server.DirectAttachedStoragePersisted != nil {
		data.DirectAttachedStoragePersisted = types.BoolValue(*server.DirectAttachedStoragePersisted)
	} else {
		data.DirectAttachedStoragePersisted = types.BoolValue(false)
	}

	// Save data into Terraform state
	tflog.Debug(ctx, "Saving updated virtual machine Terraform state ")
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vmResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data vmResourceModel

	// Read Terraform prior state data into the model
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Read API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vmResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data vmResourceModel

	// Read Terraform plan data into the model
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Update API call logic

	// Save updated data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *vmResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Reading Terraform plan data into vmResourceModel")
	var data vmResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Constructing virtual server request")
	destroyParams := virtual.DestroyServerParams{
		Id:        data.Id.ValueString(),
		Namespace: data.Namespace.ValueString(),
		Cluster:   data.Cluster.ValueString(),
	}

	tflog.Debug(ctx, "Constructing virtual machine service client")
	client := virtual.NewClient()

	tflog.Debug(ctx, "Making virtual machine deletion request")
	server, err := client.DestroyServer(ctx, &destroyParams)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting server", err.Error())
		return
	}

	tflog.Debug(ctx, "Updating virtual machine resource state")
	serverJson, err := json.MarshalIndent(server, "", "\t")
	if err != nil {
		resp.Diagnostics.AddError("Error marshaling server response", err.Error())
		return
	}
	tflog.Debug(ctx, string(serverJson))
}

// ParseVirtualServerDetailsItemWithResult parses an HTTP response into a VirtualServerDetailsItem.
// TODO: Move this back into the virtual package
func ParseVirtualServerDetailsItem(rsp *http.Response) (*virtual.VirtualServerDetailsItem, error) {
	bodyBytes, err := io.ReadAll(rsp.Body)
	defer func() { _ = rsp.Body.Close() }()
	if err != nil {
		return nil, err
	}

	// Just panic for now if we get a 4xx or higher error
	if 400 <= rsp.StatusCode {
		panic(rsp.Status)
	}

	// If direct unmarshal fails, try checking for result key
	var dest struct {
		Result virtual.VirtualServerDetailsItem `json:"result"`
	}
	if err := json.Unmarshal(bodyBytes, &dest); err != nil {
		return nil, err
	}

	return &dest.Result, nil
}
