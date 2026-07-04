package posts

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	logger           *zap.Logger
	pgx              *pgxpool.Pool
	statementBuilder squirrel.StatementBuilderType
	seedfile         string
}

func NewPostgresRepositoryFromPgxPool(connection *pgxpool.Pool, l *zap.Logger, sl string) *PostgresRepository {
	return &PostgresRepository{
		logger:           l,
		pgx:              connection,
		statementBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		seedfile:         sl,
	}
}

func (r *PostgresRepository) Initialize(ctx context.Context) error {
	return nil
}

func (r *PostgresRepository) Seed(ctx context.Context) error {
	// Check if rows already exists, if so then dont seed
	existQuery := `SELECT count(*) FROM posts;`
	row := r.pgx.QueryRow(ctx, existQuery)

	var count int
	if err := row.Scan(&count); err != nil {
		r.logger.Error("Failed to scan count of rows from posts table", zap.Error(err))
		return err
	}

	// Rows already exists, dont seed
	if count != 0 {
		r.logger.Sugar().Infof("Posts table contains %d rows, skipping seeding from file...", count)
		return nil
	}

	// Table is empty, so seed from file to table
	queryBuilder := r.statementBuilder.Insert("posts").Columns(
		"id",
		"title",
		"body",
		"author_id",
		"likes",
		"comments",
		"created_at",
		"updated_at",
	)

	file, err := os.OpenFile("./seeds/posts.json", os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open posts seeds file", zap.Error(err))
		return fmt.Errorf("Failed to open posts's seeds file. Error: %w", err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
		return fmt.Errorf("Failed to create json decoder. Returned nil pointer")
	}

	var posts []Post
	err = decoder.Decode(&posts)
	if err != nil {
		r.logger.Error("Failed to decode seeded users from file", zap.Error(err))
		return fmt.Errorf("Faield to decode seeded posts from file")
	}

	now := time.Now()
	for i := range posts {
		id, _ := uuid.NewV7()
		queryBuilder = queryBuilder.Values(id, posts[i].Title, posts[i].Body, posts[i].AuthorId, posts[i].Likes,
			posts[i].Comments, now, now)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate sql query with squirrel for seeding posts table", zap.Error(err))
		return fmt.Errorf("Failed seeding posts table")
	}

	_, err = r.pgx.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to execute query to seed posts table", zap.Error(err))
		return fmt.Errorf("Failed seeding posts table")
	}

	return nil
}

// General operations
func (r *PostgresRepository) Add(ctx context.Context, dto CreateDTO) (Post, error) {
	rawQuery := `INSERT INTO
					posts
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8);`

	id, _ := uuid.NewV7()
	likes, comments := 0, 0
	now := time.Now()

	_, err := r.pgx.Exec(ctx, rawQuery, id, dto.Title, dto.Body, dto.AuthorId, likes, comments, now, now)
	if err != nil {
		r.logger.Error("Failed to execute sql insert query", zap.Error(err))
		return Post{}, common.TranslatePostgresError(err, r.logger)
	}

	return Post{
		Id:        id,
		Title:     dto.Title,
		AuthorId:  dto.AuthorId,
		Likes:     likes,
		Comments:  comments,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context, cursor int, limit int) ([]Post, error) {
	var rows pgx.Rows
	var err error
	rawQuery := `SELECT 
					id, 
					title, 
					body, 
					author_id, 
					likes, 
					comments, 
					created_at, 
					updated_at
				FROM
					posts`
	if cursor == 0 {
		rawQuery += ` ORDER BY id 
					LIMIT $1;`

		rows, err = r.pgx.Query(ctx, rawQuery, limit)
	} else {
		rawQuery += ` WHERE 
						id > $1
					ORDER BY 
						id
					LIMIT $2;`

		rows, err = r.pgx.Query(ctx, rawQuery, cursor, limit)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var post Post
		err = rows.Scan(&post.Id, &post.Title, &post.Body, &post.AuthorId, &post.Likes, &post.Comments,
			&post.CreatedAt, &post.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("Failed to scan rows")
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (r *PostgresRepository) GetAllByAttribute(ctx context.Context, attribute string) ([]Post, error) {
	return nil, nil
}

func (r *PostgresRepository) GetById(ctx context.Context, id uuid.UUID) (Post, error) {
	rawQuery := `SELECT 
					id,
					title,
					body,
					author_id,
					likes,
					comments,
					created_at,
					updated_at
				FROM 
					posts
				WHERE
					id = $1;`

	var post Post
	row := r.pgx.QueryRow(ctx, rawQuery, id)
	err := row.Scan(&post.Id, &post.Title, &post.Body, &post.AuthorId, &post.Likes, &post.Comments,
		&post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		r.logger.Error("Failed scan row in getbyid query", zap.Error(err))
		return Post{}, common.TranslatePostgresError(err, r.logger)
	}

	return post, nil
}

func (r *PostgresRepository) GetByName(ctx context.Context, name string) (Post, error) {
	rawQuery := `SELECT 
					id,
					title,
					body,
					author_id,
					likes,
					comments,
					created_at,
					updated_at
				FROM 
					posts
				WHERE
					title = $1;`

	var post Post
	row := r.pgx.QueryRow(ctx, rawQuery, name)
	err := row.Scan(&post.Id, &post.Title, &post.Body, &post.AuthorId, &post.Likes, &post.Comments,
		&post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		r.logger.Error("Failed scan row in getbyid query", zap.Error(err))
		return Post{}, common.TranslatePostgresError(err, r.logger)
	}

	return post, nil
}

func (r *PostgresRepository) Update(ctx context.Context, dto UpdateDTO) (Post, error) {
	queryBuilder := r.statementBuilder.Update("posts").Where(squirrel.Eq{"id": dto.Id})
	for i := range dto.Updates {
		queryBuilder = queryBuilder.Set(dto.Updates[i].Field, dto.Updates[i].Value)
	}
	queryBuilder = queryBuilder.Suffix("RETURNING *")
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate update query using squirrel", zap.Error(err))
		return Post{}, common.TranslatePostgresError(err, r.logger)
	}

	tx, err := r.pgx.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		r.logger.Error("Failed to begin transaction in update query", zap.Error(err))
		return Post{}, common.TranslatePostgresError(err, r.logger)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "SELECT id FROM posts WHERE id = $1 FOR UPDATE;", dto.Id)
	if err != nil {
		r.logger.Error("Failed to lock row for update in transaction", zap.Error(err))
		return Post{}, common.TranslatePostgresError(err, r.logger)
	}

	var post Post
	row := tx.QueryRow(ctx, query, args...)
	err = row.Scan(&post.Id, &post.Title, &post.Body, &post.AuthorId, &post.Likes, &post.Comments,
		&post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		r.logger.Error("Failed to scan returned row from update query", zap.Error(err))
		return Post{}, common.TranslatePostgresError(err, r.logger)
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Error("Failed to commit transaction in update query", zap.Error(err))
		return Post{}, common.TranslatePostgresError(err, r.logger)
	}

	return post, nil
}

func (r *PostgresRepository) Replace(ctx context.Context, dto ReplaceDTO) (Post, error) {
	rawQuery := `UPDATE
					posts
				SET
					title = $2,
					body = $3,
					likes = $4,
					comments = $5,
					updated_at = $6
				WHERE
					id = $1
				RETURNING 
					*;`

	tx, err := r.pgx.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		r.logger.Error("Failed to begin transaction in replace query", zap.Error(err))
		return Post{}, common.TranslatePostgresError(err, r.logger)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, "SELECT id FROM posts WHERE id = $1 FOR UPDATE;", dto.Id)
	if err != nil {
		r.logger.Error("Failed to lock row for replace in transaction", zap.Error(err))
		return Post{}, common.TranslatePostgresError(err, r.logger)
	}

	var post Post
	likes, comments := 0, 0
	now := time.Now()

	row := tx.QueryRow(ctx, rawQuery, dto.Id, dto.Title, dto.Body, likes, comments, now)
	err = row.Scan(&post.Id, &post.Title, &post.Body, &post.AuthorId, &post.Likes, &post.Comments,
		&post.CreatedAt, &post.UpdatedAt)
	if err != nil {
		r.logger.Error("Failed to scan result of replace query", zap.Error(err))
		return Post{}, common.TranslatePostgresError(err, r.logger)
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Error("Failed to commit transaction in replace query", zap.Error(err))
		return Post{}, common.TranslatePostgresError(err, r.logger)
	}

	return post, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error) {
	rawQuery := `DELETE FROM
					posts
				WHERE
					id = $1
				RETURNING
					id;`

	row := r.pgx.QueryRow(ctx, rawQuery, dto.Id)

	var deletedId uuid.UUID
	err := row.Scan(&deletedId)
	if err != nil {
		r.logger.Error("Failed to scan returned row from delete query", zap.Error(err))
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
