package economy

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/common"
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
		r.logger.Error("Failed to execute create wallet of new player query", zap.Error(err))
		return common.TranslatePostgresError(err, r.logger)
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
		r.logger.Error("Failed to execute get wallet of user query", zap.Error(err))
		return nil, common.TranslatePostgresError(err, r.logger)
	}
	defer rows.Close()

	var wallets []WalletDTO

	for rows.Next() {
		var wallet WalletDTO
		if err := rows.Scan(
			&wallet.CurrencyID,
			&wallet.BalanceMinor,
		); err != nil {
			r.logger.Error("Failed to scan row in get wallet of user query", zap.Error(err))
			return nil, common.TranslatePostgresError(err, r.logger)
		}

		wallets = append(wallets, wallet)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error while looping through rows", zap.Error(err))
		return nil, common.TranslatePostgresError(err, r.logger)
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
		r.logger.Error("Failed to start transaction", zap.Error(err))
		return nil, err
	}
	defer func() {
		if err := transaction.Rollback(ctx); err != nil && !errors.Is(err, pgx.ErrTxClosed) {
			log.Printf("failed to rollback transaction: %v", err)
		}
	}()

	_, err = transaction.Exec(ctx, walletQuery, dto.UserID, dto.CurrencyID, dto.Amount)
	if err != nil {
		r.logger.Error("Failed to update wallet of user", zap.Error(err))
		return nil, err
	}

	dummyIdempotencyKey, _ := uuid.NewV7()
	_, err = transaction.Exec(ctx, ledgerQuery, dto.UserID, dto.CurrencyID, dto.Amount, dto.Operation, dummyIdempotencyKey, now)
	if err != nil {
		r.logger.Error("Failed to update ledger of wallet transaction", zap.Error(err))
		return nil, err
	}

	err = transaction.Commit(ctx)
	if err != nil {
		r.logger.Error("Failed to commit transaction of updating wallet", zap.Error(err))
		return nil, err
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
		r.logger.Error("Failed to execute delete wallet of player query", zap.Error(err))
		return common.TranslatePostgresError(err, r.logger)
	}

	return nil
}
