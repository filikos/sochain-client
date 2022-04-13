package util

import (
	"errors"
	"os"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const (
	EnvironmentDev  = "development"
	EnvironmentStag = "staging"
	EnvironmentProd = "production"
)

func GetParamNetwork(ctx *gin.Context, key string) (string, error) {

	v := ctx.Param(key)
	if v == "" {
		return "", errors.New("path param: Network id is missing")
	}

	if v != "btc" && v != "ltc" && v != "doge" {
		return "", errors.New("path param: Network 'id' can only be 'btc', 'eth' or 'ltc'")
	}

	return v, nil
}

func NewLogger(ciEnv string) (*zap.Logger, error) {
	if ciEnv == EnvironmentDev || ciEnv == EnvironmentStag || ciEnv == EnvironmentProd {
		return zap.NewDevelopment()
	}

	return zap.NewProduction()
}

func GetEnv(name, fallback string) string {
	if v, ok := os.LookupEnv(name); ok {
		return v
	}
	return fallback
}
