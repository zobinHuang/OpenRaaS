package service

import (
	"context"
	"errors"
	"fmt"
	"net"
	"serverd/model"
	"serverd/utils"
	"strconv"
	"strings"
)

/*
	struct: InstanceService
	description: visual machine instance service
*/
type InstanceService struct {
	instances  map[int]*model.InstanceModel // store existed wine container, {vmid: InstanceModel}
	image_name string                       // setted in FetchLayerFromDepositary, used in LaunchInstance, these two func should be called in order
}

/*
	struct: InstanceServiceConfig
	description: used for config instance of struct InstanceService
*/
type InstanceServiceConfig struct {
}

/*
	func: NewInstanceService
	description: create, config and return an instance of struct InstanceService
*/
func NewInstanceService(c *InstanceServiceConfig) model.InstanceService {
	instanceService := &InstanceService{}

	instanceService.instances = make(map[int]*model.InstanceModel)

	return instanceService

}

func getClientIp() (string, error) {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		return "", err
	}

	for _, address := range addrs {
		// Check whether the IP address is a loopback address
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String(), nil
			}

		}
	}

	return "", errors.New("cannot find the client ip address")
}

/*
	func: NewVMID
	description: find a usable vmid
*/
func (c *InstanceService) NewVMID(ctx context.Context, instanceModel *model.InstanceModel) error {
	// find out an unused vm_id
	for i := 1; i < 100; i++ {
		if _, ok := c.instances[i]; !ok {
			fmt.Printf("will use vmid %d\n", i)
			instanceModel.VMID = i
			break
		}
	}
	if instanceModel.VMID == 0 {
		return fmt.Errorf("no spare vmid")
	}
	return nil
}

/*
	func: LaunchInstance
	description: run a wine container on the host, and return a channel instance used for closing the container
*/
func (c *InstanceService) LaunchInstance(ctx context.Context, instanceModel *model.InstanceModel) chan struct{} {
	/* The current version can only be deployed on Linux */

	var execCmd string
	var params []string

	// Add Exec command
	// assuming pwd is "~/winecontainer/serverd/", and target shell at "~/winecontainer/winetools/run-wine.sh"
	execCmd = "sh"

	// Add params
	params = append(params, "../winetools/run-wine.sh")
	params = append(params, c.image_name)
	params = append(params, strconv.Itoa(instanceModel.VMID))
	if strings.HasPrefix(instanceModel.AppPath, "/") {
		params = append(params, instanceModel.AppPath)
	} else {
		params = append(params, "/"+instanceModel.AppPath)
	}
	params = append(params, instanceModel.AppFile)
	params = append(params, "'"+instanceModel.AppName+"'")
	params = append(params, instanceModel.HWKey)
	params = append(params, []string{strconv.Itoa(instanceModel.ScreenWidth), strconv.Itoa(instanceModel.ScreenHeight)}...)

	ip, err := getClientIp()
	if err != nil {
		fmt.Println(err)
		ip = "192.168.0.109"
	} else {
		fmt.Println("Found ip: ", ip)
	}
	params = append(params, ip)

	params = append(params, strconv.Itoa(instanceModel.FPS))

	vcodec := instanceModel.VCodec
	if strings.Contains(vcodec, "264") {
		vcodec = "h264"
	} else if strings.Contains(vcodec, "vp") {
		vcodec = "vpx"
	}
	params = append(params, vcodec)

	// TODO 等后台完善 option
	// instanceModel.WineOption = "test.flv"
	params = append(params, instanceModel.WineOption)

	done := utils.RunShell(execCmd, params)

	if done != nil {
		instanceModel.Done = done
		c.instances[instanceModel.VMID] = instanceModel
	}

	return done
}

/*
	func: DeleteInstance
	description: delete a target wine container with vmid
*/
func (c *InstanceService) DeleteInstance(ctx context.Context, vmid int) error {
	instanceModel, ok := c.instances[vmid]
	if !ok {
		return fmt.Errorf("cannot find target VM")
	}

	done := instanceModel.Done
	close(done)

	var err error

	// remove container
	var execCmd string
	var params []string
	execCmd = "docker"
	params = append(params, "container")
	params = append(params, "stop")
	params = append(params, "appvm"+strconv.Itoa(vmid))
	ret := utils.RunShellWithReturn(execCmd, params)
	if ret == "" {
		fmt.Printf("Failed to remove container %s\n", params[0])
		// return fmt.Errorf("Cannont umount remote dir.")
		err = fmt.Errorf("cannot umount remote dir")
	}
	fmt.Printf("Succeed to remove container %s\n", params[0])

	// umount remote dir
	params = []string{}
	execCmd = "umount"
	params = append(params, "../winetools/apps/point"+strconv.Itoa(vmid))

	ret = utils.RunShellWithReturn(execCmd, params)
	if ret == "" {
		fmt.Printf("Failed to umount %s\n", params[0])
		// return fmt.Errorf("Cannont umount remote dir.")
		err = fmt.Errorf("cannot umount remote dir")
	}
	fmt.Printf("Succeed to umount %s\n", params[0])

	delete(c.instances, vmid)

	return err
}

/*
	func: DeleteInstanceByInstanceid
	description: delete a target wine container with Instanceid
*/
func (c *InstanceService) DeleteInstanceByInstanceid(ctx context.Context, Instanceid string) error {
	vmid := -1

	// find out target instance
	for key, value := range c.instances {
		if value.Instanceid == Instanceid {
			vmid = key
		}
	}

	if vmid == -1 {
		return fmt.Errorf("cannot find target instance")
	}

	err := c.DeleteInstance(ctx, vmid)

	return err
}

/*
	func: DeleteAllInstance
	description: used when shutting down this daemon
*/
func (c *InstanceService) DeleteAllInstance(ctx context.Context) error {
	fmt.Println("invoked(DeleteAllInstance)")
	var err error
	for vmid := range c.instances {
		ret := c.DeleteInstance(ctx, vmid)
		if ret != nil {
			err = fmt.Errorf("cannot remove all instance")
		}
	}

	// make sure 'winecontainer/wintools/apps' is clear
	// umount all
	var execCmd string
	var params []string
	execCmd = "umount"
	fmt.Println("Umount apps/point* to make sure all points are clear. No problem if there is any error outputs.")
	params = append(params, "../winetools/apps/*")
	utils.RunShellWithReturn(execCmd, params)

	return err
}

/*
	func: MountFilestore
	description: mount the target cloud storage directory
*/
func (c *InstanceService) MountFilestore(ctx context.Context, vmid int, filestore model.FilestoreCore) error {
	/* The current version can only be deployed on Linux */

	var execCmd string
	var params []string

	// Add Exec command
	// assuming pwd is "~/winecontainer/serverd/", and target shell at "~/winecontainer/winetools/auto-mount.exp"
	execCmd = "expect"

	// Add params
	params = append(params, "../winetools/auto-mount.exp")
	params = append(params, strconv.Itoa(vmid))
	params = append(params, filestore.Protocal)
	params = append(params, filestore.HostAddress+":"+filestore.Port)
	params = append(params, filestore.Directory)
	params = append(params, filestore.Username)
	params = append(params, filestore.Password)

	// utils.RunShell(execCmd, params)
	ret := utils.RunShellWithReturn(execCmd, params)
	if ret == "" {
		return fmt.Errorf("cannot mount remote dir")
	}

	return nil
}

/*
	func: FetchLayerFromDepositary
	description: fetch the docker layer including some configuration of the app's installation from the target depositary server
*/
func (c *InstanceService) FetchLayerFromDepositary(ctx context.Context, vmid int, depositary model.DepositaryCore) error {
	// TODO complete model.Depositary

	var execCmd string
	var params []string

	execCmd = "docker"
	params = append(params, "pull")

	c.image_name = depositary.HostAddress + ":" + depositary.Port + "/dcwine"
	if depositary.Tag != "" {
		// TODO check this
		c.image_name = c.image_name + ":" + depositary.Tag
	}
	params = append(params, c.image_name)

	ret := utils.RunShellWithReturn(execCmd, params)

	var err error
	if ret == "" {
		fmt.Printf("Failed to pull image\n")
		err = fmt.Errorf("cannot pull image")
	}

	return err
}
