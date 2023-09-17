package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/dal"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/handler"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/service"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/service/servicecore"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

/*
func: inject
description: build layer architecture
*/
func inject(ds *dal.DataSource) (*gin.Engine, error) {
	log.Info("Injecting data sources")

	// --------------------- DAL Layer --------------------------
	rdbDAL := dal.NewRDbDAL(ds.DB)
	comsumerDAL := dal.NewConsumerDAL(&dal.ConsumerDALConfig{})
	providerDAL := dal.NewProviderDAL(&dal.ProviderDALConfig{
		DB: ds.DB,
	})
	depositoryDAL := dal.NewDepositoryDAL(&dal.DepositoryDALConfig{
		DB: ds.DB,
	})
	filestoreDAL := dal.NewFileStoreDAL(&dal.FileStoreDALConfig{
		DB: ds.DB,
	})
	instanceRoomDAL := dal.NewInstanceRoomDAL(&dal.InstanceRoomDALConfig{})
	applicationDAL := dal.NewApplicationDAL(&dal.ApplicationDALConfig{
		DB: ds.DB,
	})

	// --------------------- Service Core Layer --------------------------
	scheduleServiceCore := servicecore.NewScheduleServiceCore(&servicecore.ScheduleServiceCoreConfig{
		ConsumerDAL:     comsumerDAL,
		ProviderDAL:     providerDAL,
		DepositoryDAL:   depositoryDAL,
		FileStoreDAL:    filestoreDAL,
		InstanceRoomDAL: instanceRoomDAL,
		ApplicationDAL:  applicationDAL,
	})

	// --------------------- Service Layer --------------------------
	rdbService := service.NewRDbService(&service.RDbServiceConfig{
		RDbDAL: rdbDAL,
	})

	consumerService := service.NewConsumerService(&service.ConsumerServiceConfig{
		ICEServers:          `[{"urls":"stun:stun.l.google.com:19302"}]`,
		ScheduleServiceCore: scheduleServiceCore,
		ConsumerDAL:         comsumerDAL,
		ApplicationDAL:      applicationDAL,
		InstanceRoomDAL:     instanceRoomDAL,
	})

	providerService := service.NewProviderService(&service.ProviderServiceConfig{
		ICEServers:      `[{"urls":"stun:stun.l.google.com:19302"}]`,
		ProviderDAL:     providerDAL,
		InstanceRoomDAL: instanceRoomDAL,
		ConsumerDAL:     comsumerDAL,
	})

	depositoryService := service.NewDepositoryService(&service.DepositoryServiceConfig{
		DepositoryDAL: depositoryDAL,
	})

	fileStoreService := service.NewFileStoreService(&service.FileStoreServiceConfig{
		FileStoreDAL: filestoreDAL,
	})

	applicationService := service.NewApplicationService(&service.ApplicationServiceConfig{
		ApplicationDAL: applicationDAL,
		DepositoryDAL:  depositoryDAL,
		FileStoreDAL:   filestoreDAL,
	})

	// load RSA Private key
	pubKeyFile := os.Getenv("RSA_PUBLIC_KEY_FILE")
	pub, err := ioutil.ReadFile(pubKeyFile)
	if err != nil {
		return nil, fmt.Errorf("could not read public key pem file: %w", err)
	}
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pub)
	if err != nil {
		return nil, fmt.Errorf("could not parse public key: %w", err)
	}

	// load expiration lengths and parse as int
	idTokenExp := os.Getenv("ID_TOKEN_EXP")
	idExp, err := strconv.ParseInt(idTokenExp, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse ID_TOKEN_EXP as int: %w", err)
	}

	// Token Service
	tokenService := service.NewTokenService(&service.TSConfig{
		PubKey:           pubKey,
		IDExpirationSecs: idExp,
	})

	// --------------------- Handler Layer --------------------------
	// initialize gin router
	router := gin.Default()

	router.Use(func(timeout time.Duration) gin.HandlerFunc {
		return func(c *gin.Context) {
			ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
			defer cancel()

			c.Request = c.Request.WithContext(ctx)
			c.Next()

			// 检查是否超时
			if ctx.Err() == context.DeadlineExceeded {
				c.AbortWithStatusJSON(http.StatusGatewayTimeout, gin.H{
					"error": "context deadline exceeded",
				})
			}
		}
	}(100 * time.Second))

	// obtain base url
	baseURL := os.Getenv("SCHEDULER_API_URL")

	// handler timeout
	handlerTimeout := os.Getenv("HANDLER_TIMEOUT")
	ht, err := strconv.ParseInt(handlerTimeout, 0, 64)
	if err != nil {
		return nil, fmt.Errorf("could not parse HANDLER_TIMEOUT as int: %w", err)
	}

	// Handler
	handler.NewHandler(&handler.Config{
		R:                  router,
		RDbService:         rdbService,
		TokenService:       tokenService,
		ConsumerService:    consumerService,
		ProviderService:    providerService,
		DepositoryService:  depositoryService,
		FileStoreService:   fileStoreService,
		ApplicationService: applicationService,
		BaseURL:            baseURL,
		TimeoutDuration:    time.Duration(time.Duration(ht) * time.Second),
	})

	return router, nil
}
