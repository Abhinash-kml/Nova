package comments

import (
	"context"
	"encoding/json"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	pgx              *pgxpool.Pool
	logger           *zap.Logger
	statementBuilder squirrel.StatementBuilderType
	seedfile         string
}

func NewPostgresRepositoryFromPgxPool(connection *pgxpool.Pool, l *zap.Logger, sfp string) *PostgresRepository {
	return &PostgresRepository{
		logger:           l,
		pgx:              connection,
		statementBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		seedfile:         sfp,
	}
}

func (r *PostgresRepository) Initialize(ctx context.Context) error {
	return nil
}

func (r *PostgresRepository) Seed(ctx context.Context) error {
	// Check if rows already exists, if so then dont seed
	existQuery := `SELECT count(*) FROM comments;`
	row := r.pgx.QueryRow(ctx, existQuery)

	var count int

	// Scan returned row
	if err := row.Scan(&count); err != nil {
		r.logger.Error("Failed to scan count of rows from comments table", zap.Error(err))
		return common.ErrDbSeedingFailed
	}

	// Rows already exists, dont seed
	if count != 0 {
		r.logger.Sugar().Infof("Comments table contains %d rows, skipping seeding from file...", count)
		return nil
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
	if len(r.seedfile) == 0 || !strings.HasSuffix(r.seedfile, ".json") {
		r.logger.Error("Seeds file not provided or is not of json format. Skipping seeding")
		return common.ErrDbSeedingFailed
	}

	// open file
	file, err := os.OpenFile(r.seedfile, os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open comments seeds file", zap.String("file", r.seedfile), zap.Error(err))
		return common.ErrDbSeedingFailed
	}
	defer file.Close()

	// Create json decoder
	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
		return common.ErrDbSeedingFailed
	}

	// Decode
	var comments []Comment
	err = decoder.Decode(&comments)
	if err != nil {
		r.logger.Error("Failed to decode seeded users from file", zap.Error(err))
		return common.ErrDbSeedingFailed
	}

	for i := range comments {
		queryBuilder = queryBuilder.Values(comments[i].Id, comments[i].PostId, comments[i].AuthorId, comments[i].Body,
			comments[i].CreatedAt, comments[i].UpdatedAt)
	}

	// Generate query
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate sql query with squirrel for seeding users table", zap.Error(err))
		return common.ErrDbSeedingFailed
	}

	// Execute query
	_, err = r.pgx.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to execute query to seed users table", zap.Error(err))
		return common.ErrDbSeedingFailed
	}

	r.logger.Info("Successfully seeded comments from file", zap.Int("count", len(comments)))

	return nil
}

// General operations
func (r *PostgresRepository) Add(ctx context.Context, dto CreateDTO) (Comment, error) {
	rawQuery := `INSERT INTO
					comments(id, post_id, author_id, body, created_at, updated_at)
				VALUES($1, $2, $3, $4, $5, $6);`

	id, _ := uuid.NewV7()
	now := time.Now()

	// Execute quury
	_, err := r.pgx.Exec(ctx, rawQuery, id, dto.PostId, dto.Body, now, now)
	if err != nil {
		r.logger.Error("Failed to execute sql insert query", zap.Error(err))
		return Comment{}, common.TranslatePostgresError(err, r.logger)
	}

	return Comment{
		Id:        id,
		PostId:    dto.PostId,
		Body:      dto.Body,
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
	defer rows.Close()

	// Scan returned rows
	var comments []Comment
	for rows.Next() {
		var comment Comment
		err = rows.Scan(&comment.Id, &comment.PostId, &comment.AuthorId, &comment.Body, &comment.CreatedAt,
			&comment.UpdatedAt)
		if err != nil {
			r.logger.Error("Failed to scan returned row in getall query", zap.Error(err))
			return nil, common.TranslatePostgresError(err, r.logger)
		}

		comments = append(comments, comment)
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
	err := row.Scan(&comment.Id, &comment.PostId, &comment.AuthorId, &comment.Body, &comment.CreatedAt,
		&comment.UpdatedAt)
	if err != nil {
		r.logger.Error("Failed to scan row in getbyid query", zap.Error(err))
		return Comment{}, common.TranslatePostgresError(err, r.logger)
	}

	return comment, nil
}

func (r *PostgresRepository) Update(ctx context.Context, dto UpdateDTO) (Comment, error) {
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
	row := r.pgx.QueryRow(ctx, rawQuery, dto.Id, dto.Body, now)

	// Scan returned row
	err := row.Scan(&comment.Id, &comment.PostId, &comment.AuthorId, &comment.Body, &comment.CreatedAt,
		&comment.UpdatedAt)
	if err != nil {
		r.logger.Error("Failed to scan returned row in update query", zap.Error(err))
		return Comment{}, common.TranslatePostgresError(err, r.logger)
	}

	return comment, nil
}

// Subjected to change
func (r *PostgresRepository) Replace(ctx context.Context, dto ReplaceDTO) (Comment, error) {
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
	row := r.pgx.QueryRow(ctx, rawQuery, dto.Id, dto.Body, now)

	// Scan returned row
	err := row.Scan(&comment.Id, &comment.PostId, &comment.AuthorId, &comment.Body, &comment.CreatedAt,
		&comment.UpdatedAt)
	if err != nil {
		r.logger.Error("Failed to scan returned row in update query", zap.Error(err))
		return Comment{}, common.TranslatePostgresError(err, r.logger)
	}

	return comment, nil
}

// Subjected to change
func (r *PostgresRepository) Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error) {
	rawQuery := `DELETE FROM
					comments
				WHERE
					id = $1
				RETURNING
					id;`

	var deletedId uuid.UUID

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery, dto.Id)

	// Scan returned row
	err := row.Scan(&deletedId)
	if err != nil {
		r.logger.Error("Failed to scan returned row in delete query", zap.Error(err))
		return uuid.Nil, common.TranslatePostgresError(err, r.logger)
	}

	return deletedId, nil
}

// Bulk operations
func (r *PostgresRepository) BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error) {
	return nil, nil
}

func (r *PostgresRepository) BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error) {
	return nil, nil
}

func (r *PostgresRepository) BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error) {
	return nil, nil
}
