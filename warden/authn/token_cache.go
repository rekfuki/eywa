package authn

import (
	"encoding/json"
	"eywa/warden/db"
	"eywa/warden/types"

	"github.com/patrickmn/go-cache"
	log "github.com/sirupsen/logrus"
)

var tokenCache *cache.Cache

type payload struct {
	Action string            `json:"action"`
	Data   types.AccessToken `json:"data"`
}

// InitTokenCache initialises new token cache and sets up postgress listener
func InitTokenCache(db *db.Client) error {
	tokenCache = cache.New(cache.NoExpiration, cache.NoExpiration)
	sub, err := db.Listen("access_tokens")
	if err != nil {
		log.Fatalf("Failed to setup access token notification listener: %s", err)
	}

	allTokens, err := db.GetAllAccessTokens()
	if err != nil {
		return err
	}

	for _, t := range allTokens {
		tokenCache.Add(t.Token, t, cache.DefaultExpiration)
	}

	go listen(sub)

	return nil
}

func listen(s *db.Subscription) {
	defer s.Close()
	for {
		select {
		case err := <-s.ErrChan:
			log.Errorf("Subscription returned an error: %s", err)
			return
		case notif := <-s.Notify:
			if notif.Channel == "access_tokens" {
				var payload payload
				if err := json.Unmarshal([]byte(notif.Payload), &payload); err != nil {
					log.Errorf("Failed to unmrashal notification: %s", err)
					continue
				}

				if payload.Action == "INSERT" {
					err := tokenCache.Add(payload.Data.Token, payload.Data, cache.DefaultExpiration)
					if err != nil {
						log.Errorf("Failed to add new token to cache: %s", err)
					}
				} else if payload.Action == "DELETE" {
					tokenCache.Delete(payload.Data.Token)
				}
			}
		}
	}
}

func tokenExists(token string) bool {
	_, found := tokenCache.Get(token)

	return found
}
