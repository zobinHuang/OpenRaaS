package dal

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/zobinHuang/BrosCloud/backstage/auth/model"
	"github.com/zobinHuang/BrosCloud/backstage/auth/model/apperrors"
	"gorm.io/gorm"
)

/*
	struct: userDAL
	description: dal layer
*/
type userDAL struct {
	DB *gorm.DB
}

/*
	func: NewRDbDAL
	description: return an instance of struct rdbDAL
*/
func NewUserDAL(db *gorm.DB) model.UserDAL {
	return &userDAL{
		DB: db,
	}
}

/*
	func: FindUserProfileByEmail
	description: return model.User instance according to email
*/
func (udal *userDAL) FindUserProfileByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}

	if err := udal.DB.WithContext(ctx).Where("email = ?", email).First(user).Error; err != nil {
		log.WithFields(log.Fields{
			"Email": email,
			"error": err,
		}).Warn("Unable to find user in database with specified email")
		return nil, apperrors.NewNotFound("mobile", email)
	}

	return user, nil
}

/*
	func: CreateNewUser
	description: create a new user in database
*/
func (udal *userDAL) CreateNewUser(ctx context.Context, u *model.User) error {
	tx := udal.DB.WithContext(ctx)

	// check email conflict
	if err := tx.Where("email = ?", u.Email).First(&model.User{}).Error; err == nil {
		log.WithFields(log.Fields{
			"Email": u.Email,
			"error": err,
		}).Warn("Unable to create user in database, email already exists")
		return apperrors.NewConflict("email", u.Email)
	}

	// create a transaction
	committx := tx.Begin()

	// insert user profit
	if err := committx.Create(u).Error; err != nil {
		committx.Rollback()
		return apperrors.NewInternal()
	}

	// commit transaction
	committx.Commit()

	return nil
}

// Not used
func (udal *userDAL) GetUser(ctx context.Context, u *model.User) error {
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
