package mapbox

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

type Client struct {
	apiKey   string
	username string
}

func NewClient(apiKey, username string) *Client {
	return &Client{apiKey: apiKey, username: username}
}

// CreateOrReplace tileset creates a tileset from a file at a given path (shp files must be zips)
func (c *Client) CreateOrReplaceTileset(path, tilesetName string) error {
	// see https: //www.mapbox.com/api-documentation/#create-an-upload

	// tempCreds are the corresponding s3 creds sent back by mapbox when requesting a file upload.
	type tempCreds struct {
		AccessKeyID     string `json:"accessKeyId"`
		Bucket          string `json:"bucket"`
		Key             string `json:"key"`
		SecretAccessKey string `json:"secretAccessKey"`
		SessionToken    string `json:"sessionToken"`
		URL             string `json:"url"`
	}

	getS3Creds := func(apiKey, username string) (*tempCreds, error) {
		fmt.Println(fmt.Sprintf("https://api.mapbox.com/uploads/v1/%s/credentials?access_token=%s", username, apiKey))
		resp, err := http.Get(fmt.Sprintf("https://api.mapbox.com/uploads/v1/%s/credentials?access_token=%s", username, apiKey))
		if err != nil {
			return nil, err
		}

		var creds tempCreds
		return &creds, json.NewDecoder(resp.Body).Decode(&creds)
	}

	cpToS3Bucket := func(path string, creds *tempCreds) error {

		upload := exec.Command("aws", "s3", "cp",
			path, fmt.Sprintf("s3://%s/%s", creds.Bucket, creds.Key),
			"--region", "us-east-1")

		// set the creds
		upload.Env = append(upload.Env, []string{
			fmt.Sprintf("AWS_ACCESS_KEY_ID=%s", creds.AccessKeyID),
			fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s", creds.SecretAccessKey),
			fmt.Sprintf("AWS_SESSION_TOKEN=%s", creds.SessionToken),
		}...)

		if out, err := upload.CombinedOutput(); err != nil {
			log.Println(string(out))
			return err
		}
		return nil
	}

	startCreateOrReplaceTileset := func(tilesetName, username, apikey string, creds *tempCreds) error {
		req := struct {
			URL     string `json:"url"`
			Tileset string `json:"tileset"`
		}{
			URL:     fmt.Sprintf("http://%s.s3.amazonaws.com/%s", creds.Bucket, creds.Key),
			Tileset: fmt.Sprintf("%s.%s", username, tilesetName),
		}

		body := &bytes.Buffer{}
		if err := json.NewEncoder(body).Encode(&req); err != nil {
			return err
		}

		_, err := http.Post(fmt.Sprintf("https://api.mapbox.com/uploads/v1/%s?access_token=%s", username, apikey), "application/json", body)
		if err != nil {
			return err
		}

		return nil
	}

	tmpCreds, err := getS3Creds(c.apiKey, c.username)
	if err != nil {
		return err
	}

	if err := cpToS3Bucket(path, tmpCreds); err != nil {
		return err
	}

	if err := startCreateOrReplaceTileset(tilesetName, c.username, c.apiKey, tmpCreds); err != nil {
		return err
	}

	return nil
}
