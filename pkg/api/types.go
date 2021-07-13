package api

type ClusterParams struct {
	Namespace                                 string                 `json:"namespace"`
	ClusterID                                 string                 `json:"clusterID"`
	ExternalAPIDNSName                        string                 `json:"externalAPIDNSName"`
	ExternalAPIPort                           uint                   `json:"externalAPIPort"`
	ExternalAPIIPAddress                      string                 `json:"externalAPIAddress"`
	ExternalOauthDNSName                      string                 `json:"externalOauthDNSName"`
	ExternalOauthPort                         uint                   `json:"externalOauthPort"`
	IdentityProviders                         string                 `json:"identityProviders"`
	ServiceCIDR                               string                 `json:"serviceCIDR"`
	NamedCerts                                []NamedCert            `json:"namedCerts,omitempty"`
	PodCIDR                                   string                 `json:"podCIDR"`
	ReleaseImage                              string                 `json:"releaseImage"`
	APINodePort                               uint                   `json:"apiNodePort"`
	IngressSubdomain                          string                 `json:"ingressSubdomain"`
	OpenShiftAPIClusterIP                     string                 `json:"openshiftAPIClusterIP"`
	OauthAPIClusterIP                         string                 `json:"oauthAPIClusterIP"`
	ImageRegistryHTTPSecret                   string                 `json:"imageRegistryHTTPSecret"`
	RouterNodePortHTTP                        string                 `json:"routerNodePortHTTP"`
	RouterNodePortHTTPS                       string                 `json:"routerNodePortHTTPS"`
	BaseDomain                                string                 `json:"baseDomain"`
	NetworkType                               string                 `json:"networkType"`
	Replicas                                  string                 `json:"replicas"`
	EtcdClientName                            string                 `json:"etcdClientName"`
	OriginReleasePrefix                       string                 `json:"originReleasePrefix"`
	OpenshiftAPIServerCABundle                string                 `json:"openshiftAPIServerCABundle"`
	OauthAPIServerCABundle                    string                 `json:"oauthAPIServerCABundle"`
	CloudProvider                             string                 `json:"cloudProvider"`
	CVOSetupImage                             string                 `json:"cvoSetupImage"`
	InternalAPIPort                           uint                   `json:"internalAPIPort"`
	RouterServiceType                         string                 `json:"routerServiceType"`
	KubeAPIServerResources                    []ResourceRequirements `json:"kubeAPIServerResources"`
	OpenshiftControllerManagerResources       []ResourceRequirements `json:"openshiftControllerManagerResources"`
	ClusterVersionOperatorResources           []ResourceRequirements `json:"clusterVersionOperatorResources"`
	KubeControllerManagerResources            []ResourceRequirements `json:"kubeControllerManagerResources"`
	OpenshiftAPIServerResources               []ResourceRequirements `json:"openshiftAPIServerResources"`
	OauthAPIServerResources                   []ResourceRequirements `json:"oauthAPIServerResources"`
	KubeSchedulerResources                    []ResourceRequirements `json:"kubeSchedulerResources"`
	ControlPlaneOperatorResources             []ResourceRequirements `json:"controlPlaneOperatorResources"`
	OAuthServerResources                      []ResourceRequirements `json:"oAuthServerResources"`
	ClusterPolicyControllerResources          []ResourceRequirements `json:"clusterPolicyControllerResources"`
	AutoApproverResources                     []ResourceRequirements `json:"autoApproverResources"`
	KMSServerResources                        []ResourceRequirements `json:"kmsServerResources"`
	PortierisContainerResources               []ResourceRequirements `json:"portierisContainerResources"`
	KMSImage                                  string                 `json:"kmsImage"`
	KPInfo                                    string                 `json:"kpInfo"`
	KPRegion                                  string                 `json:"kpRegion"`
	KPAPIKey                                  string                 `json:"kpAPIKey"`
	APIServerAuditEnabled                     bool                   `json:"apiServerAuditEnabled"`
	RestartDate                               string                 `json:"restartDate"`
	ControlPlaneOperatorImage                 string                 `json:"controlPlaneOperatorImage"`
	ControlPlaneOperatorControllers           []string               `json:"controlPlaneOperatorControllers"`
	ROKSMetricsImage                          string                 `json:"roksMetricsImage"`
	ROKSMetricsSecurityContextMaster          *SecurityContext       `json:"roksMetricsSecurityContextMaster"`
	ROKSMetricsSecurityContextWorker          *SecurityContext       `json:"roksMetricsSecurityContextWorker"`
	ExtraFeatureGates                         []string               `json:"extraFeatureGates"`
	ControlPlaneOperatorSecurityContext       *SecurityContext       `json:"controlPlaneOperatorSecurityContext"`
	MasterPriorityClass                       string                 `json:"masterPriorityClass"`
	ApiserverLivenessPath                     string                 `json:"apiserverLivenessPath"`
	PortierisEnabled                          bool                   `json:"portierisEnabled"`
	PortierisImage                            string                 `json:"portierisImage"`
	KubeAPIServerSecurityContext              *SecurityContext       `json:"kubeAPIServerSecurityContext"`
	KubeSchedulerSecurityContext              *SecurityContext       `json:"kubeSchedulerSecurityContext"`
	KubeControllerManagerSecurityContext      *SecurityContext       `json:"kubeControllerManagerSecurityContext"`
	OpenshiftAPIServerSecurityContext         *SecurityContext       `json:"openshiftAPIServerSecurityContext"`
	OauthAPIServerSecurityContext             *SecurityContext       `json:"oauthAPIServerSecurityContext"`
	OpenshiftControllerManagerSecurityContext *SecurityContext       `json:"openshiftControllerManagerSecurityContext"`
	PortierisSecurityContext                  *SecurityContext       `json:"portierisSecurityContext"`
	ClusterVersionOperatorSecurityContext     *SecurityContext       `json:"clusterVersionOperatorSecurityContext"`
	KMSSecurityContext                        *SecurityContext       `json:"kmsSecurityContext"`
	ManifestBootstrapperSecurityContext       *SecurityContext       `json:"manifestBootstrapperSecurityContext"`
	OAuthServerSecurityContext                *SecurityContext       `json:"oAuthServerSecurityContext"`
	ClusterPolicyControllerSecurityContext    *SecurityContext       `json:"clusterPolicyControllerSecurityContext"`
	ClusterConfigOperatorSecurityContext      *SecurityContext       `json:"clusterConfigOperatorSecurityContext"`
	DefaultFeatureGates                       []string
	PlatformType                              string                 `json:"platformType"`
	EndpointPublishingStrategyScope           string                 `json:"endpointPublishingStrategyScope"`
	ApiserverLivenessProbe                    *Probe                 `json:"apiserverLivenessProbe,omitempty"`
	ApiserverReadinessProbe                   *Probe                 `json:"apiserverReadinessProbe,omitempty"`
	ControllerManagerLivenessProbe            *Probe                 `json:"controllerManagerLivenessProbe,omitempty"`
	SchedulerLivenessProbe                    *Probe                 `json:"schedulerLivenessProbe,omitempty"`
	KMSLivenessProbe                          *Probe                 `json:"kmsLivenessProbe,omitempty"`
	PortierisLivenessProbe                    *Probe                 `json:"portierisLivenessProbe,omitempty"`
	KubeAPIServerVerbosity                    uint                   `json:"kubeAPIServerVerbosity"`
	KonnectivityEnabled                       bool                   `json:"konnectivityEnabled"`
	KonnectivityServerImage                   string                 `json:"konnectivityServerImage"`
	KonnectivityServerSecurityContext         *SecurityContext       `json:"konnectivityServerSecurityContext"`
	KonnectivityServerContainerResources      []ResourceRequirements `json:"konnectivityServerContainerResources"`
	KonnectivityServerPort                    uint                   `json:"konnectivityServerPort"`
	KonnectivityAgentPort                     uint                   `json:"konnectivityAgentPort"`
	KonnectivityServerHealthPort              uint                   `json:"konnectivityServerHealthPort"`
	KonnectivityServerAdminPort               uint                   `json:"konnectivityServerAdminPort"`
	KonnectivityServerAgentNodePort           uint                   `json:"konnectivityServerAgentNodePort"`
}

type NamedCert struct {
	NamedCertPrefix string `json:"namedCertPrefix"`
	NamedCertDomain string `json:"namedCertDomain"`
}

type HttpGetAction struct {
	Path   string `json:"path"`
	Port   uint   `json:"port"`
	Scheme string `json:"scheme"`
}

type Probe struct {
	HttpGet             HttpGetAction `json:"httpGet"`
	InitialDelaySeconds uint          `json:"initialDelaySeconds"`
	PeriodSeconds       uint          `json:"periodSeconds"`
	TimeoutSeconds      uint          `json:"timeoutSeconds"`
	FailureThreshold    uint          `json:"failureThreshold"`
	SuccessThreshold    uint          `json:"successThreshold"`
}

type SecurityContext struct {
	RunAsUser    uint `json:"runAsUser"`
	RunAsGroup   uint `json:"runAsGroup"`
	RunAsNonRoot bool `json:"runAsNonRoot"`
	Privileged   bool `json:"privileged"`
}

type ResourceRequirements struct {
	ResourceLimit   []ResourceLimit   `json:"resourceLimit"`
	ResourceRequest []ResourceRequest `json:"resourceRequest"`
}

type ResourceLimit struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}

type ResourceRequest struct {
	CPU    string `json:"cpu"`
	Memory string `json:"memory"`
}
