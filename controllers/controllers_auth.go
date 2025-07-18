package controllers

import (
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/heronhoga/bars-be/config"
	"github.com/heronhoga/bars-be/models/entities"
	"github.com/heronhoga/bars-be/models/requests"
	"github.com/heronhoga/bars-be/utils"
)

var validate = validator.New()

func Register(c *fiber.Ctx) error {
	var req requests.RegisterRequest

		// Parse body
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Validate request
	if err := validate.Struct(req); err != nil {
		errors := make(map[string]string)
		for _, err := range err.(validator.ValidationErrors) {
			switch err.Tag() {
			case "required":
				errors[err.Field()] = "This field is required"
			case "min":
				errors[err.Field()] = "Must be at least " + err.Param() + " characters"
			default:
				errors[err.Field()] = "Invalid value"
			}
		}

		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": errors,
		})
	}

	//hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to hash password",
		})
	}

	//save new user
	var newUser entities.User

	newUser.Username = req.Username
	newUser.Password = hashedPassword
	newUser.Region = req.Region
	newUser.Discord = req.Discord

	err = config.DB.Create(&newUser).Error

	if err != nil {
		//unique constraint violation
	if strings.Contains(err.Error(), "duplicate key value") || strings.Contains(err.Error(), "UNIQUE constraint failed") {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Username already exists",
		})
		}
		//else
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create user",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "New user successfully created",
	})
}