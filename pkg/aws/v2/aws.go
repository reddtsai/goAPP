package v2

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

func LoadConfig(ctx context.Context, region string) (aws.Config, error) {
	var opts []func(*config.LoadOptions) error
	opts = append(opts, config.WithRegion(region))

	return config.LoadDefaultConfig(ctx, opts...)
}

func LoadConfigWithAKSK(ctx context.Context, region, ak, sk string) (aws.Config, error) {
	var opts []func(*config.LoadOptions) error
	opts = append(opts, config.WithRegion(region))
	opts = append(opts, config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(ak, sk, "")))

	return config.LoadDefaultConfig(ctx, opts...)
}
