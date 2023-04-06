package service

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/hack-31/point-app-backend/constant"
	"github.com/hack-31/point-app-backend/domain"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/domain/service"
	"github.com/hack-31/point-app-backend/repository"
	utils "github.com/hack-31/point-app-backend/utils/email"
	"github.com/jmoiron/sqlx"
)

type RegisterTemporaryEmail struct {
	DB    repository.Queryer
	Cache domain.Cache
	Repo  domain.UserRepo
}

func NewRegisterTemporaryEmail(db *sqlx.DB, cache domain.Cache, rep domain.UserRepo) *RegisterTemporaryEmail {
	return &RegisterTemporaryEmail{DB: db, Cache: cache, Repo: rep}
}

// メール仮登録サービス
//
// @params
// ctx コンテキスト
// email メールアドレス
//
// @returns
// temporaryEmailId 一時保存したメールを識別するID
func (r *RegisterTemporaryEmail) RegisterTemporaryEmail(ctx *gin.Context, email string) (string, error) {
	// ユーザードメインサービス
	userService := service.NewUserService(r.Repo)

	// 現在利用中のメールアドレスか確認
	existMail, err := userService.ExistByEmail(ctx, &r.DB, email)
	if err != nil {
		return "", err
	}
	if existMail {
		return "", fmt.Errorf("failed to register: %w", repository.ErrAlreadyEntry)
	}

	// キャッシュサーバーに保存するkeyの作成
	uid := uuid.New().String()
	confirmCode := model.NewConfirmCode().String()
	key := fmt.Sprintf("%s:%s", confirmCode, uid)
	// キャッシュサーバーへ保存
	err = r.Cache.Save(ctx, key, email, time.Duration(constant.ConfirmationCodeExpiration_m))
	if err != nil {
		return "", fmt.Errorf("failed to save in cache: %w", err)
	}

	// メール送信
	subject := "【ポイントアプリ】本登録を完了してください"
	body := fmt.Sprintf("ポイントアプリをご利用いただきありがとうございます。\n\n確認コードは %s です。\n\nこの確認コードの有効期限は1時間です。", confirmCode)
	_, err = utils.SendMail(email, subject, body)
	if err != nil {
		return "", fmt.Errorf("failed to send email: %w", err)
	}

	return uid, nil
}
