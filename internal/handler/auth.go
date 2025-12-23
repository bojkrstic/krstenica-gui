package handler

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"krstenica/internal/requestctx"
)

const (
	sessionCookieName       = "krstenica_session"
	contextUserKey          = "krstenica_authenticated_user"
	sessionDuration         = 30 * time.Minute
	sessionRefreshThreshold = 5 * time.Minute
	defaultRedirectPath     = "/ui"
	apiTokenDuration        = 30 * time.Minute
	adminRoleDefault        = "admin"
)

func (h *httpHandler) addAuthRoutes() {
	h.router.GET("/ui/login", h.renderLogin())
	h.router.POST("/ui/login", h.handleLogin())
	h.router.POST("/ui/logout", h.requireUIAuth(), h.handleLogout())
	h.router.POST("/"+routePrefix+"/auth/login", h.handleAPILogin())
}

func (h *httpHandler) renderLogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if _, _, ok := h.authenticateRequest(ctx); ok {
			target := sanitizeReturnURL(ctx.Query("return"))
			if target == "" {
				target = defaultRedirectPath
			}
			ctx.Redirect(http.StatusSeeOther, target)
			return
		}

		h.renderHTML(ctx, http.StatusOK, "auth/login.html", gin.H{
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
			h.renderHTML(ctx, http.StatusInternalServerError, "auth/login.html", gin.H{
				"Title":     "Пријава",
				"Error":     "Грешка при провери корисника.",
				"ReturnURL": returnURL,
			})
			return
		}
		if !ok {
			h.renderHTML(ctx, http.StatusUnauthorized, "auth/login.html", gin.H{
				"Title":     "Пријава",
				"Error":     "Погрешно корисничко име или лозинка.",
				"ReturnURL": returnURL,
			})
			return
		}

		token, err := h.createSessionToken(username)
		if err != nil {
			h.renderHTML(ctx, http.StatusInternalServerError, "auth/login.html", gin.H{
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

func (h *httpHandler) handleAPILogin() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "neispravan format zahteva"})
			return
		}

		ok, err := h.credentialsMatch(ctx, req.Username, req.Password)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "greska pri provjeri korisnika"})
			return
		}
		if !ok {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "pogresno korisnicko ime ili lozinka"})
			return
		}

		user, err := h.loadAuthenticatedUser(ctx.Request.Context(), req.Username)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "greska pri ucitavanju korisnika"})
			return
		}

		token, expiresAt, err := h.createJWTToken(user)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "greska pri generisanju tokena"})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"token":      token,
			"token_type": "Bearer",
			"expires_at": expiresAt.UTC().Format(time.RFC3339),
		})
	}
}

func (h *httpHandler) credentialsMatch(ctx *gin.Context, username, password string) (bool, error) {
	return h.service.AuthenticateUser(ctx.Request.Context(), username, password)
}

func (h *httpHandler) requireUIAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if user, expiresAt, ok := h.authenticateRequest(ctx); ok {
			if remaining := time.Until(expiresAt); remaining <= sessionRefreshThreshold {
				if token, err := h.createSessionToken(user.Username); err == nil {
					h.issueSessionCookie(ctx, token)
				}
			}
			h.attachAuthenticatedUser(ctx, user)
			ctx.Next()
			return
		}

		ctx.Redirect(http.StatusSeeOther, "/ui/login?return="+url.QueryEscape(defaultRedirectPath))
		ctx.Abort()
	}
}

func (h *httpHandler) requireAPIAuth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if user, ok := h.authenticateAPIRequest(ctx); ok {
			h.attachAuthenticatedUser(ctx, user)
			ctx.Next()
			return
		}

		if user, _, ok := h.authenticateRequest(ctx); ok {
			h.attachAuthenticatedUser(ctx, user)
			ctx.Next()
			return
		}

		ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
	}
}

func (h *httpHandler) authenticateAPIRequest(ctx *gin.Context) (*requestctx.User, bool) {
	authHeader := strings.TrimSpace(ctx.GetHeader("Authorization"))
	if authHeader == "" {
		return nil, false
	}

	const bearerPrefix = "Bearer "
	if len(authHeader) <= len(bearerPrefix) || !strings.EqualFold(authHeader[:len(bearerPrefix)], bearerPrefix) {
		return nil, false
	}

	token := strings.TrimSpace(authHeader[len(bearerPrefix):])
	if token == "" {
		return nil, false
	}

	user, err := h.parseJWTToken(token)
	if err != nil {
		return nil, false
	}
	return user, true
}

func (h *httpHandler) authenticateRequest(ctx *gin.Context) (*requestctx.User, time.Time, bool) {
	token, err := ctx.Cookie(sessionCookieName)
	if err != nil || token == "" {
		return nil, time.Time{}, false
	}
	username, expiresAt, err := h.parseSessionToken(token)
	if err != nil {
		return nil, time.Time{}, false
	}
	user, err := h.loadAuthenticatedUser(ctx.Request.Context(), username)
	if err != nil {
		return nil, time.Time{}, false
	}
	return user, expiresAt, true
}

func (h *httpHandler) createSessionToken(username string) (string, error) {
	expiresAt := time.Now().Add(sessionDuration).Unix()
	payload := username + "|" + strconv.FormatInt(expiresAt, 10)
	signature := h.signPayload(payload)
	raw := payload + "|" + signature
	return base64.RawURLEncoding.EncodeToString([]byte(raw)), nil
}

func (h *httpHandler) parseSessionToken(token string) (string, time.Time, error) {
	decoded, err := base64.RawURLEncoding.DecodeString(token)
	if err != nil {
		return "", time.Time{}, err
	}
	parts := strings.Split(string(decoded), "|")
	if len(parts) != 3 {
		return "", time.Time{}, errors.New("invalid session token")
	}
	username := parts[0]
	expiresUnix, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", time.Time{}, err
	}
	expiresAt := time.Unix(expiresUnix, 0)
	if h.signPayload(parts[0]+"|"+parts[1]) != parts[2] {
		return "", time.Time{}, errors.New("invalid session signature")
	}
	if time.Now().After(expiresAt) {
		return "", time.Time{}, errors.New("session expired")
	}
	return username, expiresAt, nil
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

func (h *httpHandler) createJWTToken(user *requestctx.User) (string, time.Time, error) {
	if user == nil {
		return "", time.Time{}, errors.New("user is required")
	}
	username := strings.TrimSpace(user.Username)
	if username == "" {
		return "", time.Time{}, errors.New("username is required")
	}
	role := strings.TrimSpace(user.Role)
	if role == "" {
		role = adminRoleDefault
	}
	city := strings.TrimSpace(user.City)
	secret := strings.TrimSpace(h.jwtSecret())
	if secret == "" {
		return "", time.Time{}, errors.New("jwt secret is not configured")
	}

	now := time.Now().UTC()
	expiresAt := now.Add(apiTokenDuration)

	header := map[string]string{
		"alg": "HS256",
		"typ": "JWT",
	}

	claims := map[string]interface{}{
		"sub":  username,
		"iat":  now.Unix(),
		"nbf":  now.Unix(),
		"exp":  expiresAt.Unix(),
		"role": role,
		"city": city,
		"uid":  user.ID,
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", time.Time{}, err
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", time.Time{}, err
	}

	segments := []string{
		base64.RawURLEncoding.EncodeToString(headerJSON),
		base64.RawURLEncoding.EncodeToString(claimsJSON),
	}

	signingInput := strings.Join(segments, ".")
	signature := base64.RawURLEncoding.EncodeToString(h.signJWT(signingInput, secret))

	token := signingInput + "." + signature
	return token, expiresAt, nil
}

func (h *httpHandler) parseJWTToken(token string) (*requestctx.User, error) {
	token = strings.TrimSpace(token)
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errors.New("invalid token")
	}

	secret := strings.TrimSpace(h.jwtSecret())
	if secret == "" {
		return nil, errors.New("jwt secret is not configured")
	}

	signingInput := strings.Join(parts[:2], ".")
	expectedSignature := h.signJWT(signingInput, secret)
	actualSignature, err := base64.RawURLEncoding.DecodeString(parts[2])
	if err != nil {
		return nil, errors.New("invalid token signature")
	}
	if !hmac.Equal(actualSignature, expectedSignature) {
		return nil, errors.New("invalid token signature")
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errors.New("invalid token payload")
	}

	var claims struct {
		Subject   string `json:"sub"`
		ExpiresAt int64  `json:"exp"`
		NotBefore int64  `json:"nbf"`
		Role      string `json:"role"`
		City      string `json:"city"`
		UserID    int64  `json:"uid"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, errors.New("invalid token payload")
	}
	if strings.TrimSpace(claims.Subject) == "" {
		return nil, errors.New("invalid token subject")
	}

	now := time.Now().UTC().Unix()
	if claims.NotBefore != 0 && now < claims.NotBefore {
		return nil, errors.New("token not yet valid")
	}
	if claims.ExpiresAt != 0 && now >= claims.ExpiresAt {
		return nil, errors.New("token expired")
	}

	role := strings.TrimSpace(claims.Role)
	if role == "" {
		role = adminRoleDefault
	}
	return &requestctx.User{
		ID:       claims.UserID,
		Username: strings.TrimSpace(claims.Subject),
		Role:     role,
		City:     strings.TrimSpace(claims.City),
	}, nil
}

func (h *httpHandler) jwtSecret() string {
	if h.conf == nil {
		return ""
	}
	secret := strings.TrimSpace(h.conf.AdminJWTSecret)
	if secret == "" {
		secret = strings.TrimSpace(h.conf.JWTSecret)
	}
	return secret
}

func (h *httpHandler) attachAuthenticatedUser(ctx *gin.Context, user *requestctx.User) {
	if ctx == nil || user == nil {
		return
	}
	ctx.Set(contextUserKey, user)
	ctx.Request = ctx.Request.WithContext(requestctx.WithUser(ctx.Request.Context(), user))
}

func (h *httpHandler) loadAuthenticatedUser(ctx context.Context, username string) (*requestctx.User, error) {
	username = strings.TrimSpace(username)
	if username == "" {
		return nil, errors.New("username is required")
	}
	modelUser, err := h.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return nil, err
	}
	role := strings.TrimSpace(modelUser.Role)
	if role == "" {
		role = adminRoleDefault
	}
	return &requestctx.User{
		ID:       modelUser.ID,
		Username: modelUser.Username,
		Role:     role,
		City:     strings.TrimSpace(modelUser.City),
	}, nil
}

func (h *httpHandler) currentUser(ctx *gin.Context) (*requestctx.User, bool) {
	if ctx == nil {
		return nil, false
	}
	value, ok := ctx.Get(contextUserKey)
	if !ok {
		return nil, false
	}
	user, ok := value.(*requestctx.User)
	return user, ok
}

func (h *httpHandler) requireRole(role string) gin.HandlerFunc {
	normalized := strings.ToLower(strings.TrimSpace(role))
	return func(ctx *gin.Context) {
		if user, ok := h.currentUser(ctx); ok && strings.ToLower(strings.TrimSpace(user.Role)) == normalized {
			ctx.Next()
			return
		}
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
	}
}

func (h *httpHandler) requireUIRole(role string) gin.HandlerFunc {
	normalized := strings.ToLower(strings.TrimSpace(role))
	return func(ctx *gin.Context) {
		if user, ok := h.currentUser(ctx); ok && strings.ToLower(strings.TrimSpace(user.Role)) == normalized {
			ctx.Next()
			return
		}
		h.renderHTML(ctx, http.StatusForbidden, "partials/error.html", gin.H{"Message": "Забрањен приступ"})
		ctx.Abort()
	}
}

func (h *httpHandler) signJWT(input, secret string) []byte {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write([]byte(input))
	return mac.Sum(nil)
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
