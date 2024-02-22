// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user.sql

package sqlc

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users
(username, login, password)
VALUES
($1, $2, $3)
RETURNING login, username, password
`

type CreateUserParams struct {
	Username pgtype.Text
	Login    string
	Password pgtype.Text
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRow(ctx, createUser, arg.Username, arg.Login, arg.Password)
	var i User
	err := row.Scan(&i.Login, &i.Username, &i.Password)
	return i, err
}

const getAllUsernames = `-- name: GetAllUsernames :many
SELECT username FROM users
`

func (q *Queries) GetAllUsernames(ctx context.Context) ([]pgtype.Text, error) {
	rows, err := q.db.Query(ctx, getAllUsernames)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []pgtype.Text
	for rows.Next() {
		var username pgtype.Text
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		items = append(items, username)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getAllUsers = `-- name: GetAllUsers :many
SELECT username, login FROM users
`

type GetAllUsersRow struct {
	Username pgtype.Text
	Login    string
}

func (q *Queries) GetAllUsers(ctx context.Context) ([]GetAllUsersRow, error) {
	rows, err := q.db.Query(ctx, getAllUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetAllUsersRow
	for rows.Next() {
		var i GetAllUsersRow
		if err := rows.Scan(&i.Username, &i.Login); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPassword = `-- name: GetPassword :one
SELECT password FROM users
WHERE login = $1
`

func (q *Queries) GetPassword(ctx context.Context, login string) (pgtype.Text, error) {
	row := q.db.QueryRow(ctx, getPassword, login)
	var password pgtype.Text
	err := row.Scan(&password)
	return password, err
}

const getUser = `-- name: GetUser :one
SELECT login, username, password FROM users
WHERE login = $1
`

func (q *Queries) GetUser(ctx context.Context, login string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, login)
	var i User
	err := row.Scan(&i.Login, &i.Username, &i.Password)
	return i, err
}
