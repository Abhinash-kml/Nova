package economy

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresWalletRepository struct {
	pgx              *pgxpool.Pool
	config           *config.Config
	logger           *zap.Logger
	statementBuilder squirrel.StatementBuilderType
}

func NewPostgresWalletRepository(pool *pgxpool.Pool, c *config.Config, l *zap.Logger) *PostgresWalletRepository {
	return &PostgresWalletRepository{
		pgx:              pool,
		config:           c,
		logger:           l,
		statementBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

func (r *PostgresWalletRepository) CreateWalletOfNewPlayer(ctx context.Context, dto CreateNewPlayerWalletDTO) error {
	rawQuery := `
		INSERT INTO
			wallets(
				user_id,
				currency_id,
				account_kind,
				balance_minor
			) SELECT
				$1::uuid AS user_id,
				id AS currency_id,
				'game' AS account_kind,
				0 AS balance_minor
			FROM
				currencies
		ON CONFLICT (user_id, currency_id)
		DO NOTHING;
	`
	_, err := r.pgx.Exec(ctx, rawQuery, dto.UserID)
	if err != nil {
		return fmt.Errorf("creating wallet of player: %w", err)
	}

	return nil
}

func (r *PostgresWalletRepository) GetWalletOfPlayer(ctx context.Context, dto GetWalletOfPlayerDTO) ([]WalletDTO, error) {
	rawQuery := `
		SELECT
			currency_id,
			balance_minor
		FROM
			wallets
		WHERE
			user_id = $1::uuid;
	`

	rows, err := r.pgx.Query(ctx, rawQuery, dto.UserID)
	if err != nil {
		return nil, fmt.Errorf("getting wallet of player: %w", err)
	}
	defer rows.Close()

	var wallets []WalletDTO

	for rows.Next() {
		var wallet WalletDTO
		if err := rows.Scan(
			&wallet.CurrencyID,
			&wallet.BalanceMinor,
		); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		wallets = append(wallets, wallet)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return wallets, nil
}

func (r *PostgresWalletRepository) UpdateWalletOfPlayer(ctx context.Context, dto UpdateWalletOfPlayerDTO) ([]WalletDTO, error) {
	walletQuery := `
		UPDATE
			wallets
		SET
			balance_minor = balance_minor + $3
		WHERE
			user_id = $1
			AND
			currency_id = $2;
	`

	ledgerQuery := `
		INSERT INTO
			wallet_ledger(
				user_id,
				currency_id,
				amount_minor,
				transaction_type,
				reference_id,
				created_at
			)
			VALUES(
				$1,
				$2,
				$3,
				$4,
				$5,
				$6
			);
	`

	now := time.Now()

	transaction, err := r.pgx.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return nil, fmt.Errorf("begin transanction: %w", err)
	}
	defer func() {
		if err := transaction.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			r.logger.Error("rolling back transaction", zap.Error(err))
		}
	}()

	_, err = transaction.Exec(ctx, walletQuery, dto.UserID, dto.CurrencyID, dto.Amount)
	if err != nil {
		return nil, fmt.Errorf("update wallet: %w", err)
	}

	dummyIdempotencyKey, _ := uuid.NewV7()
	_, err = transaction.Exec(ctx, ledgerQuery, dto.UserID, dto.CurrencyID, dto.Amount, dto.Operation, dummyIdempotencyKey, now)
	if err != nil {
		return nil, fmt.Errorf("update wallet ledger: %w", err)
	}

	err = transaction.Commit(ctx)
	if err != nil {
		return nil, fmt.Errorf("commit transaction: %w", err)
	}

	return nil, nil
}

func (r *PostgresWalletRepository) DeleteWalletOfPlayer(ctx context.Context, dto DeleteWalletOfPlayerDTO) error {
	rawQuery := `
		DELETE FROM
			wallets
		WHERE
			user_id = $1::uuid;
	`

	_, err := r.pgx.Exec(ctx, rawQuery, dto.UserID)
	if err != nil {
		return fmt.Errorf("deleting wallet of player: %w", err)
	}

	return nil
}
