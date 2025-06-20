package roster

import (
	"errors"
	"github.com/quincy/scout-events-app/src/assertions"
	"github.com/quincy/scout-events-app/src/date"
	"os"
	"testing"
	"time"
)

func Test_RosterParserFailsOnEmptyAdultFile(t *testing.T) {
	// Given the input file is empty
	emptyFile, err := os.CreateTemp("", "empty.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	}(emptyFile.Name())
	defer func(tmpFile *os.File) {
		_ = tmpFile.Close()
	}(emptyFile)

	// When I try to parse the input file
	_, err = NewCsvParser().ParseAdultRoster(emptyFile)

	// Then the RosterParser returns an error
	if err == nil {
		t.Fatalf("Expected an error")
	}

	if !errors.Is(err, EmptyRosterError) {
		t.Fatalf("Expected error to be EmptyRosterError got: %T -> %v", err, err)
	}
}

func Test_RosterParserCanParseAdultRoster(t *testing.T) {
	// Given the input file `adult-roster-example.csv`
	csv, err := os.Open("test_resources/adult-roster-example.csv")
	if err != nil {
		t.Fatalf("Failed to open roster file: %v", err)
	}
	defer func(tmpFile *os.File) {
		_ = tmpFile.Close()
	}(csv)

	// When I try to parse the input file
	actualUsers, err := NewCsvParser().ParseAdultRoster(csv)
	if err != nil {
		t.Fatalf("Failed to parse roster file: %v", err)
	}

	// Then the RosterParser returns the expected users
	expectedUsers := []AdultScoutbookUser{
		{FirstName: "Alice", LastName: "Ames", Email: "aames@example.com", Gender: "F", BsaId: 1, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification", TrainingExpiration: "03/03/2027", HealthForms: "05/06/2025(AB) | 05/06/2025 (C)", SwimClass: "", SwimClassExpiration: "", Positions: "Committee Member"},
		{FirstName: "Bob", LastName: "Brown", Email: "bbrown@example.com", Gender: "M", BsaId: 2, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification  | SCO_800 Hazardous Weather Training", TrainingExpiration: "03/03/2027 | 06/07/2025", HealthForms: "05/06/2023(AB) (Expired) | 05/06/2025 (C)", SwimClass: "", SwimClassExpiration: "", Positions: "Assistant Scoutmaster"},
		{FirstName: "Carol", LastName: "Carson", Email: "ccarson@example.com", Gender: "F", BsaId: 3, UnitNumber: "Troop 77 B", Training: "", TrainingExpiration: "", HealthForms: "05/06/2025(AB) | 05/06/2023 (C) (Expired)", SwimClass: "", SwimClassExpiration: "", Positions: "Chartered Organization Rep."},
		{FirstName: "Dan", LastName: "Dewey", Email: "ddewey@example.com", Gender: "M", BsaId: 4, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification", TrainingExpiration: "02/22/2027", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "Committee Member | Unit Advancement Chair"},
		{FirstName: "Erin", LastName: "Eckhart", Email: "eeckhart@example.com", Gender: "F", BsaId: 5, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification", TrainingExpiration: "04/04/2027", HealthForms: "02/20/2019(AB) (Expired) | 02/25/2021(C) (Expired)", SwimClass: "", SwimClassExpiration: "", Positions: "Assistant Scoutmaster | Assistant Scoutmaster | Unit Outdoors / Activities Chair"},
		{FirstName: "Frank", LastName: "Faraday", Email: "ffaraday@example.com", Gender: "M", BsaId: 6, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification", TrainingExpiration: "03/03/2027", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "Unit Scouter Reserve"},
		{FirstName: "Gertrude", LastName: "Grisham", Email: "ggrisham@example.com", Gender: "F", BsaId: 7, UnitNumber: "Troop 77 B", Training: "", TrainingExpiration: "", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "Executive Officer"},
		{FirstName: "Harold", LastName: "Hunt", Email: "hhunt@example.com", Gender: "M", BsaId: 8, UnitNumber: "Troop 77 B", Training: "", TrainingExpiration: "", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "Scoutmaster"},
		{FirstName: "Irene", LastName: "Icabod", Email: "iicabod@example.com", Gender: "F", BsaId: 9, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification", TrainingExpiration: "01/02/2027", HealthForms: "05/24/2018(AB) (Expired) | 05/24/2018(C) (Expired)", SwimClass: "", SwimClassExpiration: "", Positions: "Unit College Scouter Reserve"},
		{FirstName: "Jeff", LastName: "Jones", Email: "jjones@example.com", Gender: "M", BsaId: 10, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification  | SCO_800 Hazardous Weather Training", TrainingExpiration: "06/13/2026 | 06/05/2025", HealthForms: "03/05/2019(AB) (Expired) |", SwimClass: "Swimmer (01/01/2015)", SwimClassExpiration: "01/01/2015", Positions: "Assistant Scoutmaster | Unit Training Chair | Youth Protection Champion"},
		{FirstName: "Kristina", LastName: "Kent", Email: "kkent@example.com", Gender: "F", BsaId: 11, UnitNumber: "Troop 77 B", Training: "", TrainingExpiration: "", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "Unit Treasurer"},
		{FirstName: "Leonard", LastName: "Lewis", Email: "llewis@example.com", Gender: "M", BsaId: 12, UnitNumber: "Troop 77 B", Training: "", TrainingExpiration: "", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "Committee Chairman | Life-to-Eagle Coordinator"},
		{FirstName: "Mary", LastName: "Mumford", Email: "mmumford@example.com", Gender: "F", BsaId: 13, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification", TrainingExpiration: "02/21/2026", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "Committee Membership Coordinator | New Member Coordinator"},
	}

	if !AdultScoutbookUsers(actualUsers).ContainsExactly(expectedUsers) {
		t.Fatalf("Expected users to be\n    %v\ngot %v", expectedUsers, actualUsers)
	}
}

func Test_RosterParserFailsOnEmptyYouthFile(t *testing.T) {
	// Given the input file is empty
	emptyFile, err := os.CreateTemp("", "empty.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer func(name string) {
		err := os.Remove(name)
		if err != nil {
			t.Fatalf("Failed to remove temp file: %v", err)
		}
	}(emptyFile.Name())
	defer func(tmpFile *os.File) {
		_ = tmpFile.Close()
	}(emptyFile)

	// When I try to parse the input file
	_, err = NewCsvParser().ParseYouthRoster(emptyFile)

	// Then the RosterParser returns an error
	if err == nil {
		t.Fatalf("Expected an error")
	}

	if !errors.Is(err, EmptyRosterError) {
		t.Fatalf("Expected error to be EmptyRosterError got: %T -> %v", err, err)
	}
}

func Test_RosterParserCanParseYouthRoster(t *testing.T) {
	// Given the input file `adult-roster-example.csv`
	csv, err := os.Open("test_resources/youth-roster-example.csv")
	if err != nil {
		t.Fatalf("Failed to open roster file: %v", err)
	}
	defer func(tmpFile *os.File) {
		_ = tmpFile.Close()
	}(csv)

	// When I try to parse the input file
	actualUsers, err := NewCsvParser().ParseYouthRoster(csv)
	if err != nil {
		t.Fatalf("Failed to parse roster file: %v", err)
	}

	// Then the RosterParser returns the expected users
	expectedUsers := []YouthScoutbookUser{
		{FirstName: "Abe", LastName: "Ames", BsaId: 100, DateOfBirth: "07/04/2014", Age: 10, Gender: "M", HealthForms: "06/19/2026(AB) | 06/08/2026(C)", SwimClass: "", SwimClassExpiration: "", Positions: "Patrol Leader [ Vikings] Patrol | Scouts BSA [ Vikings] Patrol", Patrol: "Vikings", Training: "", TrainingExpiration: ""},
		{FirstName: "Billy", LastName: "Brown", BsaId: 101, DateOfBirth: "08/01/2007", Age: 17, Gender: "M", HealthForms: "06/19/2026(AB) (Expired) | 06/08/2022(C)", SwimClass: "Swimmer", SwimClassExpiration: "05/28/2019", Positions: "Scouts BSA [ Dreadnoughts] Patrol", Patrol: "Dreadnoughts", Training: "", TrainingExpiration: ""},
		{FirstName: "Charlie", LastName: "Carson", BsaId: 102, DateOfBirth: "03/11/2011", Age: 14, Gender: "M", HealthForms: "06/18/2026(AB) | 06/08/2022(C) (Expired)", SwimClass: "Nonswimmer", SwimClassExpiration: "", Positions: "Chaplain Aide | Scouts BSA [ Warthogs] Patrol", Patrol: "Warthogs", Training: "", TrainingExpiration: ""},
		{FirstName: "Daryl", LastName: "Dewey", BsaId: 103, DateOfBirth: "12/16/2008", Age: 16, Gender: "M", HealthForms: "06/18/2022(AB) (Expired) | 06/08/2022(C) (Expired)", SwimClass: "Nonswimmer", SwimClassExpiration: "", Positions: "OA Unit Representative | Scouts BSA [ Dreadnoughts] Patrol", Patrol: "Dreadnoughts", Training: "Y01 Safeguarding Youth Training Certification (Expiration Date: 03/26/2027) | SCO_800 Hazardous Weather Training (Expiration Date: 05/19/2027) |", TrainingExpiration: ""},
		{FirstName: "Ed", LastName: "Eckhart", BsaId: 104, DateOfBirth: "12/16/2009", Age: 17, Gender: "M", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "", Patrol: "", Training: "", TrainingExpiration: ""},
	}

	if !YouthScoutbookUsers(actualUsers).ContainsExactly(expectedUsers) {
		t.Fatalf("Expected users to be\n    %v\ngot %v", expectedUsers, actualUsers)
	}
}

func Test_AdultScoutbookUsersCanMapToAdultUsers(t *testing.T) {
	// Given a set of AdultScoutbookUsers
	scoutbookUser := []AdultScoutbookUser{
		{FirstName: "Alice", LastName: "Ames", Email: "aames@example.com", Gender: "F", BsaId: 1, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification", TrainingExpiration: "03/03/2027", HealthForms: "05/06/2025(AB) | 05/06/2025 (C)", SwimClass: "", SwimClassExpiration: "", Positions: "Committee Member"},
		{FirstName: "Bob", LastName: "Brown", Email: "bbrown@example.com", Gender: "M", BsaId: 2, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification  | SCO_800 Hazardous Weather Training", TrainingExpiration: "03/03/2027 | 06/07/2025", HealthForms: "05/06/2025(AB) (Expired) | 05/06/2025 (C)", SwimClass: "", SwimClassExpiration: "", Positions: "Assistant Scoutmaster"},
		{FirstName: "Carol", LastName: "Carson", Email: "ccarson@example.com", Gender: "F", BsaId: 3, UnitNumber: "Troop 77 B", Training: "", TrainingExpiration: "", HealthForms: "05/06/2025(AB) | 05/06/2023 (C) (Expired)", SwimClass: "", SwimClassExpiration: "", Positions: "Chartered Organization Rep."},
		{FirstName: "Dan", LastName: "Dewey", Email: "ddewey@example.com", Gender: "M", BsaId: 4, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification", TrainingExpiration: "02/22/2027", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "Committee Member | Unit Advancement Chair"},
		{FirstName: "Erin", LastName: "Eckhart", Email: "eeckhart@example.com", Gender: "F", BsaId: 5, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification", TrainingExpiration: "04/04/2027", HealthForms: "02/20/2019(AB) (Expired) | 02/25/2021(C) (Expired)", SwimClass: "", SwimClassExpiration: "", Positions: "Assistant Scoutmaster | Assistant Scoutmaster | Unit Outdoors / Activities Chair"},
		{FirstName: "Frank", LastName: "Faraday", Email: "ffaraday@example.com", Gender: "M", BsaId: 6, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification", TrainingExpiration: "03/03/2027", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "Unit Scouter Reserve"},
		{FirstName: "Gertrude", LastName: "Grisham", Email: "ggrisham@example.com", Gender: "F", BsaId: 7, UnitNumber: "Troop 77 B", Training: "", TrainingExpiration: "", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "Executive Officer"},
		{FirstName: "Harold", LastName: "Hunt", Email: "hhunt@example.com", Gender: "M", BsaId: 8, UnitNumber: "Troop 77 B", Training: "", TrainingExpiration: "", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "Scoutmaster"},
		{FirstName: "Irene", LastName: "Icabod", Email: "iicabod@example.com", Gender: "F", BsaId: 9, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification", TrainingExpiration: "01/02/2027", HealthForms: "05/24/2018(AB) (Expired) | 05/24/2018(C) (Expired)", SwimClass: "", SwimClassExpiration: "", Positions: "Unit College Scouter Reserve"},
		{FirstName: "Jeff", LastName: "Jones", Email: "jjones@example.com", Gender: "M", BsaId: 10, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification  | SCO_800 Hazardous Weather Training", TrainingExpiration: "06/13/2026 | 06/05/2025", HealthForms: "03/05/2019(AB) (Expired) |", SwimClass: "Swimmer (01/01/2015)", SwimClassExpiration: "01/01/2015", Positions: "Assistant Scoutmaster | Unit Training Chair | Youth Protection Champion"},
		{FirstName: "Kristina", LastName: "Kent", Email: "kkent@example.com", Gender: "F", BsaId: 11, UnitNumber: "Troop 77 B", Training: "", TrainingExpiration: "", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "Unit Treasurer"},
		{FirstName: "Leonard", LastName: "Lewis", Email: "llewis@example.com", Gender: "M", BsaId: 12, UnitNumber: "Troop 77 B", Training: "", TrainingExpiration: "", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "Committee Chairman | Life-to-Eagle Coordinator"},
		{FirstName: "Mary", LastName: "Mumford", Email: "mmumford@example.com", Gender: "F", BsaId: 13, UnitNumber: "Troop 77 B", Training: "Y01 Youth Protection Training Certification", TrainingExpiration: "02/21/2026", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "Committee Membership Coordinator | New Member Coordinator"},
	}

	var actualUsers []AdultUser
	for _, user := range scoutbookUser {
		adultUser, err := user.ToAdultUser()
		if err != nil {
			t.Fatalf("Failed to convert AdultScoutbookUser to AdultUser: %v", err)
		}
		actualUsers = append(actualUsers, adultUser)
	}

	expectedUsers := []AdultUser{
		{Name: "Alice Ames", Email: "aames@example.com", Gender: "F", BsaId: 1, UnitNumber: "Troop 77 B", Training: []UserStatusRecord{TrainingRecord("Y01 Youth Protection Training Certification", date.NewDate(2027, time.March, 3))}, HealthForms: []UserStatusRecord{HealthFormABRecord(date.NewDate(2025, time.May, 6)), HealthFormCRecord(date.NewDate(2025, time.May, 6))}, SwimClass: NonSwimmerRecord(), Positions: []string{"Committee Member"}},
		{Name: "Bob Brown", Email: "bbrown@example.com", Gender: "M", BsaId: 2, UnitNumber: "Troop 77 B", Training: []UserStatusRecord{TrainingRecord("Y01 Youth Protection Training Certification", date.NewDate(2027, time.March, 3)), TrainingRecord("SCO_800 Hazardous Weather Training", date.NewDate(2025, time.June, 7))}, HealthForms: []UserStatusRecord{HealthFormABRecord(date.NewDate(2025, time.May, 6)), HealthFormCRecord(date.NewDate(2025, time.May, 6))}, SwimClass: NonSwimmerRecord(), Positions: []string{"Assistant Scoutmaster"}},
		{Name: "Carol Carson", Email: "ccarson@example.com", Gender: "F", BsaId: 3, UnitNumber: "Troop 77 B", Training: []UserStatusRecord{}, HealthForms: []UserStatusRecord{HealthFormABRecord(date.NewDate(2025, time.May, 6)), HealthFormCRecord(date.NewDate(2023, time.May, 6))}, SwimClass: NonSwimmerRecord(), Positions: []string{"Chartered Organization Rep."}},
		{Name: "Dan Dewey", Email: "ddewey@example.com", Gender: "M", BsaId: 4, UnitNumber: "Troop 77 B", Training: []UserStatusRecord{TrainingRecord("Y01 Youth Protection Training Certification", date.NewDate(2027, time.February, 22))}, HealthForms: []UserStatusRecord{}, SwimClass: NonSwimmerRecord(), Positions: []string{"Committee Member", "Unit Advancement Chair"}},
		{Name: "Erin Eckhart", Email: "eeckhart@example.com", Gender: "F", BsaId: 5, UnitNumber: "Troop 77 B", Training: []UserStatusRecord{TrainingRecord("Y01 Youth Protection Training Certification", date.NewDate(2027, time.April, 4))}, HealthForms: []UserStatusRecord{HealthFormABRecord(date.NewDate(2019, time.February, 20)), HealthFormCRecord(date.NewDate(2021, time.February, 25))}, SwimClass: NonSwimmerRecord(), Positions: []string{"Assistant Scoutmaster", "Assistant Scoutmaster", "Unit Outdoors / Activities Chair"}},
		{Name: "Frank Faraday", Email: "ffaraday@example.com", Gender: "M", BsaId: 6, UnitNumber: "Troop 77 B", Training: []UserStatusRecord{TrainingRecord("Y01 Youth Protection Training Certification", date.NewDate(2027, time.March, 3))}, HealthForms: []UserStatusRecord{}, SwimClass: NonSwimmerRecord(), Positions: []string{"Unit Scouter Reserve"}},
		{Name: "Gertrude Grisham", Email: "ggrisham@example.com", Gender: "F", BsaId: 7, UnitNumber: "Troop 77 B", Training: []UserStatusRecord{}, HealthForms: []UserStatusRecord{}, SwimClass: NonSwimmerRecord(), Positions: []string{"Executive Officer"}},
		{Name: "Harold Hunt", Email: "hhunt@example.com", Gender: "M", BsaId: 8, UnitNumber: "Troop 77 B", Training: []UserStatusRecord{}, HealthForms: []UserStatusRecord{}, SwimClass: NonSwimmerRecord(), Positions: []string{"Scoutmaster"}},
		{Name: "Irene Icabod", Email: "iicabod@example.com", Gender: "F", BsaId: 9, UnitNumber: "Troop 77 B", Training: []UserStatusRecord{TrainingRecord("Y01 Youth Protection Training Certification", date.NewDate(2027, time.January, 2))}, HealthForms: []UserStatusRecord{HealthFormABRecord(date.NewDate(2018, time.May, 24)), HealthFormCRecord(date.NewDate(2018, time.May, 24))}, SwimClass: NonSwimmerRecord(), Positions: []string{"Unit College Scouter Reserve"}},
		{Name: "Jeff Jones", Email: "jjones@example.com", Gender: "M", BsaId: 10, UnitNumber: "Troop 77 B", Training: []UserStatusRecord{TrainingRecord("Y01 Youth Protection Training Certification", date.NewDate(2026, time.June, 13)), TrainingRecord("SCO_800 Hazardous Weather Training", date.NewDate(2025, time.June, 5))}, HealthForms: []UserStatusRecord{HealthFormABRecord(date.NewDate(2019, time.March, 5))}, SwimClass: SwimmerRecord(date.NewDate(2015, time.January, 1)), Positions: []string{"Assistant Scoutmaster", "Unit Training Chair", "Youth Protection Champion"}},
		{Name: "Kristina Kent", Email: "kkent@example.com", Gender: "F", BsaId: 11, UnitNumber: "Troop 77 B", Training: []UserStatusRecord{}, HealthForms: []UserStatusRecord{}, SwimClass: NonSwimmerRecord(), Positions: []string{"Unit Treasurer"}},
		{Name: "Leonard Lewis", Email: "llewis@example.com", Gender: "M", BsaId: 12, UnitNumber: "Troop 77 B", Training: []UserStatusRecord{}, HealthForms: []UserStatusRecord{}, SwimClass: NonSwimmerRecord(), Positions: []string{"Committee Chairman", "Life-to-Eagle Coordinator"}},
		{Name: "Mary Mumford", Email: "mmumford@example.com", Gender: "F", BsaId: 13, UnitNumber: "Troop 77 B", Training: []UserStatusRecord{TrainingRecord("Y01 Youth Protection Training Certification", date.NewDate(2026, time.February, 21))}, HealthForms: []UserStatusRecord{}, SwimClass: NonSwimmerRecord(), Positions: []string{"Committee Membership Coordinator", "New Member Coordinator"}},
	}

	if !AdultUsers(actualUsers).ContainsExactly(expectedUsers) {
		t.Fatalf("Expected users to be\n    %v\ngot %v", expectedUsers, actualUsers)
	}
}

func Test_YouthScoutbookUsersCanMapToYouthUsers(t *testing.T) {
	// Given a set of YouthScoutbookUsers
	scoutbookUser := []YouthScoutbookUser{
		{FirstName: "Abe", LastName: "Ames", BsaId: 100, DateOfBirth: "07/04/2014", Age: 10, Gender: "M", HealthForms: "06/19/2026(AB) | 06/08/2026(C)", SwimClass: "", SwimClassExpiration: "", Positions: "Patrol Leader [ Vikings] Patrol | Scouts BSA [ Vikings] Patrol", Patrol: "Vikings", Training: "", TrainingExpiration: ""},
		{FirstName: "Billy", LastName: "Brown", BsaId: 101, DateOfBirth: "08/01/2007", Age: 17, Gender: "M", HealthForms: "06/19/2026(AB) (Expired) | 06/08/2022(C)", SwimClass: "Swimmer", SwimClassExpiration: "05/28/2019", Positions: "Scouts BSA [ Dreadnoughts] Patrol", Patrol: "Dreadnoughts", Training: "", TrainingExpiration: ""},
		{FirstName: "Charlie", LastName: "Carson", BsaId: 102, DateOfBirth: "03/11/2011", Age: 14, Gender: "M", HealthForms: "06/18/2026(AB) | 06/08/2022(C) (Expired)", SwimClass: "Nonswimmer", SwimClassExpiration: "", Positions: "Chaplain Aide | Scouts BSA [ Warthogs] Patrol", Patrol: "Warthogs", Training: "", TrainingExpiration: ""},
		{FirstName: "Daryl", LastName: "Dewey", BsaId: 103, DateOfBirth: "12/16/2008", Age: 16, Gender: "M", HealthForms: "06/18/2022(AB) (Expired) | 06/08/2022(C) (Expired)", SwimClass: "Nonswimmer", SwimClassExpiration: "", Positions: "OA Unit Representative | Scouts BSA [ Dreadnoughts] Patrol", Patrol: "Dreadnoughts", Training: "Y01 Safeguarding Youth Training Certification (Expiration Date: 03/26/2027) | SCO_800 Hazardous Weather Training (Expiration Date: 05/19/2027) |", TrainingExpiration: ""},
		{FirstName: "Ed", LastName: "Eckhart", BsaId: 104, DateOfBirth: "12/16/2009", Age: 17, Gender: "M", HealthForms: "", SwimClass: "", SwimClassExpiration: "", Positions: "", Patrol: "", Training: "", TrainingExpiration: ""},
	}

	var actualUsers []YouthUser
	for _, user := range scoutbookUser {
		youthUser, err := user.ToYouthUser()
		if err != nil {
			t.Fatalf("Could not convert YouthScoutbookUser to YouthUser: %v", err)
		}
		actualUsers = append(actualUsers, youthUser)
	}

	expectedUsers := []YouthUser{
		{Name: "A. Ames", BsaId: 100, Gender: "M", DateOfBirth: date.NewDate(2014, time.July, 4), Age: 10, Patrol: "Vikings", Training: []UserStatusRecord{}, HealthForms: []UserStatusRecord{HealthFormABRecord(date.NewDate(2026, time.June, 19)), HealthFormCRecord(date.NewDate(2026, time.June, 8))}, SwimClass: NonSwimmerRecord(), Positions: []string{"Patrol Leader [ Vikings] Patrol", "Scouts BSA [ Vikings] Patrol"}},
		{Name: "B. Brown", BsaId: 101, Gender: "M", DateOfBirth: date.NewDate(2007, time.August, 1), Age: 17, Patrol: "Dreadnoughts", Training: []UserStatusRecord{}, HealthForms: []UserStatusRecord{HealthFormABRecord(date.NewDate(2026, time.June, 19)), HealthFormCRecord(date.NewDate(2022, time.June, 8))}, SwimClass: SwimmerRecord(date.NewDate(2019, time.May, 28)), Positions: []string{"Scouts BSA [ Dreadnoughts] Patrol"}},
		{Name: "C. Carson", BsaId: 102, Gender: "M", DateOfBirth: date.NewDate(2011, time.March, 11), Age: 14, Patrol: "Warthogs", Training: []UserStatusRecord{}, HealthForms: []UserStatusRecord{HealthFormABRecord(date.NewDate(2026, time.June, 18)), HealthFormCRecord(date.NewDate(2022, time.June, 8))}, SwimClass: NonSwimmerRecord(), Positions: []string{"Chaplain Aide", "Scouts BSA [ Warthogs] Patrol"}},
		{Name: "D. Dewey", BsaId: 103, Gender: "M", DateOfBirth: date.NewDate(2008, time.December, 16), Age: 16, Patrol: "Dreadnoughts", Training: []UserStatusRecord{TrainingRecord("Y01 Safeguarding Youth Training Certification", date.NewDate(2027, time.March, 26)), TrainingRecord("SCO_800 Hazardous Weather Training", date.NewDate(2027, time.May, 19))}, HealthForms: []UserStatusRecord{HealthFormABRecord(date.NewDate(2022, time.June, 18)), HealthFormCRecord(date.NewDate(2022, time.June, 8))}, SwimClass: NonSwimmerRecord(), Positions: []string{"OA Unit Representative", "Scouts BSA [ Dreadnoughts] Patrol"}},
		{Name: "E. Eckhart", BsaId: 104, Gender: "M", DateOfBirth: date.NewDate(2009, time.December, 16), Age: 17, Patrol: "", Training: []UserStatusRecord{}, HealthForms: []UserStatusRecord{}, SwimClass: NonSwimmerRecord(), Positions: []string{}},
	}

	if !YouthUsers(actualUsers).ContainsExactly(expectedUsers) {
		t.Fatalf("Expected users to be\n    %v\ngot %v", expectedUsers, actualUsers)
	}
}

type AdultScoutbookUsers = assertions.Collection[AdultScoutbookUser]
type AdultUsers = assertions.Collection[AdultUser]
type YouthScoutbookUsers = assertions.Collection[YouthScoutbookUser]
type YouthUsers = assertions.Collection[YouthUser]
