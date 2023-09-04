package infra

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"go.etcd.io/etcd/clientv3/concurrency"
	"log"
	"time"
)

type SuperGlobalEnv interface {
	Transaction(func(stm concurrency.STM) error) (bool, error)
	Put(key string, value string) error
}

var EtcdClient SuperGlobalEnv

type _superGlobalEnv struct {
	etcdClient *clientv3.Client
}

func (env *_superGlobalEnv) Transaction(f func(stm concurrency.STM) error) (bool, error) {
	txn, err := concurrency.NewSTM(env.etcdClient, func(stm concurrency.STM) error {
		return f(stm)
	})
	return txn.Succeeded, err
}

func (env *_superGlobalEnv) Put(key string, value string) error {
	_, err := env.etcdClient.Put(context.Background(), key, value)
	return err
}

//setup etcd
func SetupEtcd() error {
	//setup etcd
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
		return err
	}
	EtcdClient = &_superGlobalEnv{etcdClient: etcdClient}
	return nil
}
