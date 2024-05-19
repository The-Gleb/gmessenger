package postgres

import (
	"bytes"
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"log/slog"
	"strings"
	"time"

	"github.com/The-Gleb/gmessenger/internal/group_service/domain/entity"
	"github.com/The-Gleb/gmessenger/internal/group_service/errors"
	"github.com/The-Gleb/gmessenger/pkg/client/postgresql"
)

type groupStorage struct {
	client postgresql.Client
}

func NewGroupStorage(client postgresql.Client) *groupStorage {
	return &groupStorage{
		client: client,
	}
}

func (s *groupStorage) Create(ctx context.Context, dto entity.CreateGroupDTO) (entity.Group, error) {

	tx, err := s.client.Begin(ctx)
	if err != nil {
		slog.Error(err.Error())
		return entity.Group{}, errors.NewDomainError(errors.ErrDB, "[storage.GetOrCreateByEmail]")
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	row := tx.QueryRow(
		ctx,
		`INSERT (name, created_at) INTO groups VALUES ($1,$2) RETURNING id, created_at;`,
		dto.Name, time.Now(),
	)

	group := entity.Group{
		Name:      dto.Name,
		MemberIDs: dto.MemberIDs,
	}
	err = row.Scan(&group.ID, &group.CreatedAt)
	if err != nil {
		slog.Error(err.Error())
		return entity.Group{}, errors.NewDomainError(errors.ErrDB, "[groupStorage.Create]")
	}

	b := bytes.Buffer{}
	for _, memberID := range dto.MemberIDs {
		b.WriteString(fmt.Sprintf("(%d,%d),", group.ID, memberID))
	}
	str := b.String()
	str = strings.TrimSuffix(str, ",")

	_, err = tx.Exec(
		ctx,
		`INSERT (group_id, user_id) INTO group_user VALUES $1;`,
	)
	if err != nil {
		slog.Error(err.Error())
		return entity.Group{}, errors.NewDomainError(errors.ErrDB, "[groupStorage.Create]")
	}

	err = tx.Commit(ctx)
	if err != nil {
		slog.Error(err.Error())
		return entity.Group{}, errors.NewDomainError(errors.ErrDB, "[groupStorage.Create]")
	}
	return group, nil
}

func (s *groupStorage) IsMember(ctx context.Context, userID, groupID int64) (bool, error) {

	tx, err := s.client.Begin(ctx)
	if err != nil {
		slog.Error(err.Error())
		return false, errors.NewDomainError(errors.ErrDB, "[groupStorage.IsMember]")
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	groupExists, err := s.Exists(ctx, tx, groupID)
	if err != nil {
		slog.Error(err.Error())
		return false, errors.NewDomainError(errors.ErrDB, err.Error())
	}

	if !groupExists {
		return false, errors.NewDomainError(errors.ErrGroupNotFound, "")
	}

	row := tx.QueryRow(
		ctx,
		`SELECT CASE WHEN EXISTS (
				SELECT * FROM members
				WHERE user_id = $1 AND group_id = $2
			)
			THEN TRUE
			ELSE FALSE END;`,
	)

	var isMember bool
	err = row.Scan(&isMember)
	if err != nil {
		slog.Error(err.Error())
		return false, errors.NewDomainError(errors.ErrDB, "[groupStorage.IsMember]")
	}
	return isMember, nil
}

func (s *groupStorage) GetMemberIDs(ctx context.Context, groupID int64) ([]int64, error) {

	rows, err := s.client.Query(
		ctx,
		`SELECT user_id FROM group_user WHERE group_id = $1;`,
		groupID,
	)

	memberIDs, err := pgx.CollectRows(rows, func(row pgx.CollectableRow) (int64, error) {
		var id int64
		err := row.Scan(&id)
		return id, err
	})
	if err != nil {
		slog.Error(err.Error())
		return nil, errors.NewDomainError(errors.ErrDB, "[groupStorage.GetMemberIDs]")
	}
	return memberIDs, nil

}

func (s *groupStorage) GetGroups(ctx context.Context, userID int64, limit, offset int) ([]entity.GroupView, error) {

	tx, err := s.client.Begin(ctx)
	if err != nil {
		slog.Error(err.Error())
		return nil, errors.NewDomainError(errors.ErrDB, "[groupStorage.IsMember]")
	}
	defer func() {
		err = tx.Rollback(ctx)
		if err != nil {
			slog.Error(err.Error())
		}
	}()

	rows, err := tx.Query(
		ctx,
		`S`,
	)

}

func (s *groupStorage) Exists(ctx context.Context, tx pgx.Tx, groupID int64) (bool, error) {

	row := tx.QueryRow(
		ctx,
		`SELECT CASE WHEN EXISTS (
				SELECT * FROM groups
				WHERE id = $1
			)
			THEN TRUE
			ELSE FALSE END;`,
		groupID,
	)

	var exists bool
	err := row.Scan(&exists)
	if err != nil {
		slog.Error(err.Error())
		return false, errors.NewDomainError(errors.ErrDB, "[groupStorage.Exists]")
	}
	return exists, nil
}
