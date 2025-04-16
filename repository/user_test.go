package repository

import (
	"net/http/httptest"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/hack-31/point-app-backend/domain/model"
	"github.com/hack-31/point-app-backend/repository/entities"
	"github.com/hack-31/point-app-backend/repository/testdata"
	"github.com/hack-31/point-app-backend/testutil"
	"github.com/hack-31/point-app-backend/utils/clock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestRepository_users_GetAll(t *testing.T) {
	c := clock.FixedClocker{
		IsAsia: true,
	}
	type want struct {
		users []*entities.User
		err   error
	}
	type input struct {
		columns []string
	}

	tests := map[string]struct {
		input input
		want  want
	}{
		"全カラムのユーザー一覧を取得する": {
			input: input{
				columns: []string{},
			},
			want: want{
				users: []*entities.User{
					{
						ID:             1,
						FirstName:      "太郎",
						FirstNameKana:  "タロウ",
						FamilyName:     "本田",
						FamilyNameKana: "ホンダ",
						Email:          "honda.taro@sample.com",
						Password:       "honda.pass",
						SendingPoint:   0,
						CreatedAt:      c.Now(),
						UpdateAt:       c.Now(),
					},
					{
						ID:             2,
						FirstNameKana:  "あおい",
						FirstName:      "葵",
						FamilyName:     "斉藤",
						FamilyNameKana: "さいとう",
						Email:          "saito.aoi@example.com",
						Password:       "aoi.pass",
						SendingPoint:   100,
						CreatedAt:      c.Now(),
						UpdateAt:       c.Now(),
					},
					{
						ID:             3,
						FirstName:      "拓也",
						FirstNameKana:  "たくや",
						FamilyName:     "木村",
						FamilyNameKana: "きむら",
						Email:          "kimura.takuya@example.com",
						Password:       "kimura.pass",
						SendingPoint:   800,
						CreatedAt:      c.Now(),
						UpdateAt:       c.Now(),
					},
				},
				err: nil,
			},
		},
		"カラム指定でユーザー一覧を取得する": {
			input: input{
				columns: []string{"first_name"},
			},
			want: want{
				users: []*entities.User{
					{
						FirstName: "太郎",
					},
					{
						FirstName: "葵",
					},
					{
						FirstName: "拓也",
					},
				},
				err: nil,
			},
		},
	}
	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
			assert.NoError(t, err)

			t.Cleanup(func() {
				// 外部キー制約を有効にする
				_, err := tx.ExecContext(ctx, `SET FOREIGN_KEY_CHECKS=1`)
				assert.NoError(t, err)
				err = tx.Rollback()
				assert.NoError(t, err)
			})

			// テストデータの挿入
			testdata.Users(t, ctx, tx, func(users []*entities.User) {})

			// 実行
			r := &Repository{}
			gots, err := r.GetAll(ctx, tx, tt.input.columns...)

			// アサーション
			assert.Nil(t, err)
			assert.Equal(t, tt.want.users, gots)
		})
	}
}

func TestRepository_users_GetUserByID(t *testing.T) {
	c := clock.FixedClocker{
		IsAsia: true,
	}

	type want struct {
		user entities.User
	}
	type input struct {
		id model.UserID
	}

	tests := map[string]struct {
		input input
		want  want
	}{
		"ユーザーID_1のユーザー情報を取得する": {
			input: input{
				id: 1,
			},
			want: want{
				user: entities.User{
					ID:             1,
					FirstName:      "太郎",
					FirstNameKana:  "タロウ",
					FamilyName:     "本田",
					FamilyNameKana: "ホンダ",
					Email:          "honda.taro@sample.com",
					Password:       "honda.pass",
					SendingPoint:   0,
					CreatedAt:      c.Now(),
					UpdateAt:       c.Now(),
				},
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
			t.Cleanup(func() {
				_, err = tx.ExecContext(ctx, `SET FOREIGN_KEY_CHECKS = 1`)
				assert.NoError(t, err)
				err := tx.Rollback()
				assert.NoError(t, err)
			})

			testdata.Users(t, ctx, tx, func(users []*entities.User) {})

			// 実行
			r := &Repository{}
			gots, err := r.GetUserByID(ctx, tx, tt.input.id)

			// アサーション
			assert.Nil(t, err)
			assert.Equal(t, tt.want.user, gots)
		})
	}
}

func TestRepository_users_UpdateEmail(t *testing.T) {
	type want struct {
		email string
		err   error
	}
	type input struct {
		email string
		id    model.UserID
	}

	tests := map[string]struct {
		input input
		want  want
	}{
		"ユーザーID_1のメールを更新する": {
			input: input{
				email: "before@sample.com",
				id:    1,
			},
			want: want{
				email: "after@sample.com",
				err:   nil,
			},
		},
	}

	for n, tt := range tests {
		tt := tt
		t.Run(n, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(w)

			tx, err := testutil.OpenDBForTest(t).BeginTxx(ctx, nil)
			assert.NoError(t, err)
			t.Cleanup(func() {
				err := tx.Rollback()
				assert.NoError(t, err)
			})

			testdata.Users(t, ctx, tx, func(users []*entities.User) {
				users[0].Email = tt.input.email
			})

			// 実行
			r := &Repository{}
			err = r.UpdateEmail(ctx, tx, tt.input.id, tt.want.email)
			assert.Nil(t, err)

			// 確認
			gots, err := r.GetUserByID(ctx, tx, tt.input.id)
			assert.Nil(t, err)
			assert.Equal(t, tt.want.email, gots.Email)
		})
	}
}

func TestRepository_RegisterUser(t *testing.T) {
	t.Run("ユーザを登録し、ユーザIDを取得し、ユーザ情報ポインタにセットする", func(t *testing.T) {
		c := clock.FixedClocker{
			IsAsia: true,
		}
		// 登録データ
		okUser := &entities.User{
			FirstName:      "山田",
			FirstNameKana:  "やまだ",
			FamilyName:     "太郎",
			FamilyNameKana: "たろう",
			Password:       "pass",
			Email:          "yamada.taro@example.com",
			CreatedAt:      c.Now(),
			UpdateAt:       c.Now(),
		}

		db, mock, err := sqlmock.New()
		assert.NoError(t, err)

		t.Cleanup(func() { _ = db.Close() })

		var wantID int64 = 10
		// モックの定義
		mock.ExpectExec(
			// エスケープが必要
			`INSERT INTO users \(
			first_name, first_name_kana, family_name, family_name_kana, email, password, sending_point, created_at, update_at
		\) VALUES \(\?, \?, \?, \?, \?, \?, \?, \?, \?\)`,
		).WithArgs(
			okUser.FirstName,
			okUser.FirstNameKana,
			okUser.FamilyName,
			okUser.FamilyNameKana,
			okUser.Email,
			okUser.Password,
			okUser.SendingPoint,
			okUser.CreatedAt,
			okUser.UpdateAt,
		).WillReturnResult(sqlmock.NewResult(wantID, 1))

		w := httptest.NewRecorder()
		ctx, _ := gin.CreateTestContext(w)

		xdb := sqlx.NewDb(db, "mysql")
		r := &Repository{Clocker: c}

		err = r.RegisterUser(ctx, xdb, okUser)
		assert.NoError(t, err)
		assert.Equal(t, wantID, okUser.ID)
	})
}
