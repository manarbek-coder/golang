package v1

import (
	"net/http"
	"practice-7/internal/entity"
	"practice-7/internal/usecase"
	"practice-7/utils"

	"github.com/gin-gonic/gin"
)

type userRoutes struct {
	u usecase.UserInterface
}

func NewUserRoutes(r *gin.RouterGroup, u usecase.UserInterface) {
	handler := &userRoutes{u}

	users := r.Group("/users")
	{
		users.POST("/", handler.RegisterUser)
		users.POST("/login", handler.LoginUser)

		protected := users.Group("/")
		protected.Use(utils.JWTAuthMiddleware())
		{
			protected.GET("/me", handler.GetMe)

			admin := protected.Group("/")
			admin.Use(utils.RoleMiddleware("admin"))
			{
				admin.PATCH("/promote/:id", handler.PromoteUser)
			}
		}
	}
}

func (h *userRoutes) RegisterUser(c *gin.Context) {
	var dto entity.CreateUserDTO
	c.BindJSON(&dto)

	hash, _ := utils.HashPassword(dto.Password)

	user := entity.User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: hash,
		Role:     "user",
	}

	res, session, _ := h.u.RegisterUser(&user)

	c.JSON(http.StatusCreated, gin.H{
		"user": res,
		"sess": session,
	})
}

func (h *userRoutes) LoginUser(c *gin.Context) {
	var dto entity.LoginUserDTO
	c.BindJSON(&dto)

	token, _ := h.u.LoginUser(&dto)

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (h *userRoutes) GetMe(c *gin.Context) {
	userID := c.GetString("userID")
	user, _ := h.u.GetMe(userID)

	c.JSON(200, gin.H{
		"email": user.Email,
	})
}

func (h *userRoutes) PromoteUser(c *gin.Context) {
	id := c.Param("id")

	h.u.PromoteUser(id)

	c.JSON(200, gin.H{"message": "promoted"})
}
