package main

import (
    "context"
    "fmt"
    "os"

    arcaflow_plugin_service "arcaflow-plugin-service"
    "go.flow.arcalot.io/pluginsdk/atp"
)

func usage() {
    fmt.Println(`Usage: ./service-plugin --atp`)
}

func main() {
    if len(os.Args) != 2 {
        usage()
        os.Exit(1)
    }
    if os.Args[1] != "--atp" {
        _, _ = os.Stderr.WriteString("This plugin only supports --atp.")
        usage()
        os.Exit(1)
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    if err := atp.RunATPServer(ctx, os.Stdin, os.Stdout, arcaflow_plugin_service.Schema); err != nil {
        panic(err)
    }
}
