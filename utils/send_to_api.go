package utils

import (
    "bytes"
    "log"
    "net/http"
    "time"
)

func SendToAPI(token string, data []byte) error {
    apiURL := "http://52.23.135.169:4000/v1/message/consumer/"

    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(data))
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