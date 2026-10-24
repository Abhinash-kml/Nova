package clans

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
	config           *config.Config
	logger           *zap.Logger
	pgx              *pgxpool.Pool
	statementBuilder squirrel.StatementBuilderType
	seedfile         string
}

func NewPostgresRepositoryFromPgxPool(connection *pgxpool.Pool, c *config.Config, l *zap.Logger, sfp string) *PostgresRepository {
	return &PostgresRepository{
		config:           c,
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
	existQuery := `SELECT count(*) FROM clans;`
	row := r.pgx.QueryRow(ctx, existQuery)

	var count int

	// Scan returned row
	if err := row.Scan(&count); err != nil {
		r.logger.Error("Failed to scan count of rows from clans table", zap.Error(err))
		return common.ErrDbSeedingFailed
	}

	// Rows already exists, dont seed
	if count != 0 {
		r.logger.Sugar().Infof("Clans table contains %d rows, skipping seeding from file...", count)
		return nil
	}

	// Table is empty, so seed from file to table
	queryBuilder := r.statementBuilder.Insert("clans").Columns(
		"id",
		"name",
		"tag",
		"description",
		"leader_id",
		"coleader_ids",
		"level",
		"member_ids",
		"max_members",
		"islocked",
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
		r.logger.Error("Failed to open clans seeds file", zap.String("file", r.seedfile), zap.Error(err))
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
	var clans []Clan
	err = decoder.Decode(&clans)
	if err != nil {
		r.logger.Error("Failed to decode seeded clans from file", zap.Error(err))
		return common.ErrDbSeedingFailed
	}

	for i := range clans {
		queryBuilder = queryBuilder.Values(clans[i].Id,
			clans[i].Name,
			clans[i].Tag,
			clans[i].Description,
			clans[i].LeaderId,
			clans[i].ColeaderId,
			clans[i].Level,
			clans[i].Members,
			clans[i].MaxMembers,
			clans[i].IsLocked,
			clans[i].CreatedAt,
			clans[i].UpdatedAt)
	}

	// Generate query
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate sql query with squirrel for seeding clans table", zap.Error(err))
		return common.ErrDbSeedingFailed
	}

	// Execute query
	_, err = r.pgx.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to execute query to seed clans table", zap.Error(err))
		return common.ErrDbSeedingFailed
	}

	r.logger.Info("Successfully seeded clans from file", zap.Int("count", len(clans)))

	return nil
}

// General operations
func (r *PostgresRepository) Add(ctx context.Context, dto CreateDTO) (Clan, error) {
	createClanQuery := `INSERT INTO
					clans (
						id, 
						name, 
						tag, 
						description, 
						level, 
						max_members, 
						is_locked, 
						created_at, 
						updated_at
					)
					VALUES (
						$1, 
						$2, 
						$3, 
						$4, 
						$5, 
						$6, 
						$7, 
						$8, 
						$9
					);`

	makeLeaderQuery := `
		INSERT INTO
			clan_members (
				clan_id,
				user_id,
				role,
				joined_at
			) VALUES (
				$1,
				$2,
				'LEADER',
				$3 
			);
	`

	clanID, err := uuid.NewV7()
	if err != nil {
		return Clan{}, fmt.Errorf("generating id: %w", err)
	}
	now := time.Now()

	tx, err := r.pgx.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return Clan{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		createClanQuery,
		clanID,
		dto.Name,
		dto.Tag,
		dto.Description,
		dto.Level,
		dto.MaxMembers,
		dto.IsLocked,
		now,
		now,
	)
	if err != nil {
		return Clan{}, fmt.Errorf("creating clan: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		makeLeaderQuery,
		clanID,
		dto.LeaderID,
		now,
	)
	if err != nil {
		return Clan{}, fmt.Errorf("making leader: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return Clan{}, fmt.Errorf("commiting transaction: %w", err)
	}

	return Clan{
		Id:          clanID,
		Name:        dto.Name,
		Tag:         dto.Tag,
		Description: dto.Description,
		LeaderID:    dto.LeaderID,
		Level:       dto.Level,
		MaxMembers:  dto.MaxMembers,
		IsLocked:    dto.IsLocked,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (r *PostgresRepository) GetById(ctx context.Context, id uuid.UUID) (Clan, error) {
	rawQuery := `SELECT
					id,
					name,
					tag,
					description,
					leader_id,
					level,
					max_members,
					is_locked,
					created_at,
					updated_at
				FROM
					clans
				WHERE
					id = $1;`

	var clan Clan

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery, id)

	// Scan returned row
	err := row.Scan(
		&clan.Id,
		&clan.Name,
		&clan.Tag,
		&clan.Description,
		&clan.LeaderID,
		&clan.Level,
		&clan.Members,
		&clan.IsLocked,
		&clan.CreatedAt,
		&clan.UpdatedAt,
	)
	if err != nil {
		return Clan{}, fmt.Errorf("scanning row: %w", err)
	}

	return clan, nil
}

func (r *PostgresRepository) GetByName(ctx context.Context, name string) (Clan, error) {
	rawQuery := `SELECT
					id,
					name,
					tag,
					description,
					leader_id,
					level,
					max_members,
					is_locked,
					created_at,
					updated_at
				FROM
					clans
				WHERE
					name = $1;`

	var clan Clan

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery, name)

	// Scan returned row
	err := row.Scan(
		&clan.Id,
		&clan.Name,
		&clan.Tag,
		&clan.Description,
		&clan.LeaderID,
		&clan.Level,
		&clan.MaxMembers,
		&clan.IsLocked,
		&clan.CreatedAt,
		&clan.UpdatedAt,
	)
	if err != nil {
		return Clan{}, fmt.Errorf("scanning row: %w", err)
	}

	return clan, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context, cursor int, limit int) ([]Clan, error) {
	var rows pgx.Rows
	var err error

	// Create and execute query
	rawQuery := `SELECT 
					id,
					name,
					tag,
					description,
					leader_id,
					level,
					max_members,
					is_locked,
					created_at,
					updated_at
				FROM 
					clans`
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
	if err != nil {
		return nil, fmt.Errorf("getting clans: %w", err)
	}
	defer rows.Close()

	// Scan returned rows
	var clans []Clan
	for rows.Next() {
		var clan Clan
		err = rows.Scan(
			&clan.Id,
			&clan.Name,
			&clan.Tag,
			&clan.Description,
			&clan.LeaderID,
			&clan.Level,
			&clan.Members,
			&clan.IsLocked,
			&clan.CreatedAt,
			&clan.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		clans = append(clans, clan)
	}

	return clans, err
}

func (r *PostgresRepository) Update(ctx context.Context, dto UpdateDTO) (Clan, error) {
	return Clan{}, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error) {
	rawQuery := `DELETE FROM
					clans
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
		return uuid.Nil, fmt.Errorf("scanning row: %w", err)
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
