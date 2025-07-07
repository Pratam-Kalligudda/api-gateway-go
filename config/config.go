package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Env struct {
	PORT             string
	SECRET           string
	USER_ENDPOINT    string
	PRODUCT_ENDPOINT string
	ORDER_ENDPOINT   string
}

type ServiceConfig struct {
	ContextPath string
	TargetUrl   string
}

type Config struct {
	Services map[string]ServiceConfig
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("error while initializing env : %v", err.Error())
	}
}

func LoadEnv() Env {
	port := os.Getenv("PORT")
	if len(port) <= 0 {
		log.Fatal("error while initializing env port")
	}

	secret := os.Getenv("SECRET")
	if len(secret) <= 0 {
		log.Fatal("error while initializing env secret")
	}
	userEP := os.Getenv("USER_ENDPOINT")
	if len(userEP) <= 0 {
		log.Fatal("error while initializing env userEP")
	}
	productEP := os.Getenv("PRODUCT_ENDPOINT")
	if len(productEP) <= 0 {
		log.Fatal("error while initializing env productEP")
	}
	orderEP := os.Getenv("ORDER_ENDPOINT")
	if len(orderEP) <= 0 {
		log.Fatal("error while initializing env orderEP")
	}
	return Env{
		PORT:             port,
		SECRET:           secret,
		USER_ENDPOINT:    userEP,
		PRODUCT_ENDPOINT: productEP,
		ORDER_ENDPOINT:   orderEP,
	}
}

func LoadConfig() *Config {
	env := LoadEnv()
	config := new(Config)
	config.Services["user-service"] = ServiceConfig{
		ContextPath: "/api/users",
		TargetUrl:   env.USER_ENDPOINT,
	}
	config.Services["product-service"] = ServiceConfig{
		ContextPath: "/api/products",
		TargetUrl:   env.PRODUCT_ENDPOINT,
	}
	config.Services["order-service"] = ServiceConfig{
		ContextPath: "/api/orders",
		TargetUrl:   env.ORDER_ENDPOINT,
	}
	return config
}
