package cmdutils

// Context constants
type ContextKey string

const (
	CTX_AUTH_KEY   ContextKey = "ctxAuth"
	CXT_CONFIG_KEY ContextKey = "ctxConfig"
)

// Environment constants
type Env string

const (
	ENV_API_KEY Env = "CLI_API_KEY"
)

func (e Env) String() string {
	return string(e)
}

// Config constants
type ConfigKey string

const (
	CFG_CHUNK_SIZE         = "chunk_size"
	CFG_WORKERS            = "workers"
	CFG_DEFAULT_OUTPUT     = "default_output"
	CFG_REGION             = "region"
	CFG_ENV                = "env"
	CFG_DEBUG              = "debug"
	CFG_NO_CONFIRM         = "no_confirm"
	CFG_RAW_OUTPUT         = "raw_output"
	CFG_LANG               = "lang"
	CFG_SERVER_URL         = "server_url"
	CFG_VERSION_LAST_CHECK = "version_last_check"
)

func (c ConfigKey) String() string {
	return string(c)
}
