package output

import (
	"net/url"
	"regexp"
	"strings"

	"github.com/bitmovin/bitmovin-api-sdk-go/model"

	"github.com/bitmovin/bitmovin-api-sdk-go"
	"github.com/NYTimes/video-transcoding-api/config"

	"github.com/pkg/errors"
)

var s3Pattern = regexp.MustCompile(`^s3://`)

// New creates an output and returns an outputId and the folder path or an error
func New(destLoc string, api *bitmovin.BitmovinApi, cfg *config.Bitmovin) (outputID string, path string, err error) {
	if s3Pattern.MatchString(destLoc) {
		return s3(destLoc, api, cfg)
	}

	return "", "", errors.New("only s3 outputs are supported")
}

func s3(destLoc string, api *bitmovin.BitmovinApi, cfg *config.Bitmovin) (string, string, error) {
	bucket, folderPath, err := parseS3URL(destLoc)
	if err != nil {
		return "", "", err
	}

	output, err := api.Encoding.Outputs.S3.Create(model.S3Output{
		BucketName:  bucket,
		AccessKey:   cfg.AccessKeyID,
		SecretKey:   cfg.SecretAccessKey,
		CloudRegion: model.AwsCloudRegion(cfg.AWSStorageRegion),
	})
	if err != nil {
		return "", "", errors.Wrap(err, "creating s3 output")
	}

	return output.Id, folderPath, nil
}

func parseS3URL(s3URL string) (bucketName string, objectKey string, err error) {
	u, err := url.Parse(s3URL)
	if err != nil || u.Scheme != "s3" {
		return "", "", errors.Wrap(err, "parsing s3 url")
	}
	return u.Host, strings.TrimLeft(u.Path, "/"), nil
}
