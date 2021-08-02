package inncloud

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	// "math/big"
	// "strconv"
	// "time"

	// "github.com/hashicorp-demoapp/hashicups-client-go"
	"github.com/hashicorp/terraform-plugin-framework/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

type resourceServerType struct{}

// Server Resource schema
func (r resourceServerType) GetSchema(_ context.Context) (schema.Schema, []*tfprotov6.Diagnostic) {
	return schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"ip": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:     types.StringType,
				Required: true,
			},
			"model": {
				Type:     types.StringType,
				Required: true,
			},
			"image": {
				Type:     types.StringType,
				Required: true,
			},
			"region": {
				Type:     types.StringType,
				Required: true,
			},
			"cycle": {
				Type:     types.StringType,
				Required: true,
			},
		},
	}, nil
}

// New resource instance
func (r resourceServerType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, []*tfprotov6.Diagnostic) {
	return resourceOrder{
		p: *(p.(*provider)),
	}, nil
}

type resourceOrder struct {
	p provider
}

// Create a new resource
func (r resourceOrder) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Provider not configured",
			Detail:   "The provider hasn't been configured before apply, likely because it depends on an unknown value from another resource. This leads to weird stuff happening, so we'd prefer if you didn't do that. Thanks!",
		})
		return
	}

	var plan Server
	err := req.Plan.Get(ctx, &plan)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading plan",
			Detail:   "An unexpected error was encountered while reading the plan: " + err.Error(),
		})
		return
	}

	body, err := json.Marshal(map[string]string{
		"name":   plan.Name.Value,
		"model":  plan.Model.Value,
		"image":  plan.Image.Value,
		"region": plan.Region.Value,
		"cycle":  plan.Cycle.Value,
	})
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error making request",
			Detail:   "An unexpected error was encountered while making a request: " + err.Error(),
		})
		return
	}

	client := http.Client{}
	request, err := http.NewRequest("POST", "https://api.innatical.cloud/projects/"+r.p.projectID+"/servers", bytes.NewBuffer(body))
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error making request",
			Detail:   "An unexpected error was encountered while making a request: " + err.Error(),
		})
		return
	}

	request.Header.Set("authorization", r.p.token)
	request.Header.Set("content-type", "application/json")
	response, err := client.Do(request)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error making request",
			Detail:   "An unexpected error was encountered while making a request: " + err.Error(),
		})
		return
	}

	defer response.Body.Close()
	body, err = ioutil.ReadAll(response.Body)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error making request",
			Detail:   "An unexpected error was encountered while making a request: " + err.Error(),
		})
		return
	}

	server := make(map[string]string)

	if json.Unmarshal(body, &server) != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error making request",
			Detail:   "An unexpected error was encountered while making a request: " + err.Error(),
		})
		return
	}

	err = resp.State.Set(ctx, Server{
		ID:     types.String{Value: server["id"]},
		IP:     types.String{Value: server["ip"]},
		Name:   types.String{Value: server["name"]},
		Model:  types.String{Value: server["model"]},
		Image:  types.String{Value: server["image"]},
		Region: types.String{Value: server["region"]},
		Cycle:  types.String{Value: server["cycle"]},
	})
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error setting state",
			Detail:   "Could not set state, unexpected error: " + err.Error(),
		})
		return
	}

}

// Read resource information
func (r resourceOrder) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	var state Server
	err := req.State.Get(ctx, &state)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading state",
			Detail:   "An unexpected error was encountered while reading the state: " + err.Error(),
		})
		return
	}

	// Get order from API and then update what is in state from what the API returns
	serverID := state.ID.Value

	client := http.Client{}
	request, err := http.NewRequest("GET", "https://api.innatical.cloud/projects/"+r.p.projectID+"/servers/"+serverID, nil)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error making request",
			Detail:   "An unexpected error was encountered while making a request: " + err.Error(),
		})
		return
	}

	request.Header.Set("authorization", r.p.token)
	response, err := client.Do(request)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error making request",
			Detail:   "An unexpected error was encountered while making a request: " + err.Error(),
		})
		return
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error making request",
			Detail:   "An unexpected error was encountered while making a request: " + err.Error(),
		})
		return
	}

	var server map[string]string

	if json.Unmarshal(body, &server) != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error making request",
			Detail:   "An unexpected error was encountered while making a request: " + err.Error(),
		})
		return
	}

	// Set state
	err = resp.State.Set(ctx, Server{
		ID:     types.String{Value: server["id"]},
		IP:     types.String{Value: server["ip"]},
		Name:   types.String{Value: server["name"]},
		Model:  types.String{Value: server["model"]},
		Image:  types.String{Value: server["image"]},
		Region: types.String{Value: server["region"]},
		Cycle:  types.String{Value: server["cycle"]},
	})
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error setting state",
			Detail:   "Unexpected error encountered trying to set new state: " + err.Error(),
		})
		return
	}
}

// Update resource
func (r resourceOrder) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {

}

// Delete resource
func (r resourceOrder) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	var state Server
	err := req.State.Get(ctx, &state)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error reading state",
			Detail:   "An unexpected error was encountered while reading the state: " + err.Error(),
		})
		return
	}

	// Get order from API and then update what is in state from what the API returns
	serverID := state.ID.Value

	client := http.Client{}
	request, err := http.NewRequest("DELETE", "https://api.innatical.cloud/projects/"+r.p.projectID+"/servers/"+serverID, nil)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error making request",
			Detail:   "An unexpected error was encountered while making a request: " + err.Error(),
		})
		return
	}

	request.Header.Set("authorization", r.p.token)
	response, err := client.Do(request)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Error making request",
			Detail:   "An unexpected error was encountered while making a request: " + err.Error(),
		})
		return
	}

	defer response.Body.Close()
}
