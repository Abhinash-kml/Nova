package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type ItemsPostgresRepository struct {
	logger           *zap.Logger
	pgx              *pgxpool.Pool
	statementBuilder squirrel.StatementBuilderType
}

func NewItemsPostgresRepositoryFromPgxPool(connection *pgxpool.Pool, l *zap.Logger) *ItemsPostgresRepository {
	return &ItemsPostgresRepository{
		logger:           l,
		pgx:              connection,
		statementBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *ItemsPostgresRepository) Add(ctx context.Context, req CreateItemRequest) (PlayerItem, error) {
	rawQuery := `INSERT INTO
					player_items(
						name,
						description,
						rarity,
						is_stackable,
						meta,
						created_at
					)
				VALUES ($1, $2, $3, $4, COALESCE($5, '{}'::jsonb), $6)
				RETURNING
					id,
					name,
					description,
					rarity,
					is_stackable,
					meta,
					created_at,
					updated_at;`

	var item PlayerItem

	err := r.pgx.QueryRow(
		ctx,
		rawQuery,
		req.Name,
		req.Description,
		req.Rarity,
		req.IsStackable,
		req.Meta,
		time.Now(),
	).Scan(
		&item.Id,
		&item.Name,
		&item.Description,
		&item.Rarity,
		&item.IsStackable,
		&item.Meta,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return PlayerItem{}, fmt.Errorf("scanning row: %w", err)
	}

	return item, nil
}

func (r *ItemsPostgresRepository) GetAll(ctx context.Context, cursor int, limit int) ([]PlayerItem, error) {
	rawQuery := `SELECT
					id,
					name,
					description,
					rarity,
					is_stackable,
					meta,
					created_at,
					updated_at
				FROM
					player_items
				ORDER BY
					created_at DESC
				OFFSET $1
				LIMIT $2;`

	rows, err := r.pgx.Query(
		ctx,
		rawQuery,
		cursor,
		limit,
	)
	if err != nil {
		return nil, fmt.Errorf("getting item: %w", err)
	}
	defer rows.Close()

	var items []PlayerItem

	for rows.Next() {
		var item PlayerItem

		err = rows.Scan(
			&item.Id,
			&item.Name,
			&item.Description,
			&item.Rarity,
			&item.IsStackable,
			&item.Meta,
			&item.CreatedAt,
			&item.UpdatedAt,
		)

		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return items, nil
}

func (r *ItemsPostgresRepository) GetById(ctx context.Context, id uuid.UUID) (PlayerItem, error) {
	rawQuery := `SELECT
					id,
					name,
					description,
					rarity,
					is_stackable,
					meta,
					created_at,
					updated_at
				FROM
					player_items
				WHERE
					id = $1;`

	var item PlayerItem

	err := r.pgx.QueryRow(
		ctx,
		rawQuery,
		id,
	).Scan(
		&item.Id,
		&item.Name,
		&item.Description,
		&item.Rarity,
		&item.IsStackable,
		&item.Meta,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return PlayerItem{}, fmt.Errorf("scanning row: %w", err)
	}

	return item, nil
}

func (r *ItemsPostgresRepository) Update(ctx context.Context, req UpdateItemRequest) (PlayerItem, error) {
	queryBuilder := r.statementBuilder.
		Update("player_items").
		Where(squirrel.Eq{"id": req.Id}).
		Set("updated_at", time.Now())

	if req.Name != nil {
		queryBuilder = queryBuilder.Set("name", *req.Name)
	}

	if req.Description != nil {
		queryBuilder = queryBuilder.Set("description", *req.Description)
	}

	if req.Rarity != nil {
		queryBuilder = queryBuilder.Set("rarity", *req.Rarity)
	}

	if req.IsStackable != nil {
		queryBuilder = queryBuilder.Set("is_stackable", *req.IsStackable)
	}

	if req.Meta != nil {
		queryBuilder = queryBuilder.Set("meta", req.Meta)
	}

	queryBuilder = queryBuilder.Suffix(
		`RETURNING
			id,
			name,
			description,
			rarity,
			is_stackable,
			meta,
			created_at,
			updated_at`,
	)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return PlayerItem{}, fmt.Errorf("generate query: %w", err)
	}

	var item PlayerItem

	err = r.pgx.QueryRow(
		ctx,
		query,
		args...,
	).Scan(
		&item.Id,
		&item.Name,
		&item.Description,
		&item.Rarity,
		&item.IsStackable,
		&item.Meta,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return PlayerItem{}, fmt.Errorf("scanning row: %w", err)
	}

	return item, nil
}

func (r *ItemsPostgresRepository) Replace(ctx context.Context, req ReplaceItemRequest) (PlayerItem, error) {
	rawQuery := `UPDATE
					player_items
				SET
					name = $2,
					description = $3,
					rarity = $4,
					is_stackable = $5,
					meta = $6,
					updated_at = $7
				WHERE
					id = $1
				RETURNING
					id,
					name,
					description,
					rarity,
					is_stackable,
					meta,
					created_at,
					updated_at;`

	var item PlayerItem

	err := r.pgx.QueryRow(
		ctx,
		rawQuery,
		req.Id,
		req.Name,
		req.Description,
		req.Rarity,
		req.IsStackable,
		req.Meta,
		time.Now(),
	).Scan(
		&item.Id,
		&item.Name,
		&item.Description,
		&item.Rarity,
		&item.IsStackable,
		&item.Meta,
		&item.CreatedAt,
		&item.UpdatedAt,
	)

	if err != nil {
		return PlayerItem{}, fmt.Errorf("scanning row: %w", err)
	}

	return item, nil
}

func (r *ItemsPostgresRepository) Delete(ctx context.Context, req DeleteItemRequest) (uuid.UUID, error) {
	rawQuery := `DELETE FROM
					player_items
				WHERE
					id = $1
				RETURNING
					id;`

	var deletedId uuid.UUID

	err := r.pgx.QueryRow(
		ctx,
		rawQuery,
		req.Id,
	).Scan(&deletedId)

	if err != nil {
		return uuid.Nil, fmt.Errorf("scanning row: %w", err)
	}

	return deletedId, nil
}
