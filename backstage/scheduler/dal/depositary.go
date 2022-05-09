package dal

import (
	"context"

	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
)

/*
	@struct: DepositaryDAL
	@description: DAL layer
*/
type DepositaryDAL struct {
	DepositaryList map[string]*model.Depositary
}

/*
	@struct: DepositaryDALConfig
	@description: used for config instance of struct DepositaryDAL
*/
type DepositaryDALConfig struct {
}

/*
	@func: NewDepositaryDAL
	@description:
		create, config and return an instance of struct DepositaryDAL
*/
func NewDepositaryDAL(c *DepositaryDALConfig) model.DepositaryDAL {
	ddal := &DepositaryDAL{}

	ddal.DepositaryList = make(map[string]*model.Depositary)

	return ddal
}

/*
	@func: CreateDepositary
	@description:
		insert a new depositary to depositary list
*/
func (d *DepositaryDAL) CreateDepositary(ctx context.Context, depositary *model.Depositary) {
	d.DepositaryList[depositary.ClientID] = depositary
}

/*
	@func: DeleteDepositary
	@description:
		delete the specified depositary from depositary list
*/
func (d *DepositaryDAL) DeleteDepositary(ctx context.Context, depositaryID string) {
	delete(d.DepositaryList, depositaryID)
}
