package social

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/abhinash-kml/nova/server/realtime"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
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

func (r *PostgresRepository) AddFriend(ctx context.Context, req AddFriend_Request) error {
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

	userID, _ := uuid.Parse(req.UserID)
	targetID, _ := uuid.Parse(req.TargetID)

	_, err := r.pgx.Exec(ctx, rawQuery, userID, targetID)
	if err != nil {
		return fmt.Errorf("adding friend: %w", err)
	}

	return nil
}

func (r *PostgresRepository) RemoveFriend(ctx context.Context, req RemoveFriend_Request) error {
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

	userID, _ := uuid.Parse(req.UserID)
	targetID, _ := uuid.Parse(req.TargetID)

	result, err := r.pgx.Exec(ctx, rawQuery, userID, targetID)
	if err != nil {
		return fmt.Errorf("removing friend: %w", err)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("removing friend: no op")
	}

	return nil
}

func (r *PostgresRepository) AcceptFriendRequest(ctx context.Context, req AcceptFriendRequest_Request) error {
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

	userID, _ := uuid.Parse(req.UserID)
	targetID, _ := uuid.Parse(req.TargetID)

	result, err := r.pgx.Exec(ctx, rawQuery, targetID, userID)
	if err != nil {
		return fmt.Errorf("accept friend request: %w", err)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("accept friend request: no op")
	}

	return nil
}

func (r *PostgresRepository) RejectFriendRequest(ctx context.Context, req RejectFriendRequest_Request) error {
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

	userID, _ := uuid.Parse(req.UserID)
	targetID, _ := uuid.Parse(req.TargetID)

	result, err := r.pgx.Exec(ctx, rawQuery, targetID, userID)
	if err != nil {
		return fmt.Errorf("rejecting friend request: %w", err)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("rejecting friend request: no op")
	}

	return nil
}

func (r *PostgresRepository) GetAllIncomingFriendRequests(ctx context.Context, req GetIncomingRequestsOfUser_Request) ([]FriendRequest, error) {
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

	userID, _ := uuid.Parse(req.Id)

	rows, err := r.pgx.Query(ctx, rawQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("getting incoming requests: %w", err)
	}
	defer rows.Close()

	var requests []FriendRequest

	for rows.Next() {
		var result FriendRequest
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

func (r *PostgresRepository) GetAllOutgoingFriendRequests(ctx context.Context, req GetOutgoingRequestsOfUser_Request) ([]FriendRequest, error) {
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

	userID, _ := uuid.Parse(req.Id)

	rows, err := r.pgx.Query(ctx, rawQuery, userID)
	if err != nil {
		return nil, fmt.Errorf("outgoing friend requests: %w", err)
	}
	defer rows.Close()

	var requests []FriendRequest

	for rows.Next() {
		var result FriendRequest
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

func (r *PostgresRepository) BlockUser(ctx context.Context, req BlockUser_Request) error {
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

	userID, _ := uuid.Parse(req.UserID)
	targetID, _ := uuid.Parse(req.TargetID)

	result, err := r.pgx.Exec(ctx, rawQuery, userID, targetID)
	if err != nil {
		return fmt.Errorf("blocking user: %w", err)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("blocking user: no op")
	}

	return nil
}

func (r *PostgresRepository) UnblockUser(ctx context.Context, req UnblockUser_Request) error {
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

	userID, _ := uuid.Parse(req.UserID)
	targetID, _ := uuid.Parse(req.TargetID)

	result, err := r.pgx.Exec(ctx, rawQuery, userID, targetID)
	if err != nil {
		return fmt.Errorf("unblocking user: %w", err)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("unblocking user: no op")
	}

	return nil
}

func (r *PostgresRepository) GetAllFriends(ctx context.Context, req GetAllFriends_Request) ([]uuid.UUID, error) {
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

	userID, _ := uuid.Parse(req.Id)

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

func (r *PostgresRepository) GetAllBlocked(ctx context.Context, req GetAllBlocked_Request) ([]uuid.UUID, error) {
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

	userID, _ := uuid.Parse(req.Id)

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

func (r *PostgresRepository) GetMutualFriends(ctx context.Context, req GetMutualFriends_Request) ([]uuid.UUID, error) {
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

	userID, _ := uuid.Parse(req.UserOneID)
	friendID, _ := uuid.Parse(req.UserTwoID)

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

func (r *PostgresRepository) SendMessage(ctx context.Context, envelope realtime.Envelope) error {
	query := `
		INSERT INTO messages(
			sender_id,
			conversation_id,
			body,
			created_at
		) VALUES (
			$1,
			$2,
			$3,
			$4 
		);
	`

	var chatMessage realtime.ChatMessage
	err := json.Unmarshal(envelope.Data, &chatMessage)
	if err != nil {
		return fmt.Errorf("unmarshalling message: %w", err)
	}

	_, err = r.pgx.Exec(
		ctx,
		query,
		envelope.Header.SenderID,
		chatMessage.ChatId,
		chatMessage.Body,
		chatMessage.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("inserting message: %w", err)
	}

	return nil
}

func (r *PostgresRepository) CreateConversationSendMessage(ctx context.Context, message realtime.ChatMessage, senderID, targetID uuid.UUID) error {
	createConversationQuery := `
		INSERT INTO conversations(id, type)
		VALUES ($1, $2)
		ON CONFLICT (id) DO NOTHING;
	`

	insertMessageQuery := `
		INSERT INTO 
			messages(
				sender_id,
				conversation_id,
				body,
				created_at
			) VALUES (
				$1,
				$2,
				$3,
				$4 
			);
	`

	tx, err := r.pgx.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, createConversationQuery, message.ChatId, message.ConversationType)
	if err != nil {
		return fmt.Errorf("creating conversation: %w", err)
	}

	_, err = tx.Exec(ctx, insertMessageQuery, senderID, message.ChatId, message.Body, message.CreatedAt)
	if err != nil {
		return fmt.Errorf("insering message: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("commiting transaction: %w", err)
	}

	return nil
}

func (r *PostgresRepository) CheckIfConversationExists(ctx context.Context, conversationID uuid.UUID) (bool, error) {
	rawQuery := `
		SELECT EXISTS(
			SELECT 1 FROM conversations WHERE id = $1
		);
	`

	var exists bool

	row := r.pgx.QueryRow(ctx, rawQuery, conversationID)
	err := row.Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("scanning row: %w", err)
	}

	return exists, nil
}

func (r *PostgresRepository) CreateConversation(ctx context.Context, req CreateConversation_Request) (ConversationDetails_Response, error) {
	query := `
		INSERT INTO conversations(
			id,
			type,
			name
		)
		VALUES (
			$1,
			$2,
			$3
		)
		RETURNING
			id,
			type,
			name;
	`

	id, _ := uuid.NewV7()

	row := r.pgx.QueryRow(
		ctx,
		query,
		id,
		req.Type,
		req.Name,
	)

	var cd ConversationDetails_Response
	err := row.Scan(
		&cd.ID,
		&cd.Type,
		&cd.Name,
	)
	if err != nil {
		return ConversationDetails_Response{}, fmt.Errorf("scanning row: %w", err)
	}

	return cd, nil
}

func (r *PostgresRepository) AddConversationParticipant(ctx context.Context, req AddConversationParticipant_Request) error {
	query := `
		INSERT INTO conversation_participant(
			conversation_id,
			user_id,
			joined_at
		) VALUES (
			$1,
			$2,
			$3
		);
	`

	now := time.Now()

	_, err := r.pgx.Exec(
		ctx,
		query,
		req.ConversationID,
		req.UserID,
		now,
	)
	if err != nil {
		return fmt.Errorf("adding participant: %w", err)
	}

	return nil
}

func (r *PostgresRepository) RemoveConversationParticipant(ctx context.Context, req RemoveConversationParticipant_Request) error {
	query := `
		DELETE FROM 
			conversation_participants
		WHERE
			conversation_id = $1
			AND
			user_id IN($2);
	`

	_, err := r.pgx.Exec(
		ctx,
		query,
		req.ConversationID,
		req.UserID,
	)
	if err != nil {
		return fmt.Errorf("removing participants: %w", err)
	}

	return nil
}

func (r *PostgresRepository) GetConversationParticipants(ctx context.Context, conversationID uuid.UUID) ([]uuid.UUID, error) {
	rawQuery := `
		SELECT 
			user_id
		FROM 
			conversation_participants
		WHERE 
			conversation_id = $1;
	`

	var participants []uuid.UUID
	var err error

	rows, err := r.pgx.Query(ctx, rawQuery, conversationID)
	if err != nil {
		return nil, fmt.Errorf("ffetching rows: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var participant uuid.UUID
		err = rows.Scan(&participant)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		participants = append(participants, participant)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return participants, nil
}
