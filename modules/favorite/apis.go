package favorite

import (
	"app/modules/models"
	"app/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"strconv"
)

func Action(c *gin.Context) {
	videoId := c.DefaultQuery("video_id", "0")
	tokenString := c.DefaultQuery("token", "")
	actionType := c.DefaultQuery("action_type", "")

	// validate video_id
	videoIdInt, err := strconv.Atoi(videoId)
	if err != nil || videoIdInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code":    1,
			"status_message": "Invalid video_id",
		})
		return
	}

	userId, _ := utils.ValidateToken(tokenString)
	db := c.MustGet("db").(*gorm.DB)

	// validate action type and perform action accordingly
	switch actionType {
	case "1": // Favorite
		favorite := models.Favorite{
			UserID:  userId,
			VideoID: uint(videoIdInt),
		}
		// Check if this video has been liked by current user
		var count int64
		db.Model(&models.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoIdInt).Count(&count)
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status_code": 1,
				"status_msg":  "Already favortited",
			})
			return
		}
		// Failed to like for some reason
		if err := db.Create(&favorite).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 1,
				"status_msg":  "Failed to favorite",
			})
			return
		}
	case "2": // Un-favorite
		// If current user hasn't liked this video, unlike is not allowed
		var count int64
		db.Model(&models.Favorite{}).Where("user_id = ? AND video_id = ?", userId, videoIdInt).Count(&count)
		if count == 0 {
			c.JSON(http.StatusBadRequest, gin.H{
				"status_code": 1,
				"status_msg":  "Not favorite yet",
			})
			return
		}
		// TODO: Cannot unlike other's like

		// Failed to delete the like record for some reason
		if err := db.Where("user_id = ? AND video_id = ?", userId, videoIdInt).
			Delete(&models.Favorite{}).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"status_code": 1,
				"status_msg":  "Failed to un-favorite",
			})
			return
		}
	default:
		// If actionType is not 1 nor 2
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 1,
			"status_msg":  "Invalid action type.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status_code": 0,
		"status_msg":  "Success",
	})
}

func GetLikeVideos(c *gin.Context) {
	// Get user id string from context
	userIdStr := c.DefaultQuery("user_id", "0")
	// Validate user id
	userId, err := strconv.Atoi(userIdStr)
	if err != nil || userId <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"status_code": 1,
			"status_msg":  "Invalid user_id.",
		})
		return
	}

	db := c.MustGet("db").(*gorm.DB)

	// Get all favorite records where user_id = userId
	var favorites []models.Favorite
	db.Where("user_id = ?", userId).Order("created_at desc").Find(&favorites)

	// Get video id list
	var videoIds []uint
	for _, fav := range favorites {
		videoIds = append(videoIds, fav.VideoID)
	}

	// Get all videos liked by user id
	var videos []models.Video
	db.Preload("User").Preload("User.Profile").
		Where("id IN (?)", videoIds).Find(&videos)

	var videoResList []utils.VideoResItem
	for _, v := range videos {
		videoResList = append(videoResList, utils.VideoResItem{
			ID:            v.ID,
			PlayUrl:       v.PlayUrl,
			CoverUrl:      v.CoverUrl,
			FavoriteCount: v.FavoriteCount,
			CommentCount:  v.CommentCount,
			Title:         v.Title,
			IsFavorite:    true,
			Author: utils.Author{
				ID:             v.User.ID,
				Name:           v.User.Username,
				Avatar:         v.User.Profile.Avatar,
				Background:     v.User.Profile.Background,
				Signature:      v.User.Profile.Signature,
				FollowCount:    v.User.Profile.FollowCount,
				FollowerCount:  v.User.Profile.FollowerCount,
				TotalFavorited: strconv.Itoa(v.User.Profile.TotalFavorited),
				WorkCount:      v.User.Profile.WorkCount,
				FavoriteCount:  v.User.Profile.FavoriteCount,
			},
		})
	}

	c.JSON(http.StatusOK, utils.VideoResponse{
		StatusCode: 0,
		StatusMsg:  "Success",
		VideoList:  videoResList,
	})
}
