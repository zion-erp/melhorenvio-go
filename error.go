package melhorenvio

import "errors"

var (
	ErrClientNotInitialized = errors.New("melhor envio: client not initialized")
	ErrInvalidToken         = errors.New("melhor envio: invalid token")
)
