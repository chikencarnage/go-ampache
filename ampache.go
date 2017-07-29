package ampache

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"
)

const apiPath = "server/xml.server.php?"

type Connection struct {
	Host     string
	Username string
	Password string
	ApiKey   string
	AuthType string
	Client   *http.Client
}

func generatePassphrase(password string) string {
	now := time.Now()
	utime := now.Unix()

	key := sha256.New()
	key.Write([]byte(password))

	passphrase := sha256.New()
	passphrase.Write([]byte(fmt.Sprintf("%d%x", utime, key)))

	return fmt.Sprintf("%x", passphrase.sum(nil))
}
