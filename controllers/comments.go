package controllers

import (
	"capture-life/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type CreateCommentInput struct {
	ParentId   uint   `json:"parentId" binding:"required"`
	ParentType uint   `json:"parentType" binding:"required"`
	Author     string `json:"author" binding:"required"`
	Content    string `json:"content" binding:"required"`
}

type UpdateCommentInput struct {
	Author  string `json:"author"`
	Content string `json:"content"`
}

func CreateComment(c *gin.Context) {
	var input CreateCommentInput

	// Validate input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO Check if blog post exists or if parent comment exists

	// Create comment
	comment := models.Comment{
		ParentId:   input.ParentId,
		ParentType: models.ParentType(input.ParentType),
		Author:     input.Author,
		Content:    input.Content,
	}
	models.DB.Create(&comment)

	c.JSON(http.StatusOK, gin.H{"data": comment})
}

func GetCommentsForBlogPost(c *gin.Context) {
	var comments []models.Comment

	blogPostId, err := strconv.ParseUint(c.Param("blog-post-id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Blog Post Id malformed!"})
	}

	// Lookup Comments by BlogPostId
	// Need to get all comments for blog post
	// by searching comments where parent type == blogpost
	// and parent id matches blog post id
	if err := models.DB.Where(&models.Comment{ParentType: models.PT_BlogPost, ParentId: uint(blogPostId)}).Find(&comments).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Then for all returned comments lookup threaded comments
	// by matching parent id == comment and id of each comment
	// Do this recursively
	for i, comment := range comments {
		comments[i].SubComments = getThreadedComments(comment.ID)
	}

	c.JSON(http.StatusOK, gin.H{"data": comments})
}

func getThreadedComments(parentId uint) []models.Comment {
	var comments []models.Comment
	if err := models.DB.Where(&models.Comment{ParentType: models.PT_Comment, ParentId: parentId}).Find(&comments).Error; err != nil {
		return nil
	}
	if comments == nil {
		return nil
	}
	for i, comment := range comments {
		comments[i].SubComments = getThreadedComments(comment.ID)
	}
	return comments
}

func UpdateComment(c *gin.Context) {
	// Get model if exist
	var comment models.Comment
	if err := models.DB.Where("id = ?", c.Param("id")).First(&comment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Validate input
	var input UpdateCommentInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	models.DB.Model(&comment).Updates(input)

	c.JSON(http.StatusOK, gin.H{"data": comment})
}

func DeleteComment(c *gin.Context) {
	var comment models.Comment
	if err := models.DB.Where("id = ?", c.Param("id")).First(&comment).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Record not found!"})
		return
	}

	// Instead of deleting the comment and all sub-comments
	// just delete the comment body
	models.DB.Model(&comment).Updates(&UpdateCommentInput{
		Content: " ",
	})

	c.JSON(http.StatusOK, gin.H{"data": true})
}
