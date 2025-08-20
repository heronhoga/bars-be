package controllers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/bars-be/config"
	"github.com/heronhoga/bars-be/models/dto"
)

func GetProfile(c *fiber.Ctx) error {
	// get user from token jwt
	username := c.Locals("username").(string)
	var profileInfo dto.ProfileInformation

	tx := config.DB.Raw(`
		SELECT 
			users.username, 
			users.region, 
			users.discord,
			COUNT(DISTINCT liked_beats.id) AS likes,
			COUNT(DISTINCT beats.id) AS tracks
		FROM users
		LEFT JOIN beats ON beats.user_id = users.id
		LEFT JOIN liked_beats ON liked_beats.beat_id = beats.id
		WHERE users.username = ?
		GROUP BY users.username, users.region, users.discord
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

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Data successfully retrieved",
		"data":    beats,
	})
}
