package objectstorage

import (
	"context"
	"fmt"

	sdk "github.com/MagaluCloud/mgc-sdk-go/client"
	objSdk "github.com/MagaluCloud/mgc-sdk-go/objectstorage"
	"github.com/magaluCloud/mgccli/cmd/common/auth"
	cmdutils "github.com/magaluCloud/mgccli/cmd_utils"
)

type ObjectStorage interface {
	GetBucketService() objSdk.BucketService
}

type objectStorage struct {
	bucketService objSdk.BucketService
}

func NewObjectStorage(ctx context.Context) (ObjectStorage, error) {
	sdkCoreConfig, ok := ctx.Value(cmdutils.CTX_SDK_KEY).(sdk.CoreClient)
	if !ok {
		return nil, fmt.Errorf("sdk core client not found in context")
	}

	authCtx, ok := ctx.Value(cmdutils.CTX_AUTH_KEY).(auth.Auth)
	if !ok {
		return nil, fmt.Errorf("auth information not found in context")
	}

	accessKeyID := authCtx.GetAccessKeyID()
	secretAccessKey := authCtx.GetSecretAccessKey()

	if accessKeyID == "" || secretAccessKey == "" {
		return nil, fmt.Errorf("object storage credentials not found. Make sure you set an API key")
	}

	service, err := objSdk.New(&sdkCoreConfig, accessKeyID, secretAccessKey)
	if err != nil || service == nil {
		return nil, fmt.Errorf("failed to initialize object storage service: %w", err)
	}

	return &objectStorage{
		bucketService: service.Buckets(),
	}, nil
}

func (o *objectStorage) GetBucketService() objSdk.BucketService {
	return o.bucketService
}
