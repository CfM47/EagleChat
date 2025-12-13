package diconfig

import (
	"os"
	"time"

	"eaglechat/apps/client/internal/domain/services"
	middleware "eaglechat/apps/client/internal/middleware/domain"
	"eaglechat/apps/client/internal/middleware/infrastructure/clientconnpool"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerpool"
	"eaglechat/apps/client/internal/middleware/infrastructure/idmanagerregisterer"
	jsonmessagecache "eaglechat/apps/client/internal/middleware/infrastructure/messagecache/json"
	jsonusercache "eaglechat/apps/client/internal/middleware/infrastructure/usercache/json"
	"eaglechat/apps/client/internal/tui"
	"eaglechat/common/multicast/implementation"
	"eaglechat/common/simplecrypto/rsa"
)

type MiddlewareDeps struct {
	TUI                  *tui.TUI
	MiddlewareConnector  services.Connector
	MiddlewareRegisterer services.Registerer
	CAPubkey             rsa.PublicKey
}

func NewMiddlewareDeps(TUI *tui.TUI, connector services.Connector, registerer services.Registerer, CAPubkey rsa.PublicKey) *MiddlewareDeps {
	return &MiddlewareDeps{
		TUI:                  TUI,
		MiddlewareConnector:  connector,
		MiddlewareRegisterer: registerer,
		CAPubkey:             CAPubkey,
	}
}

func BuildMiddlewareDeps(idManagerPort string) (*MiddlewareDeps, error) {
	// Build tui
	tui := tui.New()

	// Build connector
	connector, err := buildConnector()
	if err != nil {
		return nil, err
	}

	// Build registerer
	registerer := buildRegisterer(idManagerPort)

	pk, err := getCAPubkey()
	if err != nil {
		return nil, err
	}

	return NewMiddlewareDeps(tui, connector, registerer, *pk), nil
}

func buildConnector() (services.Connector, error) {
	messageCache, err := jsonmessagecache.NewJSONMessageCache("./data/message_cache.json")
	if err != nil {
		return nil, err
	}

	userCache, err := jsonusercache.NewJSONUserCache("./data/user_cache.json", time.Second*10)
	if err != nil {
		return nil, err
	}

	return middleware.NewConnector(
		messageCache,

		clientconnpool.NewClientConnPoolBuilder(),
		idmanagerpool.NewIDManagerPoolBuilder(),

		userCache,
	), nil
}

func buildRegisterer(idManagerPort string) services.Registerer {
	return idmanagerregisterer.NewRegisterer(implementation.DefaultUDPAddress, idManagerPort, time.Second*10)
}

func getCAPubkey() (*rsa.PublicKey, error) {
	path := envOrDefault("CA_PUBLIC_KEY_PATH", "env/ca_public_key.pem")
	return rsa.PublicKeyFromFile(path)
}

func envOrDefault(env, def string) string {
	val := os.Getenv(env)
	if val == "" {
		return def
	}

	return val
}
