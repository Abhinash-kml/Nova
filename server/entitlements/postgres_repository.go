package entitlements

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	logger *zap.Logger
	pgx    *pgxpool.Pool
}

func NewPostgresRepository(pool *pgxpool.Pool, l *zap.Logger) *PostgresRepository {
	return &PostgresRepository{
		pgx:    pool,
		logger: l,
	}
}

func (r *PostgresRepository) Create(ctx context.Context, req CreateRequest) (Entitlement, error) {
	rawQuery := `
		INSERT INTO
			entitlements(
				entitlement_id,
				user_id,
				item_id,
				source,
				obtained_on,
				valid_till,
				is_active
			) VALUES (
			 	$1,
				$2,
				$3,
				$4,
				$5,
				$6,
				$7
			);
	`

	entitlementID, err := uuid.NewV7()
	if err != nil {
		return Entitlement{}, fmt.Errorf("generating id: %w", err)
	}
	userID, _ := uuid.Parse(req.UserID)
	itemID, _ := uuid.Parse(req.ItemID)

	currentTimeUTC := time.Now().UTC()

	_, err = r.pgx.Exec(
		ctx,
		rawQuery,
		entitlementID,
		userID,
		itemID,
		req.Source,
		currentTimeUTC,
		req.ValidTill,
		req.IsActive,
	)
	if err != nil {
		return Entitlement{}, fmt.Errorf("creating entitlement: %w", err)
	}

	return Entitlement{
		ID:         entitlementID,
		UserID:     userID,
		ItemID:     itemID,
		Source:     req.Source,
		ObtainedOn: currentTimeUTC,
		ValidTill:  req.ValidTill,
		IsActive:   req.IsActive,
	}, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, req DeleteRequest) (DeleteResponse, error) {
	rawQuery := `
		DELETE FROM
			entitlements
		WHERE
			entitlement_id = $1
		RETURNING
			entitlement_id;
	`

	var deletedID DeleteResponse

	row := r.pgx.QueryRow(ctx, rawQuery, req.Id)
	err := row.Scan(&deletedID.Id)
	if err != nil {
		return DeleteResponse{Id: uuid.Nil}, fmt.Errorf("deleting entitlement: %w", err)
	}

	return deletedID, nil
}

func (r *PostgresRepository) CheckIfActive(ctx context.Context, req CheckIfActiveRequest) (bool, error) {
	rawQuery := `
	SELECT EXISTS (
		SELECT
			true
		FROM
			entitlements
		WHERE
			is_active = true
			AND
			entitlement_id = $1
	);
	`

	var isActive bool
	row := r.pgx.QueryRow(ctx, rawQuery, req.Id)
	err := row.Scan(&isActive)
	if err != nil {
		return false, fmt.Errorf("checking is active: %w", err)
	}

	return isActive, nil
}

func (r *PostgresRepository) GetAllOfUser(ctx context.Context, req GetAllOfUserRequest) ([]Entitlement, error) {
	rawQuery := `
		SELECT
			entitlement_id,
			user_id,
			item_id,
			source,
			obtained_on,
			valid_till,
			is_active
		FROM
			entitlements
		WHERE
			user_id = $1;
	`

	userID, _ := uuid.Parse(req.UserID)

	rows, err := r.pgx.Query(ctx, rawQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("getting entitlements: %w", err)
	}
	defer rows.Close()

	var entitlements []Entitlement
	var ent Entitlement

	for rows.Next() {
		var valid_till sql.NullTime
		if err = rows.Scan(
			&ent.ID,
			&ent.UserID,
			&ent.ItemID,
			&ent.Source,
			&ent.ObtainedOn,
			&valid_till,
			&ent.IsActive,
		); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}
		if valid_till.Valid {
			ent.ValidTill = valid_till.Time
		}
		entitlements = append(entitlements, ent)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return entitlements, nil
}
