package waitlist

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/dioncodes/go-waitlist/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/resend/resend-go/v2"
	"gorm.io/datatypes"
)

func RegisterRoutes(rg *gin.RouterGroup) {
	auth := rg.Group("/waitlist")
	{
		auth.POST("", signUp)
		auth.POST("/opt-in", optIn)
		auth.POST("/opt-out", optOut)
	}
}

func signUp(c *gin.Context) {
	type Request struct {
		Email                 string                 `json:"email"`
		AdditionalInformation map[string]interface{} `json:"additionalInformation"`
	}

	var r Request

	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if _, err := model.GetRegistrationByEmail(r.Email); err == nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "alreadyRegistered"})
		return
	}

	registration := &model.Registration{
		Email:            r.Email,
		RegistrationDate: time.Now(),
		Confirmed:        false,
	}

	additionalInfo, err := json.Marshal(r.AdditionalInformation)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid additionalInformation"})
		return
	}
	registration.AdditionalInformation = datatypes.JSON(additionalInfo)

	registration.Save()

	client := resend.NewClient(os.Getenv("RESEND_API_KEY"))

	html, err := ParseTemplate(os.Getenv("BASE_DIR")+"/templates/opt-in.html", struct {
		AppName string
		LogoUrl string
		AccentColor string
		AdditonalText string
		Url string
	}{
		AppName: os.Getenv("APP_NAME"),
		LogoUrl: os.Getenv("EMAIL_LOGO_URL"),
		AccentColor: os.Getenv("EMAIL_ACCENT_COLOR"),
		AdditonalText: os.Getenv("EMAIL_ADDITIONAL_TEXT"),
		Url: os.Getenv("CONFIRMATION_CALLBACK_URL") + "?email=" + url.QueryEscape(registration.Email) + "&token=" + url.QueryEscape(registration.ConfirmationToken),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	params := &resend.SendEmailRequest{
		To:      []string{registration.Email},
		From:    os.Getenv("EMAIL_FROM"),
		Html:    html,
		Subject: os.Getenv("EMAIL_SUBJECT"),
	}

	_, err = client.Emails.Send(params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "emailSendingError"})
		registration.Delete()
		return
	}

	// fmt.Println(sent.Id)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}

func optIn(c *gin.Context) {
	type Request struct {
		Email             string `json:"email"`
		ConfirmationToken string `json:"confirmationToken"`
	}

	var r Request

	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registration, err := model.GetRegistrationByEmail(r.Email)

	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	if r.ConfirmationToken != registration.ConfirmationToken {
		c.Status(http.StatusBadRequest)
		return
	}

	registration.Confirmed = true
	registration.Save()

	c.Status(http.StatusNoContent)
}

func optOut(c *gin.Context) {
	type Request struct {
		Email string `json:"email"`
	}

	var r Request

	if err := c.ShouldBindJSON(&r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	registration, err := model.GetRegistrationByEmail(r.Email)

	if err != nil {
		c.Status(http.StatusNoContent)
		return
	}

	if err := registration.Delete(); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Status(http.StatusNoContent)
}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
