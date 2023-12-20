package repository

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"project0/internal/domain"
	"project0/internal/repository/dao"
	"time"
)

// 预定义错误
var (
	ErrDuplicateEmail = dao.ErrDuplicateEmail
	ErrUserNotFound   = dao.ErrRecordNotFound
)

type UserRepository struct {
	dao *dao.UserDao
}

func (repos *UserRepository) Create(ctx context.Context, u domain.User) error {
	return repos.dao.Insert(ctx, dao.User{
		Email:    u.Email,
		Password: u.Password,
	})

}

func (repos *UserRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repos.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repos.toDomain(u), nil
}

func NewUserRepository(dao *dao.UserDao) *UserRepository {

	return &UserRepository{
		dao: dao,
	}
}

func (repos *UserRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		Id:       u.Id,
		Email:    u.Email,
		Password: u.Password,
	}
}

func (repos *UserRepository) Profile(ctx context.Context, id int64) (domain.UserProfile, error) {

	profile, err := repos.dao.Profile(ctx, id)
	if err != nil {
		log.Println(err)
		return domain.UserProfile{}, err
	}

	return repos.toDomainProfile(profile), err
}

func (repos *UserRepository) toDomainProfile(u dao.UserProfile) domain.UserProfile {
	birthUnix := u.BirthDate * int64((time.Millisecond))
	t := time.Unix(0, birthUnix)

	return domain.UserProfile{
		Id:           u.Id,
		Gender:       u.Gender,
		NickName:     u.NickName,
		Introduction: u.Introduction,
		BirthDate:    t,
	}
}

func (repos *UserRepository) Edit(ctx *gin.Context, profile domain.UserProfile) error {

	err := repos.dao.Edit(ctx, profile)

	return err
}
