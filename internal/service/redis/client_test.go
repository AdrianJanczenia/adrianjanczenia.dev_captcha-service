package redis

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/go-redis/redis"
	"github.com/go-redis/redismock/v8"
)

func TestClient_Ping(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := &client{rdb: db}
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mock.ExpectPing().SetVal("PONG")
		err := client.Ping(ctx)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectPing().SetErr(errors.New("connection failed"))
		err := client.Ping(ctx)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestClient_Set(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := &client{rdb: db}
	ctx := context.Background()
	key := "test-key"
	val := "test-val"
	ttl := time.Minute

	t.Run("success", func(t *testing.T) {
		mock.ExpectSet(key, val, ttl).SetVal("OK")
		err := client.Set(ctx, key, val, ttl)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectSet(key, val, ttl).SetErr(errors.New("redis error"))
		err := client.Set(ctx, key, val, ttl)
		if err == nil {
			t.Error("expected error, got nil")
		}
	})
}

func TestClient_SetNX(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := &client{rdb: db}
	ctx := context.Background()
	for _, tc := range []struct {
		name    string
		created bool
	}{
		{name: "new seed", created: true},
		{name: "used seed", created: false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			mock.ExpectSetNX("pow:seed", "1", 5*time.Minute).SetVal(tc.created)
			created, err := client.SetNX(ctx, "pow:seed", "1", 5*time.Minute)
			if err != nil || created != tc.created {
				t.Fatalf("SetNX() = %t, %v; want %t, nil", created, err, tc.created)
			}
		})
	}
}

func TestClient_Get(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := &client{rdb: db}
	ctx := context.Background()
	key := "test-key"

	t.Run("success", func(t *testing.T) {
		mock.ExpectGet(key).SetVal("val")
		res, err := client.Get(ctx, key)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if res != "val" {
			t.Errorf("got %v, want val", res)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectGet(key).SetErr(redis.Nil)
		_, err := client.Get(ctx, key)
		if !errors.Is(err, redis.Nil) {
			t.Errorf("expected redis.Nil, got %v", err)
		}
	})
}

func TestClient_Del(t *testing.T) {
	db, mock := redismock.NewClientMock()
	client := &client{rdb: db}
	ctx := context.Background()
	key := "test-key"

	t.Run("success", func(t *testing.T) {
		mock.ExpectDel(key).SetVal(1)
		err := client.Del(ctx, key)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
}
