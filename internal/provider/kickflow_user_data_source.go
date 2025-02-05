// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"terraform-provider-kickflow/internal/client"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ datasource.DataSource = &KickflowUserDataSource{}

func NewKickflowUserDataSource() datasource.DataSource {
	return &KickflowUserDataSource{}
}

// KickflowUserDataSource defines the data source implementation.
type KickflowUserDataSource struct {
	client *client.KickflowClient
}

// KickflowUserDataSourceModel describes the data source data model.
type KickflowUserDataSourceModel struct {
	UserID    types.String `tfsdk:"user_id"`
	UserEmail types.String `tfsdk:"user_email"`

	Id            types.String `tfsdk:"id"`
	Email         types.String `tfsdk:"email"`
	Code          types.String `tfsdk:"code"`
	FirstName     types.String `tfsdk:"first_name"`
	LastName      types.String `tfsdk:"last_name"`
	FullName      types.String `tfsdk:"full_name"`
	EmployeeID    types.String `tfsdk:"employee_id"`
	Status        types.String `tfsdk:"status"`
	Locale        types.String `tfsdk:"locale"`
	CreatedAt     types.String `tfsdk:"created_at"`
	UpdatedAt     types.String `tfsdk:"updated_at"`
	DeactivatedAt types.String `tfsdk:"deactivated_at"`
}

func (d *KickflowUserDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (d *KickflowUserDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Data source for retrieving Kickflow user information",

		Attributes: map[string]schema.Attribute{
			"user_id": schema.StringAttribute{
				MarkdownDescription: "User ID to search for",
				Optional:            true,
			},
			"user_email": schema.StringAttribute{
				MarkdownDescription: "User email to search for",
				Optional:            true,
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "User ID",
				Computed:            true,
			},
			"email": schema.StringAttribute{
				MarkdownDescription: "Email address",
				Computed:            true,
			},
			"code": schema.StringAttribute{
				MarkdownDescription: "User code",
				Computed:            true,
			},
			"first_name": schema.StringAttribute{
				MarkdownDescription: "First name",
				Computed:            true,
			},
			"last_name": schema.StringAttribute{
				MarkdownDescription: "Last name",
				Computed:            true,
			},
			"full_name": schema.StringAttribute{
				MarkdownDescription: "Full name",
				Computed:            true,
			},
			"employee_id": schema.StringAttribute{
				MarkdownDescription: "Employee ID",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "User status",
				Computed:            true,
			},
			"locale": schema.StringAttribute{
				MarkdownDescription: "Locale settings",
				Computed:            true,
			},
			"created_at": schema.StringAttribute{
				MarkdownDescription: "Created at",
				Computed:            true,
			},
			"updated_at": schema.StringAttribute{
				MarkdownDescription: "Updated at",
				Computed:            true,
			},
			"deactivated_at": schema.StringAttribute{
				MarkdownDescription: "Deactivated at",
				Computed:            true,
			},
		},
	}
}

func (d *KickflowUserDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	// Prevent panic if the provider has not been configured.
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.KickflowClient)

	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected data source configuration type",
			fmt.Sprintf("Expected *client.KickflowClient but got %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *KickflowUserDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data KickflowUserDataSourceModel

	// Read Terraform configuration data into the model
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	var httpResp *http.Response
	var err error

	if !data.UserID.IsNull() {
		httpReq, err := http.NewRequest("GET", d.client.BaseURL+"/users", nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to create request",
				fmt.Sprintf("Error occurred while creating request: %s", err),
			)
			return
		}
		q := httpReq.URL.Query()
		q.Add("userId", data.UserID.ValueString())
		httpReq.URL.RawQuery = q.Encode()
		httpResp, err = d.client.DoRequest(httpReq)
	} else if !data.UserEmail.IsNull() {
		httpReq, err := http.NewRequest("GET", d.client.BaseURL+"/lookupByEmail", nil)
		if err != nil {
			resp.Diagnostics.AddError(
				"Failed to create request",
				fmt.Sprintf("Error occurred while creating request: %s", err),
			)
			return
		}
		q := httpReq.URL.Query()
		q.Add("email", data.UserEmail.ValueString())
		httpReq.URL.RawQuery = q.Encode()
		httpResp, err = d.client.DoRequest(httpReq)
	} else {
		resp.Diagnostics.AddError(
			"User ID or email address is required",
			"Please set either user_id or user_email to identify the user.",
		)
		return
	}

	if err != nil {
		resp.Diagnostics.AddError(
			"API request failed",
			fmt.Sprintf("Error occurred while retrieving user information: %s", err),
		)
		return
	}
	defer httpResp.Body.Close()

	if httpResp == nil || httpResp.StatusCode != http.StatusOK {
		resp.Diagnostics.AddError(
			"API error",
			fmt.Sprintf("Received unexpected status code from API: %d", httpResp.StatusCode),
		)
		return
	}

	var user struct {
		ID            string `json:"id"`
		Email         string `json:"email"`
		Code          string `json:"code"`
		FirstName     string `json:"first_name"`
		LastName      string `json:"last_name"`
		FullName      string `json:"full_name"`
		EmployeeID    string `json:"employee_id"`
		Status        string `json:"status"`
		Locale        string `json:"locale"`
		CreatedAt     string `json:"created_at"`
		UpdatedAt     string `json:"updated_at"`
		DeactivatedAt string `json:"deactivated_at"`
	}

	if err := json.NewDecoder(httpResp.Body).Decode(&user); err != nil {
		resp.Diagnostics.AddError(
			"Failed to parse response",
			fmt.Sprintf("Error occurred while parsing JSON: %s", err),
		)
		return
	}

	data.Id = types.StringValue(user.ID)
	data.Email = types.StringValue(user.Email)
	data.Code = types.StringValue(user.Code)
	data.FirstName = types.StringValue(user.FirstName)
	data.LastName = types.StringValue(user.LastName)
	data.FullName = types.StringValue(user.FullName)
	data.EmployeeID = types.StringValue(user.EmployeeID)
	data.Status = types.StringValue(user.Status)
	data.Locale = types.StringValue(user.Locale)
	data.CreatedAt = types.StringValue(user.CreatedAt)
	data.UpdatedAt = types.StringValue(user.UpdatedAt)
	data.DeactivatedAt = types.StringValue(user.DeactivatedAt)

	// Write logs using the tflog package
	// Documentation: https://terraform.io/plugin/log
	tflog.Trace(ctx, "read a data source")

	// Save data into Terraform state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
