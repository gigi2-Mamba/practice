package dao

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"log"
	"project0/internal/domain"
	"time"
)

// 预定义错误
var (
	ErrDuplicateEmail = errors.New("邮箱冲突")
	ErrRecordNotFound = gorm.ErrRecordNotFound //  利用中间件提供的错误码，对数据库报错信息进行反馈
)

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	//  时区
	//
	Ctime int64
	Utime int64
}

type UserProfile struct {
	Id       int64  `gorm:"primary"`
	NickName string `gorm:"unique"`
	Gender   string

	Introduction string
	BirthDate    int64
}

type UserDao struct {
	db *gorm.DB
}

func (dao *UserDao) Insert(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.Ctime = now
	u.Utime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	if me, ok := err.(*mysql.MySQLError); ok { // 驱动包的error
		const duplicateErr uint16 = 1062
		if me.Number == duplicateErr {
			//用户冲突吧
			return ErrDuplicateEmail
		}
	}
	//up := NewUserProfileDao(dao.db)
	var upInstance UserProfile
	log.Println("--------", upInstance)
	err = dao.CreateProfile(ctx, upInstance, u.Id)

	return err
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email=?", email).First(&u).Error
	return u, err
}

func NewUserDAO(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}

func (u *UserDao) CreateProfile(ctx context.Context, up UserProfile, id int64) error {
	up.Id = id
	return u.db.WithContext(ctx).Create(&up).Error

}

func (u *UserDao) Profile(ctx context.Context, id int64) (UserProfile, error) {

	var uprofile UserProfile
	err := u.db.WithContext(ctx).Where("id=?", id).First(&uprofile).Error
	return uprofile, err

}

func (dao *UserDao) Edit(ctx *gin.Context, profile domain.UserProfile) error {
	uprofile := &UserProfile{
		Id:           profile.Id,
		Gender:       profile.Gender,
		BirthDate:    profile.BirthDate.Unix(),
		Introduction: profile.Introduction,
		NickName:     profile.NickName,
	}

	return dao.db.WithContext(ctx).Updates(&uprofile).Error

}

//func (u *UserDao) Edit(ctx context.Context, up)

//func (u *UserDao) Edit(ctx context.Context, up *UserProfile) (domain.UserProfile, error) {
//
//}
