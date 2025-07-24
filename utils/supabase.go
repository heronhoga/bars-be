package utils

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func UploadToSupabase(filePath string, file io.Reader, fileSize int64, fileName string) (string, error) {
	projectURL := os.Getenv("SUPABASE_URL")
	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")
	fmt.Println(serviceKey)
	bucketName := "hg-bucket"

	// Prepare the request
	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", projectURL, bucketName, fileName)

	req, err := http.NewRequest("POST", url, file)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+serviceKey)
	req.Header.Set("Content-Type", "audio/mpeg")
	req.Header.Set("Content-Length", fmt.Sprintf("%d", fileSize))
	req.Header.Set("x-upsert", "true")

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("upload failed: %s", body)
	}

	publicURL := fmt.Sprintf("%s/storage/v1/object/public/%s/%s", projectURL, bucketName, fileName)
	return publicURL, nil
}
