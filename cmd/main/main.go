package main

import (
	"fmt"
	"os"
	"time"

	ginzap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config struct {
	ServerAddress  string `mapstructure:"SERVER_ADDRESS"`
	Port           string `mapstructure:"PORT"`
	ContextTimeout int    `mapstructure:"CONTEXT_TIMEOUT"`
	DBHost         string `mapstructure:"DB_HOST"`
	DBPort         string `mapstructure:"DB_PORT"`
	DBUser         string `mapstructure:"DB_USER"`
	DBPass         string `mapstructure:"DB_PASS"`
	DBName         string `mapstructure:"DB_NAME"`
	LogLevel       string `mapstructure:"LOG_LEVEL"`
}

func GetConfig() (Config, error) {
	cfg := Config{}

	viper.AddConfigPath("../..") // for local run
	viper.AddConfigPath("/app")  // to work with Dockerfile setup
	viper.SetConfigName("app")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return cfg, fmt.Errorf("can't find the file .env : %w", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return cfg, fmt.Errorf("environment can't be loaded: %w", err)
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.ServerAddress == "" {
		cfg.ServerAddress = "0.0.0.0"
	}

	return cfg, nil
}

func initLogger(config Config) (logger *zap.Logger, err error) {
	var zapConfig = zap.NewProductionConfig()
	if config.LogLevel != "" {
		var level zapcore.Level
		if err = level.Set(config.LogLevel); err != nil {
			return nil, err
		}
		zapConfig.Level.SetLevel(level)
	}

	logger, err = zapConfig.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

func main() {

	fmt.Println("starting Url shortener server")

	config, err := GetConfig()
	if err != nil {
		fmt.Printf("impossible to load config at startup: %s\n", err)
		os.Exit(1)
	}

	logger, err := initLogger(config)
	if err != nil {
		fmt.Printf("impossible to init logger: %s\n", err)
		os.Exit(1)
	}

	router := gin.New()
	router.Use(ginzap.Ginzap(logger, time.RFC3339, true))
	router.Use(ginzap.RecoveryWithZap(logger, true))
	router.POST("/url", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	logger.Fatal("url shortener service crashed", zap.Error(router.Run(fmt.Sprintf("%s:%s", config.ServerAddress, config.Port))))
}
