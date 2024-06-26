package common

const (
	MIND_PALACE_ROOT      = ".mind-palace"
	MIND_PALACE_MEMORIES  = "memories"
	MIND_PALACE_RESOURCES = "resources"
	MIND_PALACE_ORIGINAL  = "original"
	MIND_PALACE_CONFIG    = "config.json"
	MIND_PALACE_INFO      = ".info.json"
)

const VDB_ORIGINAL_RESOURCE_COLLECTION_NAME = "original"

const (
	PROJECT_ROOT = "PROJECT_ROOT"
	MP_ENV       = "MP_ENV"
	PROD_ENV     = "prod"
	DEV_ENV      = "dev"
	TEST_ENV     = "test"
)

var ENVS = map[string]bool{PROD_ENV: true, DEV_ENV: true, TEST_ENV: true}

const (
	_MIND_PALACE_TEST_PATH = ".tests"
	TEST_USER              = "mindpalace_test_user"
)

const (
	LOG_LEVEL   = "LOG_LEVEL"
	LEVEL_PANIC = 5
	LEVEL_FATAL = 4
	LEVEL_ERROR = 3
	LEVEL_WARN  = 2
	LEVEL_INFO  = 1
	LEVEL_DEBUG = 0
	LEVEL_TRACE = -1

	COLOR_RESET  = "\033[0m"
	COLOR_YELLOW = "\033[33m"
)

const (
	DB_DEFAULT_SCHEMA = "public"
	DB_SCHEMA_SUFFIX  = "_mindpalace"
)
