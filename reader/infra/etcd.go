package infra

import (
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"log"
	"time"
)

type ISuperGlobalEnv interface {
	Transaction(func(stm concurrency.STM) error) (bool, error)
	Put(key string, value string, option clientv3.OpOption) error
	Get(key string) (string, error)
	GetClient() *clientv3.Client
}

type SuperGlobalEnv struct {
	EtcdClient *clientv3.Client
}

func (env *SuperGlobalEnv) GetClient() *clientv3.Client {
	return env.EtcdClient
}

func (env *SuperGlobalEnv) Transaction(f func(stm concurrency.STM) error) (bool, error) {
	txn, err := concurrency.NewSTM(env.EtcdClient, func(stm concurrency.STM) error {
		return f(stm)
	})
	if err != nil {
		return false, err
	}
	return txn.Succeeded, err
}

func (env *SuperGlobalEnv) Get(key string) (string, error) {
	r, err := env.EtcdClient.Get(context.Background(), key)
	if err != nil {
		return "", err
	}
	if len(r.Kvs) != 1 {
		return "", errors.New(fmt.Sprintf("not found key: %d", len(r.Kvs)))
	}
	return string(r.Kvs[0].Value), nil
}

func (env *SuperGlobalEnv) Put(key string, value string, option clientv3.OpOption) error {
	_, err := env.EtcdClient.Put(context.Background(), key, value, option)
	return err
}

//setup etcd
func SetupEtcd() (*SuperGlobalEnv, error) {
	//setup etcd
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &SuperGlobalEnv{EtcdClient: etcdClient}, nil
}
