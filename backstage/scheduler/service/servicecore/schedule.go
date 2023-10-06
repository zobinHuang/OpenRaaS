package servicecore

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/liushuochen/gotable"
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/utils"
)

/*
@struct: ScheduleServiceCore
@description: service core layer
*/

type ScheduleServiceCore struct {
	ConsumerDAL     model.ConsumerDAL
	ProviderDAL     model.ProviderDAL
	DepositoryDAL   model.DepositoryDAL
	FileStoreDAL    model.FileStoreDAL
	InstanceRoomDAL model.InstanceRoomDAL
	ApplicationDAL  model.ApplicationDAL
}

/*
@struct: ScheduleServiceCoreConfig
@description: used for config instance of struct ScheduleServiceCore
*/
type ScheduleServiceCoreConfig struct {
	ConsumerDAL     model.ConsumerDAL
	ProviderDAL     model.ProviderDAL
	DepositoryDAL   model.DepositoryDAL
	FileStoreDAL    model.FileStoreDAL
	InstanceRoomDAL model.InstanceRoomDAL
	ApplicationDAL  model.ApplicationDAL
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
		DepositoryDAL:   c.DepositoryDAL,
		FileStoreDAL:    c.FileStoreDAL,
		InstanceRoomDAL: c.InstanceRoomDAL,
		ApplicationDAL:  c.ApplicationDAL,
	}
}

/*
@func: ScheduleStream
@description:

	core logic of scheduling stream instance is here
*/
func (sc *ScheduleServiceCore) ScheduleStream(ctx context.Context, consumer *model.Consumer, streamInstance *model.StreamInstance) (
	*model.Provider, []model.DepositoryCoreWithInst, []model.FileStoreCoreWithInst, error) {
	// 1. 粗选节点 (所有满足 APP 基本需求的节点), 并等待获取数字资产

	providersIDList := make([]string, 0)
	providers := sc.ProviderDAL.GetProvider()
	providersInRDS := make(map[string]*model.ProviderCoreWithInst)
	for i, p := range providers {
		pInfo, err := sc.ProviderDAL.GetProviderInRDSByID(ctx, p.ClientID)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("scheduler GetProviderInRDSByID err: %s, streamInstance: %+v", err.Error(), streamInstance)
		}
		providers[i].ID = p.ClientID
		providers[i].IP = pInfo.IP
		providers[i].Port = pInfo.Port
		providers[i].Processor = pInfo.Processor
		providers[i].IsContainGPU = pInfo.IsContainGPU
		providers[i].Bandwidth = pInfo.Bandwidth
		providers[i].Latency = pInfo.Latency
		providersIDList = append(providersIDList, p.ClientID[0:5])
		providersInRDS[p.ID] = pInfo
	}

	appInfo, err := sc.ApplicationDAL.GetStreamApplicationByID(ctx, streamInstance.ApplicationID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("scheduler GetStreamApplicationByID err: %s, streamInstance: %+v", err.Error(), streamInstance)
	}

	if appInfo.FileStoreList == "" {
		return nil, nil, nil, fmt.Errorf("scheduler FileStoreList is none streamInstance: %+v", streamInstance)
	}
	var fileStoreIDList []string
	if err := json.Unmarshal([]byte(appInfo.FileStoreList), &fileStoreIDList); err != nil {
		return nil, nil, nil, fmt.Errorf("scheduler unmarshal FileStoreList fail, err: %s, streamInstance: %+v", err.Error(), streamInstance)
	}

	// todo: get depositoryList from image id
	depositoryList, err := sc.DepositoryDAL.GetDepositoryInRDS(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("scheduler GetDepositoryInRDS err: %s, streamInstance: %+v", err.Error(), streamInstance)
	}
	depositoryIDList := make([]string, 0)
	for _, d := range depositoryList {
		depositoryIDList = append(depositoryIDList, d.ID[0:5])
	}

	filestoreList, err := sc.FileStoreDAL.GetFileStoreInRDSBetweenID(ctx, fileStoreIDList)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("scheduler GetFileStoreInRDS err: %s, streamInstance: %+v", err.Error(), streamInstance)
	}
	for i := 0; i < len(fileStoreIDList); i++ {
		fileStoreIDList[i] = fileStoreIDList[i][0:5]
	}

	// 打印所有满足要求的 ID 列表
	table, err := gotable.Create("节点类型", "可使用的节点索引")
	if err != nil {
		fmt.Println("Create table failed: ", err.Error())
		return nil, nil, nil, fmt.Errorf("scheduler gotable.Create err: %s, streamInstance: %+v", err.Error(), streamInstance)
	}
	table.AddRow([]string{"服务提供节点", strings.Join(providersIDList, ",")})
	table.AddRow([]string{"内容存储节点", strings.Join(fileStoreIDList, ",")})
	table.AddRow([]string{"镜像仓库节点", strings.Join(depositoryIDList, ",")})
	log.Info("正常线上节点：")
	fmt.Println("\n", table, "\n")

	// 2. 并行获得数字资产, 并等待全部完成

	var wg sync.WaitGroup

	for _, p := range providers {
		wg.Add(1)
		go func() {
			s, err := sc.GetValueFromBlockchain(p.ClientID)
			if err == nil {
				log.Infof("【服务提供节点】数字资产获取, 资产索引: %s, 资产内容: %s", providersInRDS[p.ID].ID, s)
				log.Infof("【服务提供节点】数字资产获取 (解析后):\n%s", providersInRDS[p.ID].DetailedInfo())
			}
			wg.Done()
		}()
	}

	for _, d := range depositoryList {
		wg.Add(1)
		go func() {
			s, err := sc.GetValueFromBlockchain(d.ID)
			if err == nil {
				log.Infof("【镜像仓库节点】数字资产获取, 资产索引: %s, 资产内容: %s", d.ID, s)
				log.Infof("【镜像仓库节点】数字资产获取 (解析后):\n%s", d.DetailedInfo())
			}
			wg.Done()
		}()
	}

	for _, f := range filestoreList {
		wg.Add(1)
		go func() {
			s, err := sc.GetValueFromBlockchain(f.ID)
			if err == nil {
				log.Infof("【内容存储节点】数字资产获取, 资产索引: %s, 资产内容: %s", f.ID, s)
				log.Infof("【内容存储节点】数字资产获取 (解析后):\n%s", f.DetailedInfo())
			}
			wg.Done()
		}()
	}

	// 等待所有线程执行完毕
	wg.Wait()

	if consumer.ConsumerType == "stream" {
		consumer.T2 = time.Now()
		log.Infof("%s, 开始生成业务能力动态组合方案, ID: %s", consumer.T2.Format(utils.TIME_LAYOUT), consumer.ClientID)
		log.Infof("T2 = %s", consumer.T2.Sub(consumer.T1))
	}

	// 3. 开始进一步筛选
	providersOut := make([]*model.Provider, 0)
	fileStoresOut := make([]model.FileStoreCoreWithInst, 0)
	depositoryOut := make([]model.DepositoryCoreWithInst, 0)

	// 3.1 排除异常节点
	table, err = gotable.Create("节点 ID", "节点 IP", "计算能力", "平均历史服务质量", "带宽", "时延", "是否有 GPU", "异常服务次数")
	if err != nil {
		fmt.Println("Create table failed: ", err.Error())
		return nil, nil, nil, fmt.Errorf("scheduler gotable.Create err: %s, streamInstance: %+v", err.Error(), streamInstance)
	}
	for _, p := range providers {
		providersOut = append(providersOut, p)
		if providersInRDS[p.ID].GetAbnormalHistoryTimes() == 0 {
			table.AddRow([]string{p.ID[0:5], p.IP, fmt.Sprintf("%.2f GF", p.Processor), providersInRDS[p.ID].GetMeanHistory(),
				fmt.Sprintf("%.2f Mbps", p.Bandwidth), fmt.Sprintf("%.2f ms", p.Latency), strconv.FormatBool(p.IsContainGPU),
				fmt.Sprintf("%d", providersInRDS[p.ID].GetAbnormalHistoryTimes())})
		} else {
			//log.Infof("剔除异常服务提供节点：%s，异常次数：%d，历史信息：%s", p.ID, providersInRDS[p.ID].GetAbnormalHistoryTimes(), providersInRDS[p.ID].InstHistory)
			table.AddRow([]string{p.ID[0:5] + "*", p.IP, fmt.Sprintf("%.2f GF", p.Processor), providersInRDS[p.ID].GetMeanHistory(),
				fmt.Sprintf("%.2f Mbps", p.Bandwidth), fmt.Sprintf("%.2f ms", p.Latency), strconv.FormatBool(p.IsContainGPU),
				fmt.Sprintf("%d", providersInRDS[p.ID].GetAbnormalHistoryTimes())})
		}
	}
	log.Info("服务提供节点性能表现：")
	fmt.Println("\n", table, "\n")

	table, err = gotable.Create("节点 ID", "节点 IP", "存储能力", "平均历史服务质量", "带宽", "时延", "是否支持高性能读写", "异常服务次数")
	if err != nil {
		fmt.Println("Create table failed: ", err.Error())
		return nil, nil, nil, fmt.Errorf("scheduler gotable.Create err: %s, streamInstance: %+v", err.Error(), streamInstance)
	}
	for _, f := range filestoreList {
		fileStoresOut = append(fileStoresOut, f)
		if f.GetAbnormalHistoryTimes() == 0 {
			table.AddRow([]string{f.ID[0:5], f.IP, fmt.Sprintf("%.2f GB", f.Mem), f.GetMeanHistory(),
				fmt.Sprintf("%.2f Mbps", f.Bandwidth), fmt.Sprintf("%.2f ms", f.Latency), strconv.FormatBool(f.IsContainFastNetspeed),
				fmt.Sprintf("%d", f.GetAbnormalHistoryTimes())})
		} else {
			//log.Infof("剔除异常内容存储节点：%s，异常次数：%d，历史信息：%s", f.ID, f.GetAbnormalHistoryTimes(), f.InstHistory)
			table.AddRow([]string{f.ID[0:5] + "*", f.IP, fmt.Sprintf("%.2f GB", f.Mem), f.GetMeanHistory(),
				fmt.Sprintf("%.2f Mbps", f.Bandwidth), fmt.Sprintf("%.2f ms", f.Latency), strconv.FormatBool(f.IsContainFastNetspeed),
				fmt.Sprintf("%d", f.GetAbnormalHistoryTimes())})
		}
	}
	log.Info("内容存储节点性能表现：")
	fmt.Println("\n", table, "\n")

	table, err = gotable.Create("节点 ID", "节点 IP", "存储能力", "平均历史服务质量", "带宽", "时延", "是否支持高性能读写", "异常服务次数")
	if err != nil {
		fmt.Println("Create table failed: ", err.Error())
		return nil, nil, nil, fmt.Errorf("scheduler gotable.Create err: %s, streamInstance: %+v", err.Error(), streamInstance)
	}
	for _, d := range depositoryList {
		depositoryOut = append(depositoryOut, d)
		if d.GetAbnormalHistoryTimes() == 0 {
			table.AddRow([]string{d.ID[0:5], d.IP, fmt.Sprintf("%.2f GB", d.Mem), d.GetMeanHistory(),
				fmt.Sprintf("%.2f Mbps", d.Bandwidth), fmt.Sprintf("%.2f ms", d.Latency), strconv.FormatBool(d.IsContainFastNetspeed),
				fmt.Sprintf("%d", d.GetAbnormalHistoryTimes())})
		} else {
			//log.Infof("剔除异常镜像仓库节点：%s，异常次数：%d，历史信息：%s", d.ID, d.GetAbnormalHistoryTimes(), d.InstHistory)
			table.AddRow([]string{d.ID[0:5] + "*", d.IP, fmt.Sprintf("%.2f GB", d.Mem), d.GetMeanHistory(),
				fmt.Sprintf("%.2f Mbps", d.Bandwidth), fmt.Sprintf("%.2f ms", d.Latency), strconv.FormatBool(d.IsContainFastNetspeed),
				fmt.Sprintf("%d", d.GetAbnormalHistoryTimes())})
		}
	}
	log.Info("镜像仓库节点性能表现：")
	fmt.Println("\n", table, "\n")

	// 3.2 统计总资源量和已经使用的资源量
	totalGf := 0.0
	usedGf := 0.0
	providersRemained := make(map[string]float64)
	totalMem := 0.0
	usedMem := 0.0
	fileStoresRemained := make(map[string]float64)
	for _, p := range providersOut {
		totalGf += p.Processor
		providersRemained[p.ID] = p.Processor
	}
	for _, f := range fileStoresOut {
		totalMem += f.Mem
		fileStoresRemained[f.ID] = f.Mem
	}

	consumers := sc.ConsumerDAL.GetConsumers()
	log.Infof("consumers 数量：%d", len(consumers))
	for _, c := range consumers {
		if c.Provider == nil {
			log.Infof("c.Provider is nil, c: %+v", c)
			continue
		}
		if _, ok := providersRemained[c.Provider.ID]; ok {
			if c.Provider.IsContainGPU {
				usedGf += 5.0
				providersRemained[c.Provider.ID] -= 5.0
			} else {
				usedGf += 2.0
				providersRemained[c.Provider.ID] -= 2.0
			}
		}
		if _, ok := fileStoresRemained[c.Filestore.ID]; ok {
			usedMem += 1.0
			fileStoresRemained[c.Filestore.ID] -= 1.0
		}
	}

	table, err = gotable.Create("节点类型", "已使用资源量", "总资源量")
	if err != nil {
		fmt.Println("Create table failed: ", err.Error())
		return nil, nil, nil, fmt.Errorf("scheduler gotable.Create err: %s, streamInstance: %+v", err.Error(), streamInstance)
	}
	table.AddRow([]string{"服务提供节点", fmt.Sprintf("%.2f GF", usedGf), fmt.Sprintf("%.2f GF", totalGf)})
	table.AddRow([]string{"内容存储节点", fmt.Sprintf("%.2f GB", usedMem), fmt.Sprintf("%.2f GB", totalMem)})
	log.Info("节点资源使用情况：")
	fmt.Println("\n", table, "\n")

	// 4. 特殊处理 (支撑模型更新)
	// 4.1 选择模型
	if !sc.ConsumerDAL.IsUserOverTime(consumer.UserName) {
		// 如果申请的次数小于等于2，或者距离上上次申请的时间间隔相差30min以上，则尽可能按需分配
		log.Infof("用户 %s 申请应用小于等于2或者距离上上次申请的时间间隔相差30min以上，尽可能按需分配", consumer.UserName)
		log.Infof("appinfo: %+v", appInfo)
		if !appInfo.IsProviderReqGPU {
			// 去掉高性能的，如果有低性能且正常能分配的话
			tmp := make([]*model.Provider, 0)
			flag := false
			for _, p := range providersOut {
				if !p.IsContainGPU && providersRemained[p.ID] >= 2.0 {
					tmp = append(tmp, p)
					if providersInRDS[p.ID].GetAbnormalHistoryTimes() == 0 {
						flag = true
					}
				}
			}
			if len(tmp) > 0 && flag {
				providersOut = tmp
			}
		}
		if !appInfo.IsFileStoreReqFastNetspeed {
			// 去掉高性能的，如果有低性能且正常能分配的话
			tmp := make([]model.FileStoreCoreWithInst, 0)
			flag := false
			for _, f := range fileStoresOut {
				if !f.IsContainFastNetspeed && fileStoresRemained[f.ID] >= 1.0 {
					tmp = append(tmp, f)
					if f.GetAbnormalHistoryTimes() == 0 {
						flag = true
					}
				}
			}
			if len(tmp) > 0 && flag {
				fileStoresOut = tmp
			}
		}
		if !appInfo.IsDepositoryReqFastNetspeed {
			// 去掉高性能的，如果有低性能且正常能分配的话
			tmp := make([]model.DepositoryCoreWithInst, 0)
			flag := false
			for _, d := range depositoryOut {
				if !d.IsContainFastNetspeed {
					tmp = append(tmp, d)
					if d.GetAbnormalHistoryTimes() == 0 {
						flag = true
					}
				}
			}
			if len(tmp) > 0 && flag {
				depositoryOut = tmp
			}
		}
	} else {
		log.Infof("用户 %s 申请应用数量大于2且距离上上次申请的时间间隔相差30min以内，更换为根据资源使用情况进行服务的策略", consumer.UserName)
		if usedMem*2 >= totalMem || usedGf*2 >= totalGf {
			log.Infof("资源紧缺，使用尽力服务策略")
			if time.Now().Sub(consumer.T0) <= time.Minute && (appInfo.IsProviderReqGPU || appInfo.IsFileStoreReqFastNetspeed) {
				log.Infof("用户 %s 创建应用时间间隔过短，对其性能进行限制", consumer.ClientID)
				if (usedGf/totalGf > usedMem/totalMem || !appInfo.IsFileStoreReqFastNetspeed) && appInfo.IsProviderReqGPU {
					log.Infof("对服务提供节点性能进行限制，将高性能限制为低性能")
					appInfo.IsProviderReqGPU = false
					tmp := make([]*model.Provider, 0)
					for _, p := range providersOut {
						if !p.IsContainGPU {
							tmp = append(tmp, p)
						}
					}
					providersOut = tmp
				} else {
					log.Infof("对内容存储节点性能进行限制，将高性能限制为低性能")
					appInfo.IsFileStoreReqFastNetspeed = false
					newFilestoreList := make([]model.FileStoreCoreWithInst, 0)
					for _, f := range fileStoresOut {
						if !f.IsContainFastNetspeed {
							newFilestoreList = append(newFilestoreList, f)
						}
					}
					fileStoresOut = newFilestoreList
				}
			}
		} else {
			log.Infof("资源充足，使用性能最佳策略")
		}
	}

	sc.ConsumerDAL.UserUpdateTime(consumer.UserName, consumer.T0)

	// 4.2 正常筛选
	tmp := make([]*model.Provider, 0)
	for _, p := range providersOut {
		if appInfo.IsProviderReqGPU && p.IsContainGPU && providersRemained[p.ID] >= 5.0 {
			tmp = append(tmp, p)
		} else if !appInfo.IsProviderReqGPU && providersRemained[p.ID] >= 2.0 {
			tmp = append(tmp, p)
		}
	}
	providersOut = tmp
	if appInfo.IsDepositoryReqFastNetspeed {
		newDopositoryList := make([]model.DepositoryCoreWithInst, 0)
		for _, d := range depositoryOut {
			if d.IsContainFastNetspeed {
				newDopositoryList = append(newDopositoryList, d)
			}
		}
		depositoryOut = newDopositoryList
	}
	newFilestoreList := make([]model.FileStoreCoreWithInst, 0)
	for _, f := range fileStoresOut {
		if appInfo.IsDepositoryReqFastNetspeed && f.IsContainFastNetspeed && fileStoresRemained[f.ID] >= 1.0 {
			newFilestoreList = append(newFilestoreList, f)
		} else if !appInfo.IsDepositoryReqFastNetspeed && fileStoresRemained[f.ID] >= 1.0 {
			newFilestoreList = append(newFilestoreList, f)
		}
	}
	fileStoresOut = newFilestoreList

	// 剔除异常节点，如果有正常节点的话
	for i := 0; i < len(providersOut); i++ {
		if providersInRDS[providersOut[i].ID].GetAbnormalHistoryTimes() == 0 {
			tmp := make([]*model.Provider, 0)
			for _, p := range providersOut {
				if providersInRDS[p.ID].GetAbnormalHistoryTimes() == 0 {
					tmp = append(tmp, p)
				} else {
					log.Infof("剔除异常服务提供节点：%s，异常次数：%d，历史信息：%s", p.ID, providersInRDS[p.ID].GetAbnormalHistoryTimes(), providersInRDS[p.ID].InstHistory)
				}
			}
			providersOut = tmp
			break
		}
	}
	for i := 0; i < len(depositoryOut); i++ {
		if depositoryOut[i].GetAbnormalHistoryTimes() == 0 {
			tmp := make([]model.DepositoryCoreWithInst, 0)
			for _, d := range depositoryOut {
				if d.GetAbnormalHistoryTimes() == 0 {
					tmp = append(tmp, d)
				} else {
					log.Infof("剔除异常镜像仓库节点：%s，异常次数：%d，历史信息：%s", d.ID, d.GetAbnormalHistoryTimes(), d.InstHistory)
				}
			}
			depositoryOut = tmp
			break
		}
	}
	for i := 0; i < len(fileStoresOut); i++ {
		if fileStoresOut[i].GetAbnormalHistoryTimes() == 0 {
			tmp := make([]model.FileStoreCoreWithInst, 0)
			for _, f := range fileStoresOut {
				if f.GetAbnormalHistoryTimes() == 0 {
					tmp = append(tmp, f)
				} else {
					log.Infof("剔除异常内容存储节点：%s，异常次数：%d，历史信息：%s", f.ID, f.GetAbnormalHistoryTimes(), f.InstHistory)
				}
			}
			fileStoresOut = tmp
			break
		}
	}

	// 4.3 排序
	log.Info("硬件资源编排：")
	sc.ProviderDAL.ShowInfoFromClient(providersOut, providersInRDS)
	sc.FileStoreDAL.ShowInfoFromRDS(fileStoresOut)
	sc.DepositoryDAL.ShowInfoFromRDS(depositoryOut)

	log.Infof("select info, provider: %+v, depositoryOut: %+v, fileStoresOut: %+v", providersOut, depositoryOut, fileStoresOut)
	if len(providersOut) == 0 || len(depositoryOut) == 0 || len(fileStoresOut) == 0 {
		return nil, nil, nil, fmt.Errorf("not enough resourse to schedule")
	}

	// !!!【开始】测试时间用代码
	p_num := 10
	f_num := 10
	d_num := 10
	rand.Seed(time.Now().UnixNano())
	p_randomNumbers := make([]int, p_num)
	f_randomNumbers := make([]int, f_num)
	d_randomNumbers := make([]int, d_num)
	for i := 0; i < p_num; i++ {
		p_randomNumbers[i] = rand.Intn(100)
	}
	for i := 0; i < f_num; i++ {
		f_randomNumbers[i] = rand.Intn(100)
	}
	for i := 0; i < d_num; i++ {
		d_randomNumbers[i] = rand.Intn(100)
	}
	sort.Ints(p_randomNumbers)
	sort.Ints(f_randomNumbers)
	sort.Ints(d_randomNumbers)
	// !!!【结束】测试时间用代码

	// 找到第一个不异常的节点，如果全为异常节点则选择第一个
	for i := 0; i < len(providersOut); i++ {
		if providersInRDS[providersOut[i].ID].GetAbnormalHistoryTimes() == 0 {
			providersOut[0] = providersOut[i]
			break
		}
	}
	for i := 0; i < len(depositoryOut); i++ {
		if depositoryOut[i].GetAbnormalHistoryTimes() == 0 {
			depositoryOut[0] = depositoryOut[i]
			break
		}
	}
	for i := 0; i < len(fileStoresOut); i++ {
		if fileStoresOut[i].GetAbnormalHistoryTimes() == 0 {
			fileStoresOut[0] = fileStoresOut[i]
			break
		}
	}

	table, err = gotable.Create("节点类型", "已使用资源量", "总资源量")
	if err != nil {
		fmt.Println("Create table failed: ", err.Error())
		return nil, nil, nil, fmt.Errorf("scheduler gotable.Create err: %s, streamInstance: %+v", err.Error(), streamInstance)
	}
	if providersOut[0].IsContainGPU {
		usedGf += 5.0
	} else {
		usedGf += 2.0
	}
	usedMem += 1.0
	table.AddRow([]string{"服务提供节点", fmt.Sprintf("%.2f GF", usedGf), fmt.Sprintf("%.2f GF", totalGf)})
	table.AddRow([]string{"内容存储节点", fmt.Sprintf("%.2f GB", usedMem), fmt.Sprintf("%.2f GB", totalMem)})
	log.Info("组合后节点资源使用情况：")
	fmt.Println("\n", table, "\n")

	log.Infof("【业务能力动态组合方案】")
	fmt.Printf("1. 计算 (Computation) 原子服务:\n%s", providersInRDS[providersOut[0].ID].DetailedInfo())
	fmt.Printf("2. 运行时文件 (Runtime Files) 原子服务:\n%s", fileStoresOut[0].DetailedInfo())
	fmt.Printf("3. 运行时环境 (Runtime Environment) 原子服务:\n%s", depositoryOut[0].DetailedInfo())

	consumer.Provider = providersOut[0]
	consumer.Filestore = &model.FileStore{
		FileStoreCore: model.FileStoreCore{ID: fileStoresOut[0].ID},
	}

	return providersOut[0], depositoryOut, fileStoresOut, nil
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

// GetStreamInstanceRoomByInstanceID obtain StreamInstanceRoom by instance id
func (sc *ScheduleServiceCore) GetStreamInstanceRoomByInstanceID(id string) (*model.StreamInstanceRoom, error) {
	return sc.InstanceRoomDAL.GetInstanceRoomByInstanceID(nil, id)
}

// SetValueToBlockchain set key value to blockchain
func (sc *ScheduleServiceCore) SetValueToBlockchain(key, value string) error {
	//log.Infof("SetValueToBlockchain key: %s, value: %s", key, value)
	url := "http://192.168.0.109:5001/api/set_value"
	// 准备请求的数据
	data := map[string]string{
		"key":   key,
		"value": value,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("JSON marshaling failed:", err)
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("POST request failed:", err)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("POST request failed with status:", resp.StatusCode)
		return fmt.Errorf("POST request failed with status: %+v", resp.StatusCode)
	}
	return nil
}

// GetValueFromBlockchain obtain value by key
func (sc *ScheduleServiceCore) GetValueFromBlockchain(key string) (string, error) {
	//log.Infof("GetValueFromBlockchain key: %s", key)
	url := "http://192.168.0.109:5001/api/get_value?key=" + key
	resp, err := http.Get(url)
	if err != nil {
		fmt.Println("GET request failed:", err)
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response body:", err)
		return "", err
	}

	type Response struct {
		Key     string `json:"key"`
		Message string `json:"message"`
		Value   string `json:"value"`
	}

	var response Response
	err = json.Unmarshal(body, &response)
	if err != nil {
		fmt.Println("Failed to parse JSON response:", err)
		return "", err
	}

	//log.Infof("GetValueFromBlockchain value: %s", response.Value)

	return response.Value, nil
}

// Clear delete all
func (sc *ScheduleServiceCore) Clear() {
	sc.ConsumerDAL.Clear()
	sc.ProviderDAL.Clear()
	sc.DepositoryDAL.Clear()
	sc.FileStoreDAL.Clear()
	sc.InstanceRoomDAL.Clear()
	sc.ApplicationDAL.Clear()
}
