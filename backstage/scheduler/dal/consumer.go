package dal

import (
	"context"
	"fmt"
	"time"

	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
)

type Timer struct {
	T0 time.Time
	T1 time.Time
}

/*
@struct: ConsumerDAL
@description: DAL layer
*/
type ConsumerDAL struct {
	ConsumerList map[string]*model.Consumer
	UserRecord   map[string]Timer
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

	cdal.UserRecord = make(map[string]Timer)

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

/*
@func: GetConsumerByID
@description:

	obtain a consumer by given consumer id
*/
func (d *ConsumerDAL) GetConsumerByID(ctx context.Context, consumerID string) (*model.Consumer, error) {
	consumer, ok := d.ConsumerList[consumerID]
	if !ok {
		return nil, fmt.Errorf("Failed to obtain consumer by given consumer id")
	}
	return consumer, nil
}

// Clear delete all
func (d *ConsumerDAL) Clear() {
	d.ConsumerList = make(map[string]*model.Consumer)
	d.UserRecord = make(map[string]Timer)
}

/*
@func: GetConsumers
@description:

	obtain all consumers
*/
func (d *ConsumerDAL) GetConsumers() map[string]*model.Consumer {
	return d.ConsumerList
}

func (d *ConsumerDAL) AddUser(name string) {
	d.UserRecord[name] = Timer{}
}

func (d *ConsumerDAL) HasUser(name string) bool {
	_, ok := d.UserRecord[name]
	return ok
}

func (d *ConsumerDAL) IsUserOverTime(name string) bool {
	user := d.UserRecord[name]
	return (!user.T1.IsZero() && time.Now().Sub(user.T0) < time.Minute*30) ||
		(!user.T0.IsZero() && time.Now().Sub(user.T0) >= time.Minute*30)
}

func (d *ConsumerDAL) UserUpdateTime(name string, t time.Time) {
	user := d.UserRecord[name]
	if user.T0.IsZero() {
		user.T0 = t
	} else if user.T1.IsZero() {
		user.T1 = t
	} else {
		user.T0 = user.T1
		user.T1 = t
	}
}
