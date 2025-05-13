package probe

import model "github.com/litmuschaos/litmus/chaoscenter/graphql/server/graph/model"

type GetProbeResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data GetProbeResponseData `json:"data"`
}

type GetProbeResponseData struct {
	GetProbe model.Probe `json:"getProbe"`
}

type ListProbeResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data ListProbeResponseData `json:"data"`
}

type ListProbeResponseData struct {
	Probes []model.Probe `json:"listProbes"`
}

type DeleteProbeResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data DeleteProbeResponseData `json:"data"`
}

type DeleteProbeResponseData struct {
	DeleteProbe bool `json:"deleteProbe"`
}

type GetProbeYAMLResponse struct {
	Errors []struct {
		Message string   `json:"message"`
		Path    []string `json:"path"`
	} `json:"errors"`
	Data GetProbeYAMLResponseData `json:"data"`
}

type GetProbeYAMLResponseData struct {
	GetProbeYAML string `json:"getProbeYAML"`
}

// ProbeType defines the type of probe.
type ProbeType string

const (
	ProbeTypeHTTPProbe ProbeType = "httpProbe"
	ProbeTypeCMDProbe  ProbeType = "cmdProbe"
	ProbeTypePROMProbe ProbeType = "promProbe"
	ProbeTypeK8SProbe  ProbeType = "k8sProbe"
)

// InfrastructureType defines the type of infrastructure.
type InfrastructureType string

const (
	InfrastructureTypeKubernetes InfrastructureType = "Kubernetes"
	// Add other infrastructure types if necessary
)

// KubernetesHTTPProbeRequest defines properties for Kubernetes HTTP probes.
// Aligned with UI example for addKubernetesHTTPProbe variables.request.kubernetesHTTPProperties
type KubernetesHTTPProbeRequest struct {
	ProbeTimeout       string  `json:"probeTimeout"`
	Interval           string  `json:"interval"`
	Attempt            *int    `json:"attempt,omitempty"` // Added from UI example
	URL                string  `json:"url"`
	Method             *Method `json:"method"`
	InsecureSkipVerify *bool   `json:"insecureSkipVerify,omitempty"` // Present in GQL schema, can be set
	// Removed: Retry, PollingInterval, InitialDelay, EvaluationTimeout, StopOnFailure (not in UI addProbe input for properties)
}

// Method defines the HTTP method and its properties.
// This maps to HTTPProbeInputs in the GraphQL schema.
type Method struct {
	Get  *GetMethod  `json:"get,omitempty"`
	Post *PostMethod `json:"post,omitempty"` // Kept for completeness, though UI example only shows GET
}

type GetMethod struct { // Maps to HTTPGetInput
	ResponseCode string `json:"responseCode"`
	Criteria     string `json:"criteria"`
}

type PostMethod struct { // Maps to HTTPPostInput
	Body         *string `json:"body,omitempty"`
	ContentType  *string `json:"contentType,omitempty"`
	ResponseCode string  `json:"responseCode"`
	Criteria     string  `json:"criteria"`
}

// ComparatorInput represents the comparator input for CMD probes
type ComparatorInput struct {
	Type     string `json:"type"`     // e.g., "Contains", "Equal", "NotEqual"
	Criteria string `json:"criteria"` // This field is required by the GraphQL schema
	Value    string `json:"value"`    // The value to compare against
}

// Update KubernetesCMDProbeRequest to use the corrected ComparatorInput
type KubernetesCMDProbeRequest struct {
	Command      string           `json:"command"`
	Comparator   *ComparatorInput `json:"comparator,omitempty"`
	Source       *string          `json:"source,omitempty"`
	ProbeTimeout string           `json:"probeTimeout"`
	Interval     string           `json:"interval"`
	Attempt      *int             `json:"attempt,omitempty"`
}

type K8SProbeRequest struct {
	ProbeTimeout  string  `json:"probeTimeout"`
	Interval      string  `json:"interval"`
	Attempt       *int    `json:"attempt,omitempty"` 
	Group         *string `json:"group,omitempty"` 
	Version       string  `json:"version"`
	Resource      string  `json:"resource"`
	Namespace     string  `json:"namespace"`
	ResourceNames *string `json:"resourceNames,omitempty"` 
	FieldSelector *string `json:"fieldSelector,omitempty"` 
	LabelSelector *string `json:"labelSelector,omitempty"` 
	Operation     string  `json:"operation"`             
	}

// PROMProbeRequest defines properties for Prometheus probes.
type PROMProbeRequest struct {
	Endpoint     string  `json:"endpoint"`
	Query        string  `json:"query"`
	Comparator   string  `json:"comparator"`
	Value        string  `json:"value"`
	ProbeTimeout string  `json:"probeTimeout"`
	Interval     string  `json:"interval"`
	Attempt      *int    `json:"attempt,omitempty"`
}

type ProbeRequest struct {
	Name                     string                      `json:"name"`
	Description              *string                     `json:"description,omitempty"`
	Type                     ProbeType                   `json:"type"`
	InfrastructureType       InfrastructureType          `json:"infrastructureType"`
	Tags                     []string                    `json:"tags,omitempty"`
	KubernetesHTTPProperties *KubernetesHTTPProbeRequest `json:"kubernetesHTTPProperties,omitempty"`
	KubernetesCMDProperties  *KubernetesCMDProbeRequest  `json:"kubernetesCMDProperties,omitempty"`
	K8SProperties            *K8SProbeRequest            `json:"k8sProperties,omitempty"`
	PROMProperties           *PROMProbeRequest           `json:"promProperties,omitempty"`
}

type AddProbeResponse struct {
	Errors []struct {
		Message    string                 `json:"message"`
		Path       []string               `json:"path"`
		Extensions map[string]interface{} `json:"extensions,omitempty"`
	} `json:"errors,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"` // This is the actual response
}

type Probe struct {
	Name                     string                       `json:"name"`
	Description              string                       `json:"description"` // Not a pointer in the response
	Type                     ProbeType                    `json:"type"`
	InfrastructureType       InfrastructureType          `json:"infrastructureType"`
	Tags                     []string                     `json:"tags"`
	KubernetesHTTPProperties *KubernetesHTTPProbeResponse `json:"kubernetesHTTPProperties,omitempty"`
	KubernetesCMDProperties  *KubernetesCMDProbeResponse  `json:"kubernetesCMDProperties,omitempty"`
	K8SProperties            *K8SProbeResponse            `json:"k8sProperties,omitempty"`
	PROMProperties           *PROMProbeResponse           `json:"promProperties,omitempty"`
}

type KubernetesHTTPProbeResponse struct {
	ProbeTimeout       string  `json:"probeTimeout"`
	Interval           string  `json:"interval"`
	URL                string  `json:"url"`
	Method             *Method `json:"method,omitempty"`
	InsecureSkipVerify *bool   `json:"insecureSkipVerify,omitempty"`
	Typename           string  `json:"__typename,omitempty"`
}

type ComparatorDetails struct {
	Type     string `json:"type"`
	Criteria string `json:"criteria"` 
	Value    string `json:"value"`
}

type KubernetesCMDProbeResponse struct {
	ProbeTimeout string             `json:"probeTimeout"`
	Interval     string             `json:"interval"`
	Command      string             `json:"command"`
	Comparator   *ComparatorDetails `json:"comparator,omitempty"`
}

type K8SProbeResponse struct {
	ProbeTimeout string  `json:"probeTimeout"`
	Interval     string  `json:"interval"`
	Retry        *int    `json:"retry,omitempty"` 
	Group        *string `json:"group,omitempty"` 
	Version      string  `json:"version"`
	Resource     string  `json:"resource"`
	Namespace    string  `json:"namespace"`
}

// PROMProbeResponse aligned with updated createProbeMutation (minimal)
type PROMProbeResponse struct {
	ProbeTimeout string `json:"probeTimeout"`
	Interval     string `json:"interval"`
	Endpoint     string `json:"endpoint"`
	Query        string `json:"query"`
}
