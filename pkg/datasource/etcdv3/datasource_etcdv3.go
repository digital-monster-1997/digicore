package etcdv3

import (
	"context"
	"errors"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/digital-monster-1997/digicore/pkg/client/etcdv3"
	"github.com/digital-monster-1997/digicore/pkg/dlog"
	"github.com/digital-monster-1997/digicore/pkg/utils/dgo"
	"time"
)

type etcdv3DataSourceProvider struct {
	propertyKey         string
	lastUpdatedRevision int64
	client              *etcdv3.Client
	// cancel is the func, call cancel will stop watching on the propertyKey
	cancel context.CancelFunc
	// closed indicate whether continuing to watch on the propertyKey
	// closed util.AtomicBool

	logger *dlog.Logger

	changed chan struct{}
}

// NewDataSource new a etcdv3DataSource instance.
// client is the etcdv3 client, it must be useful and should be release by User.
func NewDataSource(client *etcdv3.Client, key string) *etcdv3DataSourceProvider {
	ds := &etcdv3DataSourceProvider{
		client:      client,
		propertyKey: key,
		changed: make(chan struct{}),
	}
	go dgo.RecoverGo(ds.watch, nil)
	return ds
}

// ReadConfig ...
func (s *etcdv3DataSourceProvider) ReadConfig() ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	resp, err := s.client.Get(ctx, s.propertyKey)
	if err != nil {
		return nil, err
	}
	if resp.Count == 0 {
		return nil, errors.New("empty response")
	}
	s.lastUpdatedRevision = resp.Header.GetRevision()
	return resp.Kvs[0].Value, nil
}

// IsConfigChanged ...
func (s *etcdv3DataSourceProvider) IsConfigChanged() <-chan struct{} {
	return s.changed
}

func (s *etcdv3DataSourceProvider) handle(resp *clientv3.WatchResponse) {
	if resp.CompactRevision > s.lastUpdatedRevision {
		s.lastUpdatedRevision = resp.CompactRevision
	}
	if resp.Header.GetRevision() > s.lastUpdatedRevision {
		s.lastUpdatedRevision = resp.Header.GetRevision()
	}

	if err := resp.Err(); err != nil {
		return
	}

	for _, ev := range resp.Events {
		if ev.Type == mvccpb.PUT || ev.Type == mvccpb.DELETE {
			s.changed <- struct{}{}
		}
	}
}

func (s *etcdv3DataSourceProvider) watch() {
	ctx, cancel := context.WithCancel(context.Background())
	s.cancel = cancel
	rch := s.client.Watch(ctx, s.propertyKey, clientv3.WithCreatedNotify(), clientv3.WithRev(s.lastUpdatedRevision))
	for {
		for resp := range rch {
			s.handle(&resp)
		}
		time.Sleep(time.Second)

		ctx, cancel = context.WithCancel(context.Background())
		if s.lastUpdatedRevision > 0 {
			rch = s.client.Watch(ctx, s.propertyKey, clientv3.WithCreatedNotify(), clientv3.WithRev(s.lastUpdatedRevision))
		} else {
			rch = s.client.Watch(ctx, s.propertyKey, clientv3.WithCreatedNotify())
		}
		s.cancel = cancel
	}
}

// Close ...
func (s *etcdv3DataSourceProvider) Close() error {
	s.cancel()
	return nil
}
