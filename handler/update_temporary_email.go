package handler

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hack-31/point-app-backend/myerror"
)

type UpdateTemporaryEmail struct {
	Service UpdateTemporaryEmailService
}

func NewUpdateTemporaryEmailHandler(s UpdateTemporaryEmailService) *UpdateTemporaryEmail {
	return &UpdateTemporaryEmail{Service: s}
}

// メール仮登録ハンドラー
//
// @param ctx ginContext
func (ute *UpdateTemporaryEmail) ServeHTTP(ctx *gin.Context) {
	const mailErrTitle = "メールアドレス仮登録エラー"
	const paramErrTitle = "パラメータエラー"

	var input struct {
		Email string `json:"email"`
	}

	// ユーザーから正しいパラメータでポストデータが送られていない
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, paramErrTitle, err.Error(), err)
		return
	}

	// バリデーション検証
	err := validation.ValidateStruct(&input,
		validation.Field(
			&input.Email,
			validation.Required,
			validation.Length(1, 256),
			is.Email,
		))
	if err != nil {
		ErrResponse(ctx, http.StatusBadRequest, paramErrTitle, err.Error(), err)
		return
	}

	// サービス層にメール仮登録処理を依頼
	temporaryEmailID, err := ute.Service.UpdateTemporaryEmail(ctx, input.Email)

	// エラーレスポンスを返す
	if err != nil {
		if errors.Is(err, myerror.ErrAlreadyEntry) {
			ErrResponse(ctx, http.StatusConflict, mailErrTitle, myerror.ErrAlreadyEntry.Error(), err)
			return
		}
		ErrResponse(ctx, http.StatusInternalServerError, "サーバーエラー", err.Error(), err)
		return
	}

	// 成功時のレスポンスを返す
	rsp := struct {
		TemporaryEmailID string `json:"temporaryEmailId"`
	}{TemporaryEmailID: temporaryEmailID}
	APIResponse(ctx, http.StatusCreated, "指定のメールアドレスに確認コードを送信しました。", rsp)
}
