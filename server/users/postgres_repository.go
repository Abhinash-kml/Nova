package users

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
	logger           *zap.Logger
	pgx              *pgxpool.Pool
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
	existQuery := `SELECT count(*) FROM users;`
	row := r.pgx.QueryRow(ctx, existQuery)

	var count int

	// Scan returned row
	if err := row.Scan(&count); err != nil {
		r.logger.Error("Failed to scan count of rows from users table", zap.Error(err))
		return common.ErrDbSeedingFailed
	}

	// Rows already exists, dont seed
	if count != 0 {
		r.logger.Sugar().Infof("Users table contains %d rows, skipping seeding from file...", count)
		return nil
	}

	// Table is empty, so seed from file to table
	queryBuilder := r.statementBuilder.Insert("users").Columns(
		"id",
		"username",
		"displayname",
		"email",
		"country",
		"state",
		"avatar_url",
		"lang_tag",
		"timezone",
		"created_at",
		"updated_at",
		"verified_at",
	)

	// Check if seeds file is provided and is of json format by suffix
	if len(r.seedfile) == 0 || !strings.HasSuffix(r.seedfile, ".json") {
		r.logger.Error("Seeds file not provided or is not of json format. Skipping seeding")
		return common.ErrDbSeedingFailed
	}

	// open file
	file, err := os.OpenFile(r.seedfile, os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open users seeds file", zap.String("file", r.seedfile), zap.Error(err))
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
	var users []User
	err = decoder.Decode(&users)
	if err != nil {
		r.logger.Error("Failed to decode seeded users from file", zap.Error(err))
		return common.ErrDbSeedingFailed
	}

	for i := range users {
		queryBuilder = queryBuilder.Values(users[i].Id, users[i].Username, users[i].DisplayName, users[i].Email, users[i].Country, users[i].State,
			users[i].AvatarURL, users[i].LangTag, users[i].Timezone, users[i].CreatedAt, users[i].UpdatedAt, users[i].VerifiedAt)
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

	r.logger.Info("Successfully seeded users from file", zap.Int("count", len(users)))

	return nil
}

// General operations
func (r *PostgresRepository) Add(ctx context.Context, dto CreateDTO) (User, error) {
	rawQuery := `INSERT INTO 
					users(id, username, displayname, email, country, state, avatar_url, lang_tag, timezone, created_at, updated_at, verified_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`

	id, _ := uuid.NewV7()
	now := time.Now()
	dummyAvatarUrl := "-"

	_, err := r.pgx.Exec(ctx, rawQuery,
		id,
		dto.Username,
		dto.DisplayName,
		dto.Email,
		dto.Country,
		dto.State,
		dummyAvatarUrl,
		dto.LangTag,
		dto.Timezone,
		now,
		now,
		now)
	if err != nil {
		r.logger.Error("Failed to execute sql insert query", zap.Error(err))
		return User{}, common.TranslatePostgresError(err, r.logger)
	}

	return User{
		Id:          id,
		Username:    dto.Username,
		DisplayName: dto.DisplayName,
		Email:       dto.Email,
		Country:     dto.Country,
		State:       dto.State,
		AvatarURL:   dummyAvatarUrl,
		LangTag:     dto.LangTag,
		Timezone:    dto.Timezone,
		CreatedAt:   now,
		UpdatedAt:   now,
		VerifiedAt:  now,
	}, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context, cursor int, limit int) ([]User, error) {
	var rows pgx.Rows
	var err error

	// Create and execute query
	rawQuery := `SELECT 
					id, 
					username, 
					displayname, 
					email, 
					country, 
					state, 
					avatar_url, 
					lang_tag, 
					timezone, 
					created_at, 
					updated_at, 
					verified_at
				FROM users`
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

	// Scan returned rows
	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username, &user.DisplayName, &user.Email, &user.Country, &user.State, &user.AvatarURL, &user.LangTag, &user.Timezone,
			&user.CreatedAt, &user.UpdatedAt, &user.VerifiedAt)
		if err != nil {
			r.logger.Error("Failed to sca nreturned row in getall query", zap.Error(err))
			return nil, common.TranslatePostgresError(err, r.logger)
		}

		users = append(users, user)
	}

	return users, err
}

func (r *PostgresRepository) GetAllByAttribute(ctx context.Context, attribute string) ([]User, error) {
	return nil, nil
}

func (r *PostgresRepository) GetById(ctx context.Context, id uuid.UUID) (User, error) {
	rawQuery := `SELECT
					id, 
					username, 
					displayname, 
					email, 
					country, 
					state, 
					avatar_url, 
					lang_tag, 
					timezone, 
					created_at, 
					updated_at, 
					verified_at, 
					disabled_at
				FROM 
					users
				WHERE 
					id = $1;`

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery, id)

	var user User
	var disabledAt *time.Time

	// Scan returned row
	err := row.Scan(&user.Id, &user.Username, &user.DisplayName, &user.Email, &user.Country, &user.State, &user.AvatarURL, &user.LangTag, &user.Timezone,
		&user.CreatedAt, &user.UpdatedAt, &user.VerifiedAt, &disabledAt)
	if err != nil {
		r.logger.Error("Failed to scan row in getbyid query", zap.Error(err))
		return User{}, common.TranslatePostgresError(err, r.logger)
	}

	if disabledAt != nil && !disabledAt.IsZero() {
		user.DisabledAt = *disabledAt
	}

	return user, nil
}

func (r *PostgresRepository) GetByName(ctx context.Context, name string) (User, error) {
	rawQuery := `SELECT
					id, 
					username, 
					displayname, 
					email, 
					country, 
					state, 
					avatar_url, 
					lang_tag, 
					timezone, 
					created_at, 
					updated_at, 
					verified_at, 
					disabled_at
				FROM 
					users
				WHERE 
					username = $1;`

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery, name)

	var user User
	var disabledAt *time.Time

	// Scasn returned row
	err := row.Scan(&user.Id, &user.Username, &user.DisplayName, &user.Email, &user.Country, &user.State, &user.AvatarURL, &user.LangTag, &user.Timezone,
		&user.CreatedAt, &user.UpdatedAt, &user.VerifiedAt, &disabledAt)
	if err != nil {
		r.logger.Error("Failed to scan result in getbyname query", zap.Error(err))
		return User{}, common.TranslatePostgresError(err, r.logger)
	}

	if disabledAt != nil && !disabledAt.IsZero() {
		user.DisabledAt = *disabledAt
	}

	return user, nil
}

func (r *PostgresRepository) Update(ctx context.Context, dto UpdateDTO) (User, error) {
	// build dynamic query using squirrel
	queryBuilder := r.statementBuilder.Update("users").Where(squirrel.Eq{"id": dto.Id})
	for i := range dto.Updates {
		queryBuilder = queryBuilder.Set(dto.Updates[i].Field, dto.Updates[i].Value)
	}
	queryBuilder = queryBuilder.Suffix("RETURNING *")
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate update query using squirrel", zap.Error(err))
		return User{}, common.TranslatePostgresError(err, r.logger)
	}

	var user User
	var disabledAt *time.Time

	// Begin transaction
	tx, err := r.pgx.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		r.logger.Error("Failed to begin transaction in update query", zap.Error(err))
		return User{}, common.TranslatePostgresError(err, r.logger)
	}
	defer tx.Rollback(ctx)

	// Lock row inside transaction
	_, err = tx.Exec(ctx, "SELECT username FROM users WHERE id = $1 FOR UPDATE;", dto.Id)
	if err != nil {
		r.logger.Error("Failed to lock row for update in transaction", zap.Error(err))
		return User{}, common.TranslatePostgresError(err, r.logger)
	}

	// Execute query
	result := tx.QueryRow(ctx, query, args...)
	err = result.Scan(&user.Id, &user.Username, &user.DisplayName, &user.Email, &user.Country, &user.State, &user.AvatarURL, &user.LangTag, &user.Timezone,
		&user.CreatedAt, &user.UpdatedAt, &user.VerifiedAt, &disabledAt)
	if err != nil {
		r.logger.Error("Failed to scan returned object from update query in transaction", zap.Error(err))
		return User{}, common.TranslatePostgresError(err, r.logger)
	}

	if disabledAt != nil && !disabledAt.IsZero() {
		user.DisabledAt = *disabledAt
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Error("Failed to commit transaction in update query", zap.Error(err))
		return User{}, common.TranslatePostgresError(err, r.logger)
	}

	return user, nil
}

func (r *PostgresRepository) Replace(ctx context.Context, dto ReplaceDTO) (User, error) {
	rawQuery := `UPDATE 
					users
				SET
					username = $2,
					displayname = $3,
					email = $4,
					country = $5,
					state = $6,
					avatar_url = '-',
					lang_tag = $7,
					timezone = $8,
					updated_at = $9,
					verified_at = $10
				WHERE
					id = $1
				RETURNING 
						*;`

	var user User
	var disabledAt *time.Time
	now := time.Now()

	// Begin transaction
	tx, err := r.pgx.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return User{}, common.TranslatePostgresError(err, r.logger)
	}
	defer tx.Rollback(ctx)

	// Lock row inside transaction
	_, err = tx.Exec(ctx, "SELECT username FROM users WHERE id = $1 FOR UPDATE;", dto.Id)
	if err != nil {
		r.logger.Error("Failed to lock row for replace in transaction", zap.Error(err))
		return User{}, common.TranslatePostgresError(err, r.logger)
	}

	// Execute query
	result := tx.QueryRow(ctx, rawQuery, dto.Id, dto.Username, dto.DisplayName, dto.Email, dto.Country, dto.State, dto.LangTag, dto.Timezone, now, now)
	err = result.Scan(&user.Id, &user.Username, &user.DisplayName, &user.Email, &user.Country, &user.State, &user.AvatarURL, &user.LangTag, &user.Timezone,
		&user.CreatedAt, &user.UpdatedAt, &user.VerifiedAt, &disabledAt)
	if err != nil {
		r.logger.Error("Failed to scan result of replace query", zap.Error(err))
		return User{}, common.TranslatePostgresError(err, r.logger)
	}

	if disabledAt != nil && !disabledAt.IsZero() {
		user.DisabledAt = *disabledAt
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Error("Failed to commit transaction in replace query", zap.Error(err))
		return User{}, common.TranslatePostgresError(err, r.logger)
	}

	return user, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error) {
	deleteQuery := `DELETE FROM
						users
					WHERE 
						id = $1
					RETURNING 
						id;`
	disableQuery := `UPDATE
						users
					SET 
						disabled_at = $2
					WHERE 
						id = $1
					RETURNING 
						id;`

	now := time.Now()
	var result pgx.Row
	var deletedId uuid.UUID

	// Execute query based on delete type
	if dto.DeleteOptions.Type == "soft" {
		result = r.pgx.QueryRow(ctx, disableQuery, dto.Id, now)
	} else {
		result = r.pgx.QueryRow(ctx, deleteQuery, dto.Id)
	}

	// Scan returned row
	if err := result.Scan(&deletedId); err != nil {
		r.logger.Error("Failed to scan result in delete query", zap.Error(err))
		return uuid.Nil, common.TranslatePostgresError(err, r.logger)
	}

	return deletedId, nil
}

// Bulk operations - subject to change for improvement (dont use these queries right now)
func (r *PostgresRepository) BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error) {
	queryBuilder := r.statementBuilder.Insert("users").Columns(
		"id",
		"username",
		"displayname",
		"email",
		"country",
		"state",
		"avatar_url",
		"lang_tag",
		"timezone",
		"created_at",
		"updated_at",
		"verified_at",
	)

	var results []common.BulkOpResult

	for i := range dto.Users {
		id, _ := uuid.NewV7()
		now := time.Now()
		dummyAvatarUrl := "-"
		queryBuilder = queryBuilder.Values(id, dto.Users[i].Username, dto.Users[i].DisplayName, dto.Users[i].Email, dto.Users[i].Country, dto.Users[i].State,
			dummyAvatarUrl, dto.Users[i].LangTag, dto.Users[i].Timezone, now, now, now)

		var result common.BulkOpResult
		result.Id = id
		result.Status = 1
		result.Success = true
		result.Message = "added"
		results = append(results, result)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate bulk insert users query using squirrel", zap.Error(err))
		return nil, common.TranslatePostgresError(err, r.logger)
	}

	_, err = r.pgx.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to execute bulk insert query", zap.Error(err))
		for i := range results {
			results[i].Status = 0
			results[i].Success = false
			results[i].Message = "failed"
		}
		return results, common.TranslatePostgresError(err, r.logger)
	}

	return results, nil
}

func (r *PostgresRepository) BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error) {
	queryBuilder := r.statementBuilder.Update("users")

	var results []common.BulkOpResult

	for i := range dto.Updates {
		currentUser := dto.Updates[i]
		for j := range currentUser.Updates {
			queryBuilder = queryBuilder.Set(currentUser.Updates[j].Field, currentUser.Updates[j].Value)
		}

		var result common.BulkOpResult
		id, _ := uuid.Parse(currentUser.Id)
		result.Id = id
		result.Status = 1
		result.Success = true
		result.Message = "modified"
		results = append(results, result)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate bulk modify query with squirrel", zap.Error(err))
		return nil, common.ErrResourcesCannotBeModified
	}

	_, err = r.pgx.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to execute bulk modify query", zap.Error(err))
		for i := range results {
			results[i].Status = 0
			results[i].Success = false
			results[i].Message = "failed"
		}
		return results, common.ErrResourcesCannotBeModified
	}

	return nil, nil
}

func (r *PostgresRepository) BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error) {
	queryBuilder := r.statementBuilder.Delete("users")

	var results []common.BulkOpResult

	for i := range dto.Users {
		currentUser := dto.Users[i]
		queryBuilder = queryBuilder.Where(squirrel.Eq{"id": currentUser.Id})

		var result common.BulkOpResult
		id, _ := uuid.Parse(currentUser.Id)
		result.Id = id
		result.Status = 1
		result.Success = true
		result.Message = "deleted"
		results = append(results, result)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate bulk delete query with squirrel", zap.Error(err))
		return nil, common.ErrResourcesCannotBeDeleted
	}

	_, err = r.pgx.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to execute bulk delete query", zap.Error(err))
		for i := range results {
			results[i].Status = 0
			results[i].Message = "failed"
			results[i].Success = false
		}
		return results, common.TranslatePostgresError(err, r.logger)
	}

	return results, nil
}
