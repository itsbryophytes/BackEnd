package controller

import (
	"net/http"
	"time"

	"github.com/Caknoooo/go-gin-clean-starter/database/entities"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ArticleController interface {
	GetArticles(c *gin.Context)
	GetArticle(c *gin.Context)
	CreateArticle(c *gin.Context)
	UpdateArticle(c *gin.Context)
	DeleteArticle(c *gin.Context)
	PublishArticle(c *gin.Context)
}

type articleController struct {
	db *gorm.DB
}

func NewArticleController(db *gorm.DB) ArticleController {
	return &articleController{db: db}
}

func (ctrl *articleController) GetArticles(c *gin.Context) {
	var articles []entities.Article
	if err := ctrl.db.Order("created_at DESC").Find(&articles).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, articles)
}

func (ctrl *articleController) GetArticle(c *gin.Context) {
	var article entities.Article
	if err := ctrl.db.First(&article, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (ctrl *articleController) CreateArticle(c *gin.Context) {
	var article entities.Article
	if err := c.ShouldBindJSON(&article); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	if article.ID == uuid.Nil {
		article.ID = uuid.New()
	}

	if article.AuthorID == "" {
		if userID, exists := c.Get("user_id"); exists {
			article.AuthorID = userID.(string)
		}
	}

	if article.Status == "published" && article.PublishedAt == nil {
		now := time.Now()
		article.PublishedAt = &now
	}

	if err := ctrl.db.Create(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, article)
}

func (ctrl *articleController) UpdateArticle(c *gin.Context) {
	var article entities.Article
	if err := ctrl.db.First(&article, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	var req entities.Article
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}

	article.Title = req.Title
	article.Content = req.Content
	article.CoverImageURL = req.CoverImageURL
	
	if req.AuthorID != "" {
		article.AuthorID = req.AuthorID
	}
	
	if req.PublishedAt != nil {
		article.PublishedAt = req.PublishedAt
	} else if req.Status == "published" && article.Status != "published" {
		now := time.Now()
		article.PublishedAt = &now
	} else if req.Status == "draft" {
		article.PublishedAt = nil
	}
	article.Status = req.Status

	if err := ctrl.db.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, article)
}

func (ctrl *articleController) DeleteArticle(c *gin.Context) {
	if err := ctrl.db.Delete(&entities.Article{}, "id = ?", c.Param("id")).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

func (ctrl *articleController) PublishArticle(c *gin.Context) {
	now := time.Now()
	if err := ctrl.db.Model(&entities.Article{}).
		Where("id = ?", c.Param("id")).
		Updates(map[string]any{"status": "published", "published_at": &now}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "published"})
}
