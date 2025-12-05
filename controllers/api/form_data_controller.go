package api

import (
	"core/models"
	"core/services"
	"core/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"html"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strings"
)

var (
	CCLIST_MAX_AMOUNT      = 2
	TELEGRAM_MESSAGE_LIMIT = 4096
)

func CreateFormData(c *gin.Context) {
	data := c.Param("data")
	toUuid := utils.GetUUIDFromString(data)
	if toUuid == uuid.Nil {
		showErrorPage(c, "Token is not valid.")
		return
	}
	formToken, err := models.GetFormTokenByUuid(toUuid)
	if err != nil {
		return
	}

	origin := utils.GetRequestOrigin(c)
	if origin == "" {
		showErrorPage(c, "Request origin is not valid.")
		return
	}

	domains := services.GetDomainsName(formToken.UserID)
	if !slices.Contains(domains, origin) {
		showErrorPage(c, "This is not an allowed domain.")
		return
	}

	err = c.Request.ParseForm()
	if err != nil {
		showErrorPage(c, "Error occurred while submitting a form.")
		return
	}
	JSONData := readFormData(c.Request.Form)

	altchaValue, exists := JSONData["altcha"]
	if exists {
		altchaParam, ok := altchaValue.(string)
		if !ok || altchaParam == "" {
			showErrorPage(c, "Error occurred while submitting a form.")
			return
		}
		if !services.IsCaptchaValid(altchaParam) {
			showErrorPage(c, "Captcha is not valid.")
			return
		}
		delete(JSONData, "altcha")
	}

	go sendToTelegram(formToken, JSONData)

	if val, ok := JSONData["_next"]; ok {
		if len(strings.TrimSpace(val.(string))) > 0 {
			c.Redirect(http.StatusMovedPermanently, val.(string))
			c.Abort()
			return
		}
	}
	c.HTML(http.StatusOK, "form-verification.html", gin.H{
		"text":     "Form submitted successfully.",
		"formyUrl": os.Getenv("BASE_URL"),
	})
}

func showErrorPage(c *gin.Context, errorText string) {
	c.HTML(http.StatusOK, "form-verification.html", gin.H{
		"text":     errorText,
		"formyUrl": os.Getenv("BASE_URL"),
	})
}

func readFormData(form url.Values) map[string]interface{} {
	JSONData := map[string]interface{}{}
	for key, values := range form {
		for _, value := range values {
			JSONData[key] = html.EscapeString(strings.TrimSpace(value))
		}
	}
	return JSONData
}

func sendToTelegram(formToken *models.FormToken, formData map[string]interface{}) {
	subject := "New form submission"
	if val, ok := formData["_subject"]; ok {
		if len(strings.TrimSpace(val.(string))) > 0 {
			subject = val.(string)
		}
	}
	telegramBody := createTelegramBody(subject, formData)
	for _, message := range telegramBody {
		services.SendTelegramMessage(formToken.ChatID, message)
	}

	if val, ok := formData["_cc"]; ok {
		if len(strings.TrimSpace(val.(string))) > 0 {
			ccList := strings.SplitN(formData["_cc"].(string), ",", CCLIST_MAX_AMOUNT+1)
			for i, to := range ccList {
				if i >= CCLIST_MAX_AMOUNT {
					break
				}
				formToken, err := models.GetFormTokenByUuid(utils.GetUUIDFromString(to))
				if err != nil {
					return
				}
				for _, message := range telegramBody {
					services.SendTelegramMessage(formToken.ChatID, message)
				}
			}
		}
	}
}

func createTelegramBody(subject string, formValues map[string]interface{}) []string {
	var messageBuilder strings.Builder
	messageBuilder.WriteString(fmt.Sprintf("<b>%s</b>\n\n", subject))
	var hashtags []string
	messageBuilder.WriteString("Submitted fields:\n")
	for key := range formValues {
		if !strings.HasPrefix(key, "_") {
			hashtags = append(hashtags, "#"+strings.ReplaceAll(key, " ", "_"))
		}
	}
	messageBuilder.WriteString(strings.Join(hashtags, " ") + "\n\n")
	for key, value := range formValues {
		if strings.HasPrefix(key, "_") {
			continue
		}
		messageBuilder.WriteString(fmt.Sprintf("<b>#%s:</b>\n%v\n\n", key, value))
	}
	fullMessage := messageBuilder.String()
	var messages []string
	for len(fullMessage) > 0 {
		if len(fullMessage) > TELEGRAM_MESSAGE_LIMIT {
			messages = append(messages, fullMessage[:TELEGRAM_MESSAGE_LIMIT])
			fullMessage = fullMessage[TELEGRAM_MESSAGE_LIMIT:]
		} else {
			messages = append(messages, fullMessage)
			break
		}
	}
	return messages
}

func createHTMLBody(formValues map[string]interface{}) string {
	template := ReadMailTemplate("/views/mails/form-template.html")
	tableData := ""
	for key, value := range formValues {
		if strings.HasPrefix(key, "_") {
			continue
		}
		tableData += fmt.Sprintf("<tr><td>%s</td><td>%s</td><tr>", key, value)
	}
	return strings.Replace(template, "%s", tableData, -1)
}

func ReadMailTemplate(path string) string {
	pwd, _ := os.Getwd()
	templateFile, err := os.ReadFile(pwd + path)
	if err != nil {
		fmt.Print("Error on reading email template.")
	}
	return string(templateFile)
}
