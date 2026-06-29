package users

import (
	"context"
	"database/sql"
	"time"

	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	logger *zap.Logger
	db     *sql.DB
}

func NewPostgresRepository(connection *sql.DB, l *zap.Logger) *PostgresRepository {
	return &PostgresRepository{
		logger: l,
		db:     connection,
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

	result, err := r.db.QueryContext(ctx, rawQuery,
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
	var rows *sql.Rows
	var err error
	rawQuery := `SELECT 
				id, username, displayname, email, country, state, avatar_url, langtag, timezone, created_at, updated_at, verified_at
				FROM users`
	if cursor == 0 {
		rawQuery += ` ORDER BY id
					LIMIT $1;`

		rows, err = r.db.QueryContext(ctx, rawQuery, limit)
	} else {
		rawQuery += ` WHERE id > $1
				ORDER BY id
				LIMIT $2;`
		rows, err = r.db.QueryContext(ctx, rawQuery, cursor, limit)
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

	result, err := r.db.QueryContext(ctx, rawQuery, id)
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

	result, err := r.db.QueryContext(ctx, rawQuery, name)
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
	return User{}, nil
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

	result := r.db.QueryRowContext(ctx, rawQuery, dto.Id, dto.Username, dto.DisplayName, dto.Email, dto.Country, dto.State, dto.LangTag, dto.Timezone)
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

	var result *sql.Row
	if dto.DeleteOptions.Type == "disable" {
		result = r.db.QueryRowContext(ctx, disableQuery)
	} else {
		result = r.db.QueryRowContext(ctx, deleteQuery)
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
	return nil, nil
}

func (r *PostgresRepository) BulkModify(ctx context.Context, dto BulkModifyDTO) ([]common.BulkOpResult, error) {
	return nil, nil
}

func (r *PostgresRepository) BulkDelete(ctx context.Context, dto BulkDeleteDTO) ([]common.BulkOpResult, error) {
	return nil, nil
}
