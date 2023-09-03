package dal

import (
	"context"
	"gorm.io/gorm"

	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
)

/*
@struct: ProviderDAL
@description: DAL layer
*/
type ProviderDAL struct {
	DB           *gorm.DB
	ProviderList map[string]*model.Consumer
}

/*
@struct: ProviderDALConfig
@description: used for config instance of struct ProviderDAL
*/
type ProviderDALConfig struct {
	DB *gorm.DB
}

/*
@func: NewProviderDAL
@description:

	create, config and return an instance of struct ProviderDAL
*/
func NewProviderDAL(c *ProviderDALConfig) model.ProviderDAL {
	cdal := &ProviderDAL{}

	cdal.ProviderList = make(map[string]*model.Consumer)
	cdal.DB = c.DB

	return cdal
}

/*
@func: CreateProvider
@description:

	insert a new provider to provider list
*/
func (d *ProviderDAL) CreateProvider(ctx context.Context, provider *model.Provider) error {
	// todo
	return nil
}

/*
@func: DeleteProviderByID
@description:

	delete the specified provider from provider list
*/
func (d *ProviderDAL) DeleteProviderByID(ctx context.Context, id string) error {
	return nil
}

func (d *ProviderDAL) GetProvider(ctx context.Context) ([]model.Provider, error) {
	return nil, nil
}
func (d *ProviderDAL) GetProviderByID(ctx context.Context, id string) (*model.Provider, error) {
	return nil, nil
}
func (d *ProviderDAL) UpdateProviderByID(ctx context.Context, info *model.Provider) error {
	return nil
}
