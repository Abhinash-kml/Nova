package stats

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/common"
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
	return nil
}

func (r *PostgresRepository) Add(ctx context.Context, dto CreateDTO) (Stats, error) {
	rawQuery := `INSERT INTO
					stats(name, startvalue, created_at)
				VALUES ($1, $2, $3)
				RETURNING 
					id, 
					name, 
					startvalue, 
					created_at;`

	var stat Stats
	now := time.Now()

	err := r.pgx.QueryRow(
		ctx,
		rawQuery,
		dto.Name,
		dto.StartValue,
		now,
	).Scan(
		&stat.Id,
		&stat.Name,
		&stat.StartValue,
		&stat.CreatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to execute sql insert query", zap.Error(err))
		return Stats{}, common.TranslatePostgresError(err, r.logger)
	}

	return stat, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context, cursor int, limit int) ([]Stats, error) {
	var rows pgx.Rows
	var err error

	rawQuery := `SELECT
					id,
					name,
					startvalue,
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
		r.logger.Error("Failed to execute getall query", zap.Error(err))
		return nil, common.TranslatePostgresError(err, r.logger)
	}

	defer rows.Close()

	var stats []Stats

	for rows.Next() {
		var stat Stats

		err = rows.Scan(
			&stat.Id,
			&stat.Name,
			&stat.StartValue,
			&stat.CreatedAt,
		)

		if err != nil {
			r.logger.Error("Failed to scan returned row in getall query", zap.Error(err))
			return nil, common.TranslatePostgresError(err, r.logger)
		}

		stats = append(stats, stat)
	}

	return stats, rows.Err()
}

func (r *PostgresRepository) GetById(ctx context.Context, id int) (Stats, error) {
	rawQuery := `SELECT
					id,
					name,
					startvalue,
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
		&stat.Name,
		&stat.StartValue,
		&stat.CreatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to scan row in getbyid query", zap.Error(err))
		return Stats{}, common.TranslatePostgresError(err, r.logger)
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
		queryBuilder = queryBuilder.Set("startvalue", *dto.StartValue)
	}

	queryBuilder = queryBuilder.Suffix("RETURNING id, name, startvalue, created_at")

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate update query using squirrel", zap.Error(err))
		return Stats{}, common.TranslatePostgresError(err, r.logger)
	}

	var stat Stats

	tx, err := r.pgx.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})
	if err != nil {
		r.logger.Error("Failed to begin transaction in update query", zap.Error(err))
		return Stats{}, common.TranslatePostgresError(err, r.logger)
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		"SELECT id FROM stats WHERE id = $1 FOR UPDATE;",
		dto.Id,
	)

	if err != nil {
		r.logger.Error("Failed to lock row for update in transaction", zap.Error(err))
		return Stats{}, common.TranslatePostgresError(err, r.logger)
	}

	result := tx.QueryRow(ctx, query, args...)

	err = result.Scan(
		&stat.Id,
		&stat.Name,
		&stat.StartValue,
		&stat.CreatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to scan returned object from update query in transaction", zap.Error(err))
		return Stats{}, common.TranslatePostgresError(err, r.logger)
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Error("Failed to commit transaction in update query", zap.Error(err))
		return Stats{}, common.TranslatePostgresError(err, r.logger)
	}

	return stat, nil
}

func (r *PostgresRepository) Replace(ctx context.Context, dto ReplaceDTO) (Stats, error) {
	rawQuery := `UPDATE
					stats
				SET
					name = $2,
					startvalue = $3
				WHERE
					id = $1
				RETURNING
					id, name, startvalue, created_at;`

	var stat Stats

	tx, err := r.pgx.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		r.logger.Error("Failed to begin transaction in replace query", zap.Error(err))
		return Stats{}, common.TranslatePostgresError(err, r.logger)
	}

	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		"SELECT id FROM stats WHERE id = $1 FOR UPDATE;",
		dto.Id,
	)

	if err != nil {
		r.logger.Error("Failed to lock row for replace in transaction", zap.Error(err))
		return Stats{}, common.TranslatePostgresError(err, r.logger)
	}

	result := tx.QueryRow(
		ctx,
		rawQuery,
		dto.Id,
		dto.Name,
		dto.StartValue,
	)

	err = result.Scan(
		&stat.Id,
		&stat.Name,
		&stat.StartValue,
		&stat.CreatedAt,
	)

	if err != nil {
		r.logger.Error("Failed to scan result of replace query", zap.Error(err))
		return Stats{}, common.TranslatePostgresError(err, r.logger)
	}

	err = tx.Commit(ctx)
	if err != nil {
		r.logger.Error("Failed to commit transaction in replace query", zap.Error(err))
		return Stats{}, common.TranslatePostgresError(err, r.logger)
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
		r.logger.Error("Failed to scan result in delete query", zap.Error(err))
		return 0, common.TranslatePostgresError(err, r.logger)
	}

	return deletedId, nil
}
