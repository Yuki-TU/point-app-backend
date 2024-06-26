package handler

import (
	"net/http"
	"regexp"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hack-31/point-app-backend/myerror"
)

type Signin struct {
	Service SigninService
}

func NewSigninHandler(s SigninService) *Signin {
	return &Signin{Service: s}
}

// サインインハンドラー
//
// @param ctx ginContext
func (ru *Signin) ServeHTTP(ctx *gin.Context) {
	// ユーザのパラメータ検証
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	const errTitle = "サインインエラー"
	if err := ctx.ShouldBindJSON(&input); err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error(), err)
		return
	}
	err := validation.ValidateStruct(&input,
		validation.Field(
			&input.Email,
			validation.Length(1, 256),
			validation.Required,
			is.Email,
		),
		validation.Field(
			&input.Password,
			validation.Length(8, 50),
			// TODO: 正規表現バリデーション
			validation.Match(regexp.MustCompile(``)),
			validation.Required,
		),
	)
	if err != nil {
		ErrResponse(ctx, http.StatusBadRequest, errTitle, err.Error(), err)
		return
	}

	// サインイン処理依頼
	jwt, err := ru.Service.Signin(ctx, input.Email, input.Password)
	if err != nil {
		if errors.Is(err, myerror.ErrNotMatchLogInfo) {
			ErrResponse(ctx, http.StatusUnauthorized, errTitle, myerror.ErrNotMatchLogInfo.Error(), err)
			return
		}
		ErrResponse(ctx, http.StatusInternalServerError, errTitle, err.Error(), err)
		return
	}

	// 成功レスポンスを返す
	rsp := struct {
		Token string `json:"accessToken"`
	}{Token: jwt}
	APIResponse(ctx, http.StatusCreated, "サインイン成功しました。", rsp)
}
