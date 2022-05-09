package servicecore

import (
	"context"
	"fmt"

	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
)

/*
	@struct: ScheduleServiceCore
	@description: service core layer
*/
type ScheduleServiceCore struct {
	ConsumerDAL     model.ConsumerDAL
	ProviderDAL     model.ProviderDAL
	DepositaryDAL   model.DepositaryDAL
	FilestoreDAL    model.FilestoreDAL
	InstanceRoomDAL model.InstanceRoomDAL
}

/*
	@struct: ScheduleServiceCoreConfig
	@description: used for config instance of struct ScheduleServiceCore
*/
type ScheduleServiceCoreConfig struct {
	ConsumerDAL     model.ConsumerDAL
	ProviderDAL     model.ProviderDAL
	DepositaryDAL   model.DepositaryDAL
	FilestoreDAL    model.FilestoreDAL
	InstanceRoomDAL model.InstanceRoomDAL
}

/*
	@func: NewScheduleServiceCore
	@description:
		create, config and return an instance of struct ScheduleServiceCore
*/
func NewScheduleServiceCore(c *ScheduleServiceCoreConfig) model.ScheduleServiceCore {
	return &ScheduleServiceCore{
		ConsumerDAL:     c.ConsumerDAL,
		ProviderDAL:     c.ProviderDAL,
		DepositaryDAL:   c.DepositaryDAL,
		FilestoreDAL:    c.FilestoreDAL,
		InstanceRoomDAL: c.InstanceRoomDAL,
	}
}

/*
	@func: ScheduleStream
	@description:
		core logic of scheduling stream instance is here
*/
func (sc *ScheduleServiceCore) ScheduleStream(ctx context.Context, streamInstance *model.StreamInstance) (*model.Provider, []model.DepositaryCore, []model.FilestoreCore, error) {

	//TODO: schedule strategy

	return nil, nil, nil, fmt.Errorf("schedule strategy not implemented yet")
}

/*
	@func: CreateStreamInstanceRoom
	@description:
		create a room for the instance of stream instance
*/
func (sc *ScheduleServiceCore) CreateStreamInstanceRoom(ctx context.Context, provider *model.Provider,
	consumer *model.Consumer, streamInstance *model.StreamInstance) (*model.StreamInstanceRoom, error) {
	// initialize streamInstanceRoom instance
	streamInstanceRoom := &model.StreamInstanceRoom{
		StreamInstance: streamInstance,
		Provider:       provider,
	}

	// create consumer list, and insert our current consumer
	streamInstanceRoom.ConsumerList = make(map[string]*model.Consumer)
	streamInstanceRoom.ConsumerList[consumer.ClientID] = consumer

	// insert in dal layer
	sc.InstanceRoomDAL.CreateStreamInstanceRoom(ctx, streamInstanceRoom)

	return streamInstanceRoom, nil
}
