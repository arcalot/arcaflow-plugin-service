package arcaflow_plugin_service

import "go.flow.arcalot.io/pluginsdk/schema"

type SuccessOutput struct {
}

var successOutputSchema = schema.NewScopeSchema(
    schema.NewStructMappedObjectSchema[SuccessOutput](
        "success",
        map[string]*schema.PropertySchema{},
    ),
)

type ErrorOutput struct {
    Error string `json:"error"`
}

var errorOutputSchema = schema.NewScopeSchema(
    schema.NewStructMappedObjectSchema[ErrorOutput](
        "error",
        map[string]*schema.PropertySchema{
            "error": schema.NewPropertySchema(
                schema.NewStringSchema(nil, nil, nil),
                schema.NewDisplayValue(
                    schema.PointerTo("Error message"), nil, nil,
                ),
                true,
                nil,
                nil,
                nil,
                nil,
                nil,
            ),
        },
    ),
)
