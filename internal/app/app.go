package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	handlerCaptcha "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/handler/captcha"
	handlerPow "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/handler/pow"
	handlerVerify "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/handler/verify"
	processCaptcha "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/captcha"
	tasksCaptcha "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/captcha/task"
	processPow "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/pow"
	tasksPow "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/pow/task"
	processVerify "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/verify"
	tasksVerify "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/process/verify/task"
	"github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/registry"
	serviceRedis "github.com/AdrianJanczenia/adrianjanczenia.dev_captcha-service/internal/service/redis"
)

type App struct {
	httpServer  *http.Server
	redisClient serviceRedis.Client
}

func Build(cfg *registry.Config) (*App, error) {
	if cfg == nil || cfg.Infrastructure.Retry.MaxAttempts < 1 || cfg.Infrastructure.Retry.DelaySeconds < 0 || cfg.Server.HTTPPort == "" || cfg.Redis.URL == "" || cfg.Security.HmacSecret == "" || cfg.Security.TtlMinutes < 1 || cfg.Captcha.TtlMinutes < 1 || cfg.Captcha.MaxTries < 1 || cfg.Security.Difficulty < 1 || cfg.Security.Difficulty > 64 {
		return nil, fmt.Errorf("invalid captcha-service configuration")
	}
	maxRetries := cfg.Infrastructure.Retry.MaxAttempts
	retryDelay := cfg.Infrastructure.Retry.DelaySeconds
	var err error

	var redisClient serviceRedis.Client
	for i := range maxRetries {
		redisClient, err = serviceRedis.NewClient(cfg.Redis.URL)
		if err == nil {
			pingCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			err = redisClient.Ping(pingCtx)
			cancel()
			if err == nil {
				log.Println("INFO: successfully connected to Redis")
				break
			}
			_ = redisClient.Close()
		}
		if i+1 < maxRetries {
			log.Printf("INFO: could not connect to Redis, retrying in %v... (%d/%d)", retryDelay, i+1, maxRetries)
			time.Sleep(retryDelay)
		}
	}
	if err != nil {
		return nil, err
	}

	createSignedSeedTask := tasksPow.NewCreateSignedSeedTask(cfg.Security.HmacSecret)
	powProcess := processPow.NewProcess(createSignedSeedTask)
	powHandler := handlerPow.NewHandler(powProcess)

	validateSignatureTask := tasksCaptcha.NewValidateSignatureTask(cfg.Security.HmacSecret)
	checkSeedTimestampTask := tasksCaptcha.NewCheckSeedTimestampTask(cfg.Security.TtlMinutes)
	verifyPowTask := tasksCaptcha.NewVerifyPowTask(cfg.Security.Difficulty)
	saveUsedSeedTask := tasksCaptcha.NewSaveUsedSeedTask(redisClient, cfg.Security.TtlMinutes)
	generateCaptchaTask := tasksCaptcha.NewGenerateCaptchaTask()
	saveCaptchaTask := tasksCaptcha.NewSaveCaptchaTask(redisClient, cfg.Captcha.TtlMinutes, cfg.Captcha.MaxTries)
	captchaProcess := processCaptcha.NewProcess(validateSignatureTask, checkSeedTimestampTask, verifyPowTask, saveUsedSeedTask, generateCaptchaTask, saveCaptchaTask)
	captchaHandler := handlerCaptcha.NewHandler(captchaProcess)

	verifyCaptchaTask := tasksVerify.NewVerifyCaptchaTask(redisClient)
	verifyProcess := processVerify.NewProcess(verifyCaptchaTask)
	verifyHandler := handlerVerify.NewHandler(verifyProcess)

	mux := http.NewServeMux()
	mux.HandleFunc("/pow", powHandler.Handle)
	mux.HandleFunc("/captcha", captchaHandler.Handle)
	mux.HandleFunc("/verify", verifyHandler.Handle)

	httpServer := &http.Server{
		Addr: ":" + cfg.Server.HTTPPort,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
			mux.ServeHTTP(w, r)
		}),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		IdleTimeout:       60 * time.Second,
		MaxHeaderBytes:    32 << 10,
	}

	return &App{
		httpServer:  httpServer,
		redisClient: redisClient,
	}, nil
}

func (a *App) RunHTTP() error {
	log.Printf("INFO: HTTP server listening on %s", a.httpServer.Addr)
	return a.httpServer.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) {
	log.Println("INFO: shutting down server...")
	_ = a.httpServer.Shutdown(ctx)
	_ = a.redisClient.Close()
}
