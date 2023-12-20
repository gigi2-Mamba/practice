package web

import (
	"fmt"
	"github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"project0/internal/domain"
	"project0/internal/service"
	"time"
)

// 构建一个实例
const (
	emailRegexPattern = "^\\w+([-+.]\\w+)*@\\w+([-.]\\w+)*\\.\\w+([-.]\\w+)*$"
	// 和上面比起来，用 ` 看起来就比较清爽
	passwordRegexPattern = `^(?=.*[A-Za-z])(?=.*\d)(?=.*[$@$!%*#?&])[A-Za-z\d$@$!%*#?&]{8,}$`
)

type UserHandler struct { // 怎么选择把正则表达式放在 UserHandler里
	svc            *service.UserService // 注入了 对应的service对象
	emailRegexp    *regexp2.Regexp
	passwordRegexp *regexp2.Regexp
}

func NewUserHandler(svc *service.UserService) *UserHandler {
	return &UserHandler{
		emailRegexp:    regexp2.MustCompile(emailRegexPattern, regexp2.None),
		passwordRegexp: regexp2.MustCompile(passwordRegexPattern, regexp2.None),
		svc:            svc,
	}

}
func (u *UserHandler) RegisterRoute(server *gin.Engine) {

	ug := server.Group("/users")

	ug.POST("/signup", u.Signup)
	ug.POST("/login", u.Login)
	ug.GET("/profile", u.Profile)
	ug.POST("/edit", u.Edit)

}

// 注册 注册路由
func (u *UserHandler) Signup(ctx *gin.Context) { //  因为 type HandleFunc  func(*context){}
	// 对请求体进行验证先绑定吗
	// 构建注册请求体
	type SignUpReq struct { // 怎么把证
		Email           string `json:"email"`
		Password        string `json:"password"`
		ConfirmPassword string `json:"confirmPassword"`
	}
	var req SignUpReq
	if err := ctx.Bind(&req); err != nil { // bind要接受json数据
		log.Println(err)
		fmt.Println("req is ", req)
		ctx.String(http.StatusBadRequest, "注册参数请求错误")
		return
	}
	//
	if req.Password != req.ConfirmPassword {
		ctx.String(http.StatusOK, "两次输入密码不一致")
		return
	}
	isEmail, err := u.emailRegexp.MatchString(req.Email)

	if !isEmail {
		log.Println("EMAIL ", err)
		ctx.String(http.StatusBadRequest, "请输入合法邮箱")
	}
	isPwd, err := u.passwordRegexp.MatchString(req.Password)
	log.Println("ISpwd ,", isPwd, req.Password)
	if !isPwd {
		log.Println("pwd ", err)
		ctx.String(http.StatusBadRequest, "请至少包含数字，特殊字符，字母，整体长度8到16")
		return
	}

	err = u.svc.Signup(ctx, domain.User{
		Email:    req.Email,
		Password: req.Password,
	})

	log.Println(req)
	switch err {
	case nil:
		sessions.Default(ctx)
		ctx.String(http.StatusOK, "注册成功")
	case service.ErrDuplicateEmail:
		ctx.String(http.StatusOK, "邮箱已被使用,请更换")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
	//ctx.String(http.StatusOK, "welcome to real project")
}

func (u *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req Req
	// 接受前端数据要son格式  使用bind（）
	if err := ctx.Bind(&req); err != nil {
		return //  不用给信息，gin会给
	}
	user, err := u.svc.Login(ctx, req.Email, req.Password)

	switch err {
	case nil:
		sess := sessions.Default(ctx)
		sess.Set("userId", user.Id)
		sess.Options(sessions.Options{
			MaxAge: 900,
			//HttpOnly:
		})
		err = sess.Save() //  gin session设置要求主动保存
		if err != nil {
			ctx.String(http.StatusOK, "系统错误")
		}
		ctx.String(http.StatusOK, "登录成功")
	case service.ErrInvalidUserOrPassword:
		ctx.String(http.StatusOK, "用户名或密码不对")
	default:
		ctx.String(http.StatusOK, "系统错误")
	}
	//if err != nil {
	//	ctx.String(http.StatusOK, "登录错误")
	//	log.Println("err", err)
	//	return
	//}
	//
	//ctx.String(http.StatusOK, "登录成功")
	//http.Cookie{}
}

func (u *UserHandler) Edit(ctx *gin.Context) {

	type Req struct {
		Nickname     string `json:"nickname"`
		Gender       string `json:"gender"`
		Introduction string `json:"introduction"`
		Birthday     string `json:"birthday"`
	}
	var req Req
	if err := ctx.Bind(&req); err != nil {
		// 因为绑定的错误话，gin会自动处理
		return
	}
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.String(http.StatusOK, "日期格式错误")
		return
	}

	if req.Nickname == "" || len(req.Nickname) > 24 {
		ctx.String(http.StatusOK, "昵称为空，或者超出范围")
		return
	}

	//id,err := strconv.Atoi(ctx.Query("id"))
	// 通过session
	sess := sessions.Default(ctx)
	uid := sess.Get("userId")
	v, ok := uid.(int64)
	if !ok {
		ctx.String(http.StatusOK, "请求参数有误")
		return
	}

	err = u.svc.Edit(ctx, domain.UserProfile{
		Id:           v,
		Gender:       req.Gender,
		NickName:     req.Nickname,
		Introduction: req.Introduction,
		BirthDate:    birthday,
	})

	if err != nil {
		ctx.String(http.StatusOK, "系统错误，个人简介")
	}
	//if err == nil {
	//	log.Println("用户简介查询参数有误",err)
	//	ctx.String(http.StatusOK,"用户参数错误，简介方面")
	//}

	ctx.String(http.StatusOK, "修改成功")
	//u..Profile(ctx, id)
	//ctx.JSON(http.StatusOK)
}

func (u *UserHandler) Profile(ctx *gin.Context) {

	//birthday, err := time.Parse(time.DateOnly, req.Birthday)

	sess := sessions.Default(ctx)
	uid := sess.Get("userId")
	v, ok := uid.(int64)
	if !ok {
		ctx.String(http.StatusOK, "请求参数有误")
		return
	}
	err := u.svc.Profile(ctx, v)
	if err != nil {
		ctx.String(http.StatusOK, "系统错误，个人简介")
	}

	//ctx.String(http.StatusOK, "个人基本资料")
}
