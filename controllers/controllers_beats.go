package controllers

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/heronhoga/bars-be/config"
	"github.com/heronhoga/bars-be/models/entities"
	"github.com/heronhoga/bars-be/models/requests"
	"github.com/heronhoga/bars-be/utils"
)

func CreateNewBeat(c *fiber.Ctx) error {
	var req requests.CreateBeatRequest

	//parse body
	if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File is required",
		})
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".mp3" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Only .mp3 files are allowed",
		})
	}

	// Check file size
	const MaxFileSize = 5 * 1024 * 1024
	if file.Size > MaxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "File size must be less than or equal to 5 MB",
		})
	}

	//find existing user
	username := c.Locals("username").(string)
	var existingUser entities.User
	err = config.DB.Where("username = ?", username).First(&existingUser).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid username or password",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

	//rename file
	timestamp := time.Now().Format("20060102-150405")
	file.Filename = fmt.Sprintf("%s-%s-%s", existingUser.ID, timestamp, file.Filename)

	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to open uploaded file",
		})
	}
	defer src.Close()

	// Upload to Supabase
	fileURL, err := utils.UploadToSupabase(file.Filename, src, file.Size, file.Filename)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to upload to Supabase: " + err.Error(),
		})
	}

	//save to db
	var newBeat entities.Beat
	newBeat.UserID = existingUser.ID
	newBeat.Title = req.Title
	newBeat.Description = req.Description
	newBeat.Genre = req.Genre
	newBeat.Tags = req.Tags
	newBeat.FileURL = fileURL
	newBeat.FileSize = file.Size

	err = config.DB.Create(&newBeat).Error
	if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Database error",
		})
	}

			return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"error": "New Beat successfully created",
		})
}

func GetAllBeats(c *fiber.Ctx) error {
    limit := c.QueryInt("limit", 10)
    page := c.QueryInt("page", 1)
    title := c.Query("title")

    if page < 1 {
        page = 1
    }
    if limit < 1 {
        limit = 10
    }

    offset := (page - 1) * limit

    var beats []entities.Beat
    query := config.DB.Limit(limit).Offset(offset)

    if title != "" {
        query = query.Where("title ILIKE ?", "%"+title+"%")
    }

    if err := query.Find(&beats).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch beats",
        })
    }

    return c.JSON(beats)
}
