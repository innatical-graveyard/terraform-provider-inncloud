package inncloud

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Server struct {
	ID     types.String `tfsdk:"id"`
	IP     types.String `tfsdk:"ip"`
	Name   types.String `tfsdk:"name"`
	Model  types.String `tfsdk:"model"`
	Image  types.String `tfsdk:"image"`
	Region types.String `tfsdk:"region"`
	Cycle  types.String `tfsdk:"cycle"`
}
