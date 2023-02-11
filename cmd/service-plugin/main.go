package main

import (
	arcaflow_plugin_service "arcaflow-plugin-service"

	"go.flow.arcalot.io/pluginsdk/plugin"
)

func main() {
	plugin.Run(arcaflow_plugin_service.Schema)
}
