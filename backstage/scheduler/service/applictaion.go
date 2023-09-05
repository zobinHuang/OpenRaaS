package service

import (
	"context"

	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
)

/*
	@struct: ApplicationService
	@description: service layer
*/
type ApplicationService struct {
	ApplicationDAL model.ApplicationDAL
}

/*
	@struct: ApplicationServiceConfig
	@description: used for config instance of struct ApplicationService
*/
type ApplicationServiceConfig struct {
	ApplicationDAL model.ApplicationDAL
}

/*
	@func: NewApplicationService
	@description:
		create, config and return an instance of struct ApplicationService
*/
func NewApplicationService(c *ApplicationServiceConfig) model.ApplicationService {
	return &ApplicationService{
		ApplicationDAL: c.ApplicationDAL,
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
