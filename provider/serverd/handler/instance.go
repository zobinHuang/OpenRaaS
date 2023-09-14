package handler

import (
	"context"
	"fmt"
	"net/http"
	"serverd/model"
	"serverd/model/apperrors"
	"serverd/utils"
	"strconv"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

/*
	func: CreateInstance
	description: handler for endpoint "/api/daemon/createinstance"
*/
func (h *Handler) CreateInstance(c *gin.Context) {
	instanceModel := &model.InstanceModel{}
	log.Printf("Create instance with context: %+v", c)
	if ok := bindData(c, &instanceModel); !ok {
		log.WithFields(log.Fields{
			"Client Address": c.Request.Host,
		}).Error("Failed to bind data, abandoned")
		return
	}

	log.Printf("%+v", instanceModel)

	ctx := c.Request.Context()

	err := h.SelectFilestore(ctx, instanceModel)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": "cannot select filestore",
		})
	}

	err = h.InstanceService.NewVMID(ctx, instanceModel)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": "cannot generate new vm id",
		})
	}

	err = h.MountFilestore(ctx, instanceModel)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": "cannot mount filestore",
		})
	}

	err = h.CreateInstanceWithModel(ctx, instanceModel)
	if err != nil {
		c.JSON(apperrors.Status(err), gin.H{
			"error": "cannot create instance",
		})
	}

	// return http_ok if success
	c.JSON(http.StatusOK, gin.H{
		"vmid": strconv.Itoa(instanceModel.VMID),
	})
}

func (h *Handler) CheckInstanceByVMID(c *gin.Context) {
	vmid, ok := c.GetQuery("vmid")
	if !ok {
		log.WithFields(log.Fields{
			"Client Address": c.Request.Host,
		}).Error("Failed to extract vmid, invalid http connection request, abandoned")
		c.JSON(http.StatusBadRequest, gin.H{"error": "need vmid"})
		return
	}

	var execCmd string
	var params []string
	// docker logs --tail 10 $(docker ps -qf "name=appvm2")
	execCmd = "docker"
	params = append(params, "logs")
	params = append(params, "--tail")
	params = append(params, "10")
	params = append(params, "$(docker ps -qf 'name=appvm"+vmid+"')")

	ret, err := utils.RunShellWithReturn(execCmd, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})

	}
	// return http_ok if success
	c.JSON(http.StatusOK, gin.H{
		"log": ret,
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
		return fmt.Errorf("instance startup failed")
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
		return fmt.Errorf("no instance specified")
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
	list_len := len(instanceModel.FilestoreList)

	if list_len == 0 {
		log.Printf("Empty filestore list!")

		var filestore model.FilestoreCore
		filestore.HostAddress = "kb109.dynv6.net"
		filestore.Port = "7189"
		filestore.Protocol = "davfs"
		filestore.Username = "kb109"
		filestore.Password = "Xusir666!"
		filestore.Directory = "/public_hdd/game/PC/dcwine"
		// filestore.Directory = "/storage_ssd/6G/dcwine"

		instanceModel.TargetFilestore = filestore
	} else {
		// index := rand.Intn(list_len)
		index := 0
		instanceModel.TargetFilestore = instanceModel.FilestoreList[index]
	}

	return nil
}

/*
	func: SelectDepository
	description: check network connection status between host and depository servers, and select one as schedule target
*/
func (h *Handler) SelectDepository(ctx context.Context, instanceModel *model.InstanceModel) error {
	/* Test */

	list_len := len(instanceModel.DepositoryList)

	if list_len == 0 {
		log.Printf("Empty depository list!")

		var depository model.DepositoryCore
		depository.HostAddress = "127.0.0.1"
		depository.Port = "5000"
		depository.Tag = "latest"

		instanceModel.TargetDepository = depository
	} else {
		// index := rand.Intn(list_len)
		index := 0
		instanceModel.TargetDepository = instanceModel.DepositoryList[index]
	}

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
	func: FetchDepository
	description: fetch the docker layer including some configuration of the app's installation from the target depository server
*/
func (h *Handler) FetchDepository(ctx context.Context, instanceModel *model.InstanceModel) error {
	err := h.InstanceService.FetchLayerFromDepository(ctx, instanceModel.VMID, instanceModel.TargetDepository, instanceModel.ImageName)

	if err != nil {
		log.Printf("Failed to fetch layer: %v\n", err.Error())
	}

	return err
}
