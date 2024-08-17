package database

type table struct {
	Chunk            string
	ChunkKeyword     string
	Keyword          string
	MemoryKeyword    string
	Memory           string
	OriginalResource string
	Summary          string
}

var Table = table{
	Chunk:            "chunk",
	ChunkKeyword:     "chunk_keyword",
	Keyword:          "keyword",
	MemoryKeyword:    "memory_keyword",
	Memory:           "memory",
	OriginalResource: "original_resource",
	Summary:          "summary",
}
