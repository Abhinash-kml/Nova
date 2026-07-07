package channels

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
	existQuery := `SELECT count(*) FROM channels;`
	row := r.pgx.QueryRow(ctx, existQuery)

	var count int

	// Scan returned row
	if err := row.Scan(&count); err != nil {
		r.logger.Error("Failed to scan count of rows from channels table", zap.Error(err))
		return common.ErrDbSeedingFailed
	}

	// Rows already exists, dont seed
	if count != 0 {
		r.logger.Sugar().Infof("Channels table contains %d rows, skipping seeding from file...", count)
		return nil
	}

	// Table is empty, so seed from file to table
	queryBuilder := r.statementBuilder.Insert("channels").Columns(
		"id",
		"name",
		"persistant",
		"process_interval",
		"created_by",
		"created_at",
		"updated_by",
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
		r.logger.Error("Failed to open channels seeds file", zap.String("file", r.seedfile), zap.Error(err))
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
	var channels []ChannelDTO
	err = decoder.Decode(&channels)
	if err != nil {
		r.logger.Error("Failed to decode seeded channels from file", zap.Error(err))
		return common.ErrDbSeedingFailed
	}

	for i := range channels {
		queryBuilder = queryBuilder.Values(channels[i].Id, channels[i].Name, channels[i].IsPersistant, channels[i].CreatedBy, channels[i].CreatedAt)
	}

	// Generate query
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate sql query with squirrel for seeding channels table", zap.Error(err))
		return common.ErrDbSeedingFailed
	}

	// Execute query
	_, err = r.pgx.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to execute query to seed channels table", zap.Error(err))
		return common.ErrDbSeedingFailed
	}

	r.logger.Info("Successfully seeded channels from file", zap.Int("count", len(channels)))

	return nil
}

// General operations
func (r *PostgresRepository) GetAll(ctx context.Context, cursor int, limit int) ([]Channel, error) {
	var rows pgx.Rows
	var err error

	// Generate query based on cursor
	rawQuery := `SELECT
					id,
					name,
					process_interval,
					persistant,
					created_by,
					created_at
				FROM
					channels`
	if cursor == 0 {
		rawQuery += ` ORDER BY
						id
					LIMIT 
						$1;`

		rows, err = r.pgx.Query(ctx, rawQuery, limit)
	} else {
		rawQuery += ` WHERE 
						id > $1
					ORDER BY
						id
					LIMIT
						$2;`

		rows, err = r.pgx.Query(ctx, rawQuery, cursor, limit)
	}
	defer rows.Close()

	// Scan returned rows
	var channels []Channel
	for rows.Next() {
		var channel Channel
		err = rows.Scan(&channel.Id, &channel.Name, &channel.ProcessInterval, &channel.IsPersistant,
			&channel.CreatedBy, &channel.CreatedAt)
		if err != nil {
			r.logger.Error("Failed to scan returned row in getall query", zap.Error(err))
			return nil, common.TranslatePostgresError(err, r.logger)
		}

		channels = append(channels, channel)
	}

	return channels, nil
}

func (r *PostgresRepository) GetById(ctx context.Context, id uuid.UUID) (Channel, error) {
	rawQuery := `SELECT
					id,
					name,
					process_interval,
					persistant,
					created_by,
					created_at,
					updated_by,
					updated_at
				FROM
					channels
				WHERE
					id = $1;`

	var channel Channel

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery, id)

	// Scan returned row
	err := row.Scan(&channel.Id, &channel.Name, &channel.ProcessInterval, &channel.IsPersistant,
		&channel.CreatedBy, &channel.CreatedAt, &channel.Updatedby, &channel.UpdatedAt)
	if err != nil {
		r.logger.Error("Failed to scan returned row in getbyid query", zap.Error(err))
		return Channel{}, common.TranslatePostgresError(err, r.logger)
	}

	return channel, nil
}

func (r *PostgresRepository) Add(ctx context.Context, dto CreateDTO) (Channel, error) {
	rawQuery := `INSERT INTO
					channels(id, name, process_interval, persistant, created_by, created_at, updated_by, updated_at)
				VALUES($1, $2, $3, $4, $5, $6, $7, $8);`

	id, _ := uuid.NewV7()
	now := time.Now()

	// Execute query
	_, err := r.pgx.Exec(ctx, rawQuery, id, dto.Name, dto.ProcessInterval, dto.IsPersistant,
		dto.CreatedBy, now, uuid.Nil, now)
	if err != nil {
		r.logger.Error("Failed to execute sql insert query", zap.Error(err))
		return Channel{}, common.TranslatePostgresError(err, r.logger)
	}

	processInterval, _ := time.ParseDuration(dto.ProcessInterval)
	createdBy, _ := uuid.Parse(dto.CreatedBy)

	return Channel{
		Id:              id,
		Name:            dto.Name,
		ProcessInterval: processInterval,
		IsPersistant:    dto.IsPersistant,
		CreatedBy:       createdBy,
		CreatedAt:       now,
		Updatedby:       uuid.Nil,
		UpdatedAt:       time.Time{},
	}, nil
}

func (r *PostgresRepository) Modify(ctx context.Context, dto UpdateDTO) (Channel, error) {
	rawQuery := `UPDATE
					channels
				SET
					name = $2,
					process_interval = $3,
					persistant = $4,
					updated_by = $5,
					updated_at = $6
				WHERE
					id = $1
				RETURNING
					*;`

	var channel Channel
	now := time.Now()

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery, dto.Id, dto.Name, dto.ProcessInterval, dto.IsPersistant,
		dto.UpdatedBy, now)

	// Scan returned row
	err := row.Scan(&channel.Id, &channel.Name, &channel.ProcessInterval, &channel.IsPersistant,
		&channel.CreatedBy, &channel.CreatedAt, &channel.Updatedby, &channel.UpdatedAt)
	if err != nil {
		r.logger.Error("Failed to scan returned row in delete query", zap.Error(err))
		return Channel{}, common.TranslatePostgresError(err, r.logger)
	}

	return channel, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, dto DeleteDTO) (uuid.UUID, error) {
	deleteQuery := `DELETE FROM
						channels
					WHERE
						id = $1
					RETURNING
						id;`

	disableQuery := `UPDATE
						channels
					SET
						enabled = true
					WHERE
						id = $1
					RETURNING
						id;`

	var row pgx.Row
	var effectedChannel uuid.UUID

	// Execute query based on delete / disable
	if dto.Type == "disable" {
		row = r.pgx.QueryRow(ctx, disableQuery, dto.Id)
	} else {
		row = r.pgx.QueryRow(ctx, deleteQuery, dto.Id)
	}

	// Scan returned row
	err := row.Scan(&effectedChannel)
	if err != nil {
		r.logger.Error("Failed to scan returned row in delete query", zap.Error(err))
		return uuid.Nil, common.TranslatePostgresError(err, r.logger)
	}

	return effectedChannel, nil
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
