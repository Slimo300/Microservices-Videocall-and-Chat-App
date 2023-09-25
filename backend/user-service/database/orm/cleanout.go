package orm

import (
	"fmt"
	"time"

	"github.com/Slimo300/Microservices-Videocall-and-Chat-App/backend/user-service/models"
)

func (db *Database) CleanCodes(interval time.Duration) {
	t := time.NewTicker(interval)
	defer t.Stop()

	for range t.C {
		fmt.Println("Cleaning codes")
		db.Model(&models.ResetCode{}).Where("created < ?", time.Now().Add(-db.Config.ResetCodeDuration)).Delete(&models.ResetCode{})
		db.Model(&models.VerificationCode{}).Where("created < ?", time.Now().Add(-db.Config.VerificationCodeDuration)).Delete(&models.ResetCode{})
	}
}
