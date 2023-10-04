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
@struct: FileStoreDAL
@description: DAL layer
*/
type FileStoreDAL struct {
	DB            *gorm.DB
	FileStoreList map[string]*model.FileStore
}

/*
@struct: FileStoreDALConfig
@description: used for config instance of struct FileStoreDAL
*/
type FileStoreDALConfig struct {
	DB *gorm.DB
}

/*
@func: NewFileStoreDAL
@description:

	create, config and return an instance of struct FileStoreDAL
*/
func NewFileStoreDAL(c *FileStoreDALConfig) model.FileStoreDAL {
	ddal := &FileStoreDAL{}

	ddal.FileStoreList = make(map[string]*model.FileStore)
	ddal.DB = c.DB

	return ddal
}

/*
@func: CreateFileStore
@description:

	insert a new filestore to FileStore list
*/
func (d *FileStoreDAL) CreateFileStore(ctx context.Context, filestore *model.FileStore) {
	d.FileStoreList[filestore.ClientID] = filestore
}

/*
@func: DeleteFileStore
@description:

	delete the specified filestore from FileStore list
*/
func (d *FileStoreDAL) DeleteFileStore(ctx context.Context, filestoreID string) {
	delete(d.FileStoreList, filestoreID)
}

// CreateFileStoreInRDS create file store core info to rds
func (d *FileStoreDAL) CreateFileStoreInRDS(ctx context.Context, info *model.FileStoreCoreWithInst) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("file_store_core_with_insts").Create(info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"info":  info,
		}).Warn("Failed to create file store core info to rds")
		return err
	}

	return nil
}

// DeleteFileStoreInRDSByID delete file store core info by id in rds
func (d *FileStoreDAL) DeleteFileStoreInRDSByID(ctx context.Context, id string) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("file_store_core_with_insts").Where("id=?", id).Delete(&model.FileStoreCoreWithInst{}).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Warn("Failed to delete file store core info by id in rds")
		return err
	}

	return nil
}

// GetFileStoreInRDS obtain all file store core info from rds
func (d *FileStoreDAL) GetFileStoreInRDS(ctx context.Context) ([]model.FileStoreCoreWithInst, error) {
	var infos []model.FileStoreCoreWithInst

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("file_store_core_with_insts").Find(&infos).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to obtain all file store core info from rds")
		return nil, err
	}

	return infos, nil
}

// GetFileStoreInRDSByID get file store core info by id from rds
func (d *FileStoreDAL) GetFileStoreInRDSByID(ctx context.Context, id string) (*model.FileStoreCoreWithInst, error) {
	var info model.FileStoreCoreWithInst
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("file_store_core_with_insts").Where("id = ?", id).First(&info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    id,
		}).Warn("Failed to get file store core info by id from rds")
		return nil, err
	}
	return &info, nil
}

// UpdateFileStoreInRDSByID update file store core info by id in rds
func (d *FileStoreDAL) UpdateFileStoreInRDSByID(ctx context.Context, info *model.FileStoreCoreWithInst) error {
	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("file_store_core_with_insts").Where("id=?", info.ID).Updates(info).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"id":    info.ID,
		}).Warn("Failed to update file store core info by id in rds")
		return err
	}
	return nil
}

// GetFileStoreInRDSBetweenID get file store core info between id from rds
func (d *FileStoreDAL) GetFileStoreInRDSBetweenID(ctx context.Context, ids []string) ([]model.FileStoreCoreWithInst, error) {
	var infos []model.FileStoreCoreWithInst

	// initialize context
	tx := d.DB.WithContext(ctx)

	// query in database
	if err := tx.Table("file_store_core_with_insts").Where("id IN (?)", ids).Find(&infos).Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("Failed to obtain all file store core info from rds")
		return nil, err
	}

	return infos, nil
}

// Clear delete all
func (d *FileStoreDAL) Clear() {
	if err := d.DB.Exec("DELETE FROM public.file_store_core_with_insts").Error; err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Error("Fail to clear file_store_core_with_insts table")
	}
}

func (d *FileStoreDAL) ShowInfoFromRDS(fileStores []model.FileStoreCoreWithInst) {
	sort.Slice(fileStores, func(i, j int) bool {
		if fileStores[i].IsContainFastNetspeed {
			return true
		}
		if fileStores[j].IsContainFastNetspeed {
			return false
		}
		historyi := fileStores[i].GetMeanHistory()
		historyj := fileStores[j].GetMeanHistory()
		if historyi == "" {
			return true
		}
		if historyj == "" {
			return false
		}
		f1, err := strconv.ParseFloat(historyi[0:len(historyi)-3], 64)
		if err != nil {
			fmt.Println("ShowInfoFromRDS fileStores strconv.ParseFloat(historyi[0:len(historyi)-3], 64) 转换失败:", err)
			return false
		}
		f2, err := strconv.ParseFloat(historyj[0:len(historyj)-3], 64)
		if err != nil {
			fmt.Println("ShowInfoFromRDS fileStores strconv.ParseFloat(historyj[0:len(historyj)-3], 64) 转换失败:", err)
			return false
		}
		if f1 != f2 {
			return f1 <= f2
		}
		return 5.0/fileStores[i].Bandwidth+fileStores[i].Latency <= 5.0/fileStores[j].Bandwidth+fileStores[j].Latency
	})
	table, err := gotable.Create("节点 ID", "节点 IP", "存储能力", "平均历史服务质量", "带宽", "时延", "是否支持高性能读写", "异常服务次数")
	if err != nil {
		fmt.Println("ShowInfoFromRDS FileStoreDAL Create table failed: ", err.Error())
		return
	}
	for _, f := range fileStores {
		if f.GetAbnormalHistoryTimes() == 0 {
			table.AddRow([]string{f.ID[0:5], f.IP, fmt.Sprintf("%.2f GB", f.Mem), f.GetMeanHistory(),
				fmt.Sprintf("%.2f Mbps", f.Bandwidth), fmt.Sprintf("%.2f ms", f.Latency), strconv.FormatBool(f.IsContainFastNetspeed),
				fmt.Sprintf("%d", f.GetAbnormalHistoryTimes())})
		} else {
			table.AddRow([]string{f.ID[0:5] + "*", f.IP, fmt.Sprintf("%.2f GB", f.Mem), f.GetMeanHistory(),
				fmt.Sprintf("%.2f Mbps", f.Bandwidth), fmt.Sprintf("%.2f ms", f.Latency), strconv.FormatBool(f.IsContainFastNetspeed),
				fmt.Sprintf("%d", f.GetAbnormalHistoryTimes())})
		}
	}
	log.Info("内容存储节点性能表现：")
	fmt.Println("\n", table, "\n")
}
