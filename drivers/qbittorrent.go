package drivers

import (
	"fmt"
	"net/http"
	"strings"
)

func init() {
	registerDriver(QBitTorrentDriver{})
}

type QBitTorrentDriver struct {
}

// Name implements Driver.
func (QBitTorrentDriver) Name() string {
	return "QBitTorrent"
}

// AddMagnetURL implements Driver.
func (d QBitTorrentDriver) AddMagnetURL(config *Config, magnet string) error {
	baseURL := strings.TrimRight(config.URL, "/")

	client := newClient()

	if err := d.authenticate(config, baseURL, client); err != nil {
		return err
	}

	if err := d.addURLs(config, baseURL, client, magnet); err != nil {
		return err
	}

	return nil
}

func (d QBitTorrentDriver) authenticate(config *Config, baseURL string, client *http.Client) error {
	body := marshalFieldsURLEncoded(map[string]string{
		"username": config.Username,
		"password": config.Password,
	})

	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/v2/auth/login", body.Reader())
	if err != nil {
		return err
	}
	req.Header.Set("content-type", body.ContentType)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	cookies := resp.Cookies()
	for _, cookie := range cookies {
		if cookie.Name == "SID" {
			return nil
		}
	}

	return fmt.Errorf("could not get authentication token (status=%d)", resp.StatusCode)
}

func (d QBitTorrentDriver) addURLs(config *Config, baseURL string, client *http.Client, urls ...string) error {
	body := marshalFieldsMultipart(map[string]string{
		"urls": strings.Join(urls, "\n"),
	})

	req, err := http.NewRequest(http.MethodPost, baseURL+"/api/v2/torrents/add", body.Reader())
	if err != nil {
		return err
	}
	req.Header.Set("content-type", body.ContentType)

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("could not upload links (status=%d)", resp.StatusCode)
	}

	return nil
}
