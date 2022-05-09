package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/dgrijalva/jwt-go"
	"github.com/zobinHuang/BrosCloud/backstage/auth/dal"
	"github.com/zobinHuang/BrosCloud/backstage/auth/handler"
	"github.com/zobinHuang/BrosCloud/backstage/auth/service"

	"github.com/gin-gonic/gin"
)

/*
	func: inject
	description: build layer architecture
*/
func inject(ds *dal.DataSource) (*gin.Engine, error) {
	log.Info("Injecting data sources")

	// --------------------- DAL Layer --------------------------
	userDAL := dal.NewUserDAL(ds.DB)
	tokenDAL := dal.NewTokenDAL(ds.RedisClient)

	// --------------------- Service Layer --------------------------
	userService := service.NewUserService(&service.UserServiceConfig{
		UserDAL: userDAL,
	})

	// load RSA keys
	privKeyFile := os.Getenv("RSA_PRIVATE_KEY_FILE")
	priv, err := ioutil.ReadFile(privKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read private key pem file: %w", err)
	}
	privKey, err := jwt.ParseRSAPrivateKeyFromPEM(priv)
	if err != nil {
		return nil, fmt.Errorf("could not parse private key: %w", err)
	}

	pubKeyFile := os.Getenv("RSA_PUBLIC_KEY_FILE")
	pub, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read public key pem file: %w", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	// load HMAC secret
	refreshSecret := os.Getenv("REFRESH_SECRET")

	// load expiration lengths and parse as int
	idTokenExp := os.Getenv("ID_TOKEN_EXP")
	refreshTokenExp := os.Getenv("REFRESH_TOKEN_EXP")

	idExp, err := strconv.ParseInt(idTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse ID_TOKEN_EXP as int: %w", err)
	}

	refreshExp, err := strconv.ParseInt(refreshTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse REFRESH_TOKEN_EXP as int: %w", err)
	}

	// Token Service
	tokenService := service.NewTokenService(&service.TSConfig{
		TokenDAL:              tokenDAL,
		RefreshSecret:         refreshSecret,
		PrivKey:               privKey,
		PubKey:                pubKey,
		IDExpirationSecs:      idExp,
		RefreshExpirationSecs: refreshExp,
	})

	// --------------------- Handler Layer --------------------------
	// initialize gin router
	router := gin.Default()

	// obtain base url
	baseURL := os.Getenv("AUTH_API_URL")

	// handler timeout
	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	// Handler
	handler.NewHandler(&handler.Config{
		R:               router,
		UserService:     userService,
		TokenService:    tokenService,
		BaseURL:         baseURL,
		TimeoutDuration: time.Duration(time.Duration(ht) * time.Second),
	})

	return router, nil
}
