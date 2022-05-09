package dal

import (
	"context"

	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
)

/*
	@struct: ConsumerDAL
	@description: DAL layer
*/
type ConsumerDAL struct {
	ConsumerList map[string]*model.Consumer
}

/*
	@struct: ConsumerDALConfig
	@description: used for config instance of struct ConsumerDAL
*/
type ConsumerDALConfig struct{}

/*
	@func: NewConsumerDAL
	@description:
		create, config and return an instance of struct ConsumerDAL
*/
func NewConsumerDAL(c *ConsumerDALConfig) model.ConsumerDAL {
	cdal := &ConsumerDAL{}

	cdal.ConsumerList = make(map[string]*model.Consumer)

	return cdal
}

/*
	@func: CreateConsumer
	@description:
		insert a new consumer to consumer list
*/
func (d *ConsumerDAL) CreateConsumer(ctx context.Context, consumer *model.Consumer) {
	d.ConsumerList[consumer.ClientID] = consumer
}

/*
	@func: DeleteConsumer
	@description:
		delete the specified consumer from consumer list
*/
func (d *ConsumerDAL) DeleteConsumer(ctx context.Context, consumerID string) {
	delete(d.ConsumerList, consumerID)
}
