package dal

import (
	"context"
	"gorm.io/gorm"

	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
)

/*
@struct: DepositaryDAL
@description: DAL layer
*/
type DepositaryDAL struct {
	DB *gorm.DB
}

/*
@struct: DepositaryDALConfig
@description: used for config instance of struct DepositaryDAL
*/
type DepositaryDALConfig struct {
	DB *gorm.DB
}

/*
@func: NewDepositaryDAL
@description:

	create, config and return an instance of struct DepositaryDAL
*/
func NewDepositaryDAL(c *DepositaryDALConfig) model.DepositaryDAL {
	return &DepositaryDAL{
		DB: c.DB,
	}
}

/*
@func: CreateDepositary
@description:

	insert a new depositary to depositary list
*/
func (d *DepositaryDAL) CreateDepositary(ctx context.Context, info *model.Depositary) error {
	return nil
}

/*
@func: DeleteDepositary
@description:

	delete the specified depositary from depositary list
*/
func (d *DepositaryDAL) DeleteDepositaryByID(ctx context.Context, id string) error {
	return nil
}

func (d *DepositaryDAL) GetDepositary(ctx context.Context) ([]model.Depositary, error) {
	return nil, nil
}
func (d *DepositaryDAL) GetDepositaryByID(ctx context.Context, id string) (*model.Depositary, error) {
	return nil, nil
}
func (d *DepositaryDAL) UpdateDepositaryByID(ctx context.Context, info *model.Depositary) error {
	return nil
}
