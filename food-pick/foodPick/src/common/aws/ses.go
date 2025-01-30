package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

type emailType string

const (
	emailTypePassword         = emailType("password")
	emailTypeAuth             = emailType("authCode")
	emailTypeReport           = emailType("report")
	emailTypeSignup           = emailType("signup")
	emailTypeFoodNameReport   = emailType("foodNameReport")
	emailTypeFoodUploadReport = emailType("foodUploadReport")
)

type sesMailData struct {
	email        []string
	mailType     emailType
	failCount    uint8
	templateData string
	templateName string
}

type ReqReportSES struct {
	UserID string
	Reason string
}


func EmailSendFoodInfoEmptyReport(imageEmpty, nutrientEmpty []string) {
	templateDataMap := map[string][]string{
		"imageMissingFoods": imageEmpty,
		"nutrientMissingFoods":  nutrientEmpty,
	}
	templateDataJson, err := json.Marshal(templateDataMap)
	if err != nil {
		fmt.Println("Error marshaling template data:", err)
		return
	}

	emailSend([]string{"pkjhj485@gmail.com", "dtw7225@naver.com", "ohhyejin1213@naver.com"}, emailTypeFoodNameReport, string(templateDataJson), "foodInfoEmptyReport")
}
func EmailSendFoodUploadReport(successFoodList, failedFoodList []string) {
	templateDataMap := map[string]string{
		"successFoodList": strings.Join(successFoodList, ", "),
		"failedFoodList":  strings.Join(failedFoodList, ", "),
	}
	templateDataJson, err := json.Marshal(templateDataMap)
	if err != nil {
		fmt.Println("Error marshaling template data:", err)
		return
	}

	emailSend([]string{"pkjhj485@gmail.com", "dtw7225@naver.com", "ohhyejin1213@naver.com"}, emailTypeFoodNameReport, string(templateDataJson), "foodUploadReport")
}
func EmailSendFoodNameReport(foodNames []string) {
	currentDate := time.Now().Format("01-02")
	templateDataMap := map[string]string{
		"currentDate": currentDate,
		"foodList":    strings.Join(foodNames, ", "),
	}
	templateDataJson, err := json.Marshal(templateDataMap)
	if err != nil {
		fmt.Println("Error marshaling template data:", err)
		return
	}

	emailSend([]string{"pkjhj485@gmail.com", "dtw7225@naver.com", "ohhyejin1213@naver.com"}, emailTypeFoodNameReport, string(templateDataJson), "foodNameReport")
}
func EmailSendAuthCode(email string, validateCode string) {
	templateDataMap := map[string]string{
		"code": validateCode,
	}
	templateDataJson, err := json.Marshal(templateDataMap)
	if err != nil {
		fmt.Println("Error marshaling template data:", err)
		return
	}

	emailSend([]string{email}, emailTypeAuth, string(templateDataJson), "foodAuth")
}
func EmailSendPassword(email string, validateCode string) {
	templateDataMap := map[string]string{
		"randomValue": validateCode,
	}
	templateDataJson, err := json.Marshal(templateDataMap)
	if err != nil {
		fmt.Println("Error marshaling template data:", err)
		return
	}

	emailSend([]string{email}, emailTypePassword, string(templateDataJson), "password")
}
func EmailSendSignup(email string, validateCode string) {
	templateDataMap := map[string]string{
		"randomValue": validateCode,
	}
	templateDataJson, err := json.Marshal(templateDataMap)
	if err != nil {
		fmt.Println("Error marshaling template data:", err)
		return
	}

	emailSend([]string{email}, emailTypePassword, string(templateDataJson), "signup")
}

func EmailSendReport(email []string, req *ReqReportSES) {
	templateDataMap := map[string]string{
		"userID": req.UserID,
		"reason": req.Reason,
	}
	templateDataJson, err := json.Marshal(templateDataMap)
	if err != nil {
		fmt.Println("Error marshaling template data:", err)
		return
	}
	emailSend(email, emailTypeReport, string(templateDataJson), "foodReport")
}

func emailSend(email []string, mailType emailType, templateDataJson, templateName string) {

	mailData := sesMailData{
		email:        email,
		mailType:     mailType,
		failCount:    0,
		templateData: templateDataJson,
		templateName: templateName,
	}
	select {
	case sesMailReqChan <- mailData:
	default:
		<-sesMailReqChan
		sesMailReqChan <- mailData
	}
}

var sesMailReqChan chan sesMailData

func InitAwsSes() error {

	sesMailReqChan = make(chan sesMailData, 100)
	go func() {
		for {
			mailReq := <-sesMailReqChan
			_, err := awsClientSes.SendEmail(context.TODO(), &sesv2.SendEmailInput{
				Content: &types.EmailContent{
					Template: &types.Template{
						TemplateData: aws.String(mailReq.templateData),
						TemplateName: aws.String(mailReq.templateName),
					},
				},
				Destination: &types.Destination{
					ToAddresses: mailReq.email,
				},
				EmailTags: []types.MessageTag{{
					Name:  aws.String("type"),
					Value: aws.String(string(mailReq.mailType)),
				}},
				FromEmailAddress: aws.String("root@jokertrickster.com"),
			})
			if err != nil {
				if mailReq.failCount < 3 {
					fmt.Println("Error sending email:", err)
					mailReq.failCount += 1
					sesMailReqChan <- mailReq
				}
			}
		}
	}()
	return nil
}
