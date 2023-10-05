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
@struct: DepositoryDAL
@description: DAL layer
*/
type DepositoryDAL struct {
	DB             *gorm.DB
	DepositoryList map[string]*model.Depository
}

/*
@struct: DepositoryDALConfig
@description: used for config instance of struct DepositoryDAL
*/
type DepositoryDALConfig struct {
	DB *gorm.DB
}

/*
@func: NewDepositoryDAL
@description:

	create, config and return an instance of struct DepositoryDAL
*/
func NewDepositoryDAL(c *DepositoryDALConfig) model.DepositoryDAL {
	ddal := &DepositoryDAL{}

	ddal.DepositoryList = make(map[string]*model.Depository)
	ddal.DB = c.DB

	return ddal
}

/*
@func: CreateDepository
@description:

	insert a new depository to depository list
*/
func (d *DepositoryDAL) CreateDepository(ctx context.Context, depository *model.Depository) {
	d.DepositoryList[depository.ClientID] = depository
}

/*
@func: DeleteDepository
@description:

	delete the specified depository from depository list
*/
func (d *DepositoryDAL) DeleteDepository(ctx context.Context, depositoryID string) {
	delete(d.DepositoryList, depositoryID)
}

// CreateDepositoryInRDS create depository core info to rds
func (d *DepositoryDAL) CreateDepositoryInRDS(ctx context.Context, info *model.DepositoryCoreWithInst) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depository_core_with_insts").Create(info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"info":  info,
		}).Warn("Failed to create depository core info to rds")
		return err
	}

	return nil
}

// DeleteDepositoryInRDSByID delete depository core info by id in rds
func (d *DepositoryDAL) DeleteDepositoryInRDSByID(ctx context.Context, id string) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depository_core_with_insts").Where("id=?", id).Delete(&model.DepositoryCoreWithInst{}).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Warn("Failed to delete depository core info by id in rds")
		return err
	}
	return nil
}

// GetDepositoryInRDS obtain all depository core info from rds
func (d *DepositoryDAL) GetDepositoryInRDS(ctx context.Context) ([]model.DepositoryCoreWithInst, error) {
	var infos []model.DepositoryCoreWithInst

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depository_core_with_insts").Find(&infos).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to obtain all depository core info from rds")
		return nil, err
	}

	return infos, nil
}

// GetDepositoryInRDSByID get depository core info by id from rds
func (d *DepositoryDAL) GetDepositoryInRDSByID(ctx context.Context, id string) (*model.DepositoryCoreWithInst, error) {
	var info model.DepositoryCoreWithInst
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depository_core_with_insts").Where("id = ?", id).First(&info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Warn("Failed to get depository core info by id from rds")
		return nil, err
	}
	return &info, nil
}

// UpdateDepositoryInRDSByID update depository core info by id in rds
func (d *DepositoryDAL) UpdateDepositoryInRDSByID(ctx context.Context, info *model.DepositoryCoreWithInst) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depository_core_with_insts").Where("id=?", info.ID).Updates(info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    info.ID,
		}).Warn("Failed to update depository core info by id in rds")
		return err
	}
	return nil
}

// GetDepositoryBetweenIDInRDS get depository core info Between id from rds
func (d *DepositoryDAL) GetDepositoryBetweenIDInRDS(ctx context.Context, ids []string) ([]model.DepositoryCoreWithInst, error) {
	var infos []model.DepositoryCoreWithInst

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("depository_core_with_insts").Where("id IN (?)", ids).Find(&infos).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to obtain all depository core info from rds")
		return nil, err
	}

	return infos, nil
}

// Clear delete all
func (d *DepositoryDAL) Clear() {
	if err := d.DB.Exec("DELETE FROM public.depository_core_with_insts").Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Fail to clear depository core table")
	}
}

func (d *DepositoryDAL) ShowInfoFromRDS(depositories []model.DepositoryCoreWithInst) {
	sort.Slice(depositories, func(i, j int) bool {
		if (depositories[i].IsContainFastNetspeed && !depositories[j].IsContainFastNetspeed) ||
			(!depositories[i].IsContainFastNetspeed && depositories[j].IsContainFastNetspeed) {
			return depositories[i].IsContainFastNetspeed
		}
		if depositories[i].Bandwidth != depositories[j].Bandwidth {
			return depositories[i].Bandwidth > depositories[j].Bandwidth
		}
		historyi := depositories[i].GetMeanHistory()
		historyj := depositories[j].GetMeanHistory()
		if historyi == "" {
			return true
		}
		if historyj == "" {
			return false
		}
		f1, err := strconv.ParseFloat(historyi[0:len(historyi)-4], 64)
		if err != nil {
			fmt.Println("ShowInfoFromRDS fileStores strconv.ParseFloat(historyi[0:len(historyi)-3], 64) 转换失败:", err)
			return false
		}
		f2, err := strconv.ParseFloat(historyj[0:len(historyj)-4], 64)
		if err != nil {
			fmt.Println("ShowInfoFromRDS fileStores strconv.ParseFloat(historyj[0:len(historyj)-3], 64) 转换失败:", err)
			return false
		}
		return f1 <= f2
	})
	table, err := gotable.Create("节点 ID", "节点 IP", "存储能力", "平均历史服务质量", "带宽", "时延", "是否支持高性能读写", "异常服务次数")
	if err != nil {
		fmt.Println("ShowInfoFromRDS DepositoryDAL Create table failed: ", err.Error())
		return
	}
	for _, d := range depositories {
		if d.GetAbnormalHistoryTimes() == 0 {
			table.AddRow([]string{d.ID[0:5], d.IP, fmt.Sprintf("%.2f GB", d.Mem), d.GetMeanHistory(),
				fmt.Sprintf("%.2f Mbps", d.Bandwidth), fmt.Sprintf("%.2f ms", d.Latency), strconv.FormatBool(d.IsContainFastNetspeed),
				fmt.Sprintf("%d", d.GetAbnormalHistoryTimes())})
		} else {
			table.AddRow([]string{d.ID[0:5] + "*", d.IP, fmt.Sprintf("%.2f GB", d.Mem), d.GetMeanHistory(),
				fmt.Sprintf("%.2f Mbps", d.Bandwidth), fmt.Sprintf("%.2f ms", d.Latency), strconv.FormatBool(d.IsContainFastNetspeed),
				fmt.Sprintf("%d", d.GetAbnormalHistoryTimes())})
		}
	}
	log.Info("镜像仓库节点性能表现：")
	fmt.Println("\n", table, "\n")
}
