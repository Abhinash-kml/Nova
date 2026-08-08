package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresInventoryRepository struct {
	config           *config.Config
	logger           *zap.Logger
	pgx              *pgxpool.Pool
	statementBuilder squirrel.StatementBuilderType
}

func NewPostgresInventoryRepositoryFromPgxPool(connection *pgxpool.Pool, c *config.Config, l *zap.Logger) *PostgresInventoryRepository {
	return &PostgresInventoryRepository{
		config:           c,
		logger:           l,
		pgx:              connection,
		statementBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostgresInventoryRepository) GetInventoryOfUser(ctx context.Context, dto GetInventoryOfUserDTO) ([]PlayerInventory, error) {
	rawQuery := `SELECT
					user_id,
					item_id,
					quantity,
					updated_at,
					source,
					is_equipped
				FROM
					player_inventory
				WHERE
					user_id = $1
				ORDER BY
					updated_at DESC;`

	rows, err := r.pgx.Query(
		ctx,
		rawQuery,
		dto.Id,
	)
	if err != nil {
		return nil, fmt.Errorf("getting inventory: %w", err)
	}
	defer rows.Close()

	var inventory []PlayerInventory

	for rows.Next() {
		var item PlayerInventory

		err = rows.Scan(
			&item.UserId,
			&item.ItemId,
			&item.Quantity,
			&item.UpdatedAt,
			&item.Source,
			&item.IsEquipped,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		inventory = append(inventory, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return inventory, nil
}

func (r *PostgresInventoryRepository) DeleteInventoryOfUser(ctx context.Context, dto DeleteInventoryOfUserDTO) (UserID, error) {
	rawQuery := `DELETE FROM 
					player_inventory 
				WHERE
					user_id = $1
				RETURNING user_id;`

	var deletedID string
	userID, _ := uuid.Parse(dto.Id)
	row := r.pgx.QueryRow(ctx, rawQuery, userID)
	err := row.Scan(&deletedID)
	if err != nil {
		return UserID{}, fmt.Errorf("scanning row: %w", err)
	}

	return UserID{Id: deletedID}, nil
}

func (r *PostgresInventoryRepository) GetInventoryItemOfUser(ctx context.Context, dto GetInventoryItemOfUserDTO) (PlayerInventory, error) {
	rawQuery := `SELECT
					user_id,
					item_id,
					quantity,
					updated_at,
					source,
					is_equipped
				FROM
					player_inventory
				WHERE
					user_id = $1
					AND item_id = $2;`

	var inventory PlayerInventory

	err := r.pgx.QueryRow(
		ctx,
		rawQuery,
		dto.Id,
		dto.ItemId,
	).Scan(
		&inventory.UserId,
		&inventory.ItemId,
		&inventory.Quantity,
		&inventory.UpdatedAt,
		&inventory.Source,
		&inventory.IsEquipped,
	)

	if err != nil {
		return PlayerInventory{}, fmt.Errorf("scanning row: %w", err)
	}

	return inventory, nil
}

func (r *PostgresInventoryRepository) AddInventoryItemOfUser(ctx context.Context, dto AddInventoryItemOfUserDTO) (InventoryItem, error) {
	rawQuery := `INSERT INTO
					player_inventory(
					user_id,
					item_id,
					quantity,
					source,
					is_equipped,
					updated_at
				) VALUES ($1, $2, $3, $4, $5, $6)
				RETURNING *;`

	now := time.Now()
	var item InventoryItem
	row := r.pgx.QueryRow(
		ctx,
		rawQuery,
		dto.UserID,
		dto.ItemID,
		dto.Quantity,
		dto.Source,
		dto.IsEquipped,
		now)

	err := row.Scan(
		&item.UserID,
		&item.ItemID,
		&item.Quantity,
		&item.Source,
		&item.IsEquipped,
		&item.UpdatedAt,
	)
	if err != nil {
		return InventoryItem{}, fmt.Errorf("scanning row: %w", err)
	}

	return item, nil
}

func (r *PostgresInventoryRepository) UpdateInventoryItemOfUser(ctx context.Context, dto UpdateInventoryItemOfUserDTO) (PlayerInventory, error) {
	userID, _ := uuid.Parse(dto.UserID.Id)
	itemID, _ := uuid.Parse(dto.ItemID.ItemId)
	queryBuilder := r.statementBuilder.
		Update("player_inventory").
		Where(squirrel.Eq{
			"user_id": userID,
			"item_id": itemID,
		}).
		Set("updated_at", time.Now())

	if dto.Quantity != nil {
		queryBuilder = queryBuilder.Set("quantity", *dto.Quantity)
	}

	if dto.Source != nil {
		queryBuilder = queryBuilder.Set("source", *dto.Source)
	}

	if dto.IsEquipped != nil {
		queryBuilder = queryBuilder.Set("is_equipped", *dto.IsEquipped)
	}

	queryBuilder = queryBuilder.Suffix(
		`RETURNING
			user_id,
			item_id,
			quantity,
			updated_at,
			source,
			is_equipped`,
	)

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return PlayerInventory{}, fmt.Errorf("generate query: %w", err)
	}

	var inventory PlayerInventory

	err = r.pgx.QueryRow(
		ctx,
		query,
		args...,
	).Scan(
		&inventory.UserId,
		&inventory.ItemId,
		&inventory.Quantity,
		&inventory.UpdatedAt,
		&inventory.Source,
		&inventory.IsEquipped,
	)

	if err != nil {
		return PlayerInventory{}, fmt.Errorf("scanning row: %w", err)
	}

	return inventory, nil
}

func (r *PostgresInventoryRepository) DeleteInventoryItemOfUser(ctx context.Context, dto DeleteInventoryItemOfUserDTO) error {
	rawQuery := `DELETE FROM
					player_inventory
				WHERE
					user_id = $1
					AND item_id = $2;`

	userID, _ := uuid.Parse(dto.UserID.Id)
	itemID, _ := uuid.Parse(dto.ItemID.ItemId)
	_, err := r.pgx.Exec(
		ctx,
		rawQuery,
		userID,
		itemID,
	)
	if err != nil {
		return fmt.Errorf("deleting inventory item: %w", err)
	}

	return nil
}
