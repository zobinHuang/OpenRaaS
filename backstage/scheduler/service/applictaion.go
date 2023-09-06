package service

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/utils"
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
	log.Infof("%s, allow new application enter, id: %s", utils.GetCurrentTime(), app.ApplicationID)
	req1 := "normal"
	if app.IsProviderReqGPU {
		req1 = "powerful"
	}
	req2 := "normal"
	if app.IsProviderReqGPU {
		req2 = "powerful"
	}
	req3 := "normal"
	if app.IsProviderReqGPU {
		req3 = "powerful"
	}
	log.Infof("%s, new application id: %s, name: %s, type: %s, description: %s, image id: 50fbb73b-1979-4938-8bbe-41dd6fe066a9, launch nodes: %s, "+
		"provider request: %s, depository request: %s, fileStore request: %s",
		utils.GetCurrentTime(), app.ApplicationID, app.ApplicationName, app.HWKey, app.Description, nodeId, req1, req2, req3)
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
	log.Infof("%s, Applications Total Info, app amount: %d, launch fileStore amount: %d, launch depository amount: %d, served image amount: 1",
		utils.GetCurrentTime(), len(apps), len(fileStores), len(depositories))

	depositoriesIds := make([]string, 0, 0)
	for _, d := range depositories {
		depositoriesIds = append(depositoriesIds, d.ID)
	}
	depositoriesIdsStr, _ := json.Marshal(depositoriesIds)

	for _, app := range apps {
		req1 := "normal"
		if app.IsProviderReqGPU {
			req1 = "powerful"
		}
		req2 := "normal"
		if app.IsProviderReqGPU {
			req2 = "powerful"
		}
		req3 := "normal"
		if app.IsProviderReqGPU {
			req3 = "powerful"
		}
		log.Infof("%s, new application id: %s, name: %s, type: %s, description: %s, image id: 50fbb73b-1979-4938-8bbe-41dd6fe066a9, "+
			"launch fileStores: %s, launch depositories: %s, provider request: %s, depository request: %s, fileStore request: %s",
			utils.GetCurrentTime(), app.ApplicationID, app.ApplicationName, app.HWKey, app.Description,
			app.FileStoreList, depositoriesIdsStr, req1, req2, req3)
	}
}
