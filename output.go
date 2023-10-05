package service

import "go.flow.arcalot.io/pluginsdk/schema"

// SuccessOutput struct for name.
type SuccessOutput struct {
	Name string `json:"name"`
}

var successOutputSchema = schema.NewScopeSchema(
	schema.NewStructMappedObjectSchema[SuccessOutput](
		"success",
		map[string]*schema.PropertySchema{
			"name": schema.NewPropertySchema(
				dnsSubdomainName,
				schema.NewDisplayValue(
					schema.PointerTo("Name"),
					schema.PointerTo("Name of the service that has just been created."),
					nil,
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

// ErrorOutput struct for error.
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
