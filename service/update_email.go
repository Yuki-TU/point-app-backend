package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/auth"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/service"
	"github.com/hack-31/point-app-backend/repository"
	"github.com/jmoiron/sqlx"
)

type UpdateEmail struct {
	ExecerDB  repository.Execer
	QueryerDB repository.Queryer
	Cache     domain.Cache
	Repo      domain.UserRepo
}

func NewUpdateEmail(db *sqlx.DB, cache domain.Cache, rep domain.UserRepo) *UpdateEmail {
	return &UpdateEmail{ExecerDB: db, QueryerDB: db, Cache: cache, Repo: rep}
}

// メール本変更サービス
//
// @params temporaryEmailID 一時保存メールアドレスID
// @params confirmCode 認証コード
//
//
func (ue *UpdateEmail) UpdateEmail(ctx *gin.Context, temporaryEmailID, confirmCode string) error {
	// ユーザードメインサービス
	userService := service.NewUserService(ue.Repo)

	// コンテキストよりEmailを取得する
	email, _ := ctx.Get(auth.Email)
	stringMail := email.(string)

	// 一時メールアドレスの復元
	key := fmt.Sprintf("email:%s:%s", confirmCode, temporaryEmailID)
	temporaryEmail, err := ue.Cache.Load(ctx, key)
	if err != nil {
		return fmt.Errorf("cannot load email in cache: %w", err)
	}

	// 復元が成功したら一時メールアドレスを削除する
	if err := ue.Cache.Delete(ctx, key); err != nil {
		return fmt.Errorf("cannot delete in cache: %w", err)
	}

	// すでに登録済みのユーザーがいるか確認する
	existMail, err := userService.ExistByEmail(ctx, &ue.QueryerDB, temporaryEmail)
	if err != nil {
		return fmt.Errorf("failed to find user by mail in db: %w", err)
	}
	if existMail {
		return fmt.Errorf("exist mail address: %w", repository.ErrAlreadyEntry)
	}

	// DBに保存する
	if err := ue.Repo.UpdateEmail(ctx, ue.ExecerDB, &stringMail, &temporaryEmail); err != nil {
		return fmt.Errorf("failed to update: %w", err)
	}

	// 成功時
	return nil
}
