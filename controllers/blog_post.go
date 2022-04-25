package controllers

import (
	"capture-life/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"regexp"
)

type CreateBlogPostInput struct {
	Title   string `json:"title" binding:"required"`
	Author  string `json:"author" binding:"required"`
	Content string `json:"content" binding:"required"`
	UrlName string `json:"urlName" binding:"required"`
}

type UpdateBlogPostInput struct {
	Title   string `json:"title"`
	Author  string `json:"author"`
	Content string `json:"content"`
	UrlName string `json:"urlName"`
}

func CreateBlogPost(c *gin.Context) {
	var input CreateBlogPostInput

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// RegExp to replace certain chars from UrlName
	rexp := regexp.MustCompile(`[^\w\-\!\$\'\(\)\=\@\d_]+`)

	// Create blog post
	blogPost := models.BlogPost{
		Title:   input.Title,
		Author:  input.Author,
		Content: input.Content,
		UrlName: rexp.ReplaceAllString(input.UrlName, "-"),
	}
	models.DB.Create(&blogPost)

	c.JSON(http.StatusOK, gin.H{"data": blogPost})
}

func GetAllBlogPosts(c *gin.Context) {
	var blogPosts []models.BlogPost
	models.DB.Find(&blogPosts)
	c.JSON(http.StatusOK, gin.H{"data": blogPosts})
}

func GetBlogPost(c *gin.Context) {
	var blogPost models.BlogPost

	if err := models.DB.Where("id = ?", c.Param("id")).First(&blogPost).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": blogPost})
}

func UpdateBlogPost(c *gin.Context) {
	// Get model if exist
	var blogPost models.BlogPost
	if err := models.DB.Where("id = ?", c.Param("id")).First(&blogPost).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input UpdateBlogPostInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&blogPost).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": blogPost})
}

func DeleteBlogPost(c *gin.Context) {
	var blogPost models.BlogPost
	if err := models.DB.Where("id = ?", c.Param("id")).First(&blogPost).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	models.DB.Delete(&blogPost)

	c.JSON(http.StatusOK, gin.H{"data": true})
}
