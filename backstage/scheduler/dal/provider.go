package dal

import (
	"context"

	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
)

/*
	@struct: ProviderDAL
	@description: DAL layer
*/
type ProviderDAL struct {
	ProviderList map[string]*model.Provider
}

/*
	@struct: ProviderDALConfig
	@description: used for config instance of struct ProviderDAL
*/
type ProviderDALConfig struct{}

/*
	@func: NewProviderDAL
	@description:
		create, config and return an instance of struct ProviderDAL
*/
func NewProviderDAL(c *ProviderDALConfig) model.ProviderDAL {
	cdal := &ProviderDAL{}

	cdal.ProviderList = make(map[string]*model.Provider)

	return cdal
}

/*
	@func: CreateProvider
	@description:
		insert a new provider to provider list
*/
func (d *ProviderDAL) CreateProvider(ctx context.Context, provider *model.Provider) {
	d.ProviderList[provider.ClientID] = provider
}

/*
	@func: DeleteProvider
	@description:
		delete the specified provider from provider list
*/
func (d *ProviderDAL) DeleteProvider(ctx context.Context, providerID string) {
	delete(d.ProviderList, providerID)
}
