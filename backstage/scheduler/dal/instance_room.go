package dal

import (
	"context"

	"github.com/zobinHuang/BrosCloud/backstage/scheduler/model"
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
