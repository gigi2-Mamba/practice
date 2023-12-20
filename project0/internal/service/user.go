package service

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"project0/internal/domain"
	"project0/internal/repository"
)

var (
	ErrDuplicateEmail        = repository.ErrDuplicateEmail
	ErrInvalidUserOrPassword = errors.New("用户不存在或密码错误")
)

// 业务逻辑层调用 数据抽象层，直达数据访问层
type UserService struct {
	repos *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{
		repos: repo,
	}
}

// 登录业务逻辑 ,传了go的context 要干嘛呢
func (svc *UserService) Signup(ctx context.Context, u domain.User) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost) // 指定加密代价。
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return svc.repos.Create(ctx, u)
}

func (svc *UserService) Login(ctx *gin.Context, email string, password string) (domain.User, error) {
	u, err := svc.repos.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	// 检查密码是否正确
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, ErrInvalidUserOrPassword
	}

	return u, nil
}

func (svc *UserService) Profile(ctx *gin.Context, uid int64) error {

	profile, err := svc.repos.Profile(ctx, uid)
	// 以后再改
	ctx.JSON(http.StatusOK, profile)

	return err

}

func (svc *UserService) Edit(ctx *gin.Context, profile domain.UserProfile) error {

	return svc.repos.Edit(ctx, profile)
}
