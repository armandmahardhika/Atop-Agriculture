package datatunnel

import "go.mongodb.org/mongo-driver/bson"

// Defined DATA -> IN NODE -> MIDDLEWARE -> OUT NODES
type TMiddleware struct {
}

type Tunnel struct {
	Middlewares []TMiddleware
}

func (t *Tunnel) input(data bson.M) error {
	return nil
}
