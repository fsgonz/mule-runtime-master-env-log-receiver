package logsampler

// Constants for valid metric values
const (
	MetricNetstats = "netstats"
)

// Constants for valid output values
const (
	OutputFileLogger      = "file_logger"
	OutputPipelineEmitter = "pipeline_emitter"
)

// Constants for the logs
const (
	LastCountKey    = "LAST_COUNT"
	Format          = "v1"
	SchemaID        = "schema_id"
	NetworkSchemaId = "network_schema_id"
)

// Constants for environment variables
const (
	OrgID              = "ORG_ID"
	EnvID              = "ENV_ID"
	DeploymentID       = "DEPLOYMENT_ID"
	RootOrgID          = "ROOT_ORG_ID"
	MuleBillingEnabled = "MULE_BILLING_ENABLED"
	PodName            = "POD_NAME"
	AppName            = "APP_NAME"
)
