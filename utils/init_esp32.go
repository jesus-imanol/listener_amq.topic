package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func InitEsp32(status string, token string) error {
	err := godotenv.Load()
	if err != nil {
        return err
    }
	username := os.Getenv("USERNAME_CONSUMER")
	initEsp32Data := map[string]string{
        "status": "activate",
    }
	jsonData, err := json.Marshal(initEsp32Data)
	if err != nil {
        return err
    }
	apiURL := fmt.Sprintf("http://127.0.0.1:8081/v1/users/login?username=%s", username)
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
        return err
    }

   req.Header.Set("Content-Type", "application/json")
   req.Header.Set("Authorization", token)
   client := &http.Client{
	Timeout: time.Second * 10, 
}

resp, err := client.Do(req)
if err != nil {
	return err
}
defer resp.Body.Close()

log.Println("API response status: ", resp.Status)
return nil
}