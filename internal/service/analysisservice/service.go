package analysisservice

import "github.com/edcnogueira/valoguard-api/internal/providers/henrik"

type Service struct {
	client henrik.Client
}

func New(client henrik.Client) Service {
	return Service{
		client: client,
	}
}
