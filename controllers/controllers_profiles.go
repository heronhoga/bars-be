package controllers

import (
	"math"

	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/bars-be/config"
	"github.com/heronhoga/bars-be/models/dto"
	"github.com/heronhoga/bars-be/models/entities"
	"github.com/heronhoga/bars-be/models/requests"
	"github.com/heronhoga/bars-be/utils"
)

func GetProfile(c *fiber.Ctx) error {
	// get user from token jwt
	username := c.Locals("username").(string)
	var profileInfo dto.ProfileInformation

	tx := config.DB.Raw(`
		SELECT
			users.id, 
			users.username, 
			users.region, 
			users.discord,
			COUNT(DISTINCT liked_beats.id) AS likes,
			COUNT(DISTINCT beats.id) AS tracks
		FROM users
		LEFT JOIN beats ON beats.user_id = users.id
		LEFT JOIN liked_beats ON liked_beats.beat_id = beats.id
		WHERE users.username = ?
		GROUP BY users.username, users.region, users.discord, users.id
	`, username).Scan(&profileInfo)

	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data successfully retrieved",
		"data":    profileInfo,
	})
}

func GetBeatByUser(c *fiber.Ctx) error {
	//pagination
    limit := 5
    page := c.QueryInt("page", 1)
	    if page < 1 {
        page = 1
    }

    offset := (page - 1) * limit
	// get user from token jwt
	username := c.Locals("username").(string)
	var beats []dto.BeatByUser

	//get all beats created by user
	tx := config.DB.Raw(`
	SELECT beats.id, beats.title, beats.description, beats.genre, beats.tags, beats.file_url, beats.file_size, beats.created_at, COUNT(DISTINCT liked_beats.id) AS likes
	FROM beats
	JOIN users ON beats.user_id = users.id
	LEFT JOIN liked_beats ON liked_beats.beat_id = beats.id
	WHERE users.username = ?
	GROUP BY beats.id, beats.title, beats.description, beats.genre, beats.tags, beats.file_url, beats.file_size, beats.created_at
	ORDER BY beats.created_at DESC
	LIMIT ?
	OFFSET ?
	`, username, limit, offset).Scan(&beats)

	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	//get total beats created by user
	var total int64
	tx = config.DB.Raw(`
		SELECT COUNT(*) 
		FROM beats
		JOIN users ON beats.user_id = users.id
		WHERE users.username = ?
	`, username).Scan(&total)

		if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))


	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data successfully retrieved",
		"data":    beats,
		"totalPages": totalPages,
	})
}

func GetLikedBeatByUser(c *fiber.Ctx) error {
	//pagination
    limit := 5
    page := c.QueryInt("page", 1)
	    if page < 1 {
        page = 1
    }

    offset := (page - 1) * limit
	// get user from token jwt
	username := c.Locals("username").(string)
	
	var beats []dto.BeatByUser
	tx := config.DB.Raw(
		`SELECT
			beats.id, 
			beats.title, 
			beats.description, 
			beats.genre, 
			beats.tags, 
			beats.file_url, 
			beats.file_size, 
			beats.created_at, 
			COUNT(DISTINCT liked_beats.id) AS likes
		FROM liked_beats
		JOIN beats ON liked_beats.beat_id = beats.id
		JOIN users ON liked_beats.user_id = users.id
		WHERE users.username = ?
		GROUP BY beats.id, beats.title, beats.description, beats.genre, beats.tags, beats.file_url, beats.file_size, beats.created_at
		ORDER BY beats.created_at DESC
		LIMIT ?
		OFFSET ?
		`, username, limit, offset).Scan(&beats)

	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	//get total likes
	var total int64
	tx = config.DB.Raw(`
	SELECT COUNT(liked_beats.id)
	FROM liked_beats
	JOIN users ON liked_beats.user_id = users.id
	WHERE users.username = ?`, username).Scan(&total)

	if tx.Error != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data successfully retrieved",
		"data":    beats,
		"totalPages": totalPages,
	})
}

func EditProfile(c *fiber.Ctx) error {
	var req requests.EditProfileRequest

	//parse body
	if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	//validate request
	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": err,
		})
	}

	//get user
	var existingUser entities.User
	if err := config.DB.Where("id = ?", req.ID).First(&existingUser).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err,
		})
	}

	//map the data
	existingUser.Region = req.Region
	existingUser.Discord = req.Discord

	//update the data
	err := config.DB.Save(&existingUser).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error updating database",
		})
	}

	//save the data
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":"Update beat data successful",
	})

}

