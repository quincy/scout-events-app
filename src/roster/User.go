package roster

import (
	"github.com/quincy/scout-events-app/src/date"
)

type AdultUser struct {
	Name        string
	BsaId       int64
	Email       string
	Gender      string
	UnitNumber  string
	Training    []UserStatusRecord
	HealthForms []UserStatusRecord
	SwimClass   UserStatusRecord
	Positions   []string
}

type YouthUser struct {
	Name        string
	BsaId       int64
	Gender      string
	DateOfBirth date.Date
	Age         int
	Patrol      string
	Training    []UserStatusRecord
	HealthForms []UserStatusRecord
	SwimClass   UserStatusRecord
	Positions   []string
}
