package users

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/auth"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	config           *config.UsersConfig
	logger           *zap.Logger
	pgx              *pgxpool.Pool
	statementBuilder squirrel.StatementBuilderType
}

func NewPostgresRepositoryFromPgxPool(connection *pgxpool.Pool, c *config.UsersConfig, l *zap.Logger) *PostgresRepository {
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
	existQuery := `SELECT count(*) FROM users;`
	row := r.pgx.QueryRow(ctx, existQuery)

	var count int

	// Scan returned row
	if err := row.Scan(&count); err != nil {
		r.logger.Error("Failed to scan count of rows from users table", zap.Error(err))
	}

	// Rows already exists, dont seed
	if count != 0 {
		r.logger.Sugar().Infof("Users table contains %d rows, skipping seeding from file...", count)
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
	if len(r.config.SeedFile) == 0 || !strings.HasSuffix(r.config.SeedFile, ".json") {
		r.logger.Error("Seeds file not provided or is not of json format. Skipping seeding")
	}

	// open file
	file, err := os.OpenFile(r.config.SeedFile, os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open users seeds file", zap.String("file", r.config.SeedFile), zap.Error(err))
	}
	defer file.Close()

	// Create json decoder
	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
	}

	// Decode
	var users []User
	err = decoder.Decode(&users)
	if err != nil {
		r.logger.Error("Failed to decode seeded users from file", zap.Error(err))
	}

	for i := range users {
		queryBuilder = queryBuilder.Values(
			users[i].Id,
			users[i].Username,
			users[i].DisplayName,
			users[i].Email,
			users[i].Country,
			users[i].State,
			users[i].AvatarURL,
			users[i].LangTag,
			users[i].Timezone,
			users[i].CreatedAt,
			users[i].UpdatedAt,
			users[i].VerifiedAt,
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

	r.logger.Info("Successfully seeded users from file", zap.Int("count", len(users)))
}

// General operations
func (r *PostgresRepository) Add(ctx context.Context, req CreateRequest) (User, error) {
	rawQuery := `INSERT INTO 
					users(id, username, displayname, email, country, state, avatar_url, lang_tag, timezone, created_at, updated_at, verified_at)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`

	id, _ := uuid.NewV7()
	now := time.Now()

	_, err := r.pgx.Exec(ctx, rawQuery,
		id,
		req.Username,
		req.DisplayName,
		req.Email,
		req.Country,
		req.State,
		req.AvatarURL,
		req.LangTag,
		req.Timezone,
		now,
		now,
		now)
	if err != nil {
		return User{}, fmt.Errorf("adding user: %w", err)
	}

	return User{
		Id:          id,
		Username:    req.Username,
		DisplayName: req.DisplayName,
		Email:       req.Email,
		Country:     req.Country,
		State:       req.State,
		AvatarURL:   req.AvatarURL,
		LangTag:     req.LangTag,
		Timezone:    req.Timezone,
		CreatedAt:   now,
		UpdatedAt:   now,
		VerifiedAt:  now,
	}, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context, cursor uuid.UUID, limit int, sortBy, sortDirection string) ([]User, error) {
	// Create and execute query
	baseQuery := `SELECT 
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
				FROM users
				WHERE id > $1`

	allowedColumns := map[string]bool{
		"id":         true,
		"username":   true,
		"created_at": true,
	}

	if !allowedColumns[strings.ToLower(sortBy)] {
		return nil, fmt.Errorf("invalid sort column: %s", sortBy)
	}
	sortDirection = strings.ToUpper(sortDirection)
	if sortDirection != "ASC" && sortDirection != "DESC" {
		return nil, fmt.Errorf("invalid sort direction: %s", sortDirection)
	}

	addition := fmt.Sprintf(" ORDER BY %s %s LIMIT %d", sortBy, sortDirection, limit)
	finalQuery := baseQuery + addition

	rows, err := r.pgx.Query(ctx, finalQuery, cursor)
	if err != nil {
		return nil, fmt.Errorf("getting users: %w", err)
	}
	defer rows.Close()

	// Scan returned rows
	var users []User

	for rows.Next() {
		var user User
		err = rows.Scan(
			&user.Id,
			&user.Username,
			&user.DisplayName,
			&user.Email,
			&user.Country,
			&user.State,
			&user.AvatarURL,
			&user.LangTag,
			&user.Timezone,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.VerifiedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
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
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.Country,
		&user.State,
		&user.AvatarURL,
		&user.LangTag,
		&user.Timezone,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.VerifiedAt,
		&disabledAt,
	)
	if err != nil {
		return User{}, fmt.Errorf("scanning row: %w", err)
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
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.Country,
		&user.State,
		&user.AvatarURL,
		&user.LangTag,
		&user.Timezone,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.VerifiedAt,
		&disabledAt,
	)
	if err != nil {
		return User{}, fmt.Errorf("scanning row: %w", err)
	}

	if disabledAt != nil && !disabledAt.IsZero() {
		user.DisabledAt = *disabledAt
	}

	return user, nil
}

func (r *PostgresRepository) Update(ctx context.Context, req UpdateRequest) (User, error) {
	// build dynamic query using squirrel
	queryBuilder := r.statementBuilder.Update("users").Where(squirrel.Eq{"id": req.Id})
	if req.Username != nil {
		queryBuilder = queryBuilder.Set("username", req.Username)
	}
	if req.DisplayName != nil {
		queryBuilder = queryBuilder.Set("displayname", req.DisplayName)
	}
	if req.Email != nil {
		queryBuilder = queryBuilder.Set("email", req.Email)
	}
	if req.Country != nil {
		queryBuilder = queryBuilder.Set("country", req.Country)
	}
	if req.State != nil {
		queryBuilder = queryBuilder.Set("state", req.State)
	}
	if req.LangTag != nil {
		queryBuilder = queryBuilder.Set("langtag", req.LangTag)
	}
	if req.Timezone != nil {
		queryBuilder = queryBuilder.Set("timezone", req.Timezone)
	}
	queryBuilder = queryBuilder.Suffix("RETURNING *")
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return User{}, fmt.Errorf("building query: %w", err)
	}

	var user User
	var disabledAt *time.Time

	// Execute query
	result := r.pgx.QueryRow(ctx, query, args...)
	err = result.Scan(
		&user.Id,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.Country,
		&user.State,
		&user.AvatarURL,
		&user.LangTag,
		&user.Timezone,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.VerifiedAt,
		&disabledAt,
	)
	if err != nil {
		return User{}, fmt.Errorf("scanning row: %w", err)
	}

	if disabledAt != nil && !disabledAt.IsZero() {
		user.DisabledAt = *disabledAt
	}

	return user, nil
}

func (r *PostgresRepository) Replace(ctx context.Context, req ReplaceRequest) (User, error) {
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
		return User{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// Lock row inside transaction
	_, err = tx.Exec(ctx, "SELECT username FROM users WHERE id = $1 FOR UPDATE;", req.Id)
	if err != nil {
		return User{}, fmt.Errorf("locking row: %w", err)
	}

	// Execute query
	result := tx.QueryRow(
		ctx,
		rawQuery,
		req.Id,
		req.Username,
		req.DisplayName,
		req.Email,
		req.Country,
		req.State,
		req.LangTag,
		req.Timezone,
		now,
		now)
	err = result.Scan(
		&user.Id,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.Country,
		&user.State,
		&user.AvatarURL,
		&user.LangTag,
		&user.Timezone,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.VerifiedAt,
		&disabledAt,
	)
	if err != nil {
		return User{}, fmt.Errorf("scanning row: %w", err)
	}

	if disabledAt != nil && !disabledAt.IsZero() {
		user.DisabledAt = *disabledAt
	}

	// Commit transaction
	err = tx.Commit(ctx)
	if err != nil {
		return User{}, fmt.Errorf("commit transaction: %w", err)
	}

	return user, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error) {
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
	if req.DeleteOptions.Type == "soft" {
		result = r.pgx.QueryRow(ctx, disableQuery, req.Id, now)
	} else {
		result = r.pgx.QueryRow(ctx, deleteQuery, req.Id)
	}

	// Scan returned row
	if err := result.Scan(&deletedId); err != nil {
		return uuid.Nil, fmt.Errorf("scanning row: %w", err)
	}

	return deletedId, nil
}

// Operations
func (r *PostgresRepository) FindOrCreateUserFromSocialProfile(ctx context.Context, profile auth.UnifiedUserProfile) (uuid.UUID, bool, error) {
	// Check if user with social provider id exists
	exists, err := r.CheckIfUserWithSocialProviderIdExists(ctx, profile.UserId, profile.Provider)
	if err != nil {
		return uuid.Nil, false, fmt.Errorf("checking user with social provider: %w", err)
	}

	var user User

	if !exists {
		// 1. Create a user first in users table
		hash := sha256.Sum256([]byte(profile.UserId))
		username := fmt.Sprintf("%x", hash)
		newUser := CreateRequest{
			Username:    username,
			DisplayName: profile.DisplayName,
			Email:       profile.Email,
			Country:     "in",
			State:       "wb",
			LangTag:     "en-IN",
			Timezone:    "India/Kolkata",
			AvatarURL:   profile.AvatarUrl,
		}

		user, err := r.Add(ctx, newUser)
		if err != nil {
			return uuid.Nil, false, fmt.Errorf("creating new user: %w", err)
		}

		// 2. Add mapping of new user with the social provider id
		err = r.MapUserWithSocialProvider(ctx, user.Id, profile.UserId, profile.Provider)
		if err != nil {
			return uuid.Nil, false, fmt.Errorf("mapping user: %w", err)
		}

		// Created
		return user.Id, true, nil
	}

	// Get existing user
	user, err = r.GetBySocialProfileId(ctx, profile.Provider, profile.UserId)
	if err != nil {
		return uuid.Nil, false, fmt.Errorf("getting user by social profile: %w", err)
	}

	// Not crerated
	return user.Id, false, nil
}

func (r *PostgresRepository) CheckIfUserWithSocialProviderIdExists(ctx context.Context, socialUserid, socialProvider string) (bool, error) {
	query := `SELECT EXISTS (SELECT provider_id FROM user_mappings WHERE provider_id = $1);`

	var exists bool
	row := r.pgx.QueryRow(ctx, query, socialUserid)
	if err := row.Scan(&exists); err != nil {
		return false, fmt.Errorf("scanning row: %w", err)
	}

	return exists, nil
}

func (r *PostgresRepository) MapUserWithSocialProvider(ctx context.Context, userId uuid.UUID, providerUserId, provider string) error {
	mappingQuery := `INSERT INTO 
						user_mappings(
							user_id,
							provider_id, 
							provider
						) 
						VALUES(
							$1, 
							$2, 
							$3
						);`

	_, err := r.pgx.Exec(ctx, mappingQuery, userId, providerUserId, provider)
	if err != nil {
		return fmt.Errorf("mapping user: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetBySocialProfileId(ctx context.Context, provider, socialUserId string) (User, error) {
	query := `SELECT
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
				users u
			INNER JOIN
				user_mappings m
			ON
				u.id = m.user_id
			WHERE
				m.provider_id = $1
			AND
				m.provider = $2;`

	var user User
	var disabledAt *time.Time

	row := r.pgx.QueryRow(ctx, query, socialUserId, provider)
	err := row.Scan(
		&user.Id,
		&user.Username,
		&user.DisplayName,
		&user.Email,
		&user.Country,
		&user.State,
		&user.AvatarURL,
		&user.LangTag,
		&user.Timezone,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.VerifiedAt,
		&disabledAt,
	)
	if err != nil {
		return User{}, fmt.Errorf("scanning row: %w", err)
	}

	if disabledAt != nil && !disabledAt.IsZero() {
		user.DisabledAt = *disabledAt
	}

	return user, nil
}

func (r *PostgresRepository) CheckIfUserIsBanned(ctx context.Context, userId uuid.UUID) (bool, error) {
	rawQuery := `
		SELECT EXISTS(
			SELECT
				id
			FROM
				users_banned
			WHERE
				id = $1
		);
	`

	var banned bool

	row := r.pgx.QueryRow(
		ctx,
		rawQuery,
		userId,
	)
	err := row.Scan(&banned)
	if err != nil {
		return false, fmt.Errorf("checking for banned user: %w", err)
	}

	return banned, nil
}

// Bulk operations - subject to change for improvement (dont use these queries right now)
func (r *PostgresRepository) BulkAdd(ctx context.Context, req BulkCreateRequest) ([]common.BulkOpResult, error) {
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

	for i := range req.Users {
		id, _ := uuid.NewV7()
		now := time.Now()
		dummyAvatarUrl := "-"
		queryBuilder = queryBuilder.Values(
			id,
			req.Users[i].Username,
			req.Users[i].DisplayName,
			req.Users[i].Email,
			req.Users[i].Country,
			req.Users[i].State,
			dummyAvatarUrl,
			req.Users[i].LangTag,
			req.Users[i].Timezone,
			now,
			now,
			now,
		)

		var result common.BulkOpResult
		result.Id = id
		result.Status = 1
		result.Success = true
		result.Message = "added"
		results = append(results, result)
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("generate query: %w", err)
	}

	_, err = r.pgx.Exec(ctx, query, args...)
	if err != nil {
		for i := range results {
			results[i].Status = 0
			results[i].Success = false
			results[i].Message = "failed"
		}
		return results, fmt.Errorf("inserting user: %w", err)
	}

	return results, nil
}

func (r *PostgresRepository) BulkModify(ctx context.Context, req BulkModifyRequest) ([]common.BulkOpResult, error) {
	queryBuilder := r.statementBuilder.Update("users")

	var results []common.BulkOpResult

	// To be decided

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("generate query: %w", err)
	}

	_, err = r.pgx.Exec(ctx, query, args...)
	if err != nil {
		for i := range results {
			results[i].Status = 0
			results[i].Success = false
			results[i].Message = "failed"
		}
		return results, fmt.Errorf("updating users: %w", err)
	}

	return nil, nil
}

func (r *PostgresRepository) BulkDelete(ctx context.Context, req BulkDeleteRequest) ([]common.BulkOpResult, error) {
	queryBuilder := r.statementBuilder.Delete("users")

	var results []common.BulkOpResult

	for i := range req.Users {
		currentUser := req.Users[i]
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
		return nil, fmt.Errorf("generate query: %w", err)
	}

	_, err = r.pgx.Exec(ctx, query, args...)
	if err != nil {
		for i := range results {
			results[i].Status = 0
			results[i].Message = "failed"
			results[i].Success = false
		}
		return results, fmt.Errorf("deleting user: %w", err)
	}

	return results, nil
}

func (r *PostgresRepository) Ban(ctx context.Context, req BanRequest) (BannedUserResponse, error) {
	rawQuery := `
		INSERT INTO users_banned(
			id,
			reason,
			ban_type,
			banned_at,
			banned_till,
			banned_by
		) VALUES (
			$1,
			$2,
			$3,
			$4,
			$5,
			$6 
		);
	`

	userID, _ := uuid.Parse(req.UserID)
	currentTimeUTC := time.Now().UTC()

	_, err := r.pgx.Exec(
		ctx,
		rawQuery,
		userID,
		req.Reason,
		req.Type,
		currentTimeUTC,
		req.Duration.UTC(),
		SystemUser,
	)
	if err != nil {
		return BannedUserResponse{}, fmt.Errorf("banning user: %w", err)
	}

	return BannedUserResponse{
		UserID:     req.UserID,
		Reason:     req.Reason,
		Type:       req.Type,
		BannedBy:   "system",
		BannedTill: req.Duration,
	}, nil
}

func (r *PostgresRepository) UnBan(ctx context.Context, req UnbanRequest) error {
	rawQuery := `
		DELETE FROM
			users_banned
		WHERE
			id = $1;
	`

	userID, _ := uuid.Parse(req.UserID)

	_, err := r.pgx.Query(
		ctx,
		rawQuery,
		userID,
	)
	if err != nil {
		return fmt.Errorf("unbanning user: %w", err)
	}

	return nil
}

func (r *PostgresRepository) CheckUsername(ctx context.Context, req CheckUsernameRequest) (bool, error) {
	rawQuery := `
		SELECT EXISTS(
			SELECT 
				1 
			FROM 
				users 
			WHERE 
				username = $1
		);
	`

	row := r.pgx.QueryRow(ctx, rawQuery, req.Username)

	var valid bool

	err := row.Scan(&valid)
	if err != nil {
		return false, fmt.Errorf("scanning row: %w", err)
	}

	return valid, nil
}

func (r *PostgresRepository) GetAllUsernames(ctx context.Context, req UserId) ([]string, error) {
	rawQuery := `SELECT 
					display_name
				FROM 
					user_name_audit
				WHERE
					id = $1;`

	userID, _ := uuid.Parse(req.Id)
	rows, err := r.pgx.Query(ctx, rawQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("getting usernames: %w", err)
	}
	defer rows.Close()

	var usernames []string

	for rows.Next() {
		var username string

		err = rows.Scan(&username)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		usernames = append(usernames, username)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return usernames, nil
}
