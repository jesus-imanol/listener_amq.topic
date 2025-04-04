package utils

import (
    "bytes"
    "encoding/json"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "time"
)

// LogConsumer realiza el proceso de autenticación y retorna el Bearer Token
func LogConsumer() (string, error) {
    username := os.Getenv("USERNAME")
    password := os.Getenv("PASSWORD")

    loginData := map[string]string{
        "username": username,
        "password": password,
    }

    jsonData, err := json.Marshal(loginData)
    if err != nil {
        return "", err
    }

    apiURL := "http://127.0.0.1:8080/v1/users/login"
    req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return "", err
    }

    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{
        Timeout: time.Second * 10,
    }

    resp, err := client.Do(req)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    token := resp.Header.Get("Authorization")
    if token == "" {
        return "", err 
    }

    log.Printf("API response status: %s", resp.Status)

    body, _ := ioutil.ReadAll(resp.Body)
    log.Printf("API response body: %s", string(body))

    return token, nil
}