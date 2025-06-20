package roster

import (
	"github.com/quincy/scout-events-app/src/date"
)

type UserStatusRecord struct {
	Type           RecordType
	Name           string
	ExpirationDate date.Date
}

type RecordType int

const (
	Training RecordType = iota
	HealthForm
	SwimClass
)

func (rt RecordType) String() string {
	switch rt {
	case Training:
		return "Training"
	case HealthForm:
		return "Health Form"
	case SwimClass:
		return "Swimmer Classification"
	default:
		return "Unknown"
	}
}

func TrainingRecord(name string, expirationDate date.Date) UserStatusRecord {
	return UserStatusRecord{
		Type:           Training,
		Name:           name,
		ExpirationDate: expirationDate,
	}
}

func HealthFormABRecord(expirationDate date.Date) UserStatusRecord {
	return UserStatusRecord{
		Type:           HealthForm,
		Name:           "Health Form Parts A/B",
		ExpirationDate: expirationDate,
	}
}

func HealthFormCRecord(expirationDate date.Date) UserStatusRecord {
	return UserStatusRecord{
		Type:           HealthForm,
		Name:           "Health Form Part C",
		ExpirationDate: expirationDate,
	}
}

func SwimmerRecord(expirationDate date.Date) UserStatusRecord {
	return UserStatusRecord{
		Type:           SwimClass,
		Name:           "Swimmer",
		ExpirationDate: expirationDate,
	}
}

func NonSwimmerRecord() UserStatusRecord {
	return UserStatusRecord{
		Type: SwimClass,
		Name: "Non-Swimmer",
	}
}
