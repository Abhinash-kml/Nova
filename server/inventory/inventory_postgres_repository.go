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
	config           *config.InventoryConfig
	logger           *zap.Logger
	pgx              *pgxpool.Pool
	statementBuilder squirrel.StatementBuilderType
}

func NewPostgresInventoryRepositoryFromPgxPool(connection *pgxpool.Pool, c *config.InventoryConfig, l *zap.Logger) *PostgresInventoryRepository {
	return &PostgresInventoryRepository{
		config:           c,
		logger:           l,
		pgx:              connection,
		statementBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostgresInventoryRepository) GetInventoryOfUser(ctx context.Context, req GetInventoryOfUserRequest) ([]PlayerInventory, error) {
	rawQuery := `SELECT
					user_id,
					item_id,
					quantity,
					updated_at,
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
		req.Id,
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

func (r *PostgresInventoryRepository) DeleteInventoryOfUser(ctx context.Context, req DeleteInventoryOfUserRequest) (UserID, error) {
	rawQuery := `DELETE FROM 
					player_inventory 
				WHERE
					user_id = $1
				RETURNING 
					user_id;`

	var deletedID string
	userID, _ := uuid.Parse(req.Id)
	row := r.pgx.QueryRow(ctx, rawQuery, userID)
	err := row.Scan(&deletedID)
	if err != nil {
		return UserID{}, fmt.Errorf("scanning row: %w", err)
	}

	return UserID{Id: deletedID}, nil
}

func (r *PostgresInventoryRepository) GetInventoryItemOfUser(ctx context.Context, req GetInventoryItemOfUserRequest) (PlayerInventory, error) {
	rawQuery := `SELECT
					user_id,
					item_id,
					quantity,
					updated_at,
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
		req.Id,
		req.ItemId,
	).Scan(
		&inventory.UserId,
		&inventory.ItemId,
		&inventory.Quantity,
		&inventory.UpdatedAt,
		&inventory.IsEquipped,
	)

	if err != nil {
		return PlayerInventory{}, fmt.Errorf("scanning row: %w", err)
	}

	return inventory, nil
}

func (r *PostgresInventoryRepository) AddInventoryItemOfUser(ctx context.Context, req AddInventoryItemOfUserRequest) (InventoryItem, error) {
	rawQuery := `INSERT INTO
					player_inventory(
					user_id,
					item_id,
					quantity,
					is_equipped,
					updated_at
				) VALUES (
					$1, 
					$2, 
					$3, 
					$4, 
					$5
				)
				RETURNING 
					*;`

	now := time.Now()
	var item InventoryItem
	row := r.pgx.QueryRow(
		ctx,
		rawQuery,
		req.UserID,
		req.ItemID,
		req.Quantity,
		req.IsEquipped,
		now)

	err := row.Scan(
		&item.UserID,
		&item.ItemID,
		&item.Quantity,
		&item.IsEquipped,
		&item.UpdatedAt,
	)
	if err != nil {
		return InventoryItem{}, fmt.Errorf("scanning row: %w", err)
	}

	return item, nil
}

func (r *PostgresInventoryRepository) UpdateInventoryItemOfUser(ctx context.Context, req UpdateInventoryItemOfUserRequest) (PlayerInventory, error) {
	userID, _ := uuid.Parse(req.UserID.Id)
	itemID, _ := uuid.Parse(req.ItemID.ItemId)
	queryBuilder := r.statementBuilder.
		Update("player_inventory").
		Where(squirrel.Eq{
			"user_id": userID,
			"item_id": itemID,
		}).
		Set("updated_at", time.Now())

	if req.Quantity != nil {
		queryBuilder = queryBuilder.Set("quantity", *req.Quantity)
	}

	if req.IsEquipped != nil {
		queryBuilder = queryBuilder.Set("is_equipped", *req.IsEquipped)
	}

	queryBuilder = queryBuilder.Suffix(
		`RETURNING
			user_id,
			item_id,
			quantity,
			updated_at,
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
		&inventory.IsEquipped,
	)

	if err != nil {
		return PlayerInventory{}, fmt.Errorf("scanning row: %w", err)
	}

	return inventory, nil
}

func (r *PostgresInventoryRepository) DeleteInventoryItemOfUser(ctx context.Context, req DeleteInventoryItemOfUserRequest) error {
	rawQuery := `DELETE FROM
					player_inventory
				WHERE
					user_id = $1
					AND item_id = $2;`

	userID, _ := uuid.Parse(req.UserID.Id)
	itemID, _ := uuid.Parse(req.ItemID.ItemId)
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

func (r *PostgresInventoryRepository) BuildInventoryFromEntitlements(ctx context.Context, req GetInventoryOfUserRequest) ([]InventoryItem, error) {
	rawQuery := `
	WITH valid_entitlements AS (
		SELECT
			item_id
		FROM
			entitlements
		WHERE
			user_id = $1
			AND
			(valid_till IS NULL
			OR valid_till > CURRENT_TIMESTAMP)
	)
	SELECT
		pi.item_id,
		pi.quantity,
		pi.is_equipped
	FROM
		player_inventory pi
	JOIN
		valid_entitlements e ON e.item_id = pi.item_id
	WHERE
		pi.user_id = $1;
	`

	userID, _ := uuid.Parse(req.Id)
	rows, err := r.pgx.Query(ctx, rawQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("getting inventory: %w", err)
	}
	defer rows.Close()

	var inventory []InventoryItem
	var inventoryItem InventoryItem
	for rows.Next() {
		if err = rows.Scan(&inventoryItem.ItemID, &inventoryItem.Quantity, &inventoryItem.IsEquipped); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		inventory = append(inventory, inventoryItem)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return inventory, nil
}
