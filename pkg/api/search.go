package api

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/lashajini/mind-palace/pkg/common"
	"github.com/lashajini/mind-palace/pkg/models"
	"github.com/lashajini/mind-palace/pkg/mperrors"
	pb "github.com/lashajini/mind-palace/pkg/rpc/gen"
	"github.com/lashajini/mind-palace/pkg/rpc/loggers"
	vdbrpc "github.com/lashajini/mind-palace/pkg/rpc/vdb"
	"github.com/lashajini/mind-palace/pkg/storage/database"
)

type SearchResponse struct {
	Response []SearchResult
}

type SearchResult struct {
	MemoryID string
	Chunks   []*HighlightedChunk
}

func (s SearchResult) String() string {
	return fmt.Sprintf("\nSearchResult{\nMemoryID: %s,\nChunks: %s\n}",
		s.MemoryID,
		s.Chunks,
	)
}

type HighlightedChunk struct {
	Highlight bool
	Chunk     *models.Chunk
}

func (h HighlightedChunk) String() string {
	return fmt.Sprintf("HighlightedChunk{\n\tHighlighted: %t,\n\tChunk: %s\n}",
		h.Highlight,
		h.Chunk,
	)
}

func Search(ctx context.Context, text string, db database.DB, vdbGrpcClient vdbrpc.VDBClient) (SearchResponse, error) {
	select {
	case <-ctx.Done():
		return SearchResponse{}, ctx.Err()
	default:
		var response []SearchResult

		searchResponse, err := vdbGrpcClient.Search(ctx, text)
		if err != nil {
			return SearchResponse{}, mperrors.On(err).Wrap("failed to search in vdb")
		}

		memoryIds, highlightedChunkIds := initIdMaps(searchResponse)

		err = fetchMemoryIdsByChunkIds(db, memoryIds, highlightedChunkIds)
		if err != nil {
			return SearchResponse{}, mperrors.On(err).Wrap("failed to fetch memory ids by chunk ids")
		}

		err = fetchChunksByMemoryIds(db, memoryIds, highlightedChunkIds)
		if err != nil {
			return SearchResponse{}, mperrors.On(err).Wrap("failed to fetch chunks by memory ids")
		}

		for memoryId, chunks := range memoryIds {
			response = append(response, SearchResult{
				MemoryID: memoryId,
				Chunks:   chunks,
			})
		}

		return SearchResponse{response}, nil
	}
}

func fetchMemoryIdsByChunkIds(db database.DB, memoryIds map[string][]*HighlightedChunk, highlightedChunkIds map[string]bool) error {
	chunkIdsStr := mapKeysToStr(highlightedChunkIds)

	q := fmt.Sprintf(`SELECT DISTINCT
    memory_id
FROM
    %s.chunk
WHERE
    id IN (%s);
`, db.CurrentSchema(), chunkIdsStr)
	loggers.Log.DBInfo(context.Background(), uuid.Nil, q)

	chunkRows, err := db.Query(q)
	if err != nil {
		return mperrors.On(err).Wrap("failed to query chunk")
	}
	defer chunkRows.Close()

	for chunkRows.Next() {
		var memoryId string

		err := chunkRows.Scan(
			&memoryId,
		)

		if err != nil {
			return mperrors.On(err).Wrap("failed to scan chunk")
		}

		memoryIds[memoryId] = []*HighlightedChunk{}
	}

	return nil
}

func fetchChunksByMemoryIds(db database.DB, memoryIds map[string][]*HighlightedChunk, highlightedChunkIds map[string]bool) error {
	memoryIdsStr := mapKeysToStr(memoryIds)

	q := fmt.Sprintf(`SELECT
    m.id AS memory_id,
    
    c.id,
    c.memory_id,
    c.sequence,
    c.chunk,
    c.created_at,
    c.updated_at
FROM
    %s.memory AS m
INNER JOIN
    %s.chunk AS c ON c.memory_id = m.id
WHERE
    m.id IN (%s)
ORDER BY
    c.sequence ASC;
`, db.CurrentSchema(), db.CurrentSchema(), memoryIdsStr)
	loggers.Log.DBInfo(context.Background(), uuid.Nil, q)

	memoryRows, err := db.Query(q)
	if err != nil {
		return mperrors.On(err).Wrap("failed to query memory")
	}
	defer memoryRows.Close()

	for memoryRows.Next() {
		var memoryId string
		var chunk models.Chunk

		err := memoryRows.Scan(
			&memoryId,

			&chunk.ID,
			&chunk.MemoryID,
			&chunk.Sequence,
			&chunk.Chunk,
			&chunk.CreatedAt,
			&chunk.UpdatedAt,
		)

		if err != nil {
			return mperrors.On(err).Wrap("failed to scan memory")
		}

		memoryIds[memoryId] = append(memoryIds[memoryId], &HighlightedChunk{Highlight: highlightedChunkIds[chunk.ID.String()], Chunk: &chunk})
	}

	return nil
}

func initIdMaps(searchResponse *pb.SearchResponse) (map[string][]*HighlightedChunk, map[string]bool) {
	memoryIds := make(map[string][]*HighlightedChunk)
	highlightedChunkIds := make(map[string]bool)

	for _, row := range searchResponse.Rows {
		if row.Type == common.ROW_TYPE_WHOLE {
			memoryIds[row.Id] = nil
		} else {
			highlightedChunkIds[row.Id] = true
		}
	}

	return memoryIds, highlightedChunkIds
}

func mapKeysToStr[T any](m map[string]T) string {
	keys := []string{}
	for k := range m {
		keys = append(keys, fmt.Sprintf("'%s'", k))
	}
	return strings.Join(keys, ", ")
}
