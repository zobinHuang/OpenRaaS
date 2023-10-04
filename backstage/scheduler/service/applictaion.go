package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/liushuochen/gotable"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/utils"
	"gorm.io/gorm"
)

/*
@struct: ApplicationService
@description: service layer
*/
type ApplicationService struct {
	ApplicationDAL model.ApplicationDAL
	FileStoreDAL   model.FileStoreDAL
	DepositoryDAL  model.DepositoryDAL
}

/*
@struct: ApplicationServiceConfig
@description: used for config instance of struct ApplicationService
*/
type ApplicationServiceConfig struct {
	ApplicationDAL model.ApplicationDAL
	FileStoreDAL   model.FileStoreDAL
	DepositoryDAL  model.DepositoryDAL
}

/*
@func: NewApplicationService
@description:

	create, config and return an instance of struct ApplicationService
*/
func NewApplicationService(c *ApplicationServiceConfig) model.ApplicationService {
	return &ApplicationService{
		ApplicationDAL: c.ApplicationDAL,
		FileStoreDAL:   c.FileStoreDAL,
		DepositoryDAL:  c.DepositoryDAL,
	}
}

/*
@func: GetStreamApplicationsCount
@description:

	obtain amount of stream applications
*/
func (s *ApplicationService) GetStreamApplicationsCount(ctx context.Context) (int64, error) {
	return s.ApplicationDAL.GetStreamApplicationsCount(ctx)
}

/*
@func: GetStreamApplicationDetails
@description:

	obtain detail for specific application
*/
func (s *ApplicationService) GetStreamApplicationDetails(ctx context.Context, applicationID string) (*model.StreamApplication, error) {
	return s.ApplicationDAL.GetStreamApplicationByID(ctx, applicationID)
}

/*
@func: GetStreamApplications
@description:

	obtain application
*/
func (s *ApplicationService) GetStreamApplications(ctx context.Context, pageNumber int, pageSize int, orderBy string) ([]*model.StreamApplication, error) {
	if orderBy == model.ORDER_BY_NAME {
		// order by name
		streamApplicationList, err := s.ApplicationDAL.GetStreamApplicationsOrderedByName(ctx, pageSize, pageNumber)
		if err != nil {
			return nil, err
		}
		return streamApplicationList, nil

	} else if orderBy == model.ORDER_BY_USAGE_COUNT {
		// order by usage count
		streamApplicationList, err := s.ApplicationDAL.GetStreamApplicationsOrderedByUsageCount(ctx, pageSize, pageNumber)
		if err != nil {
			return nil, err
		}
		return streamApplicationList, nil

	} else if orderBy == model.ORDER_BY_UPDATE_TIME {
		// order by update time
		streamApplicationList, err := s.ApplicationDAL.GetStreamApplicationsOrderedByUpdateTime(ctx, pageSize, pageNumber)
		if err != nil {
			return nil, err
		}
		return streamApplicationList, nil

	} else {
		return nil, nil
	}
}

/*
	@func: GetStreamApplicationAmount
	@description:
		obtain amount of registered stream application
*/

// CreateStreamApplication register app to rds
func (s *ApplicationService) CreateStreamApplication(ctx context.Context, info *model.StreamApplication) error {
	return s.ApplicationDAL.CreateStreamApplication(ctx, info)
}

func (s *ApplicationService) AddFileStoreIDToAPPInRDS(ctx context.Context, info *model.StreamApplication, id string) error {
	app, err := s.ApplicationDAL.GetStreamApplicationByID(ctx, info.ApplicationID)
	if err == nil {
		if app.FileStoreList == "" {
			app.FileStoreList = "[\"" + id + "\"]"
		} else {
			var ids []string
			if err := json.Unmarshal([]byte(app.FileStoreList), &ids); err != nil {
				return err
			}
			ids = append(ids, id)
			idsStr, err := json.Marshal(&ids)
			if err != nil {
				return err
			}
			app.FileStoreList = string(idsStr)
		}
		return s.ApplicationDAL.UpdateStreamApplicationByID(ctx, app)
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		err := s.CreateStreamApplication(ctx, info)
		if err != nil {
			return err
		}
		return s.AddFileStoreIDToAPPInRDS(ctx, info, id)
	}
	return err
}

func (s *ApplicationService) ShowEnterInfo(ctx context.Context, app *model.StreamApplication, nodeId string) {
	log.Infof("%s, 软件上线, 软件 id: %s，上传节点 ID: %s", utils.GetCurrentTime(), app.ApplicationID, nodeId)
	log.Infof("详细信息: %s", app.DetailedInfo())
}

func (s *ApplicationService) ShowAllInfo(ctx context.Context) {
	apps, err := s.ApplicationDAL.GetStreamApplication(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("DepositoryService ShowAllInfo GetDepositoryInRDS error")
	}
	fileStores, err := s.FileStoreDAL.GetFileStoreInRDS(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("ApplicationService ShowAllInfo GetFileStoreInRDS error")
	}
	depositories, err := s.DepositoryDAL.GetDepositoryInRDS(ctx)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Error("ApplicationService ShowAllInfo GetDepositoryInRDS error")
	}

	log.Infof("整合前，软件资源信息：%+v", apps)
	log.Infof("整合后，软件资源信息：")
	log.Infof("%s, app 数量：%d, 内容存储节点数量：%d, 镜像仓库节点数量：%d",
		utils.GetCurrentTime(), len(apps), len(fileStores), len(depositories))

	//depositoriesIds := make([]string, 0, 0)
	//for _, d := range depositories {
	//	depositoriesIds = append(depositoriesIds, d.ID)
	//}
	//depositoriesIdsStr, _ := json.Marshal(depositoriesIds)

	table, err := gotable.Create("软件 ID", "软件名", "软件路径", "启动文件", "软件类型", "镜像 ID", "支持的内容存储节点", "是否需要高性能服务提供节点", "是否需要高性能内容存储节点", "是否需要高性能镜像仓库节点", "软件说明")
	if err != nil {
		fmt.Println("ShowAllInfo ApplicationService Create table failed: ", err.Error())
		return
	}
	for _, a := range apps {
		table.AddRow([]string{a.ApplicationID, a.ApplicationName, a.ApplicationPath, a.ApplicationFile, a.HWKey, a.ImageName, a.FileStoreList,
			strconv.FormatBool(a.IsProviderReqGPU), strconv.FormatBool(a.IsFileStoreReqFastNetspeed), strconv.FormatBool(a.IsDepositoryReqFastNetspeed),
			a.Description})
	}
	fmt.Println("\n", table, "\n")
}
