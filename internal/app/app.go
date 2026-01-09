package app

import (
	"context"
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
	httpServer *http.Server
}

func Build(cfg *registry.Config) (*App, error) {
	maxRetries := cfg.Infrastructure.Retry.MaxAttempts
	retryDelay := cfg.Infrastructure.Retry.DelaySeconds
	var err error

	var redisClient serviceRedis.Client
	for i := 0; i < maxRetries; i++ {
		redisClient, err = serviceRedis.NewClient(cfg.Redis.URL)
		if err == nil {
			if err = redisClient.Ping(context.Background()); err == nil {
				log.Println("INFO: successfully connected to Redis")
				break
			}
		}
		log.Printf("INFO: could not connect to Redis, retrying in %v... (%d/%d)", retryDelay, i+1, maxRetries)
		time.Sleep(retryDelay)
	}
	if err != nil {
		return nil, err
	}

	createSignedSeedTask := tasksPow.NewCreateSignedSeedTask(cfg.Security.HmacSecret)
	powProcess := processPow.NewProcess(createSignedSeedTask)
	powHandler := handlerPow.NewHandler(powProcess)

	validateSignatureTask := tasksCaptcha.NewValidateSignatureTask(cfg.Security.HmacSecret)
	checkSeedTimestampTask := tasksCaptcha.NewCheckSeedTimestampTask(cfg.Security.TtlMinutes)
	validateUsedSeedTask := tasksCaptcha.NewValidateUsedSeedTask(redisClient)
	verifyPowTask := tasksCaptcha.NewVerifyPowTask(cfg.Security.Difficulty)
	marksSeedUsedTask := tasksCaptcha.NewMarkSeedUsedTask(redisClient, cfg.Captcha.TtlMinutes)
	generateCaptchaTask := tasksCaptcha.NewGenerateCaptchaTask()
	saveCaptchaTask := tasksCaptcha.NewSaveCaptchaTask(redisClient, cfg.Captcha.TtlMinutes, cfg.Captcha.MaxTries)
	captchaProcess := processCaptcha.NewProcess(validateSignatureTask, checkSeedTimestampTask, validateUsedSeedTask, verifyPowTask, marksSeedUsedTask, generateCaptchaTask, saveCaptchaTask)
	captchaHandler := handlerCaptcha.NewHandler(captchaProcess)

	fetchCaptchaTask := tasksVerify.NewFetchCaptchaTask(redisClient)
	validateCaptchaTask := tasksVerify.NewValidateCaptchaTask(redisClient, cfg.Captcha.TtlMinutes)
	verifyProcess := processVerify.NewProcess(fetchCaptchaTask, validateCaptchaTask)
	verifyHandler := handlerVerify.NewHandler(verifyProcess)

	mux := http.NewServeMux()
	mux.HandleFunc("/pow", powHandler.Handle)
	mux.HandleFunc("/captcha", captchaHandler.Handle)
	mux.HandleFunc("/verify", verifyHandler.Handle)

	httpServer := &http.Server{
		Addr: ":" + cfg.Server.HTTPPort,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.URL.Path = strings.TrimSuffix(r.URL.Path, "/")
			mux.ServeHTTP(w, r)
		}),
	}

	return &App{
		httpServer: httpServer,
	}, nil
}

func (a *App) RunHTTP() error {
	log.Printf("INFO: HTTP server listening on %s", a.httpServer.Addr)
	return a.httpServer.ListenAndServe()
}

func (a *App) Shutdown(ctx context.Context) {
	log.Println("INFO: shutting down server...")
	_ = a.httpServer.Shutdown(ctx)
}
