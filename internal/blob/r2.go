// SPDX-FileCopyrightText: 2024 PK Lab AG <contact@pklab.io>
// SPDX-License-Identifier: MIT

package blob

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	serviceS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	s3Types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/rs/zerolog/log"
)

const (
	headFile = "head.json"
)

type R2Config struct {
	BucketName      string
	AccountId       string
	AccessKeyId     string
	AccessKeySecret string
}

func (cfg R2Config) R2Url() string {
	return fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountId)
}

type R2 struct {
	cfg    R2Config
	client *serviceS3.Client
}

func NewR2(cfg R2Config) R2 {
	return R2{
		cfg: cfg,
	}
}

func (r2 *R2) Init(ctx context.Context) error {
	r2Resolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL: r2.cfg.R2Url(),
		}, nil
	})
	cfg, err := awsConfig.LoadDefaultConfig(ctx,
		awsConfig.WithEndpointResolverWithOptions(r2Resolver),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(r2.cfg.AccessKeyId, r2.cfg.AccessKeySecret, "")),
		awsConfig.WithRegion("auto"),
		awsConfig.WithRetryMode(aws.RetryModeAdaptive),
	)
	if err != nil {
		return err
	}

	r2.client = serviceS3.NewFromConfig(cfg)
	return nil
}

func (r2 R2) Save(ctx context.Context, blob Blob) error {
	_, err := r2.client.PutObject(ctx, &serviceS3.PutObjectInput{
		Bucket:        aws.String(r2.cfg.BucketName),
		Key:           aws.String(blob.String()),
		Body:          bytes.NewBuffer(blob.Data),
		ContentType:   aws.String(blob.ContentType()),
		ContentLength: aws.Int64(int64(len(blob.Data))),
	})
	return err
}

func (r2 R2) SaveBulk(ctx context.Context, blobs []Blob) error {
	if len(blobs) == 0 {
		return nil
	}

	var wg sync.WaitGroup
	var e error
	for i := range blobs {
		wg.Add(1)
		go func(blob Blob, wg *sync.WaitGroup) {
			defer wg.Done()

			log.Info().Str("blob", blob.String()).Msg("saving blob...")
			if err := r2.Save(ctx, blob); err != nil {
				if e == nil {
					e = err
				}
			}
		}(blobs[i], &wg)
	}
	wg.Wait()

	return e
}

func (r2 R2) Head(ctx context.Context) (uint64, error) {
	output, err := r2.client.GetObject(ctx, &serviceS3.GetObjectInput{
		Bucket: aws.String(r2.cfg.BucketName),
		Key:    aws.String(headFile),
	})
	if err != nil {
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			switch apiErr.(type) {
			case *s3Types.NoSuchKey, *s3Types.NotFound:
				return 0, nil
			default:
				return 0, err
			}
		}
		return 0, err
	}

	var head uint64
	err = json.NewDecoder(output.Body).Decode(&head)
	return head, err
}

func (r2 R2) UpdateHead(ctx context.Context, head uint64) error {
	data, err := json.Marshal(head)
	if err != nil {
		return err
	}
	_, err = r2.client.PutObject(ctx, &serviceS3.PutObjectInput{
		Bucket: aws.String(r2.cfg.BucketName),
		Key:    aws.String(headFile),
		Body:   bytes.NewBuffer(data),
	})
	return err
}
