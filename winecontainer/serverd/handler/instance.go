package handler

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"serverd/model"
	"serverd/model/apperrors"
	"strconv"

	"github.com/gin-gonic/gin"
)

/*
func: CreateInstance
description: handler for endpoint "/api/daemon/createinstance"
*/
func (h *Handler) CreateInstance(c *gin.Context) {
	instanceModel := &model.InstanceModel{}
	if ok := bindData(c, &instanceModel); !ok {
		return
	}

	ctx := c.Request.Context()
	err := h.CreateInstanceWithModel(ctx, instanceModel)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err.Error(),
		})
	}

	// return http_ok if success
	c.JSON(http.StatusOK, gin.H{
		"vmid": strconv.Itoa(instanceModel.VMID),
	})
}

/*
func: CreateInstanceWithModel
description: use input instanceModel to create a new wine instance
*/
func (h *Handler) CreateInstanceWithModel(ctx context.Context, instanceModel *model.InstanceModel) error {
	// Return a tunnel to which you can pass a random parameter if you want to shut down the VM （done <- struct{}）
	done := h.InstanceService.LaunchInstance(ctx, instanceModel)
	if done == nil {
		log.Printf("Failed to start up a new instance.")
		return fmt.Errorf("Instance startup failed.")
	}
	return nil
}

/*
func: DeleteInstance
description: handler for endpoint "/api/daemon/deleteinstance"
*/
func (h *Handler) DeleteInstance(c *gin.Context) {
	// bind request
	deleteModel := &model.DeleteInstanceModel{}
	if ok := bindData(c, &deleteModel); !ok {
		return
	}

	ctx := c.Request.Context()
	err := h.DeleteInstanceWithModel(ctx, deleteModel)

	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": err,
		})
	}

	// return http_ok if success
	c.JSON(http.StatusOK, gin.H{})
}

/*
func: DeleteInstanceWithModel
description: use input deleteModel to delete the target wine instance
*/
func (h *Handler) DeleteInstanceWithModel(ctx context.Context, deleteModel *model.DeleteInstanceModel) error {
	var err error
	if deleteModel.Instanceid != "" {
		err = h.InstanceService.DeleteInstanceByInstanceid(ctx, deleteModel.Instanceid)
	} else if deleteModel.VMID != 0 {
		err = h.InstanceService.DeleteInstance(ctx, deleteModel.VMID)
	} else {
		log.Printf("No instance identification sent in.\n")
		return fmt.Errorf("No instance specified.")
	}

	if err != nil {
		log.Printf("Failed to delete instance: %v\n", err.Error())
	}
	return err
}

/*
func: SelectFilestore
description: check network connection status between host and storage servers, and select one as schedule target
*/
func (h *Handler) SelectFilestore(ctx context.Context, instanceModel *model.InstanceModel) error {
	/* Test */
	var filestore model.FilestoreCore

	filestore.HostAddress = "kb109.dynv6.net"
	filestore.Port = "7189"
	filestore.Protocol = "davfs"
	filestore.Username = "kb109"
	filestore.Password = "Xusir666!"
	filestore.Directory = "/public_hdd/game/PC/dcwine"
	// filestore.Directory = "/storage_ssd/6G/dcwine"

	instanceModel.TargetFilestore = filestore

	// TODO: complete filestore schedule process

	return nil
}

/*
func: SelectDepository
description: check network connection status between host and depository servers, and select one as schedule target
*/
func (h *Handler) SelectDepository(ctx context.Context, instanceModel *model.InstanceModel) error {
	/* Test */
	var depositary model.DepositaryCore

	depositary.HostAddress = "127.0.0.1"
	depositary.Port = "5000"
	depositary.Tag = "latest"

	instanceModel.TargetDepositary = depositary

	// TODO: complete depository schedule process

	return nil
}

/*
func: MountFilestore
description: mount the target cloud storage directory
*/
func (h *Handler) MountFilestore(ctx context.Context, instanceModel *model.InstanceModel) error {
	err := h.InstanceService.MountFilestore(ctx, instanceModel.VMID, instanceModel.TargetFilestore)

	if err != nil {
		log.Printf("Failed to mount filestore: %v\n", err.Error())
	}

	return err
}

/*
func: FetchDepositary
description: fetch the docker layer including some configuration of the app's installation from the target depositary server
*/
func (h *Handler) FetchDepositary(ctx context.Context, instanceModel *model.InstanceModel) error {
	err := h.InstanceService.FetchLayerFromDepositary(ctx, instanceModel.VMID, instanceModel.TargetDepositary)

	if err != nil {
		log.Printf("Failed to fetch layer: %v\n", err.Error())
	}

	return err
}
