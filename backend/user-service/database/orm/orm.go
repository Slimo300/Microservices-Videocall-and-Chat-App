package orm

import (
	"context"
	"fmt"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type UsersGormRepository struct {
	db *gorm.DB
}

func NewUsersGormRepository(address string) (*UsersGormRepository, error) {
	conn, err := gorm.Open(mysql.Open(fmt.Sprintf("%s?parseTime=true", address)))
	if err != nil {
		return nil, err
	}
	if err := conn.AutoMigrate(&User{}, &AuthorizationCode{}); err != nil {
		return nil, err
	}
	return &UsersGormRepository{db: conn}, nil
}

func (r *UsersGormRepository) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	var u User
	if err := r.db.WithContext(ctx).First(&u, userID).Error; err != nil {
		return nil, err
	}
	user := models.UnmarshalUserFromDatabase(u.ID, u.UserName, u.Email, u.Password, *u.HasPicture, *u.Verified)
	return user, nil
}

func (r *UsersGormRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	var u User
	if err := r.db.WithContext(ctx).Where(&User{Email: email}).First(&u).Error; err != nil {
		return nil, err
	}
	user := models.UnmarshalUserFromDatabase(u.ID, u.UserName, u.Email, u.Password, *u.HasPicture, *u.Verified)
	return user, nil
}

func (r *UsersGormRepository) RegisterUser(ctx context.Context, user *models.User, code *models.AuthorizationCode) error {
	u := marshalUser(*user)
	c := marshalCode(*code)
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.db.WithContext(ctx).Create(&u).Error; err != nil {
			return err
		}
		if err := r.db.WithContext(ctx).Create(&c).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *UsersGormRepository) CreateAuthorizationCode(ctx context.Context, code *models.AuthorizationCode) error {
	c := marshalCode(*code)
	return r.db.WithContext(ctx).Create(&c).Error
}

func (r *UsersGormRepository) UpdateUserByID(ctx context.Context, userID uuid.UUID, updateFn func(u *models.User) (*models.User, error)) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var u User
		if err := r.db.WithContext(ctx).First(&u, userID).Error; err != nil {
			return err
		}
		user := models.UnmarshalUserFromDatabase(u.ID, u.UserName, u.Email, u.Password, *u.HasPicture, *u.Verified)
		user, err := updateFn(user)
		if err != nil {
			return err
		}
		// if user is nil with no error we don't have to perform any update
		if user == nil {
			return nil
		}
		u = marshalUser(*user)
		if err := r.db.WithContext(ctx).Model(&u).Updates(u).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r *UsersGormRepository) UpdateUserByCode(ctx context.Context, code uuid.UUID, codeType models.CodeType, updateFn func(u *models.User) (*models.User, error)) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		var c AuthorizationCode
		if err := r.db.WithContext(ctx).Where(&AuthorizationCode{CodeType: codeType.String()}).First(&c, code).Error; err != nil {
			return err
		}
		var u User
		if err := r.db.WithContext(ctx).First(&u, c.UserID).Error; err != nil {
			return err
		}
		user := models.UnmarshalUserFromDatabase(u.ID, u.UserName, u.Email, u.Password, *u.HasPicture, *u.Verified)
		user, err := updateFn(user)
		if err != nil {
			return err
		}
		u = marshalUser(*user)
		if err := r.db.WithContext(ctx).Model(&u).Updates(u).Error; err != nil {
			return err
		}
		if err := r.db.WithContext(ctx).Delete(&c).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (r UsersGormRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if err := r.db.Transaction(func(tx *gorm.DB) error {
		if err := r.db.WithContext(ctx).Where("user_id = ?", userID).Delete(&AuthorizationCode{}).Error; err != nil {
			return err
		}
		if err := r.db.WithContext(ctx).Delete(&User{ID: userID}).Error; err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}
