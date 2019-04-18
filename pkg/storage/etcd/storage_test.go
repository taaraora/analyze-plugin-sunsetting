package etcd_test

import (
	"context"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/supergiant/analyze-plugin-sunsetting/pkg/storage"
	"github.com/supergiant/analyze-plugin-sunsetting/pkg/storage/etcd"

	"google.golang.org/grpc"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/embed"
)

type msg []byte

func (m msg) Payload() []byte {
	return m
}

const msg1, msg2 = "ololo1", "ololo2"

func getTestClientConfig() clientv3.Config {
	return clientv3.Config{
		Endpoints: []string{embed.DefaultAdvertiseClientURLs},
	}
}

const testsTimeout = time.Second * 5

func startEtcdServer() (func() error, error) {

	// we need to have unused tmp dir for broker
	tmpDir, err := ioutil.TempDir("", "etcd.embedded")
	if err != nil {
		return nil, errors.New("can't create temp dir")
	}

	cfg := embed.NewConfig()
	cfg.Dir = filepath.Join(tmpDir, "default.etcd")
	cfg.WalDir = filepath.Join(tmpDir, "default.wal")
	e, err := embed.StartEtcd(cfg)
	if err != nil {
		return nil, err
	}

	select {
	case <-e.Server.ReadyNotify():
		log.Printf("Server is ready!")
		return func() error {
			e.Server.Stop()
			e.Close()
			return os.RemoveAll(tmpDir)
		}, nil
	case <-time.After(60 * time.Second):
		e.Server.Stop() // trigger a shutdown
		log.Printf("Server took too long to start!")
	}

	return nil, <-e.Err()
}

func TestMain(m *testing.M) {
	stopFunc, err := startEtcdServer()
	if err != nil {
		log.Printf("server initialization error: %+v", err)
		os.Exit(1)
	}

	retCode := m.Run()

	err = stopFunc()
	if err != nil {
		log.Printf("server shutdown error: %+v", err)
		os.Exit(1)
	}
	os.Exit(retCode)
}

func TestNewETCDStorage_WithoutEndpoint(t *testing.T) {
	_, err := etcd.NewETCDStorage(clientv3.Config{}, logrus.New())
	if err == nil {
		t.Fatal("how we can connect to etcd server without uri")
	}
}

func TestNewETCDStorage_IncorrectEndpoint(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()

	_, err := etcd.NewETCDStorage(clientv3.Config{
		Endpoints:   []string{"http://wrong_host:9000"},
		DialOptions: []grpc.DialOption{grpc.WithBlock() /*grpc.WithTimeout(time.Second)*/},
		Context:     ctx,
	}, logrus.New())
	if err == nil {
		t.Fatal("how we can connect to incorrect etcd server uri")
	}
}

func TestNewETCDStorage_Close(t *testing.T) {
	stor, err := etcd.NewETCDStorage(getTestClientConfig(), logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()
	err = stor.Put(ctx, "test_get", "12345", msg([]byte(msg1)))
	if err != nil {
		t.Fatalf("can't put message into etcd %s", err.Error())
	}

	err = stor.Close()
	if err != nil {
		t.Fatalf("gracefull client close need to return nil %+v", err.Error())
	}

	_, err = stor.Get(ctx, "test_get", "12345")
	if err == nil {
		t.Fatal("it shouldn't read from etcd when client is closed")
	}
}

func TestETCDStorage_Get_Successfully(t *testing.T) {
	stor, err := etcd.NewETCDStorage(getTestClientConfig(), logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()
	err = stor.Put(ctx, "test_get", "12345", msg([]byte(msg1)))
	if err != nil {
		t.Fatalf("can't put message into etcd %s", err.Error())
	}

	m, err := stor.Get(ctx, "test_get", "12345")
	if err != nil {
		t.Fatalf("can't get message from etcd %s", err.Error())
	}

	if string(m.Payload()) != msg1 {
		t.Fatalf("message from etcd is wrong: %s", string(m.Payload()))
	}
}

func TestETCDStorage_Get_NotFound(t *testing.T) {
	stor, err := etcd.NewETCDStorage(getTestClientConfig(), logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()

	_, err = stor.Get(ctx, "test_not_found_get", "12345")
	if err != storage.ErrNotFound {
		t.Fatalf("should be ErrNotFound but got: %+v", err)
	}
}

func TestETCDStorage_Get_ClientError(t *testing.T) {
	stor, err := etcd.NewETCDStorage(clientv3.Config{
		Endpoints: []string{"http://wrong_host:9000"},
	}, logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()

	_, err = stor.Get(ctx, "test_not_found_get_client_error", "12345")
	if err == nil {
		t.Fatal("it should return error when etcd client is broken")
	}
}

func TestNewETCDStorage_GetAll_ClientError(t *testing.T) {
	stor, err := etcd.NewETCDStorage(clientv3.Config{
		Endpoints: []string{"http://wrong_host:9000"},
	}, logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()

	_, err = stor.GetAll(ctx, "test_get_all_wrong_host")
	if err == nil {
		t.Fatal("it should return error when etcd client is broken")
	}

}

func TestNewETCDStorage_GetAll_Successfully(t *testing.T) {
	stor, err := etcd.NewETCDStorage(getTestClientConfig(), logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()

	err = stor.Put(ctx, "test_get_all", "12345", msg([]byte(msg1)))
	if err != nil {
		t.Fatalf("can't put message into etcd %s", err.Error())
	}
	err = stor.Put(ctx, "test_get_all", "23456", msg([]byte(msg2)))
	if err != nil {
		t.Fatalf("can't put message into etcd %s", err.Error())
	}

	messages, err := stor.GetAll(ctx, "test_get_all")
	if err != nil {
		t.Fatalf("can't get all messages from etcd %s", err.Error())
	}
	if len(messages) != 2 {
		t.Fatalf("number of returned messages is worng, gor %v, need 2", len(messages))
	}

	result := make(map[string]struct{})
	for _, message := range messages {
		result[string(message.Payload())] = struct{}{}
	}

	_, exists := result[msg1]
	if !exists {
		t.Fatalf("message %s wasn't retirved from etcd", msg1)
	}

	_, exists = result[msg2]
	if !exists {
		t.Fatalf("message %s wasn't retirved from etcd", msg2)
	}
}

func TestETCDStorage_Put_Successfully(t *testing.T) {
	stor, err := etcd.NewETCDStorage(getTestClientConfig(), logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()
	err = stor.Put(ctx, "test_put", "012345", msg([]byte(msg1)))
	if err != nil {
		t.Fatalf("can't put message into etcd %s", err.Error())
	}

	m, err := stor.Get(ctx, "test_put", "012345")
	if err != nil {
		t.Fatalf("can't get message from etcd %s", err.Error())
	}

	if string(m.Payload()) != msg1 {
		t.Fatalf("message from etcd is wrong: %s", string(m.Payload()))
	}
}

func TestETCDStorage_Put_ClientError(t *testing.T) {
	stor, err := etcd.NewETCDStorage(clientv3.Config{
		Endpoints: []string{"http://wrong_host:9000"},
	}, logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()
	err = stor.Put(ctx, "test_put", "012345", msg([]byte(msg1)))
	if err == nil {
		t.Fatal("it should return error when client is broken")
	}
}

func TestETCDStorage_Delete_Successfully(t *testing.T) {
	stor, err := etcd.NewETCDStorage(getTestClientConfig(), logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()
	err = stor.Put(ctx, "test_delete", "012345", msg([]byte(msg1)))
	if err != nil {
		t.Fatalf("can't put message into etcd %s", err.Error())
	}

	m, err := stor.Get(ctx, "test_delete", "012345")
	if err != nil {
		t.Fatalf("can't get message from etcd %s", err.Error())
	}

	if string(m.Payload()) != msg1 {
		t.Fatalf("message from etcd is wrong: %s", string(m.Payload()))
	}

	err = stor.Delete(ctx, "test_delete", "012345")
	if err != nil {
		t.Fatalf("can't delete message from etcd %s", err.Error())
	}

	_, err = stor.Get(ctx, "test_delete", "012345")
	if err != storage.ErrNotFound {
		t.Fatalf("message from etcd need to be not found %+v", err)
	}
}

func TestETCDStorage_Delete_ErrNotFound(t *testing.T) {
	stor, err := etcd.NewETCDStorage(getTestClientConfig(), logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()

	err = stor.Delete(ctx, "test_delete_error", "012345")
	if err != storage.ErrNotFound {
		t.Fatalf("it should not delete nonexistent kv %+v", err)
	}
}

func TestETCDStorage_Delete_ClientError(t *testing.T) {
	stor, err := etcd.NewETCDStorage(clientv3.Config{
		Endpoints: []string{"http://wrong_host:9000"},
	}, logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()

	err = stor.Delete(ctx, "test_delete_timeout", "012345")
	if err == nil || err == storage.ErrNotFound {
		t.Fatalf("it should not delete kv when client connection is absent %+v", err)
	}
}

func TestETCDStorage_WatchPrefix_ReturnsExistingKVs(t *testing.T) {
	stor, err := etcd.NewETCDStorage(getTestClientConfig(), logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()

	err = stor.Put(ctx, "test_watch_prefix", "912345", msg([]byte(msg1)))
	if err != nil {
		t.Fatalf("can't put message into etcd %s", err.Error())
	}
	err = stor.Put(ctx, "test_watch_prefix", "23456", msg([]byte(msg2)))
	if err != nil {
		t.Fatalf("can't put message into etcd %s", err.Error())
	}

	events := stor.WatchPrefix(ctx, "test_watch_prefix")
	result := make(map[string]struct{})
	deadline := time.NewTimer(testsTimeout)
	for {
		if len(result) == 2 {
			break
		}
		select {
		case event := <-events:
			result[string(event.Payload())] = struct{}{}
		case <-deadline.C:
			t.Fatalf("we didn't get two messages in defined timeout %v", testsTimeout)
		}
	}

	_, exists := result[msg1]
	if !exists {
		t.Fatalf("message %s wasn't retirved from etcd watch", msg1)
	}

	_, exists = result[msg2]
	if !exists {
		t.Fatalf("message %s wasn't retirved from etcd watch", msg2)
	}
	err = stor.Close()
	if err != nil {
		t.Fatalf("gracefull client close need to return nil %+v", err.Error())
	}
}

func TestETCDStorage_WatchPrefix_ClientError(t *testing.T) {
	stor, err := etcd.NewETCDStorage(clientv3.Config{
		Endpoints: []string{"http://wrong_host:9000"},
	}, logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()

	events := stor.WatchPrefix(ctx, "test_watch_prefix")
	var result storage.WatchEvent
	deadline := time.NewTimer(testsTimeout)

	select {
	case event := <-events:
		result = event
	case <-deadline.C:
		t.Fatalf("we didn't get two messages in defined timeout %v", testsTimeout)
	}

	if result.Type() != storage.Error {
		t.Fatal("when etcd client is broken we have to get message with type error")
	}

	if result.Err() == nil {
		t.Fatal("when etcd client is broken we have to get error")
	}
}

func TestETCDStorage_WatchPrefix_ReturnsNewKVs(t *testing.T) {
	stor, err := etcd.NewETCDStorage(getTestClientConfig(), logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()

	events := stor.WatchPrefix(ctx, "test_watch_prefix_new")
	result := make(map[string]struct{})
	deadline := time.NewTimer(testsTimeout)
	readFinished := make(chan struct{})
	go func() {
		for {
			if len(result) == 2 {
				break
			}
			select {
			case event := <-events:
				result[string(event.Payload())] = struct{}{}
			case <-deadline.C:
				t.Errorf("we didn't get two messages in defined timeout %v", testsTimeout)
				return
			}
		}
		readFinished <- struct{}{}
	}()

	err = stor.Put(ctx, "test_watch_prefix_new", "912345", msg([]byte(msg1)))
	if err != nil {
		t.Fatalf("can't put message into etcd %s", err.Error())
	}
	err = stor.Put(ctx, "test_watch_prefix_new", "23456", msg([]byte(msg2)))
	if err != nil {
		t.Fatalf("can't put message into etcd %s", err.Error())
	}

	<-readFinished

	_, exists := result[msg1]
	if !exists {
		t.Fatalf("message %s wasn't retirved from etcd watch", msg1)
	}

	_, exists = result[msg2]
	if !exists {
		t.Fatalf("message %s wasn't retirved from etcd watch", msg2)
	}
	err = stor.Close()
	if err != nil {
		t.Fatalf("gracefull client close need to return nil %+v", err.Error())
	}
}

func TestETCDStorage_WatchPrefix_ReturnsUpdatedKVs(t *testing.T) {
	stor, err := etcd.NewETCDStorage(getTestClientConfig(), logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()

	events := stor.WatchPrefix(ctx, "test_watch_prefix_updated")
	result := make([]storage.WatchEvent, 0)
	deadline := time.NewTimer(testsTimeout)
	readFinished := make(chan struct{})
	go func() {
		for {
			if len(result) == 2 {
				break
			}
			select {
			case event := <-events:
				result = append(result, event)
			case <-deadline.C:
				t.Errorf("we didn't get two messages in defined timeout %v", testsTimeout)
				return
			}
		}
		readFinished <- struct{}{}
	}()

	err = stor.Put(ctx, "test_watch_prefix_updated", "912345", msg([]byte(msg1)))
	if err != nil {
		t.Fatalf("can't put message into etcd %s", err.Error())
	}
	err = stor.Put(ctx, "test_watch_prefix_updated", "912345", msg([]byte(msg2)))
	if err != nil {
		t.Fatalf("can't put message into etcd %s", err.Error())
	}

	<-readFinished

	for _, event := range result {
		if event.Type() == storage.Modified {
			return
		}
	}

	t.Fatal("Modified message type wasn't received from Watch")
}

func TestETCDStorage_WatchPrefix_ReturnsDeletedKVs(t *testing.T) {
	stor, err := etcd.NewETCDStorage(getTestClientConfig(), logrus.New())
	if err != nil {
		t.Fatalf("can't create etcd client %s", err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), testsTimeout)
	defer cancel()

	events := stor.WatchPrefix(ctx, "test_watch_prefix_delete")
	result := make([]storage.WatchEvent, 0)
	deadline := time.NewTimer(testsTimeout)
	readFinished := make(chan struct{})
	go func() {
		for {
			if len(result) == 2 {
				break
			}
			select {
			case event := <-events:
				result = append(result, event)
			case <-deadline.C:
				t.Errorf("we didn't get two messages in defined timeout %v", testsTimeout)
				return
			}
		}
		readFinished <- struct{}{}
	}()

	err = stor.Put(ctx, "test_watch_prefix_delete", "912345", msg([]byte(msg1)))
	if err != nil {
		t.Fatalf("can't put message into etcd %s", err.Error())
	}
	err = stor.Delete(ctx, "test_watch_prefix_delete", "912345")
	if err != nil {
		t.Fatalf("can't delete message into etcd %s", err.Error())
	}

	<-readFinished

	for _, event := range result {
		if event.Type() == storage.Deleted {
			return
		}
	}

	t.Fatal("Modified message type wasn't received from Watch")
}
