package controllers

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/heronhoga/bars-be/config"
	"github.com/heronhoga/bars-be/models/entities"
	"github.com/heronhoga/bars-be/models/requests"
	"github.com/heronhoga/bars-be/utils"
	"gorm.io/gorm"
)

func Like(c *fiber.Ctx) error {
	var req requests.LikeRequest

	err := c.BodyParser(&req)
	if err != nil {
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

	//get user from token jwt
	username := c.Locals("username").(string)
	var user entities.User
	err = config.DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	//parse string to uuid
	beatUUID, err := uuid.Parse(req.BeatID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid beat ID format",
		})
	}

	var existingLike entities.LikedBeat

	err = config.DB.Where("beat_id = ?", req.BeatID).Where("user_id = ?", user.ID).First(&existingLike).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newLike := entities.LikedBeat{
				BeatID: beatUUID,
				UserID: user.ID,
			}
			if err := config.DB.Create(&newLike).Error; err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Internal server error (adding like)",
				})
			}
			return c.Status(fiber.StatusOK).JSON(fiber.Map{
				"message": "Like added",
			})
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error (finding like)",
			})
		}
	} else {
		if err := config.DB.Delete(&existingLike).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Internal server error (removing like)",
			})
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Like removed",
		})
	}
}