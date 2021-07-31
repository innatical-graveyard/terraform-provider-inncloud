package main

import (
	"context"
	"terraform-provider-inncloud/inncloud"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func main() {
	tfsdk.Serve(context.Background(), inncloud.New, tfsdk.ServeOpts{
		Name: "inncloud",
	})
}
