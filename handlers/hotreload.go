package handlers

import (
	"net/http"
	"sync"
	"time"

	datastar "github.com/starfederation/datastar/sdk/go"
)

var once sync.Once

func HotReloadSSE(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)

	once.Do(func() {
		sse.ExecuteScript(
			"window.location.reload()",
			datastar.WithExecuteScriptRetryDuration(time.Second))
	})

	<- r.Context().Done()
}
