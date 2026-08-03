package social

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	pgx              *pgxpool.Pool
	config           *config.Config
	logger           *zap.Logger
	statementBuilder squirrel.StatementBuilderType
}

func NewPostgreRepository(pool *pgxpool.Pool, c *config.Config, l *zap.Logger) *PostgresRepository {
	return &PostgresRepository{
		pgx:    pool,
		config: c,
		logger: l,
	}
}

func (r *PostgresRepository) AddFriend(ctx context.Context, dto AddFriendDTO) error {
	rawQuery := `
		INSERT INTO 
			friendships (
				user_one_id,
				user_two_id,
				status,
				action_user_id
			)
		VALUES (
			LEAST($1::uuid, $2::uuid),
			GREATEST($1::uuid, $2::uuid),
			'pending',
			$1::uuid
		);`

	userID, _ := uuid.Parse(dto.UserID)
	targetID, _ := uuid.Parse(dto.TargetID)

	_, err := r.pgx.Exec(ctx, rawQuery, userID, targetID)
	if err != nil {
		return fmt.Errorf("adding friend: %w", err)
	}

	return nil
}

func (r *PostgresRepository) RemoveFriend(ctx context.Context, dto RemoveFriendDTO) error {
	rawQuery := `
		DELETE FROM 
			friendships
		WHERE 
			user_one_id = LEAST($1::uuid, $2::uuid)
			AND
			user_two_id = GREATEST($1::uuid, $2::uuid)
			AND 
			status = 'accepted';
	`

	userID, _ := uuid.Parse(dto.UserID)
	targetID, _ := uuid.Parse(dto.TargetID)

	result, err := r.pgx.Exec(ctx, rawQuery, userID, targetID)
	if err != nil {
		return fmt.Errorf("removing friend: %w", err)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("removing friend: no op")
	}

	return nil
}

func (r *PostgresRepository) AcceptFriendRequest(ctx context.Context, dto AcceptFriendRequestDTO) error {
	rawQuery := `
		UPDATE
			friendships
		SET
			status = 'accepted',
			action_user_id = $2::uuid
		WHERE
			user_one_id = LEAST($1::uuid, $2::uuid)
			AND
			user_two_id = GREATEST($1::uuid, $2::uuid)
			AND
			status = 'pending';
	`

	userID, _ := uuid.Parse(dto.UserID)
	targetID, _ := uuid.Parse(dto.TargetID)

	result, err := r.pgx.Exec(ctx, rawQuery, targetID, userID)
	if err != nil {
		return fmt.Errorf("accept friend request: %w", err)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("accept friend request: no op")
	}

	return nil
}

func (r *PostgresRepository) RejectFriendRequest(ctx context.Context, dto RejectFriendRequestDTO) error {
	rawQuery := `
		DELETE FROM 
			friendships
		WHERE
			user_one_id = LEAST($1::uuid, $2::uuid)
			AND
			user_two_id = GREATEST($1::uuid, $2::uuid)
			AND
			status = 'pending';
	`

	userID, _ := uuid.Parse(dto.UserID)
	targetID, _ := uuid.Parse(dto.TargetID)

	result, err := r.pgx.Exec(ctx, rawQuery, targetID, userID)
	if err != nil {
		return fmt.Errorf("rejecting friend request: %w", err)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("rejecting friend request: no op")
	}

	return nil
}

func (r *PostgresRepository) GetAllIncomingFriendRequests(ctx context.Context, dto GetIncomingRequestsOfUserDTO) ([]RequestDTO, error) {
	rawQuery := `
		SELECT
			action_user_id,
			created_at
		FROM
			friendships
		WHERE
			(user_one_id = $1::uuid
			OR
			user_two_id = $1::uuid)
			AND
			status = 'pending'
			AND
			action_user_id <> $1::uuid;
	`

	userID, _ := uuid.Parse(dto.Id)

	rows, err := r.pgx.Query(ctx, rawQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("getting incoming requests: %w", err)
	}
	defer rows.Close()

	var requests []RequestDTO

	for rows.Next() {
		var result RequestDTO
		if err := rows.Scan(
			&result.UserID,
			&result.InitiatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		requests = append(requests, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating row: %w", err)
	}

	return requests, nil
}

func (r *PostgresRepository) GetAllOutgoingFriendRequests(ctx context.Context, dto GetOutgoingRequestsOfUserDTO) ([]RequestDTO, error) {
	rawQuery := `
		SELECT
			CASE
				WHEN user_one_id = $1::uuid THEN user_two_id
				ELSE user_one_id
			END AS requestee_id,
			created_at
		FROM
			friendships
		WHERE
			(user_one_id = $1::uuid OR user_two_id = $1::uuid)
			AND
			action_user_id = $1::uuid
			AND
			status = 'pending';
	`

	userID, _ := uuid.Parse(dto.Id)

	rows, err := r.pgx.Query(ctx, rawQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("outgoing friend requests: %w", err)
	}
	defer rows.Close()

	var requests []RequestDTO

	for rows.Next() {
		var result RequestDTO
		if err := rows.Scan(
			&result.UserID,
			&result.InitiatedAt,
		); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		requests = append(requests, result)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return requests, nil
}

func (r *PostgresRepository) BlockUser(ctx context.Context, dto BlockUserDTO) error {
	rawQuery := `
		INSERT INTO
			friendships(
				user_one_id,
				user_two_id,
				status,
				action_user_id
			)
		VALUES(
			LEAST($1::uuid, $2::uuid),
			GREATEST($1::uuid, $2::uuid),
			'blocked',
			$1::uuid
		)
		ON CONFLICT(user_one_id, user_two_id)
		DO UPDATE SET
			status = 'blocked',
			action_user_id = $1::uuid;
	`

	userID, _ := uuid.Parse(dto.UserID)
	targetID, _ := uuid.Parse(dto.TargetID)

	result, err := r.pgx.Exec(ctx, rawQuery, userID, targetID)
	if err != nil {
		return fmt.Errorf("blocking user: %w", err)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("blocking user: no op")
	}

	return nil
}

func (r *PostgresRepository) UnblockUser(ctx context.Context, dto UnblockUserDTO) error {
	rawQuery := `
		DELETE FROM 
			friendships
		WHERE
			user_one_id = LEAST($1::uuid, $2::uuid)
			AND
			user_two_id = GREATEST($1::uuid, $2::uuid)
			AND
			action_user_id = $1::uuid
			AND
			status = 'blocked';
	`

	userID, _ := uuid.Parse(dto.UserID)
	targetID, _ := uuid.Parse(dto.TargetID)

	result, err := r.pgx.Exec(ctx, rawQuery, userID, targetID)
	if err != nil {
		return fmt.Errorf("unblocking user: %w", err)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("unblocking user: no op")
	}

	return nil
}

func (r *PostgresRepository) GetAllFriends(ctx context.Context, dto GetAllFriendsDTO) ([]uuid.UUID, error) {
	rawQuery := `
		SELECT
			CASE
				WHEN user_one_id = $1::uuid THEN user_two_id
				ELSE user_one_id
			END AS friend_id
		FROM 
			friendships
		WHERE 
			action_user_id = $1::uuid
			AND
			status = 'accepted';
	`

	userID, _ := uuid.Parse(dto.Id)

	rows, err := r.pgx.Query(ctx, rawQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("getting all friends: %w", err)
	}
	defer rows.Close()

	var friendIDs []uuid.UUID

	for rows.Next() {
		var friendID uuid.UUID
		if err := rows.Scan(&friendID); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		friendIDs = append(friendIDs, friendID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating row: %w", err)
	}

	return friendIDs, nil
}

func (r *PostgresRepository) GetAllBlocked(ctx context.Context, dto GetAllBlockedDTO) ([]uuid.UUID, error) {
	rawQuery := `
		SELECT
			CASE
				WHEN user_one_id = $1::uuid THEN user_two_id
				ELSE user_one_id
			END AS friend_id
		FROM 
			friendships
		WHERE 
			action_user_id = $1::uuid
			AND
			status = 'blocked';
	`

	userID, _ := uuid.Parse(dto.Id)

	rows, err := r.pgx.Query(ctx, rawQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("getting blocked users: %w", err)
	}
	defer rows.Close()

	var blockedIDs []uuid.UUID

	for rows.Next() {
		var blockedID uuid.UUID
		if err := rows.Scan(&blockedID); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		blockedIDs = append(blockedIDs, blockedID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return blockedIDs, nil
}

func (r *PostgresRepository) GetMutualFriends(ctx context.Context, dto GetMutualFriendsDTO) ([]uuid.UUID, error) {
	rawQuery := `
		SELECT 
    		u1.friend_id AS mutual_friend_id
		FROM (
    			SELECT 
					CASE 
						WHEN user_one_id = $1::uuid THEN user_two_id 
						ELSE user_one_id 
					END AS friend_id
    			FROM friendships 
    			WHERE 
					(user_one_id = $1::uuid OR user_two_id = $1::uuid) 
    				AND 
					status = 'accepted'
		) u1
		INNER JOIN (
		    SELECT 
				CASE 
					WHEN user_one_id = $2::uuid THEN user_two_id 
					ELSE user_one_id 
				END AS friend_id
		    FROM 
				friendships 
		    WHERE 
				(user_one_id = $2::uuid OR user_two_id = $2::uuid) 
		      	AND 
				status = 'accepted'
		) u2 ON u1.friend_id = u2.friend_id;
	`

	userID, _ := uuid.Parse(dto.UserOneID)
	friendID, _ := uuid.Parse(dto.UserTwoID)

	rows, err := r.pgx.Query(ctx, rawQuery, userID, friendID)
	if err != nil {
		return nil, fmt.Errorf("getting mutual friends: %w", err)
	}
	defer rows.Close()

	var mutualFriendIDs []uuid.UUID

	for rows.Next() {
		var mutualFriendID uuid.UUID
		if err := rows.Scan(&mutualFriendID); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		mutualFriendIDs = append(mutualFriendIDs, mutualFriendID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return mutualFriendIDs, nil
}
