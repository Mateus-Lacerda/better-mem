package repository

import (
	"better-mem/internal/core"
	"better-mem/internal/database/sqlite"
	"better-mem/internal/repository/vector"
	"context"
	"database/sql"
	"encoding/json"
	"log/slog"

	sqlite_vec "github.com/asg017/sqlite-vec-go-bindings/cgo"
)

type MemoryRepository struct {
	db *sql.DB
}

func NewMemoryRepository() *MemoryRepository {
	return &MemoryRepository{db: sqlite.GetDbVec()}
}

// Create implements [vector.MemoryVectorRepository]
func (m *MemoryRepository) Create(
	ctx context.Context,
	chatId string,
	vectors []float32,
	memoryType core.MemoryTypeEnum,
	memoryId string,
) error {
	blob, err := sqlite_vec.SerializeFloat32(vectors)
	if err != nil {
		return err
	}

	vectorsJSON, err := json.Marshal(vectors)
	if err != nil {
		return err
	}

	result, err := m.db.ExecContext(
		ctx,
		"INSERT INTO vec_memories(id, embedding, memory_type, vectors_json) VALUES (?, ?, ?, ?)",
		memoryId, blob, int(memoryType), string(vectorsJSON),
	)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	slog.Info("MemoryRepository.Create", "rowsAffected", rowsAffected)
	return err
}

// Search implements [vector.MemoryVectorRepository]
func (m *MemoryRepository) Search(
	ctx context.Context,
	chatId string,
	vector []float32,
	limit int,
	threshold float32,
) (*[]core.ScoredMemoryVector, error) {
	blob, err := sqlite_vec.SerializeFloat32(vector)
	if err != nil {
		return nil, err
	}

	query := `
		WITH knn_matches AS (
			SELECT
				id,
				distance,
				memory_type,
				vectors_json
			FROM vec_memories
			WHERE embedding MATCH ?
			  AND k = ?
		)
		SELECT
			knn.id,
			knn.distance,
			knn.memory_type,
			knn.vectors_json,
			m.chat_id,
			m.active
		FROM knn_matches knn
		LEFT JOIN (
			SELECT id, chat_id, active FROM long_term_memories
			UNION ALL
			SELECT id, chat_id, active FROM short_term_memories
		) m ON knn.id = m.id
		WHERE m.chat_id = ?
		  AND m.active = 1
		ORDER BY knn.distance
	`

	rows, err := m.db.QueryContext(ctx, query, blob, limit, chatId)
	if err != nil {
		slog.Error("Query error", "error", err)
		return nil, err
	}
	defer rows.Close()

	var memories []core.ScoredMemoryVector

	for rows.Next() {
		var id string
		var distance float64
		var memoryType int
		var vectorsJSON string
		var dbChatId string
		var active bool

		err = rows.Scan(&id, &distance, &memoryType, &vectorsJSON, &dbChatId, &active)
		if err != nil {
			slog.Error("Scan error", "error", err)
			continue
		}

		var vectors []float32
		if vectorsJSON != "" {
			err := json.Unmarshal([]byte(vectorsJSON), &vectors)
			if err != nil {
				slog.Error("Failed to unmarshal vector", "error", err)
				vectors = nil
			}
		}

		score := float32(1.0 - distance)

		slog.Info("Found row",
			"id", id,
			"distance", distance,
			"score", score,
			"chat_id", dbChatId,
			"active", active)

		if score < threshold {
			continue
		}

		memories = append(memories, core.ScoredMemoryVector{
			Id:      id,
			Vectors: vectors,
			Score:   score,
			Payload: core.MemoryPayload{
				ChatId:     dbChatId,
				MemoryType: core.MemoryTypeEnum(memoryType),
				MemoryId:   id,
				Active:     active,
			},
		})
	}

	if err = rows.Err(); err != nil {
		slog.Error("Rows iteration error", "error", err)
		return nil, err
	}

	slog.Info("Search complete", "results", len(memories), "chat_id", chatId)
	return &memories, nil
}

// DeactivateAll implements [vector.MemoryVectorRepository]
func (m *MemoryRepository) DeactivateAll(ctx context.Context, chatId string) error {
	tx, err := m.db.BeginTx(
		ctx,
		&sql.TxOptions{Isolation: 0, ReadOnly: false},
	)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(
		ctx,
		"UPDATE long_term_memories SET memory = false WHERE chat_id = ?",
		chatId,
	); err != nil {
		return err
	}
	if _, err := tx.ExecContext(
		ctx,
		"UPDATE short_term_memories SET memory = false WHERE chat_id = ?",
		chatId,
	); err != nil {
		return err
	}
	return tx.Commit()
}

// Deactivate implements [vector.MemoryVectorRepository]
func (m *MemoryRepository) Deactivate(ctx context.Context, chatId string, id string) error {
	tx, err := m.db.BeginTx(
		ctx,
		&sql.TxOptions{Isolation: 0, ReadOnly: false},
	)
	if err != nil {
		return err
	}
	if _, err := tx.ExecContext(
		ctx,
		"UPDATE long_term_memories SET memory = false WHERE chat_id = ? AND memory_id = ?",
		chatId,
		id,
	); err != nil {
		return err
	}
	if _, err := tx.ExecContext(
		ctx,
		"UPDATE short_term_memories SET memory = false WHERE chat_id = ? AND memory_id = ?",
		chatId,
		id,
	); err != nil {
		return err
	}
	return tx.Commit()
}

var _ vector.MemoryVectorRepository = (*MemoryRepository)(nil)
