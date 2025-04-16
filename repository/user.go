package repository

import (
	"context"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/go-sql-driver/mysql"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/myerror"
	"github.com/hack-31/point-app-backend/repository/entities"
)

// ユーザ情報永続化
//
// @params
// ctx コンテキスト
// db dbの値(インスタンス)
// u ユーザエンティティ
func (r *Repository) RegisterUser(ctx context.Context, db Execer, u *entities.User) error {
	u.CreatedAt = r.Clocker.Now()
	u.UpdateAt = r.Clocker.Now()

	sql := `INSERT INTO users (
			first_name, first_name_kana, family_name, family_name_kana, email, password, sending_point, created_at, update_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	result, err := db.ExecContext(ctx, sql, u.FirstName, u.FirstNameKana, u.FamilyName, u.FamilyNameKana, u.Email, u.Password, u.SendingPoint, u.CreatedAt, u.UpdateAt)
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == ErrCodeMySQLDuplicateEntry {
			return errors.Wrap(myerror.ErrAlreadyEntry, "failed to register user")
		}
		return errors.Wrap(err, "failed to register user")
	}
	id, err := result.LastInsertId()
	if err != nil {
		return errors.Wrap(err, "failed to get last insert id in user repo")
	}
	u.ID = id
	return nil
}

// メールでユーザが存在するか検索する
// @params
// ctx context
// db dbインスタンス
// email email
// columns カラム（無指定の場合は全て）
//
// @returns
// entities.User ユーザ情報
func (r *Repository) FindUserByEmail(ctx context.Context, db Queryer, email string, columns ...string) (entities.User, error) {
	formattedColumns := "*"
	if len(columns) > 0 {
		// 以下のような文字列にする
		// id, name, email, created_at, updated_at
		formattedColumns = strings.Join(columns, ", ")
		formattedColumns = strings.TrimSuffix(formattedColumns, ", ")
	}

	sql :=
		`SELECT ` + formattedColumns + `
		FROM users 
		WHERE email = ? 
		LIMIT 1;`

	var user entities.User

	if err := db.GetContext(ctx, &user, sql, email); err != nil {
		// 見つけられない時(その他のエラーも含む)
		// 見つけられない時のエラーは利用側で
		// errors.Is(err, sql.ErrNoRows)
		// で判断する
		return user, errors.Wrap(err, "failed to find user by email in user repo")
	}
	return user, nil
}

// ユーザIDでユーザが存在するか検索する
// @params
// ctx context
// db dbインスタンス
// ID ユーザID
//
// @returns
// entities.User ユーザ情報
func (r *Repository) GetUserByID(ctx context.Context, db Queryer, ID model.UserID) (entities.User, error) {
	sql := `
		SELECT *
		FROM users
		WHERE id = ?
		LIMIT 1
	`

	var user entities.User
	if err := db.GetContext(ctx, &user, sql, ID); err != nil {
		// 見つけられない時(その他のエラーも含む)
		// 見つけられない時のエラーは利用側で
		// errors.Is(err, sql.ErrNoRows)
		// で判断する
		return user, errors.Wrap(err, "failed to get user by id in user repo")
	}
	return user, nil
}

// 削除
func (r *Repository) DeleteUserByID(ctx context.Context, db Execer, ID model.UserID) (int64, error) {
	sql := "DELETE FROM `users`" +
		"WHERE `id` = ?"
	res, err := db.ExecContext(ctx, sql, ID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to delete user by id in user repo")
	}
	return res.RowsAffected()
}

// パスワードを上書き
// @params
// ctx context
// db dbインスタンス
// email email
// pass password
//
// @returns
// error
func (r *Repository) UpdatePassword(ctx context.Context, db Execer, email, pass *string) error {
	sql := `UPDATE users SET password = ? WHERE email = ?`

	_, err := db.ExecContext(ctx, sql, pass, email)
	if err != nil {
		return errors.Wrap(err, "failed to update password in user repo")
	}
	return err
}

// メールアドレスを上書きする
// @params
// ctx context
// db dbインスタンス
// userID 更新するユーザーID
// newEmail 上書きするメールアドレス
//
// @returns
// error
func (r *Repository) UpdateEmail(ctx context.Context, db Execer, userID model.UserID, newEmail string) error {
	sql := `UPDATE users SET email = ? WHERE id = ?`

	_, err := db.ExecContext(ctx, sql, newEmail, userID)
	if err != nil {
		return errors.Wrap(err, "failed to update email in user repo")
	}
	return nil
}

// アカウント情報を上書きする
// @params
// ctx context
// db dbインスタンス
// email email
// familyName familyName
// familyNameKana familyNameKana
// firstName firstName
// firstNameKana firstNameKana
//
// @returns
// error
func (r *Repository) UpdateAccount(ctx context.Context, db Execer, email, familyName, familyNameKana, firstName, firstNameKana *string) error {
	sql := `UPDATE users
						SET family_name = ?,
								family_name_kana = ?,
								first_name = ?,
								first_name_kana = ?
						WHERE email = ?`
	_, err := db.ExecContext(
		ctx,
		sql,
		familyName,
		familyNameKana,
		firstName,
		firstNameKana,
		email)
	if err != nil {
		return errors.Wrap(err, "failed to update account in user repo")
	}
	return nil
}

// ユーザ一覧
//
// @params
// ctx context
// db db
// columns 取得カラム（無指定の場合は全て）
//
// @returns
// Users ユーザ一覧
func (r *Repository) GetAll(ctx context.Context, db Queryer, columns ...string) ([]*entities.User, error) {
	formattedColumns := "*"
	if len(columns) > 0 {
		// 以下のような文字列にする
		// id, name, email, created_at, updated_at
		formattedColumns = strings.Join(columns, ", ")
		formattedColumns = strings.TrimSuffix(formattedColumns, ", ")
	}
	sql :=
		`SELECT ` + formattedColumns + `
		 FROM users;`
	var users []*entities.User
	if err := db.SelectContext(ctx, &users, sql); err != nil {
		return users, errors.Wrap(err, "failed to get all users in user repo")
	}
	return users, nil
}

type GetAllWithCursorParam struct {
	Size            int          `db:"size"`
	CursorPoint     int          `db:"point"`
	CursorUserID    model.UserID `db:"user_id"`
	CursorCreatedAt time.Time    `db:"created_at"`
}

// GetAllWithCursor ポイント順にユーザを取得する
// カーソルページネーションを使用して取得する
// ソート順は以下の通り
//
// 1. ポイントが多い順
// 2. ポイントが同じ場合は登録日が古い順
// 3. ポイントと登録日が同じ場合はユーザIDが大きい順
func (r *Repository) GetAllWithCursor(ctx context.Context, db Queryer, param GetAllWithCursorParam) ([]*entities.User, error) {
	sql := `
		WITH points AS (
			SELECT receiving_user_id AS user_id, SUM(transaction_point) AS point
			FROM transactions
			GROUP BY receiving_user_id
		)
		SELECT u.*
		FROM users AS u
		INNER JOIN points AS p
		ON u.id = p.user_id
		WHERE 
		  p.point < ? 
			OR (p.point = ? AND u.created_at > ?) 
			OR (p.point = ? AND u.created_at = ? AND u.id > ?)
		ORDER BY p.point DESC, u.created_at ASC, u.id ASC
		LIMIT ?;`

	var users []*entities.User
	err := db.SelectContext(ctx, &users, sql,
		param.CursorPoint,
		param.CursorPoint,
		param.CursorCreatedAt,
		param.CursorPoint,
		param.CursorCreatedAt,
		param.CursorUserID,
		param.Size,
	)
	if err != nil {
		return users, errors.Wrap(err, "failed to get all users in user repo")
	}

	return users, nil
}

type GetUsersParam struct {
	Size int
}

// GetUsers ポイント順にユーザを取得する
//
// ソート順は以下の通り
//
// 1. ポイントが多い順
// 2. ポイントが同じ場合は登録日が古い順
// 3. ポイントと登録日が同じ場合はユーザIDが大きい順
func (r *Repository) GetUsers(ctx context.Context, db Queryer, param GetUsersParam) ([]*entities.User, error) {
	sql := `
		WITH points AS (
			SELECT receiving_user_id AS user_id, SUM(transaction_point) AS point
			FROM transactions
			GROUP BY receiving_user_id
		)
		SELECT u.*
		FROM users AS u
		LEFT JOIN points AS p
		ON u.id = p.user_id
		ORDER BY p.point DESC, u.created_at ASC, u.id ASC
		LIMIT ?;`

	var users []*entities.User
	err := db.SelectContext(ctx, &users, sql,
		param.Size,
	)
	if err != nil {
		return users, errors.Wrap(err, "failed to get all users in user repo")
	}

	return users, nil
}
