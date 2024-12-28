package aws

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/png"
	"math/rand"
	"mime/multipart"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/disintegration/imaging"
)

var imgMeta = map[ImgType]imgMetaStruct{
	ImgTypeFood: {
		bucket:     func() string { return "dev-food-recommendation" },
		domain:     func() string { return "dev-food-recommendation.s3.ap-northeast-2.amazonaws.com" },
		path:       "images",
		width:      512,
		height:     512,
		expireTime: 2 * time.Hour,
	},
	ImgTypeCategory: {
		bucket:     func() string { return "dev-food-recommendation" },
		domain:     func() string { return "dev-food-recommendation.s3.ap-northeast-2.amazonaws.com" },
		path:       "category",
		width:      62,
		height:     62,
		expireTime: 2 * time.Hour,
	},
	ImgTypeProfile: {
		bucket:     func() string { return "dev-food-recommendation" },
		domain:     func() string { return "dev-food-recommendation.s3.ap-northeast-2.amazonaws.com" },
		path:       "profiles",
		width:      512,
		height:     512,
		expireTime: 2 * time.Hour,
	},
}

func ImageUpload(ctx context.Context, file *multipart.FileHeader, filename string, imgType ImgType) error {
	meta, ok := imgMeta[imgType]
	if !ok {
		return fmt.Errorf("not available meta info for imgType - %v", imgType)
	}
	bucket := meta.bucket()

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("fail to open file - %v", err)
	}
	defer src.Close()

	img, err := imaging.Decode(src)
	if err != nil {
		return fmt.Errorf("fail to load image - %v", err)
	}

	if meta.width < 1 || meta.height < 1 {
		if (meta.width < 1 && meta.height < img.Bounds().Size().Y) ||
			(meta.height < 1 && meta.width < img.Bounds().Size().X) {
			img = imaging.Resize(img, meta.width, meta.height, imaging.Lanczos)
		}
	} else {
		img = imaging.Fill(img, meta.width, meta.height, imaging.Center, imaging.Lanczos)
	}

	buf := new(bytes.Buffer)
	if err := imaging.Encode(buf, img, imaging.PNG, imaging.PNGCompressionLevel(png.BestCompression)); err != nil {
		return fmt.Errorf("fail to encode png image - %v", err)
	}

	_, err = awsClientS3Uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fmt.Sprintf("%s/%s", meta.path, filename)),
		Body:        buf,
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		return fmt.Errorf("fail to upload image to s3 - bucket:%s / key:%s/%s", bucket, meta.path, filename)
	}
	return nil
}

func FoodImageUpload(ctx context.Context, file *multipart.FileHeader, filename string, imgType ImgType) error {

	meta, ok := imgMeta[imgType]
	if !ok {
		return fmt.Errorf("not available meta info for imgType - %v", imgType)
	}
	bucket := meta.bucket()

	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("fail to open file - %v", err)
	}
	defer src.Close()

	img, err := imaging.Decode(src)
	if err != nil {
		return fmt.Errorf("fail to load image - %v", err)
	}

	// Step 3: 이미지 크기 가져오기
	imgWidth := img.Bounds().Dx()
	imgHeight := img.Bounds().Dy()

	// Step 4: 1.5배 확대된 크기 계산
	scaleFactor := 2.0
	zoomedWidth := int(float64(meta.width) * scaleFactor)
	zoomedHeight := int(float64(meta.height) * scaleFactor)

	// 중심 기준 자르기 영역 계산
	startX := (imgWidth - zoomedWidth) / 2
	startY := (imgHeight - zoomedHeight) / 2

	// 자르기 범위를 초과하지 않도록 보정
	if startX < 0 {
		startX = 0
		zoomedWidth = imgWidth
	}
	if startY < 0 {
		startY = 0
		zoomedHeight = imgHeight
	}
	// 확대된 영역으로 자르기
	croppedImg := imaging.Crop(img, image.Rect(startX, startY, startX+zoomedWidth, startY+zoomedHeight))

	// Step 5: 자른 이미지를 meta 크기로 다시 조정 (필요시)
	finalImg := imaging.Resize(croppedImg, meta.width, meta.height, imaging.Lanczos)

	buf := new(bytes.Buffer)
	if err := imaging.Encode(buf, finalImg, imaging.PNG, imaging.PNGCompressionLevel(png.BestCompression)); err != nil {
		return fmt.Errorf("fail to encode png image - %v", err)
	}

	_, err = awsClientS3Uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucket),
		Key:         aws.String(fmt.Sprintf("%s/%s", meta.path, filename)),
		Body:        buf,
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		return fmt.Errorf("fail to upload image to s3 - bucket:%s / key:%s/%s", bucket, meta.path, filename)
	}
	return nil
}
func ImageGetSignedURL(ctx context.Context, fileName string, imgType ImgType) (string, error) {
	meta, ok := imgMeta[imgType]
	if !ok {
		return "", fmt.Errorf("not available meta info for imgType - %v", imgType)
	}
	presignClient := s3.NewPresignClient(awsClientS3)

	key := fmt.Sprintf("%s/%s", meta.path, fileName)
	presignParams := &s3.GetObjectInput{
		Bucket: aws.String(meta.bucket()),
		Key:    aws.String(key),
	}

	presignResult, err := presignClient.PresignGetObject(ctx, presignParams, s3.WithPresignExpires(meta.expireTime))
	if err != nil {
		return "", err
	}
	return presignResult.URL, nil
}

func ImageDelete(ctx context.Context, fileName string, imgType ImgType) error {
	meta, ok := imgMeta[imgType]
	if !ok {
		return fmt.Errorf("not available meta info for imgType - %v", imgType)
	}

	bucket := meta.bucket()
	key := fmt.Sprintf("%s/%s", meta.path, fileName)

	if _, err := awsClientS3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}); err != nil {
		return fmt.Errorf("fail to delete image from s3 - bucket:%s, key:%s", bucket, key)
	}

	return nil
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func FileNameGenerateRandom() string {
	rand.Seed(time.Now().UnixNano())
	b := make([]rune, 32)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
