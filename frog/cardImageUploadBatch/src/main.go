package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strings"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

var failedList []string
var successList []string
var globalErr error

func handler(ctx context.Context, request events.CloudWatchEvent) error {
	// 서비스 계정 JSON 키 파일 경로

	serviceKey, err := AwsGetParam("food_image_upload_service_key")
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 서비스 계정 JSON 키를 byte 배열로 변환합니다.
	credentials := []byte(serviceKey)

	// Google Drive 클라이언트 생성
	config, err := google.JWTConfigFromJSON(credentials, drive.DriveScope)
	if err != nil {
		log.Fatalf("Unable to parse service account JSON: %v", err)
	}

	srv, err := drive.NewService(ctx, option.WithHTTPClient(config.Client(ctx)))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	// 원하는 Google Drive 폴더의 ID를 설정합니다.
	folderID := "1JxI_Btaeu19vMUiOg9WwgajUv47DjrGr"
	queryDate := time.Now().AddDate(0, 0, -60).In(time.FixedZone("KST", 9*60*60)).Format("2006-01-02")
	fmt.Println(queryDate)
	// query := fmt.Sprintf("'%s' in parents and mimeType contains 'image/'", folderID)
	// query := fmt.Sprintf("'%s' in parents and mimeType contains 'image/' and modifiedTime >= '%sT00:00:00Z'", folderID, queryDate)
	// 이미지 파일 검색 쿼리
	// r, err := srv.Files.List().Q(query).Fields("files(id, name, mimeType)").Do()
	// if err != nil {
	// 	log.Fatalf("Unable to retrieve files: %v", err)
	// }

	files, err := getAllFiles(srv, folderID)
	if err != nil {
		log.Fatalf("Error retrieving files: %v", err)
	}

	fmt.Printf("Total files: %d\n", len(files))

	successList = make([]string, 0)
	failedList = make([]string, 0)

	for _, file := range files {
		// 파일이 이미지인지 확인
		if strings.HasPrefix(file.MimeType, "image/") {
			fmt.Printf("Downloading image file: %s (%s)\n", file.Name, file.Id)
			result, err := downloadFile(srv, file.Id, file.Name)
			if err != nil {
				log.Fatalf("Error downloading file: %v", err)
			}
			if result {
				successList = append(successList, file.Name)
			} else {
				failedList = append(failedList, file.Name)
			}
		}
	}

	// 업로드 성공 및 실패 목록을 POST 요청으로 전송
	requestBody, err := json.Marshal(map[string]interface{}{
		"failedList":  failedList,
		"successList": successList,
	})
	if err != nil {
		log.Fatalf("JSON 변환 오류: %v", err)
	}

	req, err := http.NewRequest("POST", "https://dev-frog-api.jokertrickster.com/v0.1/games/images/upload-report", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Fatalf("POST 요청 생성 오류: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("POST 요청 오류: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("요청 성공!")
	} else {
		fmt.Printf("요청 실패: 상태 코드 %d\n", resp.StatusCode)
	}
	return nil
}

func main() {
	if err := InitAws(); err != nil {
		log.Fatalf("AWS 초기화 실패: %v", err)
	}
	lambda.Start(handler)
}

func downloadFile(srv *drive.Service, fileID string, fileName string) (bool, error) {
	resp, err := srv.Files.Get(fileID).Download()
	if err != nil {
		return false, fmt.Errorf("unable to download file: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("error downloading file: %v", resp.Status)
	}

	url := "https://dev-frog-api.jokertrickster.com/v0.1/game/cards/image"
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("image", fileName)
	if err != nil {
		return false, fmt.Errorf("unable to create form file: %v", err)
	}

	_, err = io.Copy(part, resp.Body)
	if err != nil {
		return false, fmt.Errorf("unable to copy file content: %v", err)
	}

	err = writer.Close()
	if err != nil {
		return false, fmt.Errorf("unable to close writer: %v", err)
	}

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return false, fmt.Errorf("unable to create request: %v", err)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{Timeout: 30 * time.Second}
	respUpload, err := client.Do(req)
	if err != nil {
		return false, fmt.Errorf("unable to upload file: %v", err)
	}
	defer respUpload.Body.Close()

	if respUpload.StatusCode != http.StatusOK && respUpload.StatusCode != http.StatusCreated {
		if respUpload.StatusCode == http.StatusBadRequest {
			return false, nil
		}
		return false, fmt.Errorf("error uploading file: %v", respUpload.Status)
	}

	return true, nil
}
func getAllFiles(srv *drive.Service, folderID string) ([]*drive.File, error) {
	var allFiles []*drive.File
	query := fmt.Sprintf("'%s' in parents and mimeType contains 'image/'", folderID)
	pageToken := ""

	for {
		request := srv.Files.List().Q(query).Fields("nextPageToken, files(id, name, mimeType)")
		if pageToken != "" {
			request = request.PageToken(pageToken)
		}

		r, err := request.Do()
		if err != nil {
			return nil, fmt.Errorf("unable to retrieve files: %v", err)
		}

		allFiles = append(allFiles, r.Files...)

		// 다음 페이지가 없으면 종료
		if r.NextPageToken == "" {
			break
		}

		// 다음 페이지 토큰 설정
		pageToken = r.NextPageToken
	}

	return allFiles, nil
}
