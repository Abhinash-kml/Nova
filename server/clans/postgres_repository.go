package clans

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
		queryBuilder = queryBuilder.Values(clans[i].Id, clans[i].Name, clans[i].Tag, clans[i].Description,
			clans[i].LeaderId, clans[i].ColeaderId, clans[i].Level, clans[i].MaxMembers, clans[i].IsLocked,
			clans[i].CreatedAt, clans[i].UpdatedAt)
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
	rawQuery := `INSERT INTO
					clans(id, name, tag, description, leader_id, coleader_ids, level, member_ids, max_members, islocked, created_at, updated_at)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12);`

	id, _ := uuid.NewV7()
	now := time.Now()
	_, err := r.pgx.Exec(ctx, rawQuery, id, dto.Name, dto.Tag, dto.Description, dto.LeaderId, dto.ColeaderId, dto.Level,
		dto.Members, dto.MaxMembers, dto.IsLocked, now, now)
	if err != nil {
		r.logger.Error("Failed to execute sql insert query", zap.Error(err))
		return Clan{}, common.TranslatePostgresError(err, r.logger)
	}

	return Clan{
		Id:          id,
		Name:        dto.Name,
		Tag:         dto.Tag,
		Description: dto.Description,
		LeaderId:    dto.LeaderId,
		ColeaderId:  dto.ColeaderId,
		Level:       dto.Level,
		Members:     dto.Members,
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
					coleader_ids,
					level,
					member_ids,
					max_members,
					islocked,
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
	err := row.Scan(&clan.Id, &clan.Name, &clan.Tag, &clan.Description, &clan.LeaderId, &clan.ColeaderId,
		&clan.Level, &clan.Members, &clan.MaxMembers, &clan.IsLocked, &clan.CreatedAt, &clan.UpdatedAt)
	if err != nil {
		r.logger.Error("Failed to scan returned row from getbyid query", zap.Error(err))
		return Clan{}, common.TranslatePostgresError(err, r.logger)
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
					coleader_ids,
					level,
					member_ids,
					max_members,
					islocked,
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
	err := row.Scan(&clan.Id, &clan.Name, &clan.Tag, &clan.Description, &clan.LeaderId, &clan.ColeaderId,
		&clan.Level, &clan.Members, &clan.MaxMembers, &clan.IsLocked, &clan.CreatedAt, &clan.UpdatedAt)
	if err != nil {
		r.logger.Error("Failed to scan returned row from getbyname query", zap.Error(err))
		return Clan{}, common.TranslatePostgresError(err, r.logger)
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
					coleader_ids,
					level,
					member_ids,
					max_members,
					islocked,
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
	defer rows.Close()

	// Scan returned rows
	var clans []Clan
	for rows.Next() {
		var clan Clan
		err = rows.Scan(&clan.Id, &clan.Name, &clan.Tag, &clan.Description, &clan.LeaderId, &clan.ColeaderId,
			&clan.Level, &clan.Members, &clan.MaxMembers, &clan.IsLocked, &clan.CreatedAt, &clan.UpdatedAt)
		if err != nil {
			r.logger.Error("Failed to scan returned row in getall query", zap.Error(err))
			return nil, common.TranslatePostgresError(err, r.logger)
		}

		clans = append(clans, clan)
	}

	return clans, err
}

func (r *PostgresRepository) Update(ctx context.Context, dto UpdateDTO) (Clan, error) {

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
