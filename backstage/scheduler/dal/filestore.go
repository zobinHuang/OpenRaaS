package dal

import (
	"context"

	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
)

/*
	@struct: FilestoreDAL
	@description: DAL layer
*/
type FilestoreDAL struct {
	FilestoreList map[string]*model.Filestore
}

/*
	@struct: FilestoreDALConfig
	@description: used for config instance of struct FilestoreDAL
*/
type FilestoreDALConfig struct {
}

/*
	@func: NewFilestoreDAL
	@description:
		create, config and return an instance of struct FilestoreDAL
*/
func NewFilestoreDAL(c *FilestoreDALConfig) model.FilestoreDAL {
	ddal := &FilestoreDAL{}

	ddal.FilestoreList = make(map[string]*model.Filestore)

	return ddal
}

/*
	@func: CreateFilestore
	@description:
		insert a new filestore to Filestore list
*/
func (d *FilestoreDAL) CreateFilestore(ctx context.Context, filestore *model.Filestore) {
	d.FilestoreList[filestore.ClientID] = filestore
}

/*
	@func: DeleteFilestore
	@description:
		delete the specified filestore from Filestore list
*/
func (d *FilestoreDAL) DeleteFilestore(ctx context.Context, filestoreID string) {
	delete(d.FilestoreList, filestoreID)
}
