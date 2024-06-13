package config

var MIND_PALACE_ROOT = ".mind-palace"
var MIND_PALACE_MEMORIES = "memories"
var MIND_PALACE_RESOURCES = "resources"
var MIND_PALACE_ORIGINAL = "original"
var MIND_PALACE_CONFIG = "config.json"
var MIND_PALACE_INFO = ".info.json"

var VDB_ORIGINAL_RESOURCE_COLLECTION_NAME = "original"

var PROD_ENV = "prod"
var DEV_ENV = "dev"
var TEST_ENV = "test"
var ENVS = map[string]bool{PROD_ENV: true, DEV_ENV: true, TEST_ENV: true}

var _MIND_PALACE_TEST_PATH = ".tests"
