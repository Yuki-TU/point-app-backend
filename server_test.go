package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestServer_Run(t *testing.T) {
	// 返り値を設定
	ctx, cancel := context.WithCancel(context.Background())
	eg, ctx := errgroup.WithContext(ctx)
	mux := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	})
	const port = 14280

	// テスト対象関数を呼び出し
	eg.Go(func() error {
		s := NewServer(mux, fmt.Sprintf(":%d", port))
		return s.Run(ctx)
	})

	// HACK: サーバーが起動してから、リクエストを送らないたとエラーになるた50ms待つ
	// 一時的な処理
	timer1 := time.NewTimer(50 * time.Millisecond)
	<-timer1.C

	// GETAPIをリクエスト
	in := "healthcheck"
	url := fmt.Sprintf("http://localhost:%d/%s", port, in)
	t.Logf("try request to %q", url)
	rsp, err := http.Get(url)
	assert.NoError(t, err)

	defer rsp.Body.Close()
	// レスポンス整形
	got, err := io.ReadAll(rsp.Body)
	assert.NoError(t, err, "failed to read body")

	// サーバの終了動作を検証する
	cancel()
	err = eg.Wait()
	assert.NoError(t, err)

	// 戻り値を検証する
	want := fmt.Sprintf("Hello, %s!", in)
	assert.Equal(t, want, string(got))
}
