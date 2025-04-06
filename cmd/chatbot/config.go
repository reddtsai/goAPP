package main

import (
	"github.com/spf13/viper"
)

type ServiceConfig struct {
	GenAIVendor string `mapstructure:"GEN_AI_VENDOR"`

	OpenAIApiKey string `mapstructure:"OPENAI_API_KEY"`
	OpenAIModel  string `mapstructure:"OPENAI_MODEL"`

	DeepSeekApiKey string `mapstructure:"DEEPSEEK_API_KEY"`
	DeepSeekModel  string `mapstructure:"DEEPSEEK_MODEL"`

	AwsBedrockRegion          string `mapstructure:"AWS_BEDROCK_REGION"`
	AwsBedrockModel           string `mapstructure:"AWS_BEDROCK_MODEL"`
	AwsBedrockAccessKeyID     string `mapstructure:"AWS_BEDROCK_ACCESS_KEY_ID"`
	AwsBedrockSecretAccessKey string `mapstructure:"AWS_BEDROCK_SECRET_ACCESS_KEY"`
}

func LoadConfig(path string) (*ServiceConfig, error) {
	viper.AddConfigPath(path)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	cfg := &ServiceConfig{}
	if err := viper.Unmarshal(cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
