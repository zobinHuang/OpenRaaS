package dal

import (
	"context"
	"fmt"
	"github.com/liushuochen/gotable"
	"sort"
	"strconv"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"

	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
)

/*
@struct: ProviderDAL
@description: DAL layer
*/
type ProviderDAL struct {
	DB           *gorm.DB
	ProviderList map[string]*model.Provider
}

/*
@struct: ProviderDALConfig
@description: used for config instance of struct ProviderDAL
*/
type ProviderDALConfig struct {
	DB *gorm.DB
}

/*
@func: NewProviderDAL
@description:

	create, config and return an instance of struct ProviderDAL
*/
func NewProviderDAL(c *ProviderDALConfig) model.ProviderDAL {
	pdal := &ProviderDAL{}

	pdal.ProviderList = make(map[string]*model.Provider)
	pdal.DB = c.DB

	return pdal
}

/*
@func: CreateProvider
@description:

	insert a new provider to provider list
*/
func (d *ProviderDAL) CreateProvider(ctx context.Context, provider *model.Provider) {
	d.ProviderList[provider.ClientID] = provider
}

/*
@func: DeleteProvider
@description:

	delete the specified provider from provider list
*/
func (d *ProviderDAL) DeleteProvider(ctx context.Context, providerID string) {
	delete(d.ProviderList, providerID)
}

func (d *ProviderDAL) GetProvider() []*model.Provider {
	providers := make([]*model.Provider, 0, 0)
	for _, value := range d.ProviderList {
		providers = append(providers, value)
	}
	return providers
}

// CreateProviderInRDS create provider core info to rds
func (d *ProviderDAL) CreateProviderInRDS(ctx context.Context, provider *model.ProviderCoreWithInst) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("provider_core_with_insts").Create(provider).Error; err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"provider": provider,
		}).Warn("Failed to create provider core info to rds")
		return err
	}

	return nil
}

// DeleteProviderInRDSByID
func (d *ProviderDAL) DeleteProviderInRDSByID(ctx context.Context, id string) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("provider_core_with_insts").Where("id=?", id).Delete(&model.ProviderCoreWithInst{}).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Warn("Failed to delete provider core info by id in rds")
		return err
	}

	return nil
}

// GetProviderInRDS obtain all provider core info from rds
func (d *ProviderDAL) GetProviderInRDS(ctx context.Context) ([]model.ProviderCoreWithInst, error) {
	var infos []model.ProviderCoreWithInst

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("provider_core_with_insts").Find(&infos).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to obtain all provider core info from rds")
		return nil, err
	}

	return infos, nil
}

// GetProviderInRDSByID get provider core info by id from rds
func (d *ProviderDAL) GetProviderInRDSByID(ctx context.Context, id string) (*model.ProviderCoreWithInst, error) {
	var info model.ProviderCoreWithInst
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("provider_core_with_insts").Where("id = ?", id).First(&info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Warn("Failed to get provider core info by id from rds")
		return nil, err
	}
	return &info, nil
}

// UpdateProviderInRDSByID update provider core info by id in rds
func (d *ProviderDAL) UpdateProviderInRDSByID(ctx context.Context, provider *model.ProviderCoreWithInst) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("provider_core_with_insts").Where("id=?", provider.ID).Updates(provider).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    provider.ID,
		}).Warn("Failed to update provider core info by id in rds")
		return err
	}
	return nil
}

// Clear delete all
func (d *ProviderDAL) Clear() {
	if err := d.DB.Exec("DELETE FROM public.provider_core_with_insts").Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Fail to clear provider_core_with_insts table")
	}
}

func (d *ProviderDAL) ShowInfoFromRDS(providers []model.ProviderCoreWithInst) {
	sort.Slice(providers, func(i, j int) bool {
		if (providers[i].IsContainGPU && !providers[j].IsContainGPU) ||
			(!providers[i].IsContainGPU && providers[j].IsContainGPU) {
			return providers[i].IsContainGPU
		}
		historyi := providers[i].GetMeanHistory()
		historyj := providers[j].GetMeanHistory()
		if historyi == "" {
			return true
		}
		if historyj == "" {
			return false
		}
		f1, err := strconv.ParseFloat(historyi[0:len(historyi)-3], 64)
		if err != nil {
			fmt.Println("ShowInfoFromRDS providers strconv.ParseFloat(historyi[0:len(historyi)-3], 64) 转换失败:", err)
			return false
		}
		f2, err := strconv.ParseFloat(historyj[0:len(historyj)-3], 64)
		if err != nil {
			fmt.Println("ShowInfoFromRDS providers strconv.ParseFloat(historyj[0:len(historyj)-3], 64) 转换失败:", err)
			return false
		}
		if f1 != f2 {
			return f1 < f2
		}
		return 5.0/providers[i].Bandwidth+providers[i].Latency <= 5.0/providers[j].Bandwidth+providers[j].Latency
	})
	table, err := gotable.Create("节点 ID", "节点 IP", "计算能力", "平均历史服务质量", "带宽", "时延", "是否有 GPU", "异常服务次数")
	if err != nil {
		fmt.Println("ShowInfoFromRDS providers Create table failed: ", err.Error())
		return
	}
	for _, p := range providers {
		if p.GetAbnormalHistoryTimes() == 0 {
			table.AddRow([]string{p.ID[0:5], p.IP, fmt.Sprintf("%.2f GF", p.Processor), p.GetMeanHistory(),
				fmt.Sprintf("%.2f Mbps", p.Bandwidth), fmt.Sprintf("%.2f ms", p.Latency), strconv.FormatBool(p.IsContainGPU), fmt.Sprintf("%d", p.GetAbnormalHistoryTimes())})
		} else {
			table.AddRow([]string{p.ID[0:5] + "*", p.IP, fmt.Sprintf("%.2f GF", p.Processor), p.GetMeanHistory(),
				fmt.Sprintf("%.2f Mbps", p.Bandwidth), fmt.Sprintf("%.2f ms", p.Latency), strconv.FormatBool(p.IsContainGPU), fmt.Sprintf("%d", p.GetAbnormalHistoryTimes())})
		}
	}
	log.Info("服务提供节点性能表现：")
	fmt.Println("\n", table, "\n")
}

func (d *ProviderDAL) ShowInfoFromClient(providers []*model.Provider, providersInRDS map[string]*model.ProviderCoreWithInst) {
	sort.Slice(providers, func(i, j int) bool {
		if (providers[i].IsContainGPU && !providers[j].IsContainGPU) ||
			(!providers[i].IsContainGPU && providers[j].IsContainGPU) {
			return providers[i].IsContainGPU
		}
		historyi := providersInRDS[providers[i].ID].GetMeanHistory()
		historyj := providersInRDS[providers[j].ID].GetMeanHistory()
		if historyi == "" {
			return true
		}
		if historyj == "" {
			return false
		}
		f1, err := strconv.ParseFloat(historyi[0:len(historyi)-3], 64)
		if err != nil {
			fmt.Println("ShowInfoFromClient providers strconv.ParseFloat(historyi[0:len(historyi)-3], 64) 转换失败:", err)
			return false
		}
		f2, err := strconv.ParseFloat(historyj[0:len(historyj)-3], 64)
		if err != nil {
			fmt.Println("ShowInfoFromClient providers strconv.ParseFloat(historyj[0:len(historyj)-3], 64) 转换失败:", err)
			return false
		}
		if f1 != f2 {
			return f1 < f2
		}
		return 5.0/providers[i].Bandwidth+providers[i].Latency <= 5.0/providers[j].Bandwidth+providers[j].Latency
	})
	table, err := gotable.Create("节点 ID", "节点 IP", "计算能力", "平均历史服务质量", "带宽", "时延", "是否有 GPU", "异常服务次数")
	if err != nil {
		fmt.Println("ShowInfoFromClient providers Create table failed: ", err.Error())
		return
	}
	for _, p := range providers {
		if providersInRDS[p.ID].GetAbnormalHistoryTimes() == 0 {
			table.AddRow([]string{p.ID[0:5], p.IP, fmt.Sprintf("%.2f GF", p.Processor), providersInRDS[p.ID].GetMeanHistory(),
				fmt.Sprintf("%.2f Mbps", p.Bandwidth), fmt.Sprintf("%.2f ms", p.Latency),
				strconv.FormatBool(p.IsContainGPU), fmt.Sprintf("%d", providersInRDS[p.ID].GetAbnormalHistoryTimes())})
		} else {
			table.AddRow([]string{p.ID[0:5] + "*", p.IP, fmt.Sprintf("%.2f GF", p.Processor), providersInRDS[p.ID].GetMeanHistory(),
				fmt.Sprintf("%.2f Mbps", p.Bandwidth), fmt.Sprintf("%.2f ms", p.Latency),
				strconv.FormatBool(p.IsContainGPU), fmt.Sprintf("%d", providersInRDS[p.ID].GetAbnormalHistoryTimes())})
		}
	}
	log.Info("服务提供节点性能表现：")
	fmt.Println("\n", table, "\n")
}
