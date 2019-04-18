package mock

import (
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/supergiant/analyze-plugin-sunsetting/pkg/storage"
)

type mockStorage struct {
	data     map[string][]byte
	isBroken bool
}

type mockMsg []byte

func (m mockMsg) Payload() []byte {
	return m
}

var errBroken = errors.New("internal storage error")

func GetMockStorage(t *testing.T, data map[string]string) storage.Interface {
	t.Helper()
	return NewMockStorage(data, false)
}

func GetMockBrokenStorage(t *testing.T) storage.Interface {
	t.Helper()
	return NewMockStorage(nil, true)
}

func NewMockStorage(data map[string]string, isBroken bool) storage.Interface {
	result := map[string][]byte{}
	for key, value := range data {
		result[key] = []byte(value)
	}

	return &mockStorage{
		data:     result,
		isBroken: isBroken,
	}
}

func (s *mockStorage) GetAll(ctx context.Context, prefix string) ([]storage.Message, error) {
	if s.isBroken {
		return nil, errBroken
	}
	result := []storage.Message{}
	for key := range s.data {
		if strings.Contains(key, prefix) {
			result = append(result, mockMsg(s.data[key]))
		}
	}

	return result, nil
}

func (s *mockStorage) Get(ctx context.Context, prefix string, key string) (storage.Message, error) {
	if s.isBroken {
		return nil, errBroken
	}
	v, ok := s.data[prefix+key]
	if !ok {
		return nil, storage.ErrNotFound
	}

	return mockMsg(v), nil
}

func (s *mockStorage) Put(ctx context.Context, prefix string, key string, value storage.Message) error {
	if s.isBroken {
		return errBroken
	}
	s.data[prefix+key] = value.Payload()

	return nil
}

func (s *mockStorage) Delete(ctx context.Context, prefix string, key string) error {
	if s.isBroken {
		return errBroken
	}
	delete(s.data, prefix+key)
	return nil
}

func (s *mockStorage) Close() error {
	if s.isBroken {
		return errBroken
	}
	return nil
}

func (s *mockStorage) WatchPrefix(ctx context.Context, key string) <-chan storage.WatchEvent {
	panic("not implemented")
}
