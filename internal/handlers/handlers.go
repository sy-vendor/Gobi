package handlers

import (
	"encoding/json"
	"gobi/config"
	"gobi/internal/models"
	"gobi/pkg/database"
	"gobi/pkg/errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Auth handlers
func Login(c *gin.Context) {
	var login struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&login); err != nil {
		if errors.IsValidationError(err) {
			c.Error(errors.NewBadRequestError("Invalid login request", err))
		} else {
			c.Error(errors.WrapError(err, "Invalid login request"))
		}
		return
	}

	var user models.User
	if err := database.DB.Where("username = ?", login.Username).First(&user).Error; err != nil {
		c.Error(errors.ErrInvalidCredentials)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password)); err != nil {
		c.Error(errors.ErrInvalidCredentials)
		return
	}

	cfg := config.DefaultConfig
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte(cfg.JWT.Secret))
	if err != nil {
		c.Error(errors.WrapError(err, "Could not generate token"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Register(c *gin.Context) {
	var register struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&register); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.IsValidationError(err) {
			c.Error(errors.NewBadRequestError("Invalid registration request", err))
		} else if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			c.Error(errors.NewBadRequestError("Invalid JSON format", err))
		} else {
			c.Error(errors.WrapError(err, "Invalid registration request"))
		}
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	err := database.DB.Where("username = ?", register.Username).First(&existingUser).Error
	if err == nil {
		c.Error(errors.NewConflictError("User already exists", nil))
		return
	} else if err != gorm.ErrRecordNotFound {
		c.Error(errors.WrapError(err, "Database error"))
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(register.Password), bcrypt.DefaultCost)
	if err != nil {
		c.Error(errors.WrapError(err, "Could not hash password"))
		return
	}

	user := models.User{
		Username: register.Username,
		Password: string(hashedPassword),
		Role:     "user",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not create user"))
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "User created successfully"})
}

// Query handlers
func CreateQuery(c *gin.Context) {
	var query models.Query
	if err := c.ShouldBindJSON(&query); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.IsValidationError(err) {
			c.Error(errors.NewBadRequestError("Invalid query data", err))
		} else if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			c.Error(errors.NewBadRequestError("Invalid JSON format", err))
		} else {
			c.Error(errors.WrapError(err, "Invalid query data"))
		}
		return
	}

	userID, _ := c.Get("userID")
	query.UserID = userID.(uint)

	if err := database.DB.Create(&query).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not create query"))
		return
	}

	c.JSON(http.StatusCreated, query)
}

func ListQueries(c *gin.Context) {
	var queries []models.Query
	userID, _ := c.Get("userID")
	role, _ := c.Get("role")

	query := database.DB.Model(&models.Query{})
	if role.(string) != "admin" {
		query = query.Where("user_id = ? OR is_public = ?", userID, true)
	}

	if err := query.Find(&queries).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch queries"})
		return
	}

	c.JSON(http.StatusOK, queries)
}

func GetQuery(c *gin.Context) {
	id := c.Param("id")
	var query models.Query
	if err := database.DB.First(&query, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	role, _ := c.Get("role")
	if role.(string) != "admin" && query.UserID != userID.(uint) && !query.IsPublic {
		c.Error(errors.ErrForbidden)
		return
	}

	c.JSON(http.StatusOK, query)
}

func UpdateQuery(c *gin.Context) {
	id := c.Param("id")
	var query models.Query
	if err := database.DB.First(&query, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	if query.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	if err := c.ShouldBindJSON(&query); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.IsValidationError(err) {
			c.Error(errors.NewBadRequestError("Invalid query data", err))
		} else if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			c.Error(errors.NewBadRequestError("Invalid JSON format", err))
		} else {
			c.Error(errors.WrapError(err, "Invalid query data"))
		}
		return
	}

	if err := database.DB.Save(&query).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not update query"))
		return
	}

	c.JSON(http.StatusOK, query)
}

func DeleteQuery(c *gin.Context) {
	id := c.Param("id")
	var query models.Query
	if err := database.DB.First(&query, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	if query.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	if err := database.DB.Delete(&query).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not delete query"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Query deleted successfully"})
}

// Chart handlers
func CreateChart(c *gin.Context) {
	var chart models.Chart
	if err := c.ShouldBindJSON(&chart); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.IsValidationError(err) {
			c.Error(errors.NewBadRequestError("Invalid chart data", err))
		} else if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			c.Error(errors.NewBadRequestError("Invalid JSON format", err))
		} else {
			c.Error(errors.WrapError(err, "Invalid chart data"))
		}
		return
	}

	userID, _ := c.Get("userID")
	chart.UserID = userID.(uint)

	if err := database.DB.Create(&chart).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not create chart"))
		return
	}

	c.JSON(http.StatusCreated, chart)
}

func ListCharts(c *gin.Context) {
	var charts []models.Chart
	userID, _ := c.Get("userID")

	if err := database.DB.Where("user_id = ?", userID).Find(&charts).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch charts"})
		return
	}

	c.JSON(http.StatusOK, charts)
}

func GetChart(c *gin.Context) {
	id := c.Param("id")
	var chart models.Chart
	if err := database.DB.First(&chart, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	if chart.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	c.JSON(http.StatusOK, chart)
}

func UpdateChart(c *gin.Context) {
	id := c.Param("id")
	var chart models.Chart
	if err := database.DB.First(&chart, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	if chart.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	if err := c.ShouldBindJSON(&chart); err != nil {
		var syntaxErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError
		if errors.IsValidationError(err) {
			c.Error(errors.NewBadRequestError("Invalid chart data", err))
		} else if errors.As(err, &syntaxErr) || errors.As(err, &typeErr) {
			c.Error(errors.NewBadRequestError("Invalid JSON format", err))
		} else {
			c.Error(errors.WrapError(err, "Invalid chart data"))
		}
		return
	}

	if err := database.DB.Save(&chart).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not update chart"))
		return
	}

	c.JSON(http.StatusOK, chart)
}

func DeleteChart(c *gin.Context) {
	id := c.Param("id")
	var chart models.Chart
	if err := database.DB.First(&chart, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	if chart.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	if err := database.DB.Delete(&chart).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not delete chart"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Chart deleted successfully"})
}

// Excel template handlers
func CreateTemplate(c *gin.Context) {
	file, err := c.FormFile("template")
	if err != nil {
		if errors.IsContentTypeError(err) {
			c.Error(errors.NewBadRequestError("No file uploaded", err))
		} else {
			c.Error(errors.WrapError(err, "No file uploaded"))
		}
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.Error(errors.WrapError(err, "Could not open file"))
		return
	}
	defer openedFile.Close()

	template := models.ExcelTemplate{
		Name:     file.Filename,
		UserID:   c.GetUint("userID"),
		Template: make([]byte, file.Size),
	}

	if _, err := openedFile.Read(template.Template); err != nil {
		c.Error(errors.WrapError(err, "Could not read file"))
		return
	}

	if err := database.DB.Create(&template).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not save template"))
		return
	}

	c.JSON(http.StatusCreated, template)
}

func ListTemplates(c *gin.Context) {
	var templates []models.ExcelTemplate
	userID, _ := c.Get("userID")

	if err := database.DB.Where("user_id = ?", userID).Find(&templates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch templates"})
		return
	}

	c.JSON(http.StatusOK, templates)
}

func GetTemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.ExcelTemplate
	if err := database.DB.First(&template, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	if template.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	c.JSON(http.StatusOK, template)
}

func UpdateTemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.ExcelTemplate
	if err := database.DB.First(&template, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	if template.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	file, err := c.FormFile("template")
	if err != nil {
		if errors.IsContentTypeError(err) {
			c.Error(errors.NewBadRequestError("No file uploaded", err))
		} else {
			c.Error(errors.WrapError(err, "No file uploaded"))
		}
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		c.Error(errors.WrapError(err, "Could not open file"))
		return
	}
	defer openedFile.Close()

	template.Name = file.Filename
	template.Template = make([]byte, file.Size)
	if _, err := openedFile.Read(template.Template); err != nil {
		c.Error(errors.WrapError(err, "Could not read file"))
		return
	}

	if err := database.DB.Save(&template).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not update template"))
		return
	}

	c.JSON(http.StatusOK, template)
}

func DeleteTemplate(c *gin.Context) {
	id := c.Param("id")
	var template models.ExcelTemplate
	if err := database.DB.First(&template, id).Error; err != nil {
		c.Error(errors.ErrNotFound)
		return
	}

	userID, _ := c.Get("userID")
	if template.UserID != userID.(uint) {
		c.Error(errors.ErrForbidden)
		return
	}

	if err := database.DB.Delete(&template).Error; err != nil {
		c.Error(errors.WrapError(err, "Could not delete template"))
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template deleted successfully"})
}
