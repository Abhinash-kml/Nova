package inventory

import (
	"context"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/common"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresInventoryRepository struct {
	logger           *zap.Logger
	pgx              *pgxpool.Pool
	statementBuilder squirrel.StatementBuilderType
}

func NewPostgresInventoryRepositoryFromPgxPool(connection *pgxpool.Pool, l *zap.Logger) *PostgresInventoryRepository {
	return &PostgresInventoryRepository{
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
		r.logger.Error("Failed to execute get player inventory query", zap.Error(err))
		return nil, common.TranslatePostgresError(err, r.logger)
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
			r.logger.Error("Failed to scan player inventory row", zap.Error(err))
			return nil, common.TranslatePostgresError(err, r.logger)
		}

		inventory = append(inventory, item)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Failed while iterating player inventory rows", zap.Error(err))
		return nil, common.TranslatePostgresError(err, r.logger)
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
		r.logger.Error("Failed execute delete user inventory query", zap.Error(err))
		return UserID{}, common.TranslatePostgresError(err, r.logger)
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
		r.logger.Error("Failed to scan player inventory", zap.Error(err))
		return PlayerInventory{}, common.TranslatePostgresError(err, r.logger)
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
		r.logger.Error("Failed to execute add inventory item of user query", zap.Error(err))
		return InventoryItem{}, common.TranslatePostgresError(err, r.logger)
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
		r.logger.Error("Failed to generate player inventory update query", zap.Error(err))
		return PlayerInventory{}, common.TranslatePostgresError(err, r.logger)
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
		r.logger.Error("Failed to scan updated player inventory", zap.Error(err))
		return PlayerInventory{}, common.TranslatePostgresError(err, r.logger)
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
		r.logger.Error("Failed to execute delete player inventory query", zap.Error(err))
		return common.TranslatePostgresError(err, r.logger)
	}

	return nil
}
