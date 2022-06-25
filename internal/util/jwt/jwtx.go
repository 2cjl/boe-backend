package jwtx

import (
	"boe-backend/internal/db"
	"boe-backend/internal/util/config"
	"errors"
	"log"
	"strconv"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

const (
	IdentityKey     = "id"
	OrganizationKey = "org"
	appRealm        = "bankSpike"
)

var (
	authMiddleware *jwt.GinJWTMiddleware
)

type loginForm struct {
	Phone  string `form:"phone" json:"phone" binding:"required"`
	Passwd string `form:"passwd" json:"passwd" binding:"required"`
}

type RegisterForm struct {
	Username string `form:"username" json:"username" binding:"required"`
	Phone    string `form:"phone" json:"phone" binding:"required"`
	Passwd   string `form:"passwd" json:"passwd" binding:"required"`
}

// TokenUserInfo 结构体中的数据将会编码进token
type TokenUserInfo struct {
	ID             string
	OrganizationID string
}

func GetAuthMiddleware() (*jwt.GinJWTMiddleware, error) {
	cfg := config.GetConfig()

	return getAuthMiddleware([]byte(cfg.JWT.Secret), cfg.JWT.Timeout, cfg.JWT.MaxRefresh)
}

func getAuthMiddleware(secret []byte, timeout, maxRefresh time.Duration) (*jwt.GinJWTMiddleware, error) {
	var err error

	authMiddleware, err = jwt.New(&jwt.GinJWTMiddleware{
		Realm:       appRealm,
		Key:         secret,
		Timeout:     timeout,
		MaxRefresh:  maxRefresh,
		IdentityKey: IdentityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*TokenUserInfo); ok {
				return jwt.MapClaims{
					IdentityKey:     v.ID,
					OrganizationKey: v.OrganizationID,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &TokenUserInfo{
				ID:             claims[IdentityKey].(string),
				OrganizationID: claims[OrganizationKey].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginInfo loginForm
			if err := c.ShouldBind(&loginInfo); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			phone := loginInfo.Phone
			passwd := loginInfo.Passwd

			if phone == "" || passwd == "" {
				return nil, jwt.ErrFailedAuthentication
			}
			user := db.Login(phone, passwd)
			if user == nil {
				log.Println(err)
				return nil, jwt.ErrFailedAuthentication
			}
			o := db.GetOrganizationByUser(user.ID)
			u := &TokenUserInfo{
				ID: strconv.Itoa(user.ID),
			}
			if o != nil {
				u.OrganizationID = strconv.Itoa(o.ID)
			}
			return u, nil
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"message": message,
			})
		},
		TokenLookup: "header: Authorization, query: token, cookie: jwt",

		TokenHeadName: "Bearer",

		TimeFunc: time.Now,
	})
	if err != nil {
		return nil, err
	}
	err = authMiddleware.MiddlewareInit()
	if err != nil {
		return nil, errors.New("authMiddleware.MiddlewareInit() Error:" + err.Error())
	}
	return authMiddleware, nil
}

func IsValidToken(token string) bool {
	j, err := authMiddleware.ParseTokenString(token)
	if err != nil {
		return false
	}
	return j.Valid
}
