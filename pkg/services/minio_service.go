package services

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/minio/minio-go/v6"
	"github.com/saas/hostgolang/pkg/types"
	"io"
)

const appsBinaryBucket = "applications"
const appsConfigBucket = "application-configs"

var ErrEncodeDecode = errors.New("failed to decode/encode config application config")
var ErrDataNotFound = errors.New("record not found")
var ErrPersistData = errors.New("failed to persist application data")

type StorageClient interface {
	Put(key string, data []byte) error
	Get(key string) (io.Reader, error)
	PutReleaseConfig(key string, cfg *types.ReleaseConfig) error
	GetReleaseConfig(key string) (*types.ReleaseConfig, error)
}

type minioStorageClient struct {
	minioClient *minio.Client
}

func NewMinioStorageClient(client *minio.Client) (StorageClient, error) {
	m := &minioStorageClient{}
	m.minioClient = client
	if err := m.minioClient.MakeBucket(appsBinaryBucket, "us-east-1"); err != nil {
		if exists, err := m.minioClient.BucketExists(appsBinaryBucket); err != nil || !exists {
			return nil, err
		}
	}
	if err := m.minioClient.MakeBucket(appsConfigBucket, "us-east-1"); err != nil {
		if exists, err := m.minioClient.BucketExists(appsConfigBucket); err != nil || !exists {
			return nil, err
		}
	}
	return m, nil
}

func (m *minioStorageClient) Put(key string, data []byte) error {
	l := int64(len(data))
	if _, err := m.minioClient.PutObject(appsBinaryBucket, key, bytes.NewBuffer(data), l, minio.PutObjectOptions{}); err != nil {
		return ErrPersistData
	}
	return nil
}

func (m *minioStorageClient) Get(key string) (io.Reader, error) {
	if val, err := m.minioClient.GetObject(appsBinaryBucket, key, minio.GetObjectOptions{}); err == nil {
		return val, nil
	}
	return nil, ErrDataNotFound
}

func (m *minioStorageClient) PutReleaseConfig(key string, cfg *types.ReleaseConfig) error {
	data, err := json.Marshal(cfg)
	if err != nil {
		return ErrEncodeDecode
	}
	_, err = m.minioClient.PutObject(appsConfigBucket, key, bytes.NewBuffer(data), int64(len(data)), minio.PutObjectOptions{ContentType: "application/json"})
	if err != nil {
		return ErrPersistData
	}
	return nil
}

func (m *minioStorageClient) GetReleaseConfig(key string) (*types.ReleaseConfig, error) {
	data, err := m.minioClient.GetObject(appsConfigBucket, key, minio.GetObjectOptions{})
	if err != nil {
		return nil, ErrDataNotFound
	}
	cfg := &types.ReleaseConfig{}
	if err := json.NewDecoder(data).Decode(cfg); err != nil {
		return nil, ErrDataNotFound
	}
	return cfg, nil
}
