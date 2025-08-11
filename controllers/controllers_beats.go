package controllers

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"github.com/heronhoga/bars-be/config"
	"github.com/heronhoga/bars-be/models/dto"
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

	//validate request
	if err := utils.ValidateStruct(req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"errors": err,
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
			"message": "New Beat successfully created",
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

    // Get username from middleware
    username := c.Locals("username").(string)
    var existingUser entities.User
    err := config.DB.Where("username = ?", username).First(&existingUser).Error
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

    // Query beats
    var beats []dto.FullBeatAndUser
    rawQuery := `
        SELECT 
            beats.id,
            users.username,
            beats.title,
            beats.description,
            beats.genre,
            beats.tags,
            beats.file_url,
            beats.file_size,
            beats.created_at,
            COUNT(DISTINCT liked_beats_all.id) AS likes,
            CASE 
                WHEN COUNT(DISTINCT liked_beats_user.id) > 0 THEN '1'
                ELSE '0'
            END AS is_liked
        FROM beats
        JOIN users ON beats.user_id = users.id
        LEFT JOIN liked_beats AS liked_beats_all 
            ON liked_beats_all.beat_id = beats.id
        LEFT JOIN liked_beats AS liked_beats_user 
            ON liked_beats_user.beat_id = beats.id 
            AND liked_beats_user.user_id = ?
        WHERE (? = '' OR beats.title ILIKE '%' || ? || '%')
        GROUP BY beats.id, users.username, beats.title, beats.description, beats.genre, beats.tags, beats.file_url, beats.file_size, beats.created_at
        ORDER BY beats.created_at DESC
        LIMIT ? OFFSET ?
    `

    if err := config.DB.Raw(rawQuery, existingUser.ID, title, title, limit, offset).Scan(&beats).Error; err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": "Failed to fetch beats",
        })
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Beats successfully retrieved",
        "data":    beats,
    })
}


func DeleteBeat (c *fiber.Ctx) error {
	fmt.Println("deleting beat..")
	var req requests.DeleteBeatRequest

	//parse body
	if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	//find existingBeat
	var existingBeat dto.BeatAndUser
	query := `
		SELECT beats.id, beats.file_url, users.username
		FROM beats
		JOIN users ON beats.user_id = users.id
		WHERE beats.id = ?
		LIMIT 1;
	`

	err := config.DB.Raw(query, req.BeatID).Scan(&existingBeat).Error

	fmt.Println(existingBeat)

	if err != nil {
					return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	usernameFromToken, ok := c.Locals("username").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid token",
		})
	}

	fmt.Println("username in context:", usernameFromToken)
	fmt.Println("username from db:", existingBeat.Username)

	if existingBeat.Username != usernameFromToken {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}

	//delete beat from supabase
	err = utils.DeleteSupabaseFile(existingBeat.FileURL)
	if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error (deleting file from supabase)",
		})
	}

	//delete beat from db
	err = config.DB.Where("id = ?", existingBeat.ID).Delete(&entities.Beat{}).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to delete beat",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Beat deleted successfully",
	})

}

func EditBeat (c *fiber.Ctx) error {
	var req requests.EditBeatRequest

	beatId := c.Params("beatid")

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

	var existingBeat entities.Beat

	if err := config.DB.First(&existingBeat, "id = ?", beatId).Error; err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Beat not found"})
	}

	//update
	existingBeat.Title = req.Title
	existingBeat.Genre = req.Genre
	existingBeat.Description = req.Description
	existingBeat.Tags = req.Tags

	err := config.DB.Save(&existingBeat).Error
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Internal server error updating database"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":"Update beat data successful",
	})
}

func GetFavoriteBeats (c *fiber.Ctx) error {
	var beats []dto.BeatAndLike

	err := config.DB.Raw(`
		SELECT 
			beats.id,
			beats.title, 
			beats.file_url, 
			users.username, 
			COUNT(liked_beats.id) AS likes
		FROM beats 
		JOIN users ON beats.user_id = users.id
		LEFT JOIN liked_beats ON liked_beats.beat_id = beats.id
		GROUP BY beats.id, beats.title, beats.file_url, users.username
		ORDER BY likes DESC
		LIMIT 5
	`).Scan(&beats).Error

	if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Beats successfully retrieved",
		"data": beats,
	})
}
