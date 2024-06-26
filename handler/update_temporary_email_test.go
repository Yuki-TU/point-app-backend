package handler

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/myerror"
	"github.com/hack-31/point-app-backend/testutil"
)

func TestUpdateTemporaryEmail(t *testing.T) {
	t.Parallel()
	type want struct {
		status  int
		rspFile string
	}

	tests := map[string]struct {
		reqFile string
		want    want
	}{
		"正しいリクエストの時は201となる": {
			reqFile: "testdata/update_temporary_email/201_req.json.golden",
			want: want{
				status:  http.StatusCreated,
				rspFile: "testdata/update_temporary_email/201_rsp.json.golden",
			},
		},
		"リクエストデータが正しくない場合は400エラーを返す": {
			reqFile: "testdata/update_temporary_email/400_req.json.golden",
			want: want{
				status:  http.StatusBadRequest,
				rspFile: "testdata/update_temporary_email/400_rsp.json.golden",
			},
		},
		"登録済みのメールアドレスは409エラーを返す": {
			reqFile: "testdata/update_temporary_email/409_req.json.golden",
			want: want{
				status:  http.StatusConflict,
				rspFile: "testdata/update_temporary_email/409_rsp.json.golden",
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			// テスト並列実行する
			t.Parallel()

			// サービス層のモック定義
			moq := &UpdateTemporaryEmailServiceMock{}
			moq.UpdateTemporaryEmailFunc = func(ctx *gin.Context, email string) (string, error) {
				if tt.want.status == http.StatusCreated {
					return "8e8d0f82-89a1-4cc6-ba25-13c864ad09db", nil
				}

				if tt.want.status == http.StatusConflict {
					return "", myerror.ErrAlreadyEntry
				}

				if tt.want.status == http.StatusBadRequest {
					return "", errors.New("email: must be a valid email address.")
				}

				return "", errors.New("error from mock")
			}

			// テストデータを挿入
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)
			c.Request, _ = http.NewRequest("POST", "/temporary_email", bytes.NewReader(testutil.LoadFile(t, tt.reqFile)))
			ute := NewUpdateTemporaryEmailHandler(moq)

			// リクエスト送信
			ute.ServeHTTP(c)

			// レスポンス
			resp := w.Result()
			// 検証
			testutil.AssertResponse(t, resp, tt.want.status, testutil.LoadFile(t, tt.want.rspFile))
		})
	}
}
