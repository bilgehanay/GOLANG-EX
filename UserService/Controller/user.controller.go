package controller

import (
	"deneme.com/bng-go/Middleware"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"

	"deneme.com/bng-go/Model"
	service "deneme.com/bng-go/Service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UserController struct {
	UserService service.UserService
}

func New(userservice service.UserService) UserController {
	return UserController{UserService: userservice}
}

func (uc *UserController) CreateUser(ctx *gin.Context) {
	var user Model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := uc.UserService.CreateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) GetUser(ctx *gin.Context) {
	userid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}
	user, err := uc.UserService.GetUser(&userid)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success", "user": user})
}

func (uc *UserController) GetUsers(ctx *gin.Context) {
	users, err := uc.UserService.GetUsers()
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success", "users": users})
}

func (uc *UserController) UpdateUser(ctx *gin.Context) {
	var user Model.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	err := uc.UserService.UpdateUser(&user)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) DeleteUser(ctx *gin.Context) {
	userid, err := uuid.Parse(ctx.Param("id"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid user ID"})
		return
	}
	err = uc.UserService.DeleteUser(&userid)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success"})
}

func (uc *UserController) Login(ctx *gin.Context) {
	var login_req struct {
		Email    string `json:"email" bson:"email,omitempty"`
		Password string `json:"password" bson:"password,omitempty"`
	}
	if err := ctx.ShouldBindJSON(&login_req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	id, err := uc.UserService.LoginUser(login_req.Email, login_req.Password)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	userClaims := service.UserClaims{
		Id:       id,
		Email:    login_req.Email,
		Password: login_req.Password,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * time.Duration(1))),
		},
	}

	token, err := service.NewAccessToken(userClaims)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "success", "token": token})

}

func (uc *UserController) VerifyToken(ctx *gin.Context) {
	token := ctx.Param("token")
	if len(token) == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid token"})
		return
	}

	UserClaims, err := service.ParseAccessToken(token)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "success", "user": UserClaims})

}

func (uc *UserController) RegisterUserRoutes(rg *gin.RouterGroup) {
	userroute := rg.Group("/user")
	userroute.Use(Middleware.RateLimit())
	userroute.POST("", uc.CreateUser)
	userroute.GET("/:id", uc.GetUser)
	userroute.GET("", uc.GetUsers)
	userroute.PUT("", uc.UpdateUser)
	userroute.DELETE("/:id", uc.DeleteUser)

	userroute.POST("/login", uc.Login)
	userroute.GET("/verify/:token", uc.VerifyToken)
}
