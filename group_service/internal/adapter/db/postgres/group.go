package postgres

import (
	"context"
	"log/slog"

	"github.com/The-Gleb/gmessenger/app/pkg/client/postgresql"
	"github.com/The-Gleb/gmessenger/group_service/internal/adapter/db/sqlc"
	"github.com/The-Gleb/gmessenger/group_service/internal/domain/entity"
	"github.com/The-Gleb/gmessenger/group_service/internal/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type groupStorage struct {
	client postgresql.Client
	sqlc   *sqlc.Queries
}

func NewGroupStorage(client postgresql.Client) *groupStorage {
	return &groupStorage{
		client: client,
		sqlc:   sqlc.New(client),
	}
}

func (s *groupStorage) Create(ctx context.Context, group entity.GroupCreate) (entity.Group, error) {

	dbTx, err := s.client.Begin(ctx)
	if err != nil {
		return entity.Group{}, err
	}
	defer dbTx.Rollback(ctx) //nolint:all

	sqlcTx := s.sqlc.WithTx(dbTx)

	sqlcGroup, err := sqlcTx.CreateGroup(ctx, sqlc.CreateGroupParams{
		Name: pgtype.Text{
			String: group.Name,
			Valid:  true,
		},
		CreatedAt: pgtype.Timestamp{
			Time:  group.CreatedAt,
			Valid: true,
		},
	})
	if err != nil {
		return entity.Group{}, errors.NewDomainError(errors.ErrDB, err.Error())
	}

	for _, member := range group.MembersLogins {
		err := sqlcTx.AddMember(ctx, sqlc.AddMemberParams{
			MemberLogin: pgtype.Text{
				String: member,
				Valid:  true,
			},
			GroupID: pgtype.Int8{
				Int64: sqlcGroup.ID,
				Valid: true,
			},
		})

		if err != nil {
			return entity.Group{}, errors.NewDomainError(errors.ErrDB, err.Error())
		}
	}

	err = dbTx.Commit(ctx)
	if err != nil {
		slog.Error(err.Error())
		return entity.Group{}, errors.NewDomainError(errors.ErrDB, err.Error())
	}

	return entity.Group{
		ID:        sqlcGroup.ID,
		Name:      sqlcGroup.Name.String,
		CreatedAt: sqlcGroup.CreatedAt.Time,
	}, nil
}

func (s *groupStorage) IsMember(ctx context.Context, userLogin string, groupID int64) (bool, error) {

	groupExists, err := s.sqlc.Exists(ctx, groupID)
	if err != nil {
		slog.Error(err.Error())
		return false, errors.NewDomainError(errors.ErrDB, err.Error())
	}

	if !groupExists {
		return false, errors.NewDomainError(errors.ErrGroupNotFound, "")
	}

	isMember, err := s.sqlc.IsMember(ctx, sqlc.IsMemberParams{
		GroupID: pgtype.Int8{
			Int64: groupID,
			Valid: true,
		},
		MemberLogin: pgtype.Text{
			String: userLogin,
			Valid:  true,
		},
	})

	if err != nil {
		slog.Error(err.Error())
		return false, errors.NewDomainError(errors.ErrDB, err.Error())
	}

	if !isMember {
		return false, nil
	}

	return true, nil

}

func (s *groupStorage) GetMembers(ctx context.Context, groupID int64) ([]string, error) {

	sqlcMembers, err := s.sqlc.GetMembers(ctx, pgtype.Int8{
		Int64: groupID,
		Valid: true,
	})
	if err != nil {
		slog.Error(err.Error())
		return nil, errors.NewDomainError(errors.ErrDB, err.Error())
	}

	memberLogins := make([]string, len(sqlcMembers))

	for i, sqlcMember := range sqlcMembers {
		memberLogins[i] = sqlcMember.String
	}

	return memberLogins, nil

}

func (s *groupStorage) GetGroups(ctx context.Context, userLogin string, limit, offset int) ([]entity.Group, error) {

	sqlcGroups, err := s.sqlc.GetGroups(ctx, pgtype.Text{
		String: userLogin,
		Valid:  true,
	})
	if err != nil && err != pgx.ErrNoRows {
		slog.Error(err.Error())
		return nil, errors.NewDomainError(errors.ErrDB, err.Error())
	}

	groups := make([]entity.Group, len(sqlcGroups))

	for i, sqlcGroup := range sqlcGroups {
		groups[i] = entity.Group{
			ID:        sqlcGroup.ID,
			Name:      sqlcGroup.Name.String,
			CreatedAt: sqlcGroup.CreatedAt.Time,
		}
	}

	return groups, nil

}

func (s *groupStorage) Exists(ctx context.Context, groupID int64) (bool, error) {

	groupExists, err := s.sqlc.Exists(ctx, groupID)
	if err != nil {
		slog.Error(err.Error())
		return false, errors.NewDomainError(errors.ErrDB, err.Error())
	}

	if !groupExists {
		return false, nil
	}

	return true, nil

}
