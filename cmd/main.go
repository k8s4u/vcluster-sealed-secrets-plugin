package main

import (
	"fmt"

	"github.com/k8s4u/vcluster-sealed-secrets-plugin/hooks"
	"github.com/loft-sh/vcluster-sdk/plugin"
)

func main() {
	fmt.Println("Starting vcluster-sealed-secrets-plugin")
	_ = plugin.MustInit("sealed-secrets")
	plugin.MustRegister(hooks.NewSecretHook())
	plugin.MustStart()
}
