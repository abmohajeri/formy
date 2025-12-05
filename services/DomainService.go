package services

import (
	"core/config"
	"core/models"
	"log"
)

func CreateUserAllowedDomain(user *models.User, domain string) (bool, string) {
	var allowedDomain models.AllowedDomain
	errForm := config.GetDB().Where("user_id = ? and name = ?", user.ID, domain).First(&allowedDomain)
	if errForm.RowsAffected == 0 {
		allowedDomain = models.AllowedDomain{
			Name:   domain,
			UserID: user.ID,
		}
		err := allowedDomain.Save()
		if err != nil {
			return false, `Error occurred\! Please try again\.`
		}
	} else {
		return false, `Domain is exist for your user\! Please try another domain\.`
	}
	return true, ""
}

func GetDomains(userId uint64) []models.AllowedDomain {
	var domains []models.AllowedDomain
	err := config.GetDB().Where("user_id = ?", userId).Find(&domains).Error
	if err != nil {
		log.Println("Error fetching domains:", err)
		return nil
	}
	return domains
}

func GetDomainsName(userId uint64) []string {
	domains := GetDomains(userId)
	var domainsName []string
	for _, domain := range domains {
		domainsName = append(domainsName, domain.Name)
	}
	return domainsName
}
