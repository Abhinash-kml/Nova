package economy

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresCurrencyRepository struct {
	pgx              *pgxpool.Pool
	config           *config.Config
	logger           *zap.Logger
	statementbuilder squirrel.StatementBuilderType
}

func NewPostgresCurrencyRepository(pool *pgxpool.Pool, c *config.Config, l *zap.Logger) *PostgresCurrencyRepository {
	return &PostgresCurrencyRepository{
		pgx:              pool,
		config:           c,
		logger:           l,
		statementbuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostgresCurrencyRepository) Get(ctx context.Context, dto GetCurrencyDTO) (CurrencyDTO, error) {
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
	var currency CurrencyDTO

	row := r.pgx.QueryRow(ctx, rawQuery, dto.ID)
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
		r.logger.Error("Failed to scan row in get currency query")
		return CurrencyDTO{}, common.TranslatePostgresError(err, r.logger)
	}

	return currency, nil
}

func (r *PostgresCurrencyRepository) GetAll(ctx context.Context, dto GetAllCurrencyDTO) ([]CurrencyDTO, error) {
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

	rows, err := r.pgx.Query(ctx, rawQuery, dto.Cursor, dto.Limit)
	if err != nil {
		r.logger.Error("Failed to execute get all currency query", zap.Error(err))
		return nil, common.TranslatePostgresError(err, r.logger)
	}
	defer rows.Close()

	var currencies []CurrencyDTO

	for rows.Next() {
		var currency CurrencyDTO

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
			r.logger.Error("Failed to scan row in get currency query")
			return nil, common.TranslatePostgresError(err, r.logger)
		}

		currencies = append(currencies, currency)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error while iterating through rows", zap.Error(err))
		return nil, common.TranslatePostgresError(err, r.logger)
	}

	return currencies, nil
}

func (r *PostgresCurrencyRepository) Create(ctx context.Context, dto CreateCurrencyDTO) (CurrencyDTO, error) {
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

	var currency CurrencyDTO

	row := r.pgx.QueryRow(
		ctx,
		rawQuery,
		dto.Name,
		dto.Code,
		dto.Kind,
		dto.IsPurchasable,
		dto.IsPurchasable,
		dto.ExipiresAfterDays,
		dto.MaxBalanceMinor,
		dto.SortOrder)
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
		r.logger.Error("Failed to scan row in create currency query")
		return CurrencyDTO{}, common.TranslatePostgresError(err, r.logger)
	}

	return currency, nil
}

func (r *PostgresCurrencyRepository) Update(ctx context.Context, dto UpdateCurrencyDTO) (CurrencyDTO, error) {
	builder := r.statementbuilder.Update("currencies")
	if dto.Name != "" {
		builder = builder.Set("name", dto.Name)
	}
	if dto.Code != "" {
		builder = builder.Set("code", dto.Code)
	}
	if dto.Kind != "" {
		builder = builder.Set("kind", dto.Kind)
	}
	if dto.IsPurchasable != nil {
		builder = builder.Set("is_purchasable", *dto.IsPurchasable)
	}
	if dto.IsTradable != nil {
		builder = builder.Set("is_tradable", *dto.IsTradable)
	}
	if dto.ExipiresAfterDays != nil {
		builder = builder.Set("expires_after_days", *dto.ExipiresAfterDays)
	}
	if dto.MaxBalanceMinor != nil {
		builder = builder.Set("max_balance_minor", *dto.MaxBalanceMinor)
	}
	if dto.SortOrder != nil {
		builder = builder.Set("sort_order", *dto.SortOrder)
	}
	builder = builder.Suffix("RETURNING *")

	query, args, err := builder.ToSql()
	if err != nil {
		r.logger.Error("Failed to build sql query for updating currency", zap.Error(err))
		return CurrencyDTO{}, common.TranslatePostgresError(err, r.logger)
	}

	row := r.pgx.QueryRow(ctx, query, args...)

	var currency CurrencyDTO

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
		r.logger.Error("Failed to scan row in update currency query")
		return CurrencyDTO{}, common.TranslatePostgresError(err, r.logger)
	}

	return currency, nil
}

func (r *PostgresCurrencyRepository) Delete(ctx context.Context, dto DeleetCurrencyDTO) (int, error) {
	rawQuery := `
		DELETE FROM
			currencies
		WHERE
			id = $1
		RETURNING
			id;
	`

	row := r.pgx.QueryRow(ctx, rawQuery, dto.ID)

	var deletedCurrencyID int

	if err := row.Scan(&deletedCurrencyID); err != nil {
		r.logger.Error("Failed to scan row in delete currency query", zap.Error(err))
		return 0, common.TranslatePostgresError(err, r.logger)
	}

	return deletedCurrencyID, nil
}
