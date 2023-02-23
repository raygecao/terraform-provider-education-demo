package main

import (
	"context"
	"terraform-provider-education-demo/education"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	providerserver.Serve(context.TODO(), education.NewProvider, providerserver.ServeOpts{
		Address: "raygecao.cn/test/education-demo",
	})
}
