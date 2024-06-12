package database

type table struct {
	Keyword          string
	MemoryKeyword    string
	Memory           string
	OriginalResource string
	Summary          string
}

var Table = table{
	Keyword:          "keyword",
	MemoryKeyword:    "memory_keyword",
	Memory:           "memory",
	OriginalResource: "original_resource",
	Summary:          "summary",
}
