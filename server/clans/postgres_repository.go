package clans

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/abhinash-kml/nova/server/config"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresRepository struct {
	config           *config.Config
	logger           *zap.Logger
	pgx              *pgxpool.Pool
	statementBuilder squirrel.StatementBuilderType
	seedfile         string
}

func NewPostgresRepositoryFromPgxPool(connection *pgxpool.Pool, c *config.Config, l *zap.Logger, sfp string) *PostgresRepository {
	return &PostgresRepository{
		config:           c,
		logger:           l,
		pgx:              connection,
		statementBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		seedfile:         sfp,
	}
}

func (r *PostgresRepository) Initialize(ctx context.Context) error {
	return nil
}

func (r *PostgresRepository) Seed(ctx context.Context) {
	// Check if rows already exists, if so then dont seed
	existQuery := `SELECT count(*) FROM clans;`
	row := r.pgx.QueryRow(ctx, existQuery)

	var count int

	// Scan returned row
	if err := row.Scan(&count); err != nil {
		r.logger.Error("scanning row", zap.Error(err))
	}

	// Rows already exists, dont seed
	if count != 0 {
		r.logger.Sugar().Infof("Clans table contains %d rows, skipping seeding from file...", count)
	}

	// Table is empty, so seed from file to table
	queryBuilder := r.statementBuilder.Insert("clans").Columns(
		"id",
		"name",
		"tag",
		"description",
		"leader_id",
		"level",
		"max_members",
		"is_locked",
		"created_at",
		"updated_at",
	)

	// Check if seeds file is provided and is of json format by suffix
	if len(r.seedfile) == 0 || !strings.HasSuffix(r.seedfile, ".json") {
		r.logger.Error("Seeds file not provided or is not of json format. Skipping seeding")
	}

	// open file
	file, err := os.OpenFile(r.seedfile, os.O_RDONLY, 0755)
	if err != nil {
		r.logger.Error("Failed to open clans seeds file", zap.String("file", r.seedfile), zap.Error(err))
	}
	defer file.Close()

	// Create json decoder
	decoder := json.NewDecoder(file)
	if decoder == nil {
		r.logger.Error("Failed to create json decoder. Returned nil pointer")
	}

	// Decode
	var clans []Clan
	err = decoder.Decode(&clans)
	if err != nil {
		r.logger.Error("Failed to decode seeded clans from file", zap.Error(err))
	}

	for i := range clans {
		queryBuilder = queryBuilder.Values(clans[i].Id,
			clans[i].Name,
			clans[i].Tag,
			clans[i].Description,
			clans[i].LeaderID,
			clans[i].Level,
			clans[i].MaxMembers,
			clans[i].IsLocked,
			clans[i].CreatedAt,
			clans[i].UpdatedAt)
	}

	// Generate query
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		r.logger.Error("Failed to generate sql query with squirrel for seeding clans table", zap.Error(err))
	}

	// Execute query
	_, err = r.pgx.Exec(ctx, query, args...)
	if err != nil {
		r.logger.Error("Failed to execute query to seed clans table", zap.Error(err))
	}

	r.logger.Info("Successfully seeded clans from file", zap.Int("count", len(clans)))
}

// General operations
func (r *PostgresRepository) Add(ctx context.Context, req CreateRequest) (Clan, error) {
	createClanQuery := `INSERT INTO
							clans (
								id, 
								name, 
								tag, 
								description,
								leader_id, 
								level, 
								max_members, 
								is_locked, 
								created_at, 
								updated_at
							)
							VALUES (
								$1, 
								$2, 
								$3, 
								$4, 
								$5, 
								$6, 
								$7, 
								$8, 
								$9,
								$10
							);`

	makeLeaderQuery := `
		INSERT INTO
			clan_members (
				clan_id,
				user_id,
				role,
				joined_at
			) VALUES (
				$1,
				$2,
				'LEADER',
				$3 
			);
	`

	clanID, err := uuid.NewV7()
	if err != nil {
		return Clan{}, fmt.Errorf("generating id: %w", err)
	}
	now := time.Now()

	tx, err := r.pgx.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return Clan{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(
		ctx,
		createClanQuery,
		clanID,
		req.Name,
		req.Tag,
		req.Description,
		req.LeaderID,
		req.Level,
		req.MaxMembers,
		req.IsLocked,
		now,
		now,
	)
	if err != nil {
		return Clan{}, fmt.Errorf("creating clan: %w", err)
	}

	_, err = tx.Exec(
		ctx,
		makeLeaderQuery,
		clanID,
		req.LeaderID,
		now,
	)
	if err != nil {
		return Clan{}, fmt.Errorf("making leader: %w", err)
	}
	err = tx.Commit(ctx)
	if err != nil {
		return Clan{}, fmt.Errorf("commiting transaction: %w", err)
	}

	return Clan{
		Id:          clanID,
		Name:        req.Name,
		Tag:         req.Tag,
		Description: req.Description,
		LeaderID:    req.LeaderID,
		Level:       req.Level,
		MaxMembers:  req.MaxMembers,
		IsLocked:    req.IsLocked,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

func (r *PostgresRepository) GetById(ctx context.Context, id uuid.UUID) (Clan, error) {
	rawQuery := `SELECT
					id,
					name,
					tag,
					description,
					leader_id,
					level,
					max_members,
					is_locked,
					created_at,
					updated_at
				FROM
					clans
				WHERE
					id = $1;`

	var clan Clan

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery, id)

	// Scan returned row
	err := row.Scan(
		&clan.Id,
		&clan.Name,
		&clan.Tag,
		&clan.Description,
		&clan.LeaderID,
		&clan.Level,
		&clan.MaxMembers,
		&clan.IsLocked,
		&clan.CreatedAt,
		&clan.UpdatedAt,
	)
	if err != nil {
		return Clan{}, fmt.Errorf("scanning row: %w", err)
	}

	return clan, nil
}

func (r *PostgresRepository) GetByName(ctx context.Context, name string) (Clan, error) {
	rawQuery := `SELECT
					id,
					name,
					tag,
					description,
					leader_id,
					level,
					max_members,
					is_locked,
					created_at,
					updated_at
				FROM
					clans
				WHERE
					name = $1;`

	var clan Clan

	// Execute query
	row := r.pgx.QueryRow(ctx, rawQuery, name)

	// Scan returned row
	err := row.Scan(
		&clan.Id,
		&clan.Name,
		&clan.Tag,
		&clan.Description,
		&clan.LeaderID,
		&clan.Level,
		&clan.MaxMembers,
		&clan.IsLocked,
		&clan.CreatedAt,
		&clan.UpdatedAt,
	)
	if err != nil {
		return Clan{}, fmt.Errorf("scanning row: %w", err)
	}

	return clan, nil
}

func (r *PostgresRepository) GetAll(ctx context.Context, cursor uuid.UUID, limit int) ([]Clan, error) {
	// Create and execute query
	baseQuery := `SELECT 
					id,
					name,
					tag,
					description,
					leader_id,
					level,
					max_members,
					is_locked,
					created_at,
					updated_at
				FROM 
					clans
				WHERE
					id > $1`

	addition := fmt.Sprintf(" ORDER BY name ASC LIMIT %d", limit)
	finalQuery := baseQuery + addition

	rows, err := r.pgx.Query(ctx, finalQuery, cursor)
	if err != nil {
		return nil, fmt.Errorf("getting clans: %w", err)
	}
	defer rows.Close()

	// Scan returned rows
	var clans []Clan
	for rows.Next() {
		var clan Clan
		err = rows.Scan(
			&clan.Id,
			&clan.Name,
			&clan.Tag,
			&clan.Description,
			&clan.LeaderID,
			&clan.Level,
			&clan.MaxMembers,
			&clan.IsLocked,
			&clan.CreatedAt,
			&clan.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		clans = append(clans, clan)
	}

	return clans, err
}

func (r *PostgresRepository) Update(ctx context.Context, req UpdateRequest) (Clan, error) {
	builder := r.statementBuilder.Update("clans").Where("id", req.ClanId).Suffix("RETURNING *")
	if req.Name != nil {
		builder = builder.Set("name", *req.Name)
	}
	if req.Description != nil {
		builder = builder.Set("description", *req.Description)
	}
	if req.Tag != nil {
		builder = builder.Set("tag", *req.Tag)
	}
	if req.IsLocked != nil {
		builder = builder.Set("is_locked", *req.IsLocked)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return Clan{}, fmt.Errorf("building query: %w", err)
	}

	var clan Clan

	row := r.pgx.QueryRow(ctx, query, args...)
	err = row.Scan(
		&clan.Id,
		&clan.Name,
		&clan.Tag,
		&clan.Description,
		&clan.LeaderID,
		&clan.Level,
		&clan.IsLocked,
		&clan.CreatedAt,
		&clan.UpdatedAt,
	)
	if err != nil {
		return Clan{}, fmt.Errorf("scanning row: %w", err)
	}

	return clan, nil
}

func (r *PostgresRepository) Delete(ctx context.Context, req DeleteRequest) (DeleteResponse, error) {
	deleteClanQuery := `DELETE FROM
							clans
						WHERE
							id = $1
						RETURNING
							id;`

	deleteMembersQuery := `
		DELETE FROM
			clan_members
		WHERE
			clan_id = $1
	`

	tx, err := r.pgx.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return DeleteResponse{DeletedClanID: uuid.Nil}, fmt.Errorf("begin transanction: %w", err)
	}
	defer tx.Rollback(ctx)

	var deletedID uuid.UUID
	row := tx.QueryRow(ctx, deleteClanQuery, req.ClanId)
	err = row.Scan(&deletedID)
	if err != nil {
		return DeleteResponse{DeletedClanID: uuid.Nil}, fmt.Errorf("scanning row: %w", err)
	}

	_, err = tx.Exec(ctx, deleteMembersQuery)
	if err != nil {
		return DeleteResponse{DeletedClanID: uuid.Nil}, fmt.Errorf("deleting clan members: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return DeleteResponse{DeletedClanID: uuid.Nil}, fmt.Errorf("commit transaction: %w", err)
	}

	return DeleteResponse{DeletedClanID: deletedID}, nil
}

func (r *PostgresRepository) AddMember(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	rawQuery := `
		INSERT INTO
			clan_members(
				clan_id,
				user_id,
				role
			) VALUES (
				$1,
				$2,
				'MEMBER' 
			);
	`

	_, err := r.pgx.Exec(ctx, rawQuery, req.ClanID, req.TargetID)
	if err != nil {
		return ActionResponse{}, fmt.Errorf("adding member: %w", err)
	}

	return ActionResponse{Status: 1, Message: "added"}, nil
}

func (r *PostgresRepository) RemoveMember(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	rawQuery := `
		DELETE FROM 
			clan_members 
		WHERE
			clan_id = $1
			AND
			user_id = $2; 
	`
	_, err := r.pgx.Exec(ctx, rawQuery, req.ClanID, req.TargetID)
	if err != nil {
		return ActionResponse{}, fmt.Errorf("removing member: %w", err)
	}

	return ActionResponse{Status: 1, Message: "removed"}, nil
}

func (r *PostgresRepository) ListAllOutogingInvites(ctx context.Context, req ActionRequest) ([]OutgoingInvite, error) {
	rawQuery := `
		SELECT
			initiator_id,
			target_id
		FROM
			clan_operations
		WHERE
			clan_id = $1;
	`

	var invites []OutgoingInvite

	rows, err := r.pgx.Query(ctx, rawQuery, req.ClanID)
	if err != nil {
		return nil, fmt.Errorf("getting invites: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var invite OutgoingInvite
		err = rows.Scan(&invite.UserID, &invite.TargetID)
		if err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		invites = append(invites, invite)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return invites, nil
}

func (r *PostgresRepository) InviteMember(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	rawQuery := `
		INSERT INTO 
			clan_operations(
				clan_id,
				initiator_id,
				target_id
			) VALUES (
			 	$1,
				$2,
				$3
			);
	`

	_, err := r.pgx.Exec(ctx, rawQuery, req.ClanID, req.InitiatorID, req.TargetID)
	if err != nil {
		return ActionResponse{Status: 0, Message: "failed"}, fmt.Errorf("inviting member: %w", err)
	}

	return ActionResponse{Status: 1, Message: "invited"}, nil
}

func (r *PostgresRepository) AcceptInvite(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	removeOperationQuery := `
		DELETE FROM
			clan_operations
		WHERE
			clan_id = $1
			AND
			target_id = $2;
	`

	addMemberQuery := `
		INSERT INTO
			clan_members(
				clan_id,
				user_id,
				role,
				joined_at
			) VALUES (
				$1,
				$2,
				'MEMBER',
				CURRENT_TIMESTAMP 
			);
	`

	tx, err := r.pgx.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return ActionResponse{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = r.pgx.Exec(ctx, removeOperationQuery, req.ClanID, req.TargetID)
	if err != nil {
		return ActionResponse{}, fmt.Errorf("removing operation: %w", err)
	}

	_, err = r.pgx.Exec(ctx, addMemberQuery, req.ClanID, req.TargetID)
	if err != nil {
		return ActionResponse{}, fmt.Errorf("adding member: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return ActionResponse{}, fmt.Errorf("comming transaction: %w", err)
	}

	return ActionResponse{Status: 1, Message: "acceped"}, nil
}

func (r *PostgresRepository) RejectInvite(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	rawQuery := `
		DELETE FROM
			clan_operations
		WHERE
			clan_id = $1
			AND
			target_id = $2;
	`

	_, err := r.pgx.Exec(ctx, rawQuery, req.ClanID, req.TargetID)
	if err != nil {
		return ActionResponse{}, fmt.Errorf("rejecting invite: %w", err)
	}

	return ActionResponse{Status: 1, Message: "rejected"}, nil
}

func (r *PostgresRepository) ListAllJoinRequests(ctx context.Context, req ActionRequest) ([]uuid.UUID, error) {
	rawQuery := `
		SELECT
			initiator_id
		FROM
			clan_operations
		WHERE
			clan_id = $1;
	`

	rows, err := r.pgx.Query(ctx, rawQuery, req.ClanID)
	if err != nil {
		return nil, fmt.Errorf("getting join requests: %w", err)
	}
	defer rows.Close()

	var requests []uuid.UUID

	for rows.Next() {
		var request uuid.UUID
		if err = rows.Scan(&request); err != nil {
			return nil, fmt.Errorf("scanning row: %w", err)
		}

		requests = append(requests, request)
	}

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("iterating rows: %w", err)
	}

	return requests, nil
}

func (r *PostgresRepository) SendJoinRequest(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	rawQuery := `
		INSERT INTO 
			clan_operations(
				clan_id,
				initiator_id,
				target_id
			) VALUES (
				$1,
				$2,
				'01a1498d-fb97-721b-90f4-b134122fcd23' 
			);
	`

	_, err := r.pgx.Exec(ctx, rawQuery, req.ClanID, req.InitiatorID)
	if err != nil {
		return ActionResponse{}, fmt.Errorf("sending request: %w", err)
	}

	return ActionResponse{Status: 1, Message: "sent-request"}, nil
}

func (r *PostgresRepository) AcceptJoinRequest(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	removeOperationQuery := `
		DELETE FROM
			clan_operations
		WHERE
			clan_id = $1
			AND
			initiator_id = $2;
	`

	addMemberQuery := `
		INSERT INTO 
			clan_members(
				clan_id,
				user_id,
				role,
				joined_at
			) VALUES (
				$1,
				$2,
				'MEMBER',
				CURRENT_TIMESTAMP 
			);
	`

	tx, err := r.pgx.BeginTx(ctx, pgx.TxOptions{IsoLevel: pgx.ReadCommitted})
	if err != nil {
		return ActionResponse{}, fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	_, err = r.pgx.Exec(ctx, removeOperationQuery, req.ClanID, req.InitiatorID)
	if err != nil {
		return ActionResponse{}, fmt.Errorf("deleting operation: %w", err)
	}

	_, err = r.pgx.Exec(ctx, addMemberQuery, req.ClanID, req.InitiatorID)
	if err != nil {
		return ActionResponse{}, fmt.Errorf("adding member: %w", err)
	}

	err = tx.Commit(ctx)
	if err != nil {
		return ActionResponse{}, fmt.Errorf("commiting transaction: %w", err)
	}

	return ActionResponse{Status: 1, Message: "request-accepted"}, nil
}

func (r *PostgresRepository) RejectJoinRequest(ctx context.Context, req ActionRequest) (ActionResponse, error) {
	rawQuery := `
		DELETE FROM
			clan_operations
		WHERE
			clan_id = $1
			AND
			initiator_id = $2;
	`

	_, err := r.pgx.Exec(ctx, rawQuery, req.ClanID, req.InitiatorID)
	if err != nil {
		return ActionResponse{}, fmt.Errorf("rejecting request: %w", err)
	}

	return ActionResponse{Status: 1, Message: "rejected"}, nil
}
