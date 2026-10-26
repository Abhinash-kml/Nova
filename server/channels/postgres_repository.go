package channels

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
	config           *config.ChannelsConfig
	logger           *zap.Logger
	pgx              *pgxpool.Pool
	statementBuilder squirrel.StatementBuilderType
}

func NewPostgresRepositoryFromPgxPool(connection *pgxpool.Pool, c *config.ChannelsConfig, l *zap.Logger) *PostgresRepository {
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
	existQuery := `SELECT count(*) FROM channels;`
	row := r.pgx.QueryRow(ctx, existQuery)

	var count int

	// Scan returned row
	if err := row.Scan(&count); err != nil {
		r.logger.Error("Failed to scan count of rows from channels table", zap.Error(err))
	}

	// Rows already exists, dont seed
	if count != 0 {
		r.logger.Sugar().Infof("Channels table contains %d rows, skipping seeding from file...", count)
		return
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
	if len(r.config.SeedFile) == 0 || !strings.HasSuffix(r.config.SeedFile, ".json") {
		r.logger.Error("Seeds file not provided or is not of json format. Skipping seeding")
		return
	}

	// open file
	file, err := os.OpenFile(r.config.SeedFile, os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open channels seeds file", zap.String("file", r.config.SeedFile), zap.Error(err))
		return
	}
	defer file.Close()

	// Create json decoder
	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
		return
	}

	// Decode
	var channels []Channel
	err = decoder.Decode(&channels)
	if err != nil {
		r.logger.Error("Failed to decode seeded channels from file", zap.Error(err))
		return
	}

	for i := range channels {
		queryBuilder = queryBuilder.Values(
			channels[i].Id,
			channels[i].Name,
			channels[i].IsPersistant,
			channels[i].ProcessInterval,
			channels[i].CreatedBy,
			channels[i].CreatedAt,
			channels[i].Updatedby,
			channels[i].UpdatedAt,
		)
	}

	// Generate query
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate sql query with squirrel for seeding channels table", zap.Error(err))
		return
	}

	// Execute query
	_, err = r.pgx.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to execute query to seed channels table", zap.Error(err))
		return
	}

	r.logger.Info("Successfully seeded channels from file", zap.Int("count", len(channels)))
}

// General operations
func (r *PostgresRepository) GetAll(ctx context.Context, cursor uuid.UUID, limit int) ([]Channel, error) {
	rawQuery := `SELECT
					id,
					name,
					process_interval,
					persistant,
					created_by,
					created_at
				FROM
					channels
				WHERE
					id >= $1
				ORDER BY
					name
				LIMIT
					$2;`

	rows, err := r.pgx.Query(ctx, rawQuery, cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("getting channels: %w", err)
	}
	defer rows.Close()

	var channels []Channel
	var channel Channel

	for rows.Next() {
		err = rows.Scan(
			&channel.Id,
			&channel.Name,
			&channel.ProcessInterval,
			&channel.IsPersistant,
			&channel.CreatedBy,
			&channel.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		channels = append(channels, channel)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
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
	err := row.Scan(
		&channel.Id,
		&channel.Name,
		&channel.ProcessInterval,
		&channel.IsPersistant,
		&channel.CreatedBy,
		&channel.CreatedAt,
		&channel.Updatedby,
		&channel.UpdatedAt,
	)
	if err != nil {
		return Channel{}, fmt.Errorf("scanning row: %w", err)
	}

	return channel, nil
}

func (r *PostgresRepository) Add(ctx context.Context, req CreateRequest) (Channel, error) {
	rawQuery := `INSERT INTO
					channels (
						id,
						name,
						process_interval,
						persistant, 
						created_by, 
						created_at, 
						updated_by, 
						updated_at
					) VALUES (
					 	$1, 
						$2, 
						$3, 
						$4, 
						$5, 
						$6, 
						$7, 
						$8
					);`

	id, err := uuid.NewV7()
	if err != nil {
		return Channel{}, fmt.Errorf("generating id: %w", err)
	}
	currentTime := time.Now()

	// Execute query
	_, err = r.pgx.Exec(
		ctx,
		rawQuery,
		id,
		req.Name,
		req.ProcessInterval,
		req.IsPersistant,
		req.CreatedBy,
		currentTime,
		uuid.Nil,
		currentTime,
	)
	if err != nil {
		return Channel{}, fmt.Errorf("inserting channel: %w", err)
	}

	processInterval, _ := time.ParseDuration(req.ProcessInterval)
	createdBy, _ := uuid.Parse(req.CreatedBy)

	return Channel{
		Id:              id,
		Name:            req.Name,
		ProcessInterval: processInterval,
		IsPersistant:    req.IsPersistant,
		CreatedBy:       createdBy,
		CreatedAt:       currentTime,
		Updatedby:       uuid.Nil,
		UpdatedAt:       time.Time{},
	}, nil
}

func (r *PostgresRepository) Modify(ctx context.Context, req UpdateRequest) (Channel, error) {
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
	currentTime := time.Now()

	// Execute query
	row := r.pgx.QueryRow(
		ctx,
		rawQuery,
		req.Id,
		req.Name,
		req.ProcessInterval,
		req.IsPersistant,
		req.UpdatedBy,
		currentTime,
	)

	// Scan returned row
	err := row.Scan(
		&channel.Id,
		&channel.Name,
		&channel.ProcessInterval,
		&channel.IsPersistant,
		&channel.CreatedBy,
		&channel.CreatedAt,
		&channel.Updatedby,
		&channel.UpdatedAt,
	)
	if err != nil {
		return Channel{}, fmt.Errorf("scanning row: %w", err)
	}

	return channel, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, req DeleteRequest) (uuid.UUID, error) {
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
	if req.Type == "disable" {
		row = r.pgx.QueryRow(ctx, disableQuery, req.Id)
	} else {
		row = r.pgx.QueryRow(ctx, deleteQuery, req.Id)
	}

	// Scan returned row
	err := row.Scan(&effectedChannel)
	if err != nil {
		return uuid.Nil, fmt.Errorf("scanning row: %w", err)
	}

	return effectedChannel, nil
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
