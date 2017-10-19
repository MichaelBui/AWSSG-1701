package types

type (
	Configs map[string]interface{}

	DBConfig struct {
		Type      string
		Dsn       string
		Initiated bool
	}

	AwsS3Config struct {
		Bucket string
	}

	AwsETConfig struct {
		Pipeline string
		Preset   string
	}

	CloudConfig struct {
		Type    string
		Configs Configs
	}

	AppConfigs struct {
		DB    *DBConfig
		Cloud *CloudConfig
	}
)
