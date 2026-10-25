package comments

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	pgx              *pgxpool.Pool
	config           *config.CommentsConfig
	logger           *zap.Logger
	statementBuilder squirrel.StatementBuilderType
}

func NewPostgresRepositoryFromPgxPool(connection *pgxpool.Pool, c *config.CommentsConfig, l *zap.Logger) *PostgresRepository {
	return &PostgresRepository{
		config:           c,
		logger:           l,
		pgx:              connection,
		statementBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostgresRepository) Initialize(ctx context.Context) error {
	return nil
}

func (r *PostgresRepository) Seed(ctx context.Context) {
	// Check if rows already exists, if so then dont seed
	existQuery := `SELECT count(*) FROM comments;`
	row := r.pgx.QueryRow(ctx, existQuery)

	var count int

	// Scan returned row
	if err := row.Scan(&count); err != nil {
		r.logger.Error("Failed to scan count of rows from comments table", zap.Error(err))
	}

	// Rows already exists, dont seed
	if count != 0 {
		r.logger.Sugar().Infof("Comments table contains %d rows, skipping seeding from file...", count)
	}

	// Table is empty, so seed from file to table
	queryBuilder := r.statementBuilder.Insert("comments").Columns(
		"id",
		"post_id",
		"author_id",
		"body",
		"created_at",
		"updated_at",
	)

	// Check if seeds file is provided and is of json format by suffix
	if len(r.config.SeedFile) == 0 || !strings.HasSuffix(r.config.SeedFile, ".json") {
		r.logger.Error("Seeds file not provided or is not of json format. Skipping seeding")
	}

	// open file
	file, err := os.OpenFile(r.config.SeedFile, os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open comments seeds file", zap.String("file", r.config.SeedFile), zap.Error(err))
	}
	defer file.Close()

	// Create json decoder
	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
	}

	// Decode
	var comments []Comment
	err = decoder.Decode(&comments)
	if err != nil {
		r.logger.Error("Failed to decode seeded users from file", zap.Error(err))
	}

	for i := range comments {
		queryBuilder = queryBuilder.Values(
			comments[i].Id,
			comments[i].PostId,
			comments[i].AuthorId,
			comments[i].Body,
			comments[i].CreatedAt,
			comments[i].UpdatedAt,
		)
	}

	// Generate query
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate sql query with squirrel for seeding users table", zap.Error(err))
	}

	// Execute query
	_, err = r.pgx.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to execute query to seed users table", zap.Error(err))
	}

	r.logger.Info("Successfully seeded comments from file", zap.Int("count", len(comments)))
}

// General operations
func (r *PostgresRepository) Add(ctx context.Context, req CreateRequest) (Comment, error) {
	rawQuery := `INSERT INTO
					comments(id, post_id, author_id, body, created_at, updated_at)
				VALUES($1, $2, $3, $4, $5, $6);`

	id, _ := uuid.NewV7()
	now := time.Now()

	// Execute quury
	_, err := r.pgx.Exec(ctx, rawQuery, id, req.PostId, req.Body, now, now)
	if err != nil {
		return Comment{}, fmt.Errorf("adding comment: %w", err)
	}

	return Comment{
		Id:        id,
		PostId:    req.PostId,
		Body:      req.Body,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context, cursor int, limit int) ([]Comment, error) {
	// Generate and execute query based on cursor value
	var rows pgx.Rows
	var err error
	rawQuery := `SELECT 
					id,
					post_id,
					author_id,
					body,
					created_at,
					updated_at
				FROM 
					comments`
	if cursor == 0 {
		rawQuery += ` LIMIT 
						$1;`
		rows, err = r.pgx.Query(ctx, rawQuery, limit)
	} else {
		rawQuery := ` WHERE
						id > $1
					LIMIT 
						$2;`
		rows, err = r.pgx.Query(ctx, rawQuery, cursor, limit)
	}
	if err != nil {
		return nil, fmt.Errorf("getting comment: %w", err)
	}
	defer rows.Close()

	// Scan returned rows
	var comments []Comment
	for rows.Next() {
		var comment Comment

		err = rows.Scan(
			&comment.Id,
			&comment.PostId,
			&comment.AuthorId,
			&comment.Body,
			&comment.CreatedAt,
			&comment.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return comments, nil
}

func (r *PostgresRepository) GetAllByAttribute(ctx context.Context, attribute string) ([]Comment, error) {
	return nil, nil
}

func (r *PostgresRepository) GetById(ctx context.Context, id uuid.UUID) (Comment, error) {
	rawQuery := `SELECT
					id,
					post_id,
					author_id,
					body,
					created_at,
					updated_at
				FROM
					comments
				WHERE
					id = $1;`

	var comment Comment

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery, id)

	// Scan returned row
	err := row.Scan(
		&comment.Id,
		&comment.PostId,
		&comment.AuthorId,
		&comment.Body,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return Comment{}, fmt.Errorf("scanning row: %w", err)
	}

	return comment, nil
}

func (r *PostgresRepository) Update(ctx context.Context, req UpdateRequest) (Comment, error) {
	rawQuery := `UPDATE
					comments
				SET
					body = $2,
					updated_at = $3
				WHERE
					id = $1
				RETURNING
					*;`

	var comment Comment
	now := time.Now()

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery, req.Id, req.Body, now)

	// Scan returned row
	err := row.Scan(
		&comment.Id,
		&comment.PostId,
		&comment.AuthorId,
		&comment.Body,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return Comment{}, fmt.Errorf("scanning row: %w", err)
	}

	return comment, nil
}

// Subjected to change
func (r *PostgresRepository) Replace(ctx context.Context, req ReplaceRequest) (Comment, error) {
	rawQuery := `UPDATE
					comments
				SET
					body = $2,
					updated_at = $3
				WHERE
					id = $1
				RETURNING
					*;`

	var comment Comment
	now := time.Now()

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery, req.Id, req.Body, now)

	// Scan returned row
	err := row.Scan(
		&comment.Id,
		&comment.PostId,
		&comment.AuthorId,
		&comment.Body,
		&comment.CreatedAt,
		&comment.UpdatedAt,
	)
	if err != nil {
		return Comment{}, fmt.Errorf("scanning row: %w", err)
	}

	return comment, nil
}

// Subjected to change
func (r *PostgresRepository) Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error) {
	rawQuery := `DELETE FROM
					comments
				WHERE
					id = $1
				RETURNING
					id;`

	var deletedId uuid.UUID

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery, req.Id)

	// Scan returned row
	err := row.Scan(&deletedId)
	if err != nil {
		return uuid.Nil, fmt.Errorf("scanning row: %w", err)
	}

	return deletedId, nil
}

// Bulk operations
func (r *PostgresRepository) BulkAdd(ctx context.Context, req BulkCreateRequest) ([]common.BulkOpResult, error) {
	return nil, nil
}

func (r *PostgresRepository) BulkModify(ctx context.Context, req BulkModifyRequest) ([]common.BulkOpResult, error) {
	return nil, nil
}

func (r *PostgresRepository) BulkDelete(ctx context.Context, req BulkDeleteRequest) ([]common.BulkOpResult, error) {
	return nil, nil
}
