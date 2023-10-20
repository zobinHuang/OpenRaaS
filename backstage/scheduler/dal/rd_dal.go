package dal

import (
	"context"

	"github.com/zobinHuang/OpenRaaS/backstage/scheduler/model"
	"gorm.io/gorm"
)

/*
	struct: rdbDAL
	description: dal layer
*/
type rdbDAL struct {
	DB *gorm.DB
}

/*
	func: NewRDbDAL
	description: return an instance of struct rdbDAL
*/
func NewRDbDAL(db *gorm.DB) model.RDbDAL {
	return &rdbDAL{
		DB: db,
	}
}

func (r *rdbDAL) GetRDbModel(ctx context.Context, rdbm *model.RDbModel) error {
	// initialize context
	// tx := r.DB.WithContext(ctx)

	// retrieve
	/*
		if err := tx.Where("user_name = ?", rdbm.UserName).First(&model.UserProfile{}).Error; err != nil {
			log.Printf("Unable to find user with username %v in database\n", up.UserName)
			return apperrors.NewNotFound("user_name", rdbm.UserName)
		}
	*/

	// create
	/*
		committx := tx.Begin()
		if err := committx.Create(rdbm).Error; err != nil {
			committx.Rollback()
			return apperrors.NewInternal()
		}
		committx.Commit()
	*/

	// update
	/*
		rdbModel := &model.RDbModel{}
		committx := tx.Begin()
		if err := committx.Model(&rdbModel).Where(
				"user_name = ?", rdbm.UserName).Updates(rdbm).Error; err != nil {
			committx.Rollback()
			return apperrors.NewInternal()
		}
	*/

	return nil
}
