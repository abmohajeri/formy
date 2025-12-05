package services

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"time"

	"github.com/altcha-org/altcha-lib-go"
)

var altchaHMACKey = os.Getenv("ALTCHA_HMAC_KEY")

func AltchaHandler(c *gin.Context) {
	expires := time.Now().Add(2 * time.Minute)
	challenge, err := altcha.CreateChallenge(altcha.ChallengeOptions{
		HMACKey:   altchaHMACKey,
		MaxNumber: 100_000,
		Expires:   &expires,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to create challenge: %s", err)})
		return
	}

	c.JSON(http.StatusOK, challenge)
}

func IsCaptchaValid(altchaParam string) bool {
	verified, err := altcha.VerifySolution(altchaParam, altchaHMACKey, true)
	if err != nil || !verified {
		return false
	}

	return true
}
