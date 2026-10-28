//go:build integration

package redis

import (
	"context"
	"encoding/json"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestCaptchaScriptAgainstRedis(t *testing.T) {
	redisURL := os.Getenv("REDIS_TEST_URL")
	if redisURL == "" {
		t.Skip("REDIS_TEST_URL is not set")
	}
	clientInterface, err := NewClient(redisURL)
	if err != nil {
		t.Fatal(err)
	}
	defer clientInterface.Close()
	client := clientInterface.(*client)
	ctx := context.Background()
	if err := client.Ping(ctx); err != nil {
		t.Fatal(err)
	}

	id := uuid.NewString()
	key := "captcha:" + id
	seedKey := "pow:" + id
	defer client.rdb.Del(ctx, key, seedKey)
	if err := client.Set(ctx, key, `{"value":"ABC","triesLeft":3,"solved":false}`, 5*time.Second); err != nil {
		t.Fatal(err)
	}
	ttlBefore := client.rdb.PTTL(ctx, key).Val()
	status, err := client.VerifyCaptcha(ctx, key, "wrong")
	if err != nil || status != "invalid_value" {
		t.Fatalf("wrong answer: status=%q err=%v", status, err)
	}
	stored, err := client.Get(ctx, key)
	if err != nil {
		t.Fatal(err)
	}
	var state struct {
		TriesLeft int  `json:"triesLeft"`
		Solved    bool `json:"solved"`
	}
	if err := json.Unmarshal([]byte(stored), &state); err != nil || state.TriesLeft != 2 || state.Solved {
		t.Fatalf("wrong answer changed state incorrectly: %+v err=%v", state, err)
	}
	if ttlAfter := client.rdb.PTTL(ctx, key).Val(); ttlAfter <= 0 || ttlAfter > ttlBefore {
		t.Fatalf("CAPTCHA TTL changed unexpectedly: before=%v after=%v", ttlBefore, ttlAfter)
	}
	status, err = client.VerifyCaptcha(ctx, key, "ABC")
	if err != nil || status != "solved" {
		t.Fatalf("correct answer: status=%q err=%v", status, err)
	}
	status, err = client.VerifyCaptcha(ctx, key, "wrong")
	if err != nil || status != "invalid_value" {
		t.Fatalf("wrong answer to solved challenge: status=%q err=%v", status, err)
	}

	created, err := client.SetNX(ctx, seedKey, "1", 5*time.Second)
	if err != nil || !created {
		t.Fatalf("first seed use: created=%v err=%v", created, err)
	}
	created, err = client.SetNX(ctx, seedKey, "1", 5*time.Second)
	if err != nil || created {
		t.Fatalf("reused seed: created=%v err=%v", created, err)
	}
}

func TestCaptchaAttemptsAreAtomicAgainstRedis(t *testing.T) {
	redisURL := os.Getenv("REDIS_TEST_URL")
	if redisURL == "" {
		t.Skip("REDIS_TEST_URL is not set")
	}
	client, err := NewClient(redisURL)
	if err != nil {
		t.Fatal(err)
	}
	defer client.Close()
	ctx := context.Background()
	key := "captcha:" + uuid.NewString()
	defer client.Del(ctx, key)
	if err := client.Set(ctx, key, `{"value":"ABC","triesLeft":3,"solved":false}`, 5*time.Second); err != nil {
		t.Fatal(err)
	}
	const workers = 16
	results := make(chan string, workers)
	var wg sync.WaitGroup
	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			status, err := client.VerifyCaptcha(ctx, key, "wrong")
			if err != nil {
				results <- "error"
				return
			}
			results <- status
		}()
	}
	wg.Wait()
	close(results)
	counts := map[string]int{}
	for status := range results {
		counts[status]++
	}
	if counts["invalid_value"] != 2 || counts["no_tries"] != 1 || counts["not_found"] != workers-3 {
		t.Fatalf("unexpected concurrent attempt results: %+v", counts)
	}
}
