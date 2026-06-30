package users

import (
	"context"
	"database/sql"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	logger *zap.Logger
	db     *sql.DB
	pgx    *pgx.Conn
}

func NewPostgresRepository(connection *sql.DB, l *zap.Logger) *PostgresRepository {
	return &PostgresRepository{
		logger: l,
		db:     connection,
	}
}

func NewPostgresRepositoryFromPGX(connection *pgx.Conn, l *zap.Logger) *PostgresRepository {
	return &PostgresRepository{
		logger: l,
		pgx:    connection,
	}
}

func (r *PostgresRepository) Initialize() error {
	return nil
}

func (r *PostgresRepository) Seed() error {
	return nil
}

// General operations
func (r *PostgresRepository) Add(ctx context.Context, dto CreateDTO) (User, error) {
	rawQuery := `INSERT INTO 
				users(username, displayname, email, country, state, avatar_url, langtag, timezone, created_at, updated_at, verified_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) 
				RETURNING *;`

	now := time.Now()
	dummyAvatarUrl := "-"

	result, err := r.pgx.Query(ctx, rawQuery,
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
		return User{}, err
	}

	var user User
	err = result.Scan(&user.Id, &user.Username, &user.DisplayName, &user.Email, &user.Country, &user.State, &user.AvatarURL, &user.LangTag, &user.Timezone,
		&user.CreatedAt, &user.UpdatedAt, &user.VerifiedAt)

	return user, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context, cursor int, limit int) ([]User, error) {
	var rows pgx.Rows
	var err error
	rawQuery := `SELECT 
				id, username, displayname, email, country, state, avatar_url, langtag, timezone, created_at, updated_at, verified_at
				FROM users`
	if cursor == 0 {
		rawQuery += ` ORDER BY id
					LIMIT $1;`

		rows, err = r.pgx.Query(ctx, rawQuery, limit)
	} else {
		rawQuery += ` WHERE id > $1
				ORDER BY id
				LIMIT $2;`
		rows, err = r.pgx.Query(ctx, rawQuery, cursor, limit)
	}

	var users []User
	for rows.Next() {
		var user User
		err = rows.Scan(&user.Id, &user.Username, &user.DisplayName, &user.Email, &user.Country, &user.State, &user.AvatarURL, &user.LangTag, &user.Timezone,
			&user.CreatedAt, &user.UpdatedAt, &user.VerifiedAt)
		if err != nil {

		}

		users = append(users, user)
	}
	defer rows.Close()

	return users, err
}

func (r *PostgresRepository) GetAllByAttribute(ctx context.Context, attribute string) ([]User, error) {
	return nil, nil
}

func (r *PostgresRepository) GetById(ctx context.Context, id uuid.UUID) (User, error) {
	rawQuery := `SELECT
				id, username, displayname, email, country, state, avatar_url, langtag, timezone, created_at, updated_at, verified_at
				FROM users
				WHERE id = $1;`

	result, err := r.pgx.Query(ctx, rawQuery, id)
	if err != nil {
		r.logger.Error("Failed to execute egtbyid query", zap.Error(err))
		return User{}, err
	}

	var user User
	err = result.Scan(&user.Id, &user.Username, &user.DisplayName, &user.Email, &user.Country, &user.State, &user.AvatarURL, &user.LangTag, &user.Timezone,
		&user.CreatedAt, &user.UpdatedAt, &user.VerifiedAt)
	if err != nil {
		r.logger.Error("Failed to scan result in query", zap.Error(err))
		return User{}, err
	}

	return user, nil
}

func (r *PostgresRepository) GetByName(ctx context.Context, name string) (User, error) {
	rawQuery := `SELECT
				id, username, displayname, email, country, state, avatar_url, langtag, timezone, created_at, updated_at, verified_at
				FROM users
				WHERE username = $1;`

	result, err := r.pgx.Query(ctx, rawQuery, name)
	if err != nil {
		r.logger.Error("Failed to execute getbyid query", zap.Error(err))
		return User{}, err
	}

	var user User
	err = result.Scan(&user.Id, &user.Username, &user.DisplayName, &user.Email, &user.Country, &user.State, &user.AvatarURL, &user.LangTag, &user.Timezone,
		&user.CreatedAt, &user.UpdatedAt, &user.VerifiedAt)
	if err != nil {
		r.logger.Error("Failed to scan result in query", zap.Error(err))
		return User{}, err
	}

	return user, nil
}

func (r *PostgresRepository) Update(ctx context.Context, dto UpdateDTO) (User, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	builder := psql.Update("users").Where(squirrel.Eq{"id": dto.Id})
	for i := range dto.Updates {
		builder.Set(dto.Updates[i].Field, dto.Updates[i].Value)
	}
	builder.Suffix("RETURNING *")
	query, args, err := builder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate update query using squirrel", zap.Error(err))
		return User{}, err
	}

	var user User
	result := r.pgx.QueryRow(ctx, query, args...)
	err = result.Scan(&user.Id, &user.Username, &user.DisplayName, &user.Email, &user.Country, &user.State, &user.AvatarURL, &user.LangTag, &user.Timezone,
		&user.CreatedAt, &user.UpdatedAt, &user.VerifiedAt)
	if err != nil {
		r.logger.Error("Failed to scan returned object from update query", zap.Error(err))
		return User{}, err
	}

	return user, nil
}

func (r *PostgresRepository) Replace(ctx context.Context, dto ReplaceDTO) (User, error) {
	rawQuery := `UPDATE users
				SET
				username = $2,
				displayname = $3,
				email = $4,
				country = $5,
				state = $6,
				avatar_url = '-',
				langtag = $8,
				timezone = $9,
				created_at = now(),
				updated_at = now(),
				verified_at = now(),
				WHERE
				id = $1
				RETURNING *;`

	result := r.pgx.QueryRow(ctx, rawQuery, dto.Id, dto.Username, dto.DisplayName, dto.Email, dto.Country, dto.State, dto.LangTag, dto.Timezone)
	var user User
	err := result.Scan(&user.Id, &user.Username, &user.DisplayName, &user.Email, &user.Country, &user.State, &user.AvatarURL, &user.LangTag, &user.Timezone,
		&user.CreatedAt, &user.UpdatedAt, &user.VerifiedAt)
	if err != nil {
		r.logger.Error("Failed to scan result of replace query", zap.Error(err))
		return User{}, err
	}

	return user, err
}

func (r *PostgresRepository) Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error) {
	deleteQuery := `DELETE FROM
					users
					WHERE id = $1
					RETURNING id;`
	disableQuery := `UPDATE
					users
					SET disabled = true
					WHERE id = $1
					RETURNING id;`

	var result pgx.Row
	if dto.DeleteOptions.Type == "disable" {
		result = r.pgx.QueryRow(ctx, disableQuery, dto.Id)
	} else {
		result = r.pgx.QueryRow(ctx, deleteQuery, dto.Id)
	}

	var deletedId uuid.UUID
	if err := result.Scan(&deletedId); err != nil {
		r.logger.Error("Failed to scan result in delete query", zap.Error(err))
		return uuid.Nil, err
	}

	return deletedId, nil
}

// Bulk operations
func (r *PostgresRepository) BulkAdd(ctx context.Context, dto BulkCreateDTO) ([]common.BulkOpResult, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	builder := psql.Insert("users").Columns(
		"id",
		"username",
		"displayname",
		"email",
		"country",
		"state",
		"avatar_url",
		"langtag",
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
		builder.Values(id, dto.Users[i].Username, dto.Users[i].DisplayName, dto.Users[i].Email, dto.Users[i].Country, dto.Users[i].State,
			dummyAvatarUrl, dto.Users[i].LangTag, dto.Users[i].Timezone, now, now, now)

		var result common.BulkOpResult
		result.Id = id
		result.Status = 200
		result.Success = true
		result.Message = "added"
		results = append(results, result)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate bulk insert users query using squirrel", zap.Error(err))
		return nil, common.ErrResourceNotAdded
	}

	_, err = r.pgx.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to execute bulk insert query", zap.Error(err))
		return nil, common.ErrResourceNotAdded
	}

	return results, nil
}

func (r *PostgresRepository) BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error) {
	psql := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	builder := psql.Update("users")

	for i := range dto.Updates {
		currentUser := dto.Updates[i]
		for j := range currentUser.Updates {
			builder = builder.Set(currentUser.Updates[j].Field, currentUser.Updates[j].Value)
		}
	}

	query, args, err := builder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate bulk modify query with squirrel", zap.Error(err))
		return nil, common.ErrResourcesCannotBeModified
	}

	_, err = r.pgx.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to execute bulk modify query", zap.Error(err))
		return nil, common.ErrResoucesCannotBeModified
	}

	return nil, nil
}

func (r *PostgresRepository) BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error) {
	return nil, nil
}
