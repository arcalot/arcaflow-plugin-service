package arcaflow_plugin_service

import (
    v1 "k8s.io/api/core/v1"
    metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Input contains the structure for creating or removing a service.
type Input struct {
    Connection Connection `json:"connection,omitempty" yaml:"connection,omitempty"`
    Service    Service    `json:"service,omitempty" yaml:"service,omitempty"`
}

// Connection describes how to connect to the Kubernetes API.
type Connection struct {
    Host    string `json:"host,omitempty" yaml:"host,omitempty"`
    APIPath string `json:"path,omitempty" yaml:"path,omitempty"`

    Username string `json:"username,omitempty" yaml:"username,omitempty"`
    Password string `json:"password,omitempty" yaml:"password,omitempty"`

    ServerName string `json:"serverName,omitempty" yaml:"serverName,omitempty"`

    CertData string `json:"cert,omitempty" yaml:"cert,omitempty"`
    KeyData  string `json:"key,omitempty" yaml:"key,omitempty"`
    CAData   string `json:"cacert,omitempty" yaml:"cacert,omitempty"`

    BearerToken string `json:"bearerToken,omitempty" yaml:"bearerToken,omitempty"`

    QPS   float64 `json:"qps,omitempty" yaml:"qps,omitempty"`
    Burst int64   `json:"burst,omitempty" yaml:"burst,omitempty"`
}

// Service describes the service parameters.
type Service struct {
    Metadata metav1.ObjectMeta `json:"metadata,omitempty" yaml:"metadata,omitempty"`
    Spec     v1.ServiceSpec    `json:"spec" yaml:"spec"`
}
