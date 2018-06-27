package ampache

import (
	"crypto/sha256"
	"encoding/xml"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"time"
)

const (
	apiuri     string = "server/xml.server.php"
	apiversion int64  = 350001
)

// Connection holds the general connection parameters
type Connection struct {
	Host       string
	APIVersion int64
	auth       Auth
	client     *http.Client
	passphrase passphrase
}

// Auth holds the information parsed from the servers handshake response
type Auth struct {
	Token    string `xml:"auth"`           // Authentication token
	Version  string `xml:"api"`            // API version
	Expire   string `xml:"session_expire"` // ISO8601 date when session expires
	Update   string `xml:"update"`         // ISO 8601 date of last update
	Add      string `xml:"add"`            // ISO 8601 date of last add
	Clean    string `xml:"clean"`          // ISO 8601 date of last clean
	Songs    string `xml:"songs"`          // number of songs
	Artists  string `xml:"artists"`        // number of artists
	Albums   string `xml:"albums"`         // number of albums
	Catalogs string `xml:"catalogs"`       // number of tags
	Videos   string `xml:"videos"`         // number of videos
	Error    string `xml:"error"`          // authentication errors
}

type passphrase struct {
	hash string
	time int64
}

// NewConnection will return a Connection object specifying APIVersion to use and an http client
func NewConnection(host string) *Connection {
	// create the full api url
	u, err := url.Parse(host)
	// TODO (David Splittberger) handle this for real
	if err != nil {
		log.Println("um, fix me")
	}
	u.Path = path.Join(u.Path, apiuri)
	return &Connection{Host: u.String(), APIVersion: apiversion, client: makeHTTPClient()}
}

// PasswordAuth authenticates with the host defined in *Connection using usename/password
func (c *Connection) PasswordAuth(username, password string) error {
	hashinfo, err := generatePassphrase(password)
	if err != nil {
		return err
	}

	req := makeHTTPRequest(c.Host, map[string]string{
		"action":    "handshake",
		"auth":      hashinfo.hash,
		"timestamp": strconv.FormatInt(hashinfo.time, 10),
		"version":   strconv.FormatInt(c.APIVersion, 10),
		"user":      username,
	})

	response, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		// TODO (David Splittberger) Put a better error message here. Include the HTTP status message
		return fmt.Errorf("Did not get 200 response code. Got %d", response.StatusCode)
	}

	err = xml.NewDecoder(response.Body).Decode(&c.auth)
	if err != nil {
		// TODO (David Splitberger) If we failed to decode, we probably got an xml error from the server
		// We need to try and unmarshal to a xml error struct
		log.Printf("failed to unmarshal\n%s", err)
		return err
	}

	if c.auth.Error != "" {
		return errors.New(c.auth.Error)
	}

	return nil
}

// APIAuth authenticates with the host defined in *Connection using an APIKey
func (c *Connection) APIAuth(apiKey string) error {
	req := makeHTTPRequest(c.Host, map[string]string{
		"action":  "handshake",
		"auth":    apiKey,
		"version": strconv.FormatInt(c.APIVersion, 10),
	})

	response, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return fmt.Errorf("Did not get 200 response code. Got %d", response.StatusCode)
	}

	err = xml.NewDecoder(response.Body).Decode(&c.auth)
	if err != nil {
		log.Printf("failed to unmarshal\n%s", err)
		return err
	}

	return nil
}

// Ping will prolong a session by contacting the ampache server using the ping method
func (c *Connection) Ping() {
	// create the http request
	req := makeHTTPRequest(c.Host, map[string]string{
		"action": "ping",
		"auth":   c.auth.Token,
	})

	response, err := c.client.Do(req)
	if err != nil {
		response.Body.Close()
		panic("Ping failed! Please fix me to fail gracefully!")
	}

	// TODO (David Splittberger) Get the updated session info from the ping response
	response.Body.Close()
}

func generatePassphrase(password string) (passphrase, error) {
	utime := time.Now().Unix()
	info := passphrase{time: utime}

	hashOne := sha256.New()
	_, err := hashOne.Write([]byte(password))
	if err != nil {
		return info, err
	}

	hashTwo := sha256.New()
	_, err = hashTwo.Write([]byte(fmt.Sprintf("%d%x", utime, hashOne.Sum(nil))))
	if err != nil {
		return info, err
	}

	info.hash = fmt.Sprintf("%x", hashTwo.Sum(nil))
	return info, nil
}

func makeHTTPClient() *http.Client {
	client := &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 20,
		},
		Timeout: time.Duration(10) * time.Second,
	}

	return client
}

func makeHTTPRequest(host string, params map[string]string) *http.Request {
	req, _ := http.NewRequest("GET", host, nil)

	query := req.URL.Query()
	for k, v := range params {
		query.Add(k, v)
	}

	req.URL.RawQuery = query.Encode()
	return req
}
