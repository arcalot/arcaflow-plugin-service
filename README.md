# Kubernetes service plugin for Arcaflow

This plugin creates a Kubernetes service for the given specification. You can use it with the following workflow specification:

```yaml
input:
  root: RootObject
  objects:
    RootObject:
      id: RootObject
      properties:
        kubeconfig:
          display:
            description: The complete kubeconfig file as a string
            name: Kubeconfig file contents
          type:
            type_id: string
        name:
          type:
            type_id: string
steps:
  read_kubeconfig:
    plugin: quay.io/arcalot/arcaflow-plugin-kubeconfig:latest
    input:
      kubeconfig: !expr $.input.kubeconfig
  create_service:
    plugin: quay.io/arcalot/arcaflow-plugin-service:latest
    input:
      connection: !expr $.steps.read_kubeconfig.outputs.success.connection
      service:
        metadata:
          namespace: default
          # Service metadata here
        spec:
          # Service spec here
```