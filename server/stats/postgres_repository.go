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
	return nil
}

func (r *PostgresRepository) Add(ctx context.Context, dto CreateDTO) (Stats, error) {
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
		dto.Key,
		dto.Name,
		dto.Description,
		dto.StartValue,
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
	var rows pgx.Rows
	var err error

	rawQuery := `SELECT
					id,
					key,
					name,
					description,
					start_value,
					created_at
				FROM stats`

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

func (r *PostgresRepository) Update(ctx context.Context, dto UpdateDTO) (Stats, error) {
	queryBuilder := r.statementBuilder.
		Update("stats").
		Where(squirrel.Eq{"id": dto.Id})

	if dto.Name != nil {
		queryBuilder = queryBuilder.Set("name", *dto.Name)
	}

	if dto.StartValue != nil {
		queryBuilder = queryBuilder.Set("start_value", *dto.StartValue)
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
		dto.Id,
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

func (r *PostgresRepository) Replace(ctx context.Context, dto ReplaceDTO) (Stats, error) {
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
		dto.Id,
	)
	if err != nil {
		return Stats{}, fmt.Errorf("locking row: %w", err)
	}

	result := tx.QueryRow(
		ctx,
		rawQuery,
		dto.Id,
		dto.Name,
		dto.Key,
		dto.Description,
		dto.StartValue,
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

func (r *PostgresRepository) Delete(ctx context.Context, dto DeleteDTO) (int, error) {
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
		dto.Id,
	).Scan(&deletedId)
	if err != nil {
		return 0, fmt.Errorf("scanning row: %w", err)
	}

	return deletedId, nil
}

func (r *PostgresRepository) GetPlayerStats(ctx context.Context, dto GetPlayerStatDTO) (PlayerStatsResponseDTO, error) {

	var response PlayerStatsResponseDTO

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
		dto.Id,
	)
	if err != nil {
		return PlayerStatsResponseDTO{}, fmt.Errorf("getting player stats: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var stat StatResponseDTO
		if err := rows.Scan(&stat.Id, &stat.Value); err != nil {
			return PlayerStatsResponseDTO{}, fmt.Errorf("scanning row: %w", err)
		}

		response.Stats = append(response.Stats, stat)
	}

	if err := rows.Err(); err != nil {
		return PlayerStatsResponseDTO{}, fmt.Errorf("iterating rows: %w", err)
	}

	return response, nil
}

func (r *PostgresRepository) UpdatePlayerStats(ctx context.Context, dto UpdatePlayerStatDTO) error {
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
		dto.Id,
		dto.StatId,
		dto.Value,
	)
	if err != nil {
		return fmt.Errorf("updating stats: %w", err)
	}

	return nil
}

func (r *PostgresRepository) DeletePlayerStats(ctx context.Context, dto DeletePlayerStatDTO) error {
	query := `DELETE FROM 
				player_stats_progression
			WHERE 
				player_id = $1`

	_, err := r.pgx.Exec(
		ctx,
		query,
		dto.Id,
	)
	if err != nil {
		return fmt.Errorf("delete stats: %w", err)
	}

	return nil
}

func (r *PostgresRepository) DeletePlayerStatSpecific(ctx context.Context, dto DeletePlayerStatSpecificDTO) error {
	query := `DELETE FROM 
				player_stats_progression
			WHERE 
				player_id = $1
			AND 
				stat_id = $2`

	statId, _ := uuid.Parse(dto.StatsId)

	_, err := r.pgx.Exec(
		ctx,
		query,
		dto.Id,
		statId,
	)
	if err != nil {
		r.logger.Error("Failed to exec deleteplayerstatspecific query", zap.Error(err))
		return fmt.Errorf("delete stat specififc: %w", err)
	}

	return nil
}
