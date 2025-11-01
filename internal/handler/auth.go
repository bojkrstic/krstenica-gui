package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	sessionCookieName   = "krstenica_session"
	contextUserKey      = "krstenica_authenticated_user"
	sessionDuration     = 24 * time.Hour
	defaultRedirectPath = "/ui"
)

func (h *httpHandler) addAuthRoutes() {
	h.router.GET("/ui/login", h.renderLogin())
	h.router.POST("/ui/login", h.handleLogin())
	h.router.POST("/ui/logout", h.requireUIAuth(), h.handleLogout())
}

func (h *httpHandler) renderLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, ok := h.authenticateRequest(ctx); ok {
			target := sanitizeReturnURL(ctx.Query("return"))
			if target == "" {
				target = defaultRedirectPath
			}
			ctx.Redirect(http.StatusSeeOther, target)
			return
		}

		ctx.HTML(http.StatusOK, "auth/login.html", gin.H{
			"Title":     "Пријава",
			"ReturnURL": sanitizeReturnURL(ctx.Query("return")),
		})
	}
}

func (h *httpHandler) handleLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		username := strings.TrimSpace(ctx.PostForm("username"))
		password := ctx.PostForm("password")
		returnURL := sanitizeReturnURL(ctx.PostForm("return"))

		ok, err := h.credentialsMatch(ctx, username, password)
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "auth/login.html", gin.H{
				"Title":     "Пријава",
				"Error":     "Грешка при провери корисника.",
				"ReturnURL": returnURL,
			})
			return
		}
		if !ok {
			ctx.HTML(http.StatusUnauthorized, "auth/login.html", gin.H{
				"Title":     "Пријава",
				"Error":     "Погрешно корисничко име или лозинка.",
				"ReturnURL": returnURL,
			})
			return
		}

		token, err := h.createSessionToken(username)
		if err != nil {
			ctx.HTML(http.StatusInternalServerError, "auth/login.html", gin.H{
				"Title": "Пријава",
				"Error": "Грешка при креирању сесије. Покушајте поново.",
			})
			return
		}

		h.issueSessionCookie(ctx, token)
		if returnURL == "" {
			returnURL = defaultRedirectPath
		}
		ctx.Redirect(http.StatusSeeOther, returnURL)
	}
}

func (h *httpHandler) handleLogout() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h.clearSessionCookie(ctx)
		ctx.Redirect(http.StatusSeeOther, "/ui/login")
	}
}

func (h *httpHandler) credentialsMatch(ctx *gin.Context, username, password string) (bool, error) {
	return h.service.AuthenticateUser(ctx.Request.Context(), username, password)
}

func (h *httpHandler) requireUIAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if username, ok := h.authenticateRequest(ctx); ok {
			ctx.Set(contextUserKey, username)
			ctx.Next()
			return
		}

		target := ctx.Request.URL.Path
		if raw := ctx.Request.URL.RawQuery; raw != "" {
			target += "?" + raw
		}
		ctx.Redirect(http.StatusSeeOther, "/ui/login?return="+url.QueryEscape(target))
		ctx.Abort()
	}
}

func (h *httpHandler) requireAPIAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if username, ok := h.authenticateRequest(ctx); ok {
			ctx.Set(contextUserKey, username)
			ctx.Next()
			return
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
}

func (h *httpHandler) authenticateRequest(ctx *gin.Context) (string, bool) {
	token, err := ctx.Cookie(sessionCookieName)
	if err != nil || token == "" {
		return "", false
	}
	username, err := h.parseSessionToken(token)
	if err != nil {
		return "", false
	}
	return username, true
}

func (h *httpHandler) createSessionToken(username string) (string, error) {
	expiresAt := time.Now().Add(sessionDuration).Unix()
	payload := username + "|" + strconv.FormatInt(expiresAt, 10)
	signature := h.signPayload(payload)
	raw := payload + "|" + signature
	return base64.RawURLEncoding.EncodeToString([]byte(raw)), nil
}

func (h *httpHandler) parseSessionToken(token string) (string, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", err
	}
	parts := strings.Split(string(decoded), "|")
	if len(parts) != 3 {
		return "", errors.New("invalid session token")
	}
	username := parts[0]
	expiresUnix, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", err
	}
	if h.signPayload(parts[0]+"|"+parts[1]) != parts[2] {
		return "", errors.New("invalid session signature")
	}
	if time.Now().Unix() > expiresUnix {
		return "", errors.New("session expired")
	}
	return username, nil
}

func (h *httpHandler) issueSessionCookie(ctx *gin.Context, token string) {
	ctx.SetCookie(sessionCookieName, token, int(sessionDuration.Seconds()), "/", "", h.isSecureRequest(ctx), true)
}

func (h *httpHandler) clearSessionCookie(ctx *gin.Context) {
	ctx.SetCookie(sessionCookieName, "", -1, "/", "", h.isSecureRequest(ctx), true)
}

// isSecureRequest inspects the request/proxy headers to decide if cookies need the Secure flag.
func (h *httpHandler) isSecureRequest(ctx *gin.Context) bool {
	if ctx.Request.TLS != nil {
		return true
	}

	if proto := strings.TrimSpace(ctx.GetHeader("X-Forwarded-Proto")); strings.EqualFold(proto, "https") {
		return true
	}

	if forwarded := ctx.GetHeader("Forwarded"); forwarded != "" {
		for _, entry := range strings.Split(forwarded, ",") {
			for _, part := range strings.Split(entry, ";") {
				fragment := strings.TrimSpace(part)
				if fragment == "" {
					continue
				}
				keyValue := strings.SplitN(fragment, "=", 2)
				if len(keyValue) != 2 {
					continue
				}
				if strings.EqualFold(strings.TrimSpace(keyValue[0]), "proto") &&
					strings.EqualFold(strings.Trim(strings.TrimSpace(keyValue[1]), `"`), "https") {
					return true
				}
			}
		}
	}

	host := strings.TrimSpace(h.conf.Host)
	return strings.HasPrefix(strings.ToLower(host), "https")
}

func (h *httpHandler) signPayload(payload string) string {
	secret := strings.TrimSpace(h.conf.Auth.SessionSecret)
	if secret == "" {
		secret = strings.TrimSpace(h.conf.AdminJWTSecret)
	}
	if secret == "" {
		secret = "krstenica-session-secret"
	}
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(payload))
	return hex.EncodeToString(mac.Sum(nil))
}

func sanitizeReturnURL(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if strings.HasPrefix(value, "//") {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(value), "http://") || strings.HasPrefix(strings.ToLower(value), "https://") {
		return ""
	}
	if !strings.HasPrefix(value, "/") {
		return ""
	}
	return value
}
