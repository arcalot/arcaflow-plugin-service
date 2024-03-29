// Package main for arcaflow_plugin_service.
package main

import (
	service "arcaflow-plugin-service"

	"go.flow.arcalot.io/pluginsdk/plugin"
)

func main() {
	plugin.Run(service.Schema)
}
