package dal

import (
	"context"
	"fmt"

	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
)

/*
	@struct: InstanceRoomDAL
	@description: DAL layer
*/
type InstanceRoomDAL struct {
	StreamInstanceRoomList map[string]*model.StreamInstanceRoom
}

/*
	@struct: InstanceRoomDALConfig
	@description: used for config instance of struct InstanceRoomDAL
*/
type InstanceRoomDALConfig struct{}

/*
	@func: NewInstanceRoomDAL
	@description:
		create, config and return an instance of struct InstanceRoomDAL
*/
func NewInstanceRoomDAL(c *InstanceRoomDALConfig) model.InstanceRoomDAL {
	irdal := &InstanceRoomDAL{}

	irdal.StreamInstanceRoomList = make(map[string]*model.StreamInstanceRoom)

	return irdal
}

/*
	@func: CreateStreamInstanceRoom
	@description:
		insert a new stream instance room to the list
*/
func (d *InstanceRoomDAL) CreateStreamInstanceRoom(ctx context.Context, streamInstanceRoom *model.StreamInstanceRoom) {
	d.StreamInstanceRoomList[streamInstanceRoom.InstanceID] = streamInstanceRoom
}

/*
	@func: DeleteStreamInstanceRoom
	@description:
		delete a specified stream instance room from the list
*/
func (d *InstanceRoomDAL) DeleteStreamInstanceRoom(ctx context.Context, instanceID string) {
	delete(d.StreamInstanceRoomList, instanceID)
}

/*
	@func: GetConsumerMapByInstanceID
	@description:
		obtain consumer map by given instance room id
*/
func (d *InstanceRoomDAL) GetConsumerMapByInstanceID(ctx context.Context, instanceID string) (map[string]*model.Consumer, error) {
	instanceRoom, ok := d.StreamInstanceRoomList[instanceID]
	if !ok {
		return nil, fmt.Errorf("No instance founded by given instance id")
	}
	return instanceRoom.ConsumerList, nil
}

/*
	@func: GetProviderByInstanceID
	@description:
		obtain provider by given instance room id
*/
func (d *InstanceRoomDAL) GetProviderByInstanceID(ctx context.Context, instanceID string) (*model.Provider, error) {
	instanceRoom, ok := d.StreamInstanceRoomList[instanceID]
	if !ok {
		return nil, fmt.Errorf("No instance founded by given instance id")
	}
	return instanceRoom.Provider, nil
}
