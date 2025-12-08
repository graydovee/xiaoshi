package config

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"time"
)

var (
	ErrMissingCredentials = errors.New("missing S3 credentials in environment variables")
)

// LoadConfigFromS3 从S3/MinIO加载配置文件
// bucket: 存储桶名称
// path:   文件路径（如 "configs/app/prod.yaml"）
// 返回: 文件内容字节数组或错误
func LoadConfigFromS3() ([]byte, error) {

	path := os.Getenv("S3_PATH")
	if path == "" {
		return nil, errors.New("S3_PATH is not set")
	}

	// 1. 从环境变量获取认证信息
	accessKey := os.Getenv("S3_ACCESS_KEY")
	secretKey := os.Getenv("S3_SECRET")
	endpoint := os.Getenv("S3_ENDPOINT")
	bucket := os.Getenv("S3_BUCKET")

	if accessKey == "" || secretKey == "" {
		return nil, ErrMissingCredentials
	}

	customTransport := http.DefaultTransport.(*http.Transport).Clone()
	customTransport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}

	// 2. 创建AWS配置（兼容MinIO）
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("auto"),
		config.WithBaseEndpoint(endpoint),
		config.WithHTTPClient(&http.Client{Transport: customTransport}),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKey,
			secretKey,
			"",
		)),
	)
	if err != nil {
		return nil, err
	}

	// 4. 创建S3客户端
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.UsePathStyle = true
	})

	// 5. 设置超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err := client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 7. 读取文件内容
	return io.ReadAll(resp.Body)
}

func LoadConfigFromPath() (*Config, error) {
	path := os.Getenv("CONFIG_FILE_PATH")
	if path == "" {
		return nil, errors.New("CONFIG_FILE_PATH is not set")
	}

	configData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err = yaml.Unmarshal(configData, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
