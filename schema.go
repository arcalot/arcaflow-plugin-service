// Package service package for arcaflow_plugin_service.
package service

import (
	"context"
	"fmt"
	"regexp"
	"time"

	"go.flow.arcalot.io/pluginsdk/schema"
	core "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/internalversion/scheme"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
)

var connectionSchema = schema.NewStructMappedObjectSchema[Connection](
	"Connection",
	map[string]*schema.PropertySchema{
		"host": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Host"),
				schema.PointerTo("Host name and port of the Kubernetes server"),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`"kubernetes.default.svc"`),
			nil,
		).TreatEmptyAsDefaultValue(),
		"path": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Path"),
				schema.PointerTo("Path to the API server."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`"/api"`),
			nil,
		).TreatEmptyAsDefaultValue(),
		"username": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Username"),
				schema.PointerTo("Username for basic authentication."),
				nil,
			),
			false,
			[]string{"password"},
			nil,
			nil,
			nil,
			nil,
		),
		"password": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Password"),
				schema.PointerTo("Password for basic authentication."),
				nil,
			),
			false,
			[]string{"username"},
			nil,
			nil,
			nil,
			nil,
		),
		"serverName": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("TLS server name"),
				schema.PointerTo("Expected TLS server name to verify in the certificate."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		),
		"cacert": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, regexp.MustCompile(`^\s*-----BEGIN CERTIFICATE-----(\s*.*\s*)*-----END CERTIFICATE-----\s*$`)),
			schema.NewDisplayValue(
				schema.PointerTo("CA certificate"),
				schema.PointerTo("CA certificate in PEM format to verify Kubernetes server certificate against."),
				nil,
			),
			false,
			[]string{"cert", "key"},
			nil,
			nil,
			nil,
			nil,
		),
		"cert": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, regexp.MustCompile(`^\s*-----BEGIN CERTIFICATE-----(\s*.*\s*)*-----END CERTIFICATE-----\s*$`)),
			schema.NewDisplayValue(
				schema.PointerTo("Client certificate"),
				schema.PointerTo("Client certificate in PEM format to authenticate against Kubernetes with."),
				nil,
			),
			false,
			[]string{"key"},
			nil,
			nil,
			nil,
			nil,
		),
		"key": schema.NewPropertySchema(
			schema.NewStringSchema(schema.IntPointer(1), nil, regexp.MustCompile(`^\s*-----BEGIN ([A-Z]+) PRIVATE KEY-----(\s*.*\s*)*-----END ([A-Z]+) PRIVATE KEY-----\s*$`)),
			schema.NewDisplayValue(
				schema.PointerTo("Client key"),
				schema.PointerTo("Client private key in PEM format to authenticate against Kubernetes with."),
				nil,
			),
			false,
			[]string{"cert"},
			nil,
			nil,
			nil,
			nil,
		),
		"bearerToken": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Bearer token"),
				schema.PointerTo("Bearer token to authenticate against the Kubernetes API with."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		),
		"qps": schema.NewPropertySchema(
			schema.NewFloatSchema(
				schema.PointerTo(0.0),
				nil,
				schema.NewUnits(schema.NewUnit(
					"q",
					"q",
					"query",
					"queries",
				), nil)),
			schema.NewDisplayValue(
				schema.PointerTo("QPS"),
				schema.PointerTo("Queries Per Second allowed against the API."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`5.0`),
			nil,
		).TreatEmptyAsDefaultValue(),
		"burst": schema.NewPropertySchema(
			schema.NewIntSchema(schema.IntPointer(0), nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Burst"),
				schema.PointerTo("Burst value for query throttling."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo(`10`),
			nil,
		).TreatEmptyAsDefaultValue(),
	},
)

var serviceSchema = schema.NewStructMappedObjectSchema[Service](
	"Service",
	map[string]*schema.PropertySchema{
		"metadata": schema.NewPropertySchema(
			schema.NewRefSchema("ObjectMeta", nil),
			schema.NewDisplayValue(
				schema.PointerTo("Metadata"),
				schema.PointerTo("Service metadata."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		),
		"spec": schema.NewPropertySchema(
			schema.NewStructMappedObjectSchema[v1.ServiceSpec](
				"ServiceSpec",
				map[string]*schema.PropertySchema{
					"ports": schema.NewPropertySchema(
						schema.NewListSchema(
							schema.NewStructMappedObjectSchema[v1.ServicePort](
								"ServicePort",
								map[string]*schema.PropertySchema{
									"name": schema.NewPropertySchema(
										dnsSubdomainName,
										schema.NewDisplayValue(
											schema.PointerTo("Name"),
											schema.PointerTo("The name of this port within the service."),
											nil,
										),
										false,
										nil,
										nil,
										nil,
										nil,
										nil,
									).TreatEmptyAsDefaultValue(),
									"protocol": schema.NewPropertySchema(
										schema.NewStringEnumSchema(
											map[string]*schema.DisplayValue{
												"TCP":  schema.NewDisplayValue(schema.PointerTo("TCP"), nil, nil),
												"UDP":  schema.NewDisplayValue(schema.PointerTo("UDP"), nil, nil),
												"SCTP": schema.NewDisplayValue(schema.PointerTo("SCTP"), nil, nil),
											},
										),
										schema.NewDisplayValue(
											schema.PointerTo("Protocol"),
											schema.PointerTo("Protocol for this service."),
											nil,
										),
										true,
										nil,
										nil,
										nil,
										nil,
										nil,
									).TreatEmptyAsDefaultValue(),
									"appProtocol": schema.NewPropertySchema(
										labelValue,
										schema.NewDisplayValue(
											schema.PointerTo("App protocol"),
											schema.PointerTo("The application protocol for this port. See RFC6335."),
											nil,
										),
										false,
										nil,
										nil,
										nil,
										nil,
										nil,
									).TreatEmptyAsDefaultValue(),
									"port": schema.NewPropertySchema(
										schema.NewIntSchema(schema.IntPointer(1), schema.IntPointer(65535), nil),
										schema.NewDisplayValue(
											schema.PointerTo("Port"),
											schema.PointerTo("Port number that will be exposed."),
											nil,
										),
										false,
										nil,
										nil,
										nil,
										nil,
										nil,
									).TreatEmptyAsDefaultValue(),
								},
							),
							nil,
							nil,
						),
						schema.NewDisplayValue(
							schema.PointerTo("Ports"),
							schema.PointerTo("Ports for this service."),
							nil,
						),
						false,
						nil,
						nil,
						nil,
						nil,
						nil,
					),
					"selector": schema.NewPropertySchema(
						schema.NewMapSchema(
							labelName,
							labelValue,
							nil,
							nil,
						),
						schema.NewDisplayValue(
							schema.PointerTo("Selector"),
							schema.PointerTo("Target selector."),
							nil,
						),
						false,
						nil,
						nil,
						nil,
						nil,
						nil,
					),
				},
			),
			schema.NewDisplayValue(
				schema.PointerTo("Specification"),
				schema.PointerTo("Service specification."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		),
	},
)

var objectMeta = schema.NewStructMappedObjectSchema[metav1.ObjectMeta](
	"ObjectMeta",
	map[string]*schema.PropertySchema{
		"name": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Name"),
				schema.PointerTo("Resource name."),
				nil,
			),
			false,
			nil,
			nil,
			[]string{
				"generateName",
			},
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"generateName": schema.NewPropertySchema(
			schema.NewStringSchema(nil, nil, nil),
			schema.NewDisplayValue(
				schema.PointerTo("Name prefix"),
				schema.PointerTo("Name prefix to generate pod names from."),
				nil,
			),
			false,
			nil,
			nil,
			[]string{
				"name",
			},
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"namespace": schema.NewPropertySchema(
			dnsSubdomainName,
			schema.NewDisplayValue(
				schema.PointerTo("Namespace"),
				schema.PointerTo("Kubernetes namespace to deploy in."),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			schema.PointerTo("\"default\""),
			nil,
		).TreatEmptyAsDefaultValue(),
		"labels": schema.NewPropertySchema(
			schema.NewMapSchema(
				labelName,
				labelValue,
				nil,
				nil,
			),
			schema.NewDisplayValue(
				schema.PointerTo("Labels"),
				schema.PointerTo(
					"Kubernetes labels to appy. See https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/ for details.",
				),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
		"annotations": schema.NewPropertySchema(
			schema.NewMapSchema(
				labelName,
				labelValue,
				nil,
				nil,
			),
			schema.NewDisplayValue(
				schema.PointerTo("Annotations"),
				schema.PointerTo(
					"Kubernetes annotations to appy. See https://kubernetes.io/docs/concepts/overview/working-with-objects/annotations/ for details.",
				),
				nil,
			),
			false,
			nil,
			nil,
			nil,
			nil,
			nil,
		).TreatEmptyAsDefaultValue(),
	},
)

var labelName = schema.NewStringSchema(
	nil,
	nil,
	regexp.MustCompile(`^(|([a-zA-Z](|[a-zA-Z\-.]{0,251}[a-zA-Z0-9]))/)([a-zA-Z](|[a-zA-Z\\-]{0,61}[a-zA-Z0-9]))$`),
)
var labelValue = schema.NewStringSchema(
	nil,
	schema.IntPointer(63),
	regexp.MustCompile(`^(|[a-zA-Z0-9]+(|[-_.][a-zA-Z0-9]+)*[a-zA-Z0-9])$`),
)
var dnsSubdomainName = schema.NewStringSchema(
	nil,
	schema.IntPointer(253),
	regexp.MustCompile(`^[a-z0-9]($|[a-z0-9\-_]*[a-z0-9])$`),
)

var inputSchema = schema.NewScopeSchema(
	schema.NewStructMappedObjectSchema[Input](
		"input",
		map[string]*schema.PropertySchema{
			"connection": schema.NewPropertySchema(
				schema.NewRefSchema("Connection", nil),
				schema.NewDisplayValue(
					schema.PointerTo("Kubernetes"),
					schema.PointerTo("Kubernetes connection parameters."),
					nil,
				),
				true,
				nil,
				nil,
				nil,
				nil,
				nil,
			),
			"service": schema.NewPropertySchema(
				schema.NewRefSchema("Service", nil),
				schema.NewDisplayValue(
					schema.PointerTo("Service"),
					schema.PointerTo("Service to create."),
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
	connectionSchema,
	serviceSchema,
	objectMeta,
)

// Schema creates a pointer to new callable schema.
var Schema = schema.NewCallableSchema(
	schema.NewCallableStep[Input](
		"create",
		inputSchema,
		map[string]*schema.StepOutputSchema{
			"success": schema.NewStepOutputSchema(
				successOutputSchema,
				schema.NewDisplayValue(
					schema.PointerTo("Success"),
					schema.PointerTo("Service successfully created."),
					nil,
				),
				false,
			),
			"error": schema.NewStepOutputSchema(
				errorOutputSchema,
				schema.NewDisplayValue(
					schema.PointerTo("Error"),
					schema.PointerTo("Service creation failed."),
					nil,
				),
				true,
			),
		},
		schema.NewDisplayValue(
			schema.PointerTo("Create service"),
			schema.PointerTo("Create a Kubernetes service with the given specification."),
			nil,
		),
		func(ctx context.Context, input Input) (string, any) {
			connectionConfig := restclient.Config{
				Host:    input.Connection.Host,
				APIPath: input.Connection.APIPath,
				ContentConfig: restclient.ContentConfig{
					GroupVersion:         &core.SchemeGroupVersion,
					NegotiatedSerializer: scheme.Codecs.WithoutConversion(),
				},
				Username:    input.Connection.Username,
				Password:    input.Connection.Password,
				BearerToken: input.Connection.BearerToken,
				Impersonate: restclient.ImpersonationConfig{},
				TLSClientConfig: restclient.TLSClientConfig{
					ServerName: input.Connection.ServerName,
					CertData:   []byte(input.Connection.CertData),
					KeyData:    []byte(input.Connection.KeyData),
					CAData:     []byte(input.Connection.CAData),
				},
				UserAgent: "Arcaflow",
				QPS:       float32(input.Connection.QPS),
				Burst:     int(input.Connection.Burst),
				Timeout:   time.Second * 15,
			}

			cli, err := kubernetes.NewForConfig(&connectionConfig)
			if err != nil {
				return "error", ErrorOutput{
					fmt.Sprintf("failed to create Kubernetes config (%v)", err),
				}
			}

			if input.Service.Metadata.Name == "" && input.Service.Metadata.GenerateName == "" {
				input.Service.Metadata.GenerateName = "arcaflow-svc-plugin-"
			}

			srv, err := cli.CoreV1().Services(input.Service.Metadata.Namespace).Create(
				context.Background(),
				&v1.Service{
					ObjectMeta: input.Service.Metadata,
					Spec:       input.Service.Spec,
				},
				metav1.CreateOptions{},
			)
			if err != nil {
				return "error", ErrorOutput{
					fmt.Sprintf("failed to create service (%v)", err),
				}
			}
			return "success", SuccessOutput{
				srv.Name,
			}
		},
	),
)
