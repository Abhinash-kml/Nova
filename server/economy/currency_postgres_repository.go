package economy

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresCurrencyRepository struct {
	pgx              *pgxpool.Pool
	config           *config.CurrencyConfig
	logger           *zap.Logger
	statementbuilder squirrel.StatementBuilderType
}

func NewPostgresCurrencyRepository(pool *pgxpool.Pool, c *config.CurrencyConfig, l *zap.Logger) *PostgresCurrencyRepository {
	return &PostgresCurrencyRepository{
		pgx:              pool,
		config:           c,
		logger:           l,
		statementbuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostgresCurrencyRepository) Get(ctx context.Context, req GetCurrencyRequest) (Currency, error) {
	rawQuery := `
		SELECT
			id,
			name,
			code,
			kind,
			is_purchasable,
			is_tradable,
			expires_after_days,
			max_balance_minor,
			sort_order
		FROM
			currencies
		WHERE
			id = $1;
	`
	var currency Currency

	row := r.pgx.QueryRow(ctx, rawQuery, req.ID)
	if err := row.Scan(
		&currency.ID,
		&currency.Name,
		&currency.Code,
		&currency.Kind,
		&currency.IsPurchasable,
		&currency.IsTradable,
		&currency.ExipiresAfterDays,
		&currency.MaxBalanceMinor,
		&currency.SortOrder,
	); err != nil {
		return Currency{}, fmt.Errorf("getting currency (%s): %w", req.ID, err)
	}

	return currency, nil
}

func (r *PostgresCurrencyRepository) GetAll(ctx context.Context, req GetAllCurrencyRequest) ([]Currency, error) {
	rawQuery := `
		SELECT
			id,
			name,
			code,
			kind,
			is_purchasable,
			is_tradable,
			expires_after_days,
			max_balance_minor,
			sort_order
		FROM
			currencies
		WHERE
			id > $1
		LIMIT 
			$2;
	`

	rows, err := r.pgx.Query(ctx, rawQuery, req.Cursor, req.Limit)
	if err != nil {
		return nil, fmt.Errorf("getting all currency: %w", err)
	}
	defer rows.Close()

	var currencies []Currency
	var currency Currency

	for rows.Next() {
		if err := rows.Scan(
			&currency.ID,
			&currency.Name,
			&currency.Code,
			&currency.Kind,
			&currency.IsPurchasable,
			&currency.IsTradable,
			&currency.ExipiresAfterDays,
			&currency.MaxBalanceMinor,
			&currency.SortOrder,
		); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		currencies = append(currencies, currency)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return currencies, nil
}

func (r *PostgresCurrencyRepository) Create(ctx context.Context, req CreateCurrencyRequest) (Currency, error) {
	rawQuery := `
		INSERT INTO
			currencies(
				name,
				code,
				kind,
				is_purchasable,
				is_tradable,
				expires_after_days,
				max_balance_minor,
				sort_order
			) VALUES(
				$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7,
				$8
			) 
		RETURNING *;
	`

	var currency Currency

	row := r.pgx.QueryRow(
		ctx,
		rawQuery,
		req.Name,
		req.Code,
		req.Kind,
		req.IsPurchasable,
		req.IsPurchasable,
		req.ExipiresAfterDays,
		req.MaxBalanceMinor,
		req.SortOrder)
	if err := row.Scan(
		&currency.ID,
		&currency.Name,
		&currency.Code,
		&currency.Kind,
		&currency.IsPurchasable,
		&currency.IsTradable,
		&currency.ExipiresAfterDays,
		&currency.MaxBalanceMinor,
		&currency.SortOrder,
	); err != nil {
		return Currency{}, fmt.Errorf("scanning row: %w", err)
	}

	return currency, nil
}

func (r *PostgresCurrencyRepository) Update(ctx context.Context, req UpdateCurrencyRequest) (Currency, error) {
	builder := r.statementbuilder.Update("currencies")
	if req.Name != "" {
		builder = builder.Set("name", req.Name)
	}
	if req.Code != "" {
		builder = builder.Set("code", req.Code)
	}
	if req.Kind != "" {
		builder = builder.Set("kind", req.Kind)
	}
	if req.IsPurchasable != nil {
		builder = builder.Set("is_purchasable", *req.IsPurchasable)
	}
	if req.IsTradable != nil {
		builder = builder.Set("is_tradable", *req.IsTradable)
	}
	if req.ExipiresAfterDays != nil {
		builder = builder.Set("expires_after_days", *req.ExipiresAfterDays)
	}
	if req.MaxBalanceMinor != nil {
		builder = builder.Set("max_balance_minor", *req.MaxBalanceMinor)
	}
	if req.SortOrder != nil {
		builder = builder.Set("sort_order", *req.SortOrder)
	}
	builder = builder.Suffix("RETURNING *")

	query, args, err := builder.ToSql()
	if err != nil {
		r.logger.Error("Failed to build sql query for updating currency", zap.Error(err))
		return Currency{}, common.TranslatePostgresError(err, r.logger)
	}

	row := r.pgx.QueryRow(ctx, query, args...)
	var currency Currency

	if err := row.Scan(
		&currency.ID,
		&currency.Name,
		&currency.Code,
		&currency.Kind,
		&currency.IsPurchasable,
		&currency.IsTradable,
		&currency.ExipiresAfterDays,
		&currency.MaxBalanceMinor,
		&currency.SortOrder,
	); err != nil {
		return Currency{}, fmt.Errorf("scanning row: %w", err)
	}

	return currency, nil
}

func (r *PostgresCurrencyRepository) Delete(ctx context.Context, req DeleteCurrencyRequest) (int, error) {
	rawQuery := `
		DELETE FROM
			currencies
		WHERE
			id = $1
		RETURNING
			id;
	`

	row := r.pgx.QueryRow(ctx, rawQuery, req.ID)
	var deletedCurrencyID int

	if err := row.Scan(&deletedCurrencyID); err != nil {
		return 0, fmt.Errorf("scanning row: %w", err)
	}

	return deletedCurrencyID, nil
}
