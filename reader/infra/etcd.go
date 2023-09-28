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
	Put(key string, value string) error
	Get(key string) (string, error)
}

type SuperGlobalEnv struct {
	etcdClient *clientv3.Client
}

func (env *SuperGlobalEnv) Transaction(f func(stm concurrency.STM) error) (bool, error) {
	txn, err := concurrency.NewSTM(env.etcdClient, func(stm concurrency.STM) error {
		return f(stm)
	})
	if err != nil {
		return false, err
	}
	return txn.Succeeded, err
}

func (env *SuperGlobalEnv) Get(key string) (string, error) {
	r, err := env.etcdClient.Get(context.Background(), key)
	if err != nil {
		return "", err
	}
	if len(r.Kvs) != 1 {
		return "", errors.New(fmt.Sprintf("not found key: %d", len(r.Kvs)))
	}
	return string(r.Kvs[0].Value), nil
}

func (env *SuperGlobalEnv) Put(key string, value string) error {
	_, err := env.etcdClient.Put(context.Background(), key, value)
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
	return &SuperGlobalEnv{etcdClient: etcdClient}, nil
}
