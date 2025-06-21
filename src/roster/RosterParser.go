package roster

import (
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/quincy/scout-events-app/src/date"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const noLimit = -1
const paddedPipe = `\s*\|\s*`

var paddedPipePattern = regexp.MustCompile(paddedPipe)

type Parser interface {
	ParseAdultRoster(input *os.File) ([]AdultScoutbookUser, error)
	ParseYouthRoster(input *os.File) ([]YouthScoutbookUser, error)
}

type csvParser struct{}

func NewCsvParser() Parser {
	return &csvParser{}
}

func (p *csvParser) ParseAdultRoster(input *os.File) ([]AdultScoutbookUser, error) {
	r := csv.NewReader(input)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) <= 2 { // Check for header-only or empty file
		return nil, EmptyRosterError
	}

	// Skip header rows (index 0 and 1) and process remaining records
	var users []AdultScoutbookUser
	for _, record := range records[2:] {
		user := AdultScoutbookUser{
			FirstName:           parseString(record[1]),
			LastName:            parseString(record[2]),
			Email:               parseString(record[3]),
			Gender:              parseString(record[4]),
			BsaId:               parseInt64(record[5]),
			UnitNumber:          parseString(record[6]),
			Training:            parseString(record[7]),
			TrainingExpiration:  parseString(record[8]),
			HealthForms:         parseString(record[9]),
			SwimClass:           parseString(record[10]),
			SwimClassExpiration: parseString(record[11]),
			Positions:           parseString(record[12]),
		}
		users = append(users, user)
	}

	return users, nil
}

func (p *csvParser) ParseYouthRoster(input *os.File) ([]YouthScoutbookUser, error) {
	r := csv.NewReader(input)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(records) <= 2 { // Check for header-only or empty file
		return nil, EmptyRosterError
	}

	// Skip header rows (index 0 and 1) and process remaining records
	var users []YouthScoutbookUser
	for _, record := range records[2:] {
		user := YouthScoutbookUser{
			FirstName:           parseString(record[1]),
			LastName:            parseString(record[2]),
			BsaId:               parseInt64(record[3]),
			DateOfBirth:         parseString(record[4]),
			Age:                 parseInt(record[5]),
			Gender:              parseString(record[6]),
			HealthForms:         parseString(record[7]),
			SwimClass:           parseString(record[8]),
			SwimClassExpiration: parseString(record[9]),
			Positions:           parseString(record[10]),
			Patrol:              parseString(record[11]),
			Training:            parseString(record[12]),
			TrainingExpiration:  parseString(record[13]),
		}
		users = append(users, user)
	}

	return users, nil
}

const SPACE = " "

func parseString(value string) string {
	return strings.Trim(value, SPACE)
}

func parseInt(value string) int {
	trimmed := strings.Trim(value, SPACE)
	bsaId, err := strconv.Atoi(trimmed)
	if err != nil {
		return 0
	}
	return bsaId
}

func parseInt64(value string) int64 {
	trimmed := strings.Trim(value, SPACE)
	bsaId, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return 0
	}
	return bsaId
}

func parseYouthTraining(training string, expiration string) ([]UserStatusRecord, error) {
	// examples
	// Training: "Y01 Safeguarding Youth Training Certification (Expiration Date: 03/26/2027) | SCO_800 Hazardous Weather Training (Expiration Date: 05/19/2027) |"
	// TrainingExpiration: ""

	if expiration != "" {
		return nil, errors.New("training expiration date is not expected. Scoutbook has started sending expiration dates for youth training")
	}

	if training == "" {
		return []UserStatusRecord{}, nil
	}

	var trainingPattern = regexp.MustCompile(`(?P<name>.*?)\s*\(Expiration Date:\s*(?P<date>\d{2}/\d{2}/\d{4})\)`)
	trainings := paddedPipePattern.Split(training, noLimit)

	var records []UserStatusRecord
	for _, t := range trainings {
		if t == "" {
			continue
		}

		matches := trainingPattern.FindStringSubmatch(t)
		if matches == nil {
			continue
		}

		name := matches[trainingPattern.SubexpIndex("name")]
		dateStr := matches[trainingPattern.SubexpIndex("date")]

		expirationDate, err := date.ParseDate(dateStr)
		if err != nil {
			return nil, fmt.Errorf("invalid date format in training: %s", t)
		}

		records = append(records, TrainingRecord(strings.TrimSpace(name), expirationDate))
	}

	return records, nil
}

func parseAdultTraining(training string, expiration string) ([]UserStatusRecord, error) {
	if training == "" {
		return []UserStatusRecord{}, nil
	}

	trainings := paddedPipePattern.Split(training, noLimit)
	expirations := paddedPipePattern.Split(expiration, noLimit)

	var records []UserStatusRecord
	for i, t := range trainings {
		if t == "" {
			continue
		}

		var exp string
		if len(expirations) > i {
			exp = expirations[i]
		}

		expirationDate, err := date.ParseDate(exp)
		if err != nil {
			return nil, fmt.Errorf("invalid date format in training: %s", t)
		}

		records = append(records, TrainingRecord(t, expirationDate))
	}

	return records, nil
}

func parseHealthForms(healthForms string) ([]UserStatusRecord, error) {
	// 06/19/2026(AB) (Expired) | 06/08/2022(C)
	var healthFormPattern = regexp.MustCompile(`(?sm)(?P<date>\d{2}/\d{2}/\d{4})\s*\((?P<type>AB|C)\)(?:\s*\(Expired\))?`)
	forms := paddedPipePattern.Split(healthForms, noLimit)

	records := []UserStatusRecord{}
	for _, form := range forms {
		matches := healthFormPattern.FindStringSubmatch(form)
		if matches == nil {
			continue
		}

		dateStr := matches[healthFormPattern.SubexpIndex("date")]
		formType := matches[healthFormPattern.SubexpIndex("type")]

		if formType == "AB" {
			expirationDate, err := date.ParseDate(dateStr)
			if err == nil {
				records = append(records, HealthFormABRecord(expirationDate))
			}
		} else if formType == "C" {
			expirationDate, err := date.ParseDate(dateStr)
			if err == nil {
				records = append(records, HealthFormCRecord(expirationDate))
			}
		} else {
			return []UserStatusRecord{}, fmt.Errorf("unknown health form type: [%s] in '%s'", formType, healthForms)
		}
	}

	return records, nil
}

func parseSwimClass(swimClass string, swimClassExpiration string) (UserStatusRecord, error) {
	// Swimmer (01/01/2015)
	if swimClass == "" {
		return NonSwimmerRecord(), nil
	}

	// Check if they're a swimmer
	if strings.Contains(swimClass, "Swimmer") {
		// Parse expiration date if available
		if swimClassExpiration != "" {
			expirationDate, err := date.ParseDate(swimClassExpiration)
			if err == nil {
				return SwimmerRecord(expirationDate), nil
			}
		}

		// If date parsing fails or no expiration, use empty date
		return UserStatusRecord{}, errors.New("swim class expiration is invalid")
	}

	// Default to non-swimmer for all other cases
	return NonSwimmerRecord(), nil
}

// AdultScoutbookUser is a placeholder to put all the values parsed from the Scoutbook CSV
// before it's mapped to an AdultUser
type AdultScoutbookUser struct {
	FirstName           string
	LastName            string
	Email               string
	Gender              string
	BsaId               int64
	UnitNumber          string
	Training            string
	TrainingExpiration  string
	HealthForms         string
	SwimClass           string
	SwimClassExpiration string
	Positions           string
}

func (u *AdultScoutbookUser) ToAdultUser() (AdultUser, error) {
	healthForms, err := parseHealthForms(u.HealthForms)
	if err != nil {
		return AdultUser{}, err
	}

	training, err := parseAdultTraining(u.Training, u.TrainingExpiration)
	if err != nil {
		return AdultUser{}, err
	}

	swimClass, err := parseSwimClass(u.SwimClass, u.SwimClassExpiration)
	if err != nil {
		return AdultUser{}, err
	}

	return AdultUser{
		Name:        u.FirstName + " " + u.LastName,
		BsaId:       u.BsaId,
		Email:       u.Email,
		Gender:      u.Gender,
		UnitNumber:  u.UnitNumber,
		Training:    training,
		HealthForms: healthForms,
		SwimClass:   swimClass,
		Positions:   strings.Split(u.Positions, " | "),
	}, nil
}

// YouthScoutbookUser is a placeholder to put all the values parsed from the Scoutbook CSV
// before it's mapped to a YouthUser
type YouthScoutbookUser struct {
	FirstName           string
	LastName            string
	BsaId               int64
	Email               string
	DateOfBirth         string
	Age                 int
	Gender              string
	HealthForms         string
	SwimClass           string
	SwimClassExpiration string
	Positions           string
	Patrol              string
	Training            string
	TrainingExpiration  string
}

func (u *YouthScoutbookUser) ToYouthUser() (YouthUser, error) {
	healthForms, err := parseHealthForms(u.HealthForms)
	if err != nil {
		return YouthUser{}, err
	}

	bday, err := date.ParseDate(u.DateOfBirth)
	if err != nil {
		return YouthUser{}, nil
	}

	training, err := parseYouthTraining(u.Training, u.TrainingExpiration)
	if err != nil {
		return YouthUser{}, err
	}

	swimClass, err := parseSwimClass(u.SwimClass, u.SwimClassExpiration)
	if err != nil {
		return YouthUser{}, err
	}

	var positions []string
	if u.Positions == "" {
		positions = []string{}
	} else {
		positions = paddedPipePattern.Split(u.Positions, noLimit)
	}

	return YouthUser{
		Name:        string(u.FirstName[0]) + ". " + u.LastName,
		BsaId:       u.BsaId,
		Gender:      u.Gender,
		DateOfBirth: bday,
		Age:         u.Age,
		Patrol:      u.Patrol,
		Training:    training,
		HealthForms: healthForms,
		SwimClass:   swimClass,
		Positions:   positions,
	}, nil
}

var EmptyRosterError = errors.New("roster file is empty")
