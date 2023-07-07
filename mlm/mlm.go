package mlm

import (
	"regexp"
	"strings"
)

type MedicalLogicModule struct {
	Title       string
	Name        string
	Arden       string
	Version     string
	Date        string
	Usage       int
	Status      int
	Description string
	Content     string

	logic         string // content without comments
	evoke         string
	evoked_events []string
}

// remove all comments
func (mlm *MedicalLogicModule) parseLogic() {
	// remove all c++-style comments from data using regex
	re := regexp.MustCompile(`//.*`)
	logic := re.ReplaceAllString(mlm.Content, "")

	// remove block level comments from text using regular expression
	re = regexp.MustCompile(`(?sU)(/\*.*\*/)`)
	logic = re.ReplaceAllString(logic, "")

	mlm.logic = logic
}

// parse evoke slot
func (mlm *MedicalLogicModule) parseEvokeSlot() {
	re := regexp.MustCompile(`(?isU);;\s+evoke:(.*);;`)
	match := re.FindStringSubmatch(mlm.logic)

	if match != nil {
		mlm.evoke = match[1]
	}
}

func (mlm *MedicalLogicModule) parseTitle() {
	re := regexp.MustCompile(`(?siU)maintenance:.*title:(.*);;`)
	match := re.FindStringSubmatch(mlm.logic)
	if match != nil {
		mlm.Title = strings.Trim(match[1], " ")
	}
}

func (mlm *MedicalLogicModule) parseName() {
	re := regexp.MustCompile(`(?siU)maintenance:.*title:.*;;.*(?:mlmname|filename):(.*);;`)
	match := re.FindStringSubmatch(mlm.logic)
	if match != nil {
		mlm.Name = strings.ToUpper(strings.Trim(match[1], " "))
	}
}

// parse arden version
func (mlm *MedicalLogicModule) parseArden() {
	re := regexp.MustCompile(`(?siU)maintenance:.*title:.*;;.*(?:mlmname|filename):.*;;.*arden:\s*version\s*(.*);;`)
	match := re.FindStringSubmatch(mlm.logic)
	if match != nil {
		mlm.Arden = strings.Trim(match[1], " ")
	}
}

// parse mlm version
func (mlm *MedicalLogicModule) parseVersion() {
	re := regexp.MustCompile(`(?siU)maintenance:.*title:.*;;.*(?:mlmname|filename):.*;;.*arden:.*;;.*version:(.*);;`)
	match := re.FindStringSubmatch(mlm.logic)
	if match != nil {
		mlm.Version = strings.Trim(match[1], " ")
	}
}

// parse date
func (mlm *MedicalLogicModule) parseDate() {
	re := regexp.MustCompile(`(?siU)maintenance:.*title:.*;;.*(?:mlmname|filename):.*;;.*arden:.*;;.*version:.*;;.*date:(.*);;`)
	match := re.FindStringSubmatch(mlm.logic)
	if match != nil {
		mlm.Date = strings.Trim(match[1], " ")
	}
}

// parse usage
func (mlm *MedicalLogicModule) parseUsage() {
	evoked_events := (regexp.MustCompile((`(?isU)([^\s]+)\s*;`))).FindAllStringSubmatch(mlm.evoke, -1)

	if evoked_events == nil {
		mlm.Usage = 2 // no evoked events
		return
	}

	for _, evoked_event := range evoked_events {
		x := evoked_event[1]
		mlm.evoked_events = append(mlm.evoked_events, x)

		if regexp.MustCompile(`(?isU)` + x + `\s*:=\s*event\s*{\s*ActivateApplication User UserInfo.*}`).MatchString(mlm.logic) {
			mlm.Usage = 1 // 1 = activate application
			return
		}
	}

	mlm.Usage = 0 // event present but not activate application
}

// func (mlm *MedicalLogicModule) String() string {
// 	return fmt.Sprintf("%+v\n", mlm)
// }

func (mlm *MedicalLogicModule) Initialize() {
	mlm.parseLogic()
	mlm.parseEvokeSlot()

	mlm.parseTitle()
	mlm.parseName()
	mlm.parseUsage()
	mlm.parseArden()
	mlm.parseVersion()
	mlm.parseDate()
	mlm.Status = 4 // 4 = active
}

func New(content string) *MedicalLogicModule {
	mlm := &MedicalLogicModule{
		Content: strings.ReplaceAll(content, "{{{SINGLE-QUOTE}}}}}", "'"),
	}
	mlm.Initialize()

	return mlm
}
