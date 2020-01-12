package session

import (
	"encoding/json"
	"github.com/go-redis/redis"
	"github.com/saas/hostgolang/pkg/types"
)
//go:generate mockgen -destination=mocks/session_store.go -package=mocks github.com/saas/hostgolang/pkg/session Store
type Store interface {
	Put(token string, account *types.Account) error
	Get(token string) (*types.Account, error)
}

type redisSessionStore struct {
	client *redis.Client
}

func (r *redisSessionStore) Put(token string, account *types.Account) error {
	data, err := json.Marshal(account)
	if err != nil {
		return err
	}
	return r.client.Set(token, data, 0).Err()
}

func (r *redisSessionStore) Get(token string) (*types.Account, error) {
	data, err := r.client.Get(token).Bytes()
	if err != nil {
		return nil, err
	}
	account := &types.Account{}
	if err := json.Unmarshal(data, account); err != nil {
		return nil, err
	}
	return account, nil
}

func NewRedisSessionStore(client *redis.Client) Store {
	return &redisSessionStore{client:client}
}
