package inncloud

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

func New() tfsdk.Provider {
	return &provider{}
}

type provider struct {
	configured bool
	token      string
	projectID  string
}

// GetSchema
func (p *provider) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"token": {
				Type:      types.StringType,
				Required:  true,
				Sensitive: true,
			},
			"project_id": {
				Type:      types.StringType,
				Required:  true,
				Sensitive: true,
			},
		},
	}, nil
}

// Provider schema struct
type providerData struct {
	Token     types.String `tfsdk:"token"`
	ProjectID types.String `tfsdk:"project_id"`
}

func (p *provider) Configure(ctx context.Context, req tfsdk.ConfigureProviderRequest, resp *tfsdk.ConfigureProviderResponse) {
	// Retrieve provider data from configuration
	var config providerData
	err := req.Config.Get(ctx, &config)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error parsing configuration",
			Detail:   "Error parsing the configuration, this is an error in the provider. Please report the following to the provider developer:\n\n" + err.Error(),
		})
		return
	}

	// User must specify a token
	var token string
	if config.Token.Unknown {
		// Cannot connect to inncloud with an unknown value
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityWarning,
			Summary:  "Unable to create client",
			Detail:   "Cannot use unknown value as token",
		})
		return
	}

	if config.Token.Null {
		token = os.Getenv("INNCLOUD_TOKEN")
	} else {
		token = config.Token.Value
	}

	if token == "" {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			// Error vs warning - empty value must stop execution
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to find token",
			Detail:   "Token cannot be an empty string",
		})
	}

	// User must specify a token
	var projectID string
	if config.ProjectID.Unknown {
		// Cannot connect to inncloud with an unknown value
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityWarning,
			Summary:  "Unable to create client",
			Detail:   "Cannot use unknown value as project_id",
		})
		return
	}

	if config.Token.Null {
		projectID = os.Getenv("INNCLOUD_PROJECT_ID")
	} else {
		projectID = config.ProjectID.Value
	}

	if token == "" {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			// Error vs warning - empty value must stop execution
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Unable to find token",
			Detail:   "project_id cannot be an empty string",
		})
	}

	p.projectID = projectID
	p.token = token
	p.configured = true
}

// GetResources - Defines provider resources
func (p *provider) GetResources(_ context.Context) (map[string]tfsdk.ResourceType, []*tfprotov6.Diagnostic) {
	return map[string]tfsdk.ResourceType{
		"inncloud_server": resourceServerType{},
	}, nil
}

// GetDataSources - Defines provider data sources
func (p *provider) GetDataSources(_ context.Context) (map[string]tfsdk.DataSourceType, []*tfprotov6.Diagnostic) {
	return map[string]tfsdk.DataSourceType{}, nil
}
