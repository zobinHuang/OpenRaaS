package dal

import (
	"context"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model/apperrors"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

/*
@struct: ApplicationDAL
@description: DAL layer
*/
type ApplicationDAL struct {
	DB *gorm.DB
}

/*
@struct: ApplicationDALConfig
@description: used for config instance of struct ApplicationDAL
*/
type ApplicationDALConfig struct {
	DB *gorm.DB
}

/*
@func: NewApplicationDAL
@description:

	create, config and return an instance of struct ApplicationDAL
*/
func NewApplicationDAL(c *ApplicationDALConfig) model.ApplicationDAL {
	return &ApplicationDAL{
		DB: c.DB,
	}
}

/*
@func: GetStreamApplicationByID
@description:

	obtain stream application according to given application id
*/
func (d *ApplicationDAL) GetStreamApplicationByID(ctx context.Context, applicationID string) (*model.StreamApplication, error) {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// instantiate stream application
	streamApplication := &model.StreamApplication{}

	// retrieve
	if err := tx.Where("application_id = ?", applicationID).First(streamApplication).Error; err != nil {
		log.WithFields(log.Fields{
			"Given Application ID": applicationID,
		}).Warn("Unable to obtain stream application with given application id")
		return nil, apperrors.NewNotFound("application_id", applicationID)
	}

	return streamApplication, nil
}

/*
@func: GetStreamApplicationsOrderedByUpdateTime
@description:

	obtain stream application list, ordered by update time
*/
func (d *ApplicationDAL) GetStreamApplicationsOrderedByUpdateTime(ctx context.Context, listLength int, listID int) ([]*model.StreamApplication, error) {
	// initialize application list
	streamApplicationList := make([]*model.StreamApplication, listLength)

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Limit(listLength).Offset(listLength * (listID - 1)).Order("updated_at").Find(&streamApplicationList).Error; err != nil {
		log.WithFields(log.Fields{
			"Given List Length": listLength,
			"Given List ID":     listID,
			"error":             err,
		}).Warn("Failed to obtain stream application list based on update time by given list metadata")
		return nil, err
	}

	return streamApplicationList, nil
}

/*
@func: GetStreamApplicationsOrderedByName
@description:

	obtain stream application list, ordered by application name
*/
func (d *ApplicationDAL) GetStreamApplicationsOrderedByName(ctx context.Context, listLength int, listID int) ([]*model.StreamApplication, error) {
	// initialize application list
	streamApplicationList := make([]*model.StreamApplication, listLength)

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Limit(listLength).Offset(listLength * (listID - 1)).Order("application_name").Find(&streamApplicationList).Error; err != nil {
		log.WithFields(log.Fields{
			"Given List Length": listLength,
			"Given List ID":     listID,
			"error":             err,
		}).Warn("Failed to obtain stream application list based on name by given list metadata")
		return nil, err
	}

	return streamApplicationList, nil
}

/*
@func: GetStreamApplicationsOrderedByUsageCount
@description:

	obtain stream application list, ordered by usage count
*/
func (d *ApplicationDAL) GetStreamApplicationsOrderedByUsageCount(ctx context.Context, listLength int, listID int) ([]*model.StreamApplication, error) {
	// initialize application list
	streamApplicationList := make([]*model.StreamApplication, listLength)

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Limit(listLength).Offset(listLength * (listID - 1)).Order("usage_count").Find(&streamApplicationList).Error; err != nil {
		log.WithFields(log.Fields{
			"Given List Length": listLength,
			"Given List ID":     listID,
			"error":             err,
		}).Warn("Failed to obtain stream application list based on usage count by given list metadata")
		return nil, err
	}

	return streamApplicationList, nil
}

/*
@func: GetStreamApplicationsCount
@description:

	obtain total count of stream applications
*/
func (d *ApplicationDAL) GetStreamApplicationsCount(ctx context.Context) (int64, error) {
	var count int64

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("stream_applications").Count(&count).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to obtain count of stream applications")
		return 0, err
	}

	return count, nil
}

// GetStreamApplication get all stream applications from rds
func (d *ApplicationDAL) GetStreamApplication(ctx context.Context) ([]model.StreamApplication, error) {
	var apps []model.StreamApplication

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("stream_applications").Find(&apps).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to obtain all stream applications from rds")
		return nil, err
	}

	return apps, nil
}

// DeleteStreamApplicationByID delete stream application by id in rds
func (d *ApplicationDAL) DeleteStreamApplicationByID(ctx context.Context, id string) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("stream_applications").Where("application_id=?", id).Delete(&model.StreamApplication{}).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Warn("Failed to delete stream application by id in rds")
		return err
	}
	return nil
}

// UpdateStreamApplicationByID update stream application to rds by id
func (d *ApplicationDAL) UpdateStreamApplicationByID(ctx context.Context, info *model.StreamApplication) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("stream_applications").Where("application_id=?", info.ApplicationID).Updates(info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"app":   info,
		}).Warn("Failed to update stream application to rds by id")
		return err
	}
	return nil
}

// CreateStreamApplication create stream application to rds
func (d *ApplicationDAL) CreateStreamApplication(ctx context.Context, info *model.StreamApplication) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("stream_applications").Create(info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"app":   info,
		}).Warn("Failed to create stream application to rds")
		return err
	}

	return nil
}

// Clear delete all
func (d *ApplicationDAL) Clear() {
	d.DB.Delete(&model.StreamApplication{})
}
