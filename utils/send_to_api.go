package utils

import (
    "bytes"
    "log"
    "net/http"
    "time"
)

// SendToAPI envía datos JSON a la API
func SendToAPI(data []byte) error {
    apiURL := "http://127.0.0.1:4000/v1/message/"

    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(data))
    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{
        Timeout: time.Second * 10,
    }

    resp, err := client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    log.Printf("API response status: %s", resp.Status)
    return nil
}