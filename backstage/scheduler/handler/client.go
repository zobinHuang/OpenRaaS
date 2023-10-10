package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
)

var upGrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

/*
@func: WSConnect
@description:

	handler for endpoint "/api/scheduler/wsconnect"
*/
func (h *Handler) WSConnect(c *gin.Context) {
	// extract client type from url
	clientType, ok := c.GetQuery("type")
	if !ok {
		log.WithFields(
			log.Fields{
				"Client Address": c.Request.Host,
			}).Warn("Failed to extract client type, invalid websocket connection request, abandoned")
		return
	}
	if clientType != model.CLIENT_TYPE_PROVIDER &&
		clientType != model.CLIENT_TYPE_CONSUMER &&
		clientType != model.CLIENT_TYPE_DEPOSITARY &&
		clientType != model.CLIENT_TYPE_FILESTORE {
		log.WithFields(log.Fields{
			"Given Client Type": clientType,
			"Client Address":    c.Request.Host,
		}).Warn("Unknown client type, abandoned")
		return
	}

	// upgrade to websocket connection
	ws, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.WithFields(log.Fields{
			"Client Address": c.Request.Host,
			"error":          err,
		}).Warn("Failed to upgrade to websocket connection, abandoned")
		return
	}

	ctx := c.Request.Context()

	switch clientType {
	case model.CLIENT_TYPE_CONSUMER:
		// create consumer instance and start to serve it
		consumer, err := h.ConsumerService.CreateConsumer(ctx, ws)
		if err != nil {
			return
		}

		// register receive callbacks based on websocket type
		h.ConsumerService.InitRecvRoute(ctx, consumer)

	case model.CLIENT_TYPE_PROVIDER:
		uuid, ok := c.GetQuery("uuid")
		if !ok {
			log.WithFields(log.Fields{
				"Client Address": c.Request.Host,
			}).Error("Failed to extract uuid, invalid websocket connection request, abandoned")
			return
		}
		// create provider instance and start to serve it
		provider, err := h.ProviderService.CreateProvider(ctx, ws, uuid)
		if err != nil {
			return
		}

		// register receive callbacks based on websocket type
		h.ProviderService.InitRecvRoute(ctx, provider)

	case model.CLIENT_TYPE_DEPOSITARY:
		// not need now

	case model.CLIENT_TYPE_FILESTORE:
		// not need now

	default:
		// leave empty
	}
}

// NodeOnline register node to rds
func (h *Handler) NodeOnline(c *gin.Context) {
	nodeType, ok := c.GetQuery("type")
	if !ok {
		log.WithFields(
			log.Fields{
				"Node Address": c.Request.Host,
			}).Warn("Failed to extract node type, invalid websocket connection request, abandoned")
		return
	}
	if nodeType != model.CLIENT_TYPE_PROVIDER &&
		nodeType != model.CLIENT_TYPE_DEPOSITARY &&
		nodeType != model.CLIENT_TYPE_FILESTORE {
		log.WithFields(log.Fields{
			"Given Node Type": nodeType,
			"Node Address":    c.Request.Host,
		}).Warn("Unknown Node type, abandoned")
		return
	}

	log.Info("type: ", nodeType)

	ctx := c.Request.Context()

	var sc model.ScheduleServiceCore
	sc = *h.ConsumerService.GetScheduleServiceCore()

	switch nodeType {
	case model.CLIENT_TYPE_PROVIDER:
		var headerData model.ProviderCoreWithInst
		if err := c.BindJSON(&headerData); err != nil {
			log.WithFields(
				log.Fields{
					"error": err.Error(),
				}).Error("Failed to extract provider json data")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := h.ProviderService.CreateProviderInRDS(ctx, &headerData); err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Error("Failed to register provider to rds")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		h.ProviderService.ShowEnterInfo(ctx, &headerData)
		h.ProviderService.ShowAllInfo(ctx)

		go func() {
			providerCoreWithInst := &headerData
			if s, err := json.Marshal(providerCoreWithInst); err != nil {
				log.WithFields(
					log.Fields{
						"error":      err.Error(),
						"headerData": headerData,
					}).Warn("Failed to Marshal providerCoreWithInst")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			} else {
				err := sc.SetValueToBlockchain(headerData.ID, string(s))
				if err != nil {
					log.WithFields(
						log.Fields{
							"error":      err.Error(),
							"headerData": headerData,
						}).Warn("providerCoreWithInst Failed to SetValueToBlockchain")
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				} else {
					log.Infof("服务提供节点数字资产发布, id: %s, value: %s", headerData.ID, s)
				}
			}
		}()

	case model.CLIENT_TYPE_DEPOSITARY:
		var headerData model.DepositoryCoreWithInst
		if err := c.BindJSON(&headerData); err != nil {
			log.WithFields(
				log.Fields{
					"error": err.Error(),
				}).Warn("Failed to extract depository json data")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := h.DepositoryService.CreateDepositoryInRDS(ctx, &headerData); err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("Failed to register depository to rds")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		h.DepositoryService.ShowEnterInfo(ctx, &headerData)
		h.DepositoryService.ShowAllInfo(ctx)

		go func() {
			depositoryCoreWithInst := &headerData
			if s, err := json.Marshal(depositoryCoreWithInst); err != nil {
				log.WithFields(
					log.Fields{
						"error":      err.Error(),
						"headerData": headerData,
					}).Warn("Failed to Marshal depositoryCoreWithInst")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			} else {
				err := sc.SetValueToBlockchain(headerData.ID, string(s))
				if err != nil {
					log.WithFields(
						log.Fields{
							"error":      err.Error(),
							"headerData": headerData,
						}).Warn("depositoryCoreWithInst Failed to SetValueToBlockchain")
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				log.Infof("镜像仓库节点数字资产发布, id: %s, value: %s", headerData.ID, s)
			}
		}()

	case model.CLIENT_TYPE_FILESTORE:
		var headerData model.FileStoreCoreWithInst
		if err := c.BindJSON(&headerData); err != nil {
			log.WithFields(
				log.Fields{
					"error": err.Error(),
				}).Warn("Failed to extract FileStore json data")
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		if err := h.FileStoreService.CreateFileStoreInRDS(ctx, &headerData); err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("Failed to register FileStore to rds")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		h.FileStoreService.ShowEnterInfo(ctx, &headerData)
		h.FileStoreService.ShowAllInfo(ctx)

		go func() {
			fileStoreCoreWithInst := &headerData
			if s, err := json.Marshal(fileStoreCoreWithInst); err != nil {
				log.WithFields(
					log.Fields{
						"error":      err.Error(),
						"headerData": headerData,
					}).Warn("Failed to Marshal fileStoreCoreWithInst")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			} else {
				err := sc.SetValueToBlockchain(headerData.ID, string(s))
				if err != nil {
					log.WithFields(
						log.Fields{
							"error":      err.Error(),
							"headerData": headerData,
						}).Warn("fileStoreCoreWithInst Failed to SetValueToBlockchain")
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				log.Infof("内容存储节点数字资产发布, id: %s, value: %s", headerData.ID, s)
			}
		}()

	default:
		// leave empty
	}
}

// ApplicationOnline register application to rds
func (h *Handler) ApplicationOnline(c *gin.Context) {
	ctx := c.Request.Context()
	type newReqData struct {
		model.StreamApplication
		NewFileStoreID string `json:"new_file_store_id"`
	}
	var headerData newReqData
	if err := c.BindJSON(&headerData); err != nil {
		log.WithFields(
			log.Fields{
				"error": err.Error(),
			}).Warn("Failed to extract stream application json data")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	appCore := model.StreamApplication{
		StreamApplicationCore: model.StreamApplicationCore{},
		ApplicationCore: model.ApplicationCore{
			ApplicationName: headerData.ApplicationName,
			ApplicationID:   headerData.ApplicationID,
			ApplicationPath: headerData.ApplicationPath,
			ApplicationFile: headerData.ApplicationFile,
			HWKey:           headerData.HWKey,
			OperatingSystem: headerData.OperatingSystem,
			CreateUser:      headerData.CreateUser,
			Description:     headerData.Description,
			UsageCount:      headerData.UsageCount,
		},
		AppInfoAttach: model.AppInfoAttach{
			FileStoreList:               headerData.FileStoreList,
			ImageName:                   headerData.ImageName,
			IsProviderReqGPU:            headerData.IsProviderReqGPU,
			IsFileStoreReqFastNetspeed:  headerData.IsFileStoreReqFastNetspeed,
			IsDepositoryReqFastNetspeed: headerData.IsDepositoryReqFastNetspeed,
		},
	}
	if headerData.NewFileStoreID == "" {
		if err := h.ApplicationService.CreateStreamApplication(ctx, &appCore); err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("CreateStreamApplication, Failed to register stream application to rds")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else {
		if err := h.ApplicationService.AddFileStoreIDToAPPInRDS(ctx, &appCore, headerData.NewFileStoreID); err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("AddFileStoreIDToAPPInRDS, Failed to register stream application to rds")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}
	h.ApplicationService.ShowEnterInfo(ctx, &appCore, headerData.NewFileStoreID)
	h.ApplicationService.ShowAllInfo(ctx)
}

func (h *Handler) Clear(c *gin.Context) {
	h.ConsumerService.Clear()
}

func (h *Handler) LearnNetWork(c *gin.Context) {
	var scheduleServiceCore model.ScheduleServiceCore
	scheduleServiceCore = *h.ConsumerService.GetScheduleServiceCore()
	ctx := c.Request.Context()
	appID, ok := c.GetQuery("app_id")
	if !ok {
		log.WithFields(log.Fields{
			"Client Address": c.Request.Host,
		}).Error("Failed to extract app_id, invalid http connection request, abandoned")
		c.JSON(http.StatusBadRequest, gin.H{"error": "need app_id"})
		return
	}

	var consumer = &model.Consumer{}
	consumer.ConsumerType = "terminal"
	consumer.UserName = "python3"
	provider, depositoryList, filestoreList, err := scheduleServiceCore.ScheduleStream(ctx, consumer, &model.StreamInstance{
		StreamApplication: &model.StreamApplication{
			ApplicationCore: model.ApplicationCore{
				ApplicationID: appID,
			},
		},
	})
	log.Printf("%+v\n", provider)
	if err != nil {
		log.WithFields(log.Fields{
			"Client Address": c.Request.Host,
			"error":          err.Error(),
		}).Error("Failed to ScheduleStream, abandoned")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	req := struct {
		ProviderCore   model.ProviderCore             `json:"provider_core"`
		DepositoryList []model.DepositoryCoreWithInst `json:"depository_list"`
		FileStoreList  []model.FileStoreCoreWithInst  `json:"filestore_list"`
	}{
		ProviderCore:   provider.ProviderCore,
		DepositoryList: depositoryList,
		FileStoreList:  filestoreList,
	}
	log.Printf("Before marshal: %+v\n", req)
	reqStr, err := json.Marshal(&req)
	log.Printf("After marshal: %+v\n", string(reqStr))
	if err != nil {
		log.WithFields(log.Fields{
			"Client Address": c.Request.Host,
			"error":          err.Error(),
		}).Error("Failed to Marshal, abandoned")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"info": string(reqStr),
	})
}

func (h *Handler) RecordHistory(c *gin.Context) {
	type ReqData struct {
		InstanceID string `json:"instance_id"`
		Latency    string `json:"latency"`
	}
	var headerData ReqData
	if err := c.BindJSON(&headerData); err != nil {
		log.WithFields(
			log.Fields{
				"error": err.Error(),
			}).Error("Failed to extract RecordHistory json data")
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	fmt.Printf("获得来自前端的服务反馈,实例 ID: %s, 业务层时延: %s\n", headerData.InstanceID, headerData.Latency)

	var sc model.ScheduleServiceCore
	sc = *h.ConsumerService.GetScheduleServiceCore()
	instRoom, err := sc.GetStreamInstanceRoomByInstanceID(headerData.InstanceID)
	if err != nil {
		log.WithFields(
			log.Fields{
				"error":      err.Error(),
				"headerData": headerData,
			}).Warn("Failed to GetStreamInstanceRoomByInstanceID")
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	go func() {
		s, err := sc.GetValueFromBlockchain(instRoom.Provider.ID)
		if err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("Failed to GetValueFromBlockchain(instRoom.Provider.ID)")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var p model.ProviderCoreWithInst
		p.InstHistory = make(map[string]string)
		err = json.Unmarshal([]byte(s), &p)
		if err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("Failed to ProviderCoreWithInst Unmarshal")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if p.InstHistory == nil {
			p.InstHistory = make(map[string]string)
		}
		p.InstHistory[headerData.InstanceID] = headerData.Latency
		if s, err := json.Marshal(&p); err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("Failed to Marshal providerCoreWithInst")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			err := sc.SetValueToBlockchain(p.ID, string(s))
			if err != nil {
				log.WithFields(
					log.Fields{
						"error":      err.Error(),
						"headerData": headerData,
					}).Warn("providerCoreWithInst Failed to SetValueToBlockchain")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			log.Infof("【服务提供节点】数字资产更新 (解析后):\n%s", p.DetailedInfo())
			log.Infof("【服务提供节点】数字资产更新, 资产索引: %s, 资产内容: %s", p.ID, s)
			err = h.ProviderService.UpdateProviderInRDS(context.TODO(), &p)
			if err != nil {
				log.Infof("RecordHistory h.ProviderService.UpdateProviderInRDS error: %s", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		s, err = sc.GetValueFromBlockchain(instRoom.SelectedDepository.ID)
		if err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("Failed to GetValueFromBlockchain(instRoom.SelectedDepository.ID)")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var d model.DepositoryCoreWithInst
		err = json.Unmarshal([]byte(s), &d)
		if err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("Failed to DepositoryCoreWithInst Unmarshal")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if d.InstHistory == nil {
			d.InstHistory = make(map[string]string)
		}
		d.InstHistory[headerData.InstanceID] = instRoom.SelectedDepositoryBandWidth
		if s, err := json.Marshal(&d); err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("Failed to Marshal DepositoryCoreWithInst")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			err := sc.SetValueToBlockchain(d.ID, string(s))
			if err != nil {
				log.WithFields(
					log.Fields{
						"error":      err.Error(),
						"headerData": headerData,
					}).Warn("DepositoryCoreWithInst Failed to SetValueToBlockchain")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			log.Infof("【镜像仓库节点】数字资产更新 (解析后):\n%s", d.DetailedInfo())
			log.Infof("【镜像仓库节点】数字资产更新, 资产索引: %s, 资产内容: %s", d.ID, s)
			err = h.DepositoryService.UpdateFileStoreInRDS(context.TODO(), &d)
			if err != nil {
				log.Infof("RecordHistory h.DepositoryService.UpdateFileStoreInRDS error: %s", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}

		s, err = sc.GetValueFromBlockchain(instRoom.SelectedFileStore.ID)
		if err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("Failed to GetValueFromBlockchain(instRoom.SelectedFileStore.ID)")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		var f model.FileStoreCoreWithInst
		err = json.Unmarshal([]byte(s), &f)
		if err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("Failed to FileStoreCoreWithInst Unmarshal")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if f.InstHistory == nil {
			f.InstHistory = make(map[string]string)
		}
		f.InstHistory[headerData.InstanceID] = instRoom.SelectedFileStoreLatency
		if s, err := json.Marshal(&f); err != nil {
			log.WithFields(
				log.Fields{
					"error":      err.Error(),
					"headerData": headerData,
				}).Warn("Failed to Marshal FileStoreCoreWithInst")
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		} else {
			err := sc.SetValueToBlockchain(f.ID, string(s))
			if err != nil {
				log.WithFields(
					log.Fields{
						"error":      err.Error(),
						"headerData": headerData,
					}).Warn("FileStoreCoreWithInst Failed to SetValueToBlockchain")
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			log.Infof("【内容存储节点】数字资产更新 (解析后):\n%s", f.DetailedInfo())
			log.Infof("【内容存储节点】数字资产更新, 资产索引: %s, 资产内容: %s", f.ID, s)
			err = h.FileStoreService.UpdateFileStoreInRDS(context.TODO(), &f)
			if err != nil {
				log.Infof("RecordHistory h.FileStoreService.UpdateFileStoreInRDS error: %s", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
		}
	}()
}
