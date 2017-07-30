package ampache

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	apipath    string = "server/xml.server.php?"
	apiversion uint32 = 350001
)

// Connection holds the general connection parameters
type Connection struct {
	Host       string
	APIVersion uint32
	Auth       Auth
	client     *http.Client
}

// Auth object contains the handshake response information
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
}

type passphrase struct {
	hash string
	time string
}

// NewConnection will return a Connection object specifying APIVersion to use and an http client
func NewConnection(url string) *Connection {
	return &Connection{Host: url, APIVersion: apiversion, client: makeHTTPClient()}
}

// PasswordAuth authenticates with the host defined in *Connection using usename/password
func (c *Connection) PasswordAuth(username, password string) error {
	hashinfo, err := generatePassphrase(password)
	response, err := c.client.Get(fmt.Sprintf("%s/%saction=handshake&auth=%s&timestamp=%s&version=%d&user=%s", c.Host, apipath, hashinfo.hash, hashinfo.time, c.APIVersion, username))
	defer response.Body.Close()
	if err != nil {
		return err
	} else if response.StatusCode != 200 {
		return fmt.Errorf("Did not get 200 response code. Got %d", response.StatusCode)
	}

	err = xml.NewDecoder(response.Body).Decode(&c.Auth)
	if err != nil {
		log.Printf("failed to unmarshal\n%s", err)
	}

	go c.ping()

	return nil
}

// APIAuth authenticates with the host defined in *Connection using an APIKey
func (c *Connection) APIAuth(apiKey string) error {
	response, err := c.client.Get(fmt.Sprintf("%s/%saction=handshake&auth=%s&version=%d", c.Host, apipath, apiKey, c.APIVersion))
	defer response.Body.Close()
	if err != nil {
		return err
	} else if response.StatusCode != 200 {
		return fmt.Errorf("Did not get 200 response code. Got %d", response.StatusCode)
	}

	err = xml.NewDecoder(response.Body).Decode(&c.Auth)
	if err != nil {
		log.Printf("failed to unmarshal")
	}

	go c.ping()

	return nil
}

func (c *Connection) ping() {
	// TODO (David Splittberger): base this on the expiration value sent by server
	for {
		response, err := c.client.Get(fmt.Sprintf("%s/%saction=ping&auth=%s", c.Host, apipath, c.Auth.Token))
		if err != nil {
			response.Body.Close()
			panic("Ping failed! Please fix me to fail gracefully!")
		}
		response.Body.Close()
		time.Sleep(300 * time.Second)
	}
}

func generatePassphrase(password string) (passphrase, error) {
	utime := int32(time.Now().Unix())
	info := passphrase{time: fmt.Sprintf("%d", utime)}

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
