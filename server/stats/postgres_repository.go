package stats

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	config           *config.StatsConfig
	logger           *zap.Logger
	pgx              *pgxpool.Pool
	statementBuilder squirrel.StatementBuilderType
}

func NewPostgresRepositoryFromPgxPool(connection *pgxpool.Pool, c *config.StatsConfig, l *zap.Logger) *PostgresRepository {
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
}

func (r *PostgresRepository) Add(ctx context.Context, req CreateRequest) (Stats, error) {
	rawQuery := `INSERT INTO
					stats(key, name, description, start_value, created_at)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING 
					id,
					key, 
					name,
					description, 
					start_value, 
					created_at;`

	var stat Stats
	now := time.Now()

	err := r.pgx.QueryRow(
		ctx,
		rawQuery,
		req.Key,
		req.Name,
		req.Description,
		req.StartValue,
		now,
	).Scan(
		&stat.Id,
		&stat.Key,
		&stat.Name,
		&stat.Description,
		&stat.StartValue,
		&stat.CreatedAt,
	)

	if err != nil {
		return Stats{}, fmt.Errorf("creating stat: %w", err)
	}

	return stat, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context, cursor int, limit int) ([]Stats, error) {
	rawQuery := `SELECT
					id,
					key,
					name,
					description,
					start_value,
					created_at
				FROM 
					stats
				WHERE
					id >= $1
				ORDER BY
					name
				LIMIT 
					$2;`

	rows, err := r.pgx.Query(ctx, rawQuery, cursor, limit)
	if err != nil {
		return nil, fmt.Errorf("getting all stats: %w", err)
	}
	defer rows.Close()

	var stats []Stats

	for rows.Next() {
		var stat Stats

		err = rows.Scan(
			&stat.Id,
			&stat.Key,
			&stat.Name,
			&stat.Description,
			&stat.StartValue,
			&stat.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		stats = append(stats, stat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return stats, nil
}

func (r *PostgresRepository) GetById(ctx context.Context, id int) (Stats, error) {
	rawQuery := `SELECT
					id,
					key,
					name,
					description,
					start_value,
					created_at
				FROM 
					stats
				WHERE
					id = $1;`

	var stat Stats

	err := r.pgx.QueryRow(
		ctx,
		rawQuery,
		id,
	).Scan(
		&stat.Id,
		&stat.Key,
		&stat.Name,
		&stat.Description,
		&stat.StartValue,
		&stat.CreatedAt,
	)

	if err != nil {
		return Stats{}, fmt.Errorf("scanning row: %w", err)
	}

	return stat, nil
}

func (r *PostgresRepository) Update(ctx context.Context, req UpdateRequest) (Stats, error) {
	queryBuilder := r.statementBuilder.
		Update("stats").
		Where(squirrel.Eq{"id": req.Id})

	if req.Name != nil {
		queryBuilder = queryBuilder.Set("name", *req.Name)
	}

	if req.StartValue != nil {
		queryBuilder = queryBuilder.Set("start_value", *req.StartValue)
	}

	queryBuilder = queryBuilder.Suffix("RETURNING id, key, name, description, start_value, created_at")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return Stats{}, fmt.Errorf("generate query: %w", err)
	}

	var stat Stats

	tx, err := r.pgx.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		return Stats{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		"SELECT id FROM stats WHERE id = $1 FOR UPDATE;",
		req.Id,
	)
	if err != nil {
		return Stats{}, fmt.Errorf("locking row: %w", err)
	}

	result := tx.QueryRow(ctx, query, args...)

	err = result.Scan(
		&stat.Id,
		&stat.Key,
		&stat.Name,
		&stat.Description,
		&stat.StartValue,
		&stat.CreatedAt,
	)
	if err != nil {
		return Stats{}, fmt.Errorf("scanning row: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return Stats{}, fmt.Errorf("commit transaction: %w", err)
	}

	return stat, nil
}

func (r *PostgresRepository) Replace(ctx context.Context, req ReplaceRequest) (Stats, error) {
	rawQuery := `UPDATE
					stats
				SET
					name = $2,
					key = $3,
					description = $4,
					start_value = $5
				WHERE
					id = $1
				RETURNING
					id, 
					name, 
					key, 
					description, 
					start_value, 
					created_at;`

	var stat Stats

	tx, err := r.pgx.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		r.logger.Error("Failed to begin transaction in replace query", zap.Error(err))
		return Stats{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		"SELECT id FROM stats WHERE id = $1 FOR UPDATE;",
		req.Id,
	)
	if err != nil {
		return Stats{}, fmt.Errorf("locking row: %w", err)
	}

	result := tx.QueryRow(
		ctx,
		rawQuery,
		req.Id,
		req.Name,
		req.Key,
		req.Description,
		req.StartValue,
	)
	err = result.Scan(
		&stat.Id,
		&stat.Name,
		&stat.Key,
		&stat.Description,
		&stat.StartValue,
		&stat.CreatedAt,
	)
	if err != nil {
		return Stats{}, fmt.Errorf("scanning row: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return Stats{}, fmt.Errorf("commit transaction: %w", err)
	}

	return stat, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, req DeleteRequest) (int, error) {
	deleteQuery := `DELETE FROM
						stats
					WHERE
						id = $1
					RETURNING
						id;`

	var deletedId int

	err := r.pgx.QueryRow(
		ctx,
		deleteQuery,
		req.Id,
	).Scan(&deletedId)
	if err != nil {
		return 0, fmt.Errorf("scanning row: %w", err)
	}

	return deletedId, nil
}

func (r *PostgresRepository) GetPlayerStats(ctx context.Context, req GetPlayerStatRequest) (PlayerStatsResponse, error) {
	query := `SELECT
				stat_id,
				current_value
			FROM 
				player_stats_progression
			WHERE 
				player_id = $1;`

	rows, err := r.pgx.Query(
		ctx,
		query,
		req.Id,
	)
	if err != nil {
		return PlayerStatsResponse{}, fmt.Errorf("getting player stats: %w", err)
	}
	defer rows.Close()

	var response PlayerStatsResponse
	for rows.Next() {
		var stat StatResponse
		if err := rows.Scan(&stat.Id, &stat.Value); err != nil {
			return PlayerStatsResponse{}, fmt.Errorf("scanning row: %w", err)
		}

		response.Stats = append(response.Stats, stat)
	}

	if err := rows.Err(); err != nil {
		return PlayerStatsResponse{}, fmt.Errorf("iterating rows: %w", err)
	}

	return response, nil
}

func (r *PostgresRepository) UpdatePlayerStats(ctx context.Context, req UpdatePlayerStatRequest) error {
	query := `INSERT INTO
				player_stats_progression (
					player_id,
					stat_id,
					current_value,
					updated_at
				)
				VALUES (
					$1,
					$2,
					$3,
					CURRENT_TIMESTAMP
				)
				ON CONFLICT (player_id, stat_id)
				DO UPDATE SET
					current_value = EXCLUDED.current_value,
					updated_at = CURRENT_TIMESTAMP`

	_, err := r.pgx.Exec(
		ctx,
		query,
		req.Id,
		req.StatId,
		req.Value,
	)
	if err != nil {
		return fmt.Errorf("updating stats: %w", err)
	}

	return nil
}

func (r *PostgresRepository) DeletePlayerStats(ctx context.Context, req DeletePlayerStatRequest) error {
	query := `DELETE FROM 
				player_stats_progression
			WHERE 
				player_id = $1`

	_, err := r.pgx.Exec(
		ctx,
		query,
		req.Id,
	)
	if err != nil {
		return fmt.Errorf("delete stats: %w", err)
	}

	return nil
}

func (r *PostgresRepository) DeletePlayerStatSpecific(ctx context.Context, req DeletePlayerStatSpecificRequest) error {
	query := `DELETE FROM 
				player_stats_progression
			WHERE 
				player_id = $1
			AND 
				stat_id = $2`

	statId, _ := uuid.Parse(req.StatsId)

	_, err := r.pgx.Exec(
		ctx,
		query,
		req.Id,
		statId,
	)
	if err != nil {
		r.logger.Error("Failed to exec deleteplayerstatspecific query", zap.Error(err))
		return fmt.Errorf("delete stat specififc: %w", err)
	}

	return nil
}
