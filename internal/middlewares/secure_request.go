package middlewares

import (
	"net/http"
	"time"
	"fmt"
	"github.com/spartaut/utils/middlewares"
	"github.com/spartaut/utils/crypto"
	"github.com/yudzmaestro/test-api-microservice/pkg/types"
)

func MakeSecureHttpRequestMiddleware(authToken *types.AuthToken) middlewares.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ts := time.Now().Format("2006-01-02 15:04:05")
			r.Header.Set("X-IP-CHECK", "1")
			r.Header.Set("Authorization", fmt.Sprintf("token=%s", authToken.Token))
			r.Header.Set("datetime", ts)
			r.Header.Set("signature", generateSignature(authToken.Token, ts, authToken.Key))

			next.ServeHTTP(w, r)

			r.Header.Del("signature")
		})
	}
}

func generateSignature(token string, datetime string, key string) string {
	return crypto.EncodeSHA256HMAC(key, token, datetime)
}