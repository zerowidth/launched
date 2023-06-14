package main

import (
	"bytes"
	"encoding/json"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	"howett.net/plist"
)

var validate = validator.New()
var translate ut.Translator

func init() {
	validate.RegisterValidation("minute-cron", validateMinuteCron)
	validate.RegisterValidation("hour-cron", validateHourCron)
	validate.RegisterValidation("day-of-month-cron", validateDayOfMonthCron)
	validate.RegisterValidation("month-cron", validateMonthCron)
	validate.RegisterValidation("weekday-cron", validateWeekdayCron)

	en := en.New()
	uni := ut.New(en, en)
	translate, _ = uni.GetTranslator("en")
	en_translations.RegisterDefaultTranslations(validate, translate)

	for _, tag := range []string{"minute-cron", "hour-cron", "day-of-month-cron", "month-cron", "weekday-cron"} {
		validate.RegisterTranslation(tag, translate,
			func(ut ut.Translator) error {
				return ut.Add(tag, "is not a valid cron expression", true)
			},
			func(ut ut.Translator, fe validator.FieldError) string {
				t, _ := ut.T(tag, fe.Field())
				return t
			})
	}
	validate.RegisterTranslation("required", translate,
		func(ut ut.Translator) error {
			return ut.Add("required", "is required", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("required", fe.Field())
			return t
		})
	validate.RegisterTranslation("number", translate,
		func(ut ut.Translator) error {
			return ut.Add("number", "must be a number", true)
		},
		func(ut ut.Translator, fe validator.FieldError) string {
			t, _ := ut.T("number", fe.Field())
			return t
		})
}

// JSON encoded launchd plist, for encoded use in URL path.
// Uses string types to easily allow for empty values.
type LaunchdPlist struct {
	ID                string // set when loaded
	Name              string `redis:"name" form:"name" validate:"required"`
	Command           string `redis:"command" form:"command" validate:"required"`
	StartInterval     string `redis:"start_interval,omitempty" form:"start_interval" validate:"omitempty,number"`
	Minute            string `redis:"minute,omitempty" form:"minute" validate:"minute-cron"`
	Hour              string `redis:"hour,omitempty" form:"hour" validate:"hour-cron"`
	DayOfMonth        string `redis:"day_of_month,omitempty" form:"day_of_month" validate:"day-of-month-cron"`
	Month             string `redis:"month,omitempty" form:"month" validate:"month-cron"`
	Weekday           string `redis:"weekday,omitempty" form:"weekday" validate:"weekday-cron"`
	RunAtLoad         string `redis:"run_at_load,omitempty" form:"run_at_load"`
	RestartOnCrash    string `redis:"restart_on_crash,omitempty" form:"restart_on_crash"`
	StartOnMount      string `redis:"start_on_mount,omitempty" form:"start_on_mount"`
	QueueDirectories  string `redis:"queue_directories,omitempty" form:"queue_directories"`
	Environment       string `redis:"environment,omitempty" form:"environment"`
	User              string `redis:"user,omitempty" form:"user"`
	Group             string `redis:"group,omitempty" form:"group"`
	WorkingDirectory  string `redis:"working_directory,omitempty" form:"working_directory"`
	RootDirectory     string `redis:"root_directory,omitempty" form:"root_directory"`
	StandardOutPath   string `redis:"standard_out_path,omitempty" form:"standard_out_path"`
	StandardErrorPath string `redis:"standard_error_path,omitempty" form:"standard_error_path"`
	CreatedAt         string `redis:"created_at"` // written to when stored
}

func NewPlistFromForm(values url.Values) LaunchdPlist {
	plist := LaunchdPlist{}
	_ = decoder.Decode(&plist, values)
	return plist
}

func (p LaunchdPlist) JSONIndent() string {
	encoded, _ := json.MarshalIndent(p, "", "  ")
	return string(encoded)
}

func (p LaunchdPlist) PlistXML() string {

	buf := new(bytes.Buffer)
	xml := map[string]interface{}{
		"Label":            p.Label(),
		"ProgramArguments": []string{"sh", "-c", p.Command},
	}
	if p.StartInterval != "" {
		start, _ := strconv.Atoi(p.StartInterval)
		xml["StartInterval"] = start
	}
	crons := p.CronIntervals()
	if len(crons) > 0 {
		xml["StartCalendarInterval"] = crons
	}
	if p.RunAtLoad != "" {
		xml["RunAtLoad"] = true
	}
	if p.RestartOnCrash != "" {
		xml["KeepAlive"] = map[string]interface{}{
			"Crashed": true,
		}
	}
	if p.StartOnMount != "" {
		xml["StartOnMount"] = true
	}
	if p.QueueDirectories != "" {
		xml["QueueDirectories"] = strings.Split(p.QueueDirectories, ",")
	}
	if p.Environment != "" {
		xml["EnvironmentVariables"] = p.EnvironmentMap()
	}
	if p.User != "" {
		xml["UserName"] = p.User
	}
	if p.Group != "" {
		xml["GroupName"] = p.Group
	}
	if p.WorkingDirectory != "" {
		xml["WorkingDirectory"] = p.WorkingDirectory
	}
	if p.RootDirectory != "" {
		xml["RootDirectory"] = p.RootDirectory
	}
	if p.StandardOutPath != "" {
		xml["StandardOutPath"] = p.StandardOutPath
	}
	if p.StandardErrorPath != "" {
		xml["StandardErrorPath"] = p.StandardErrorPath
	}

	encoder := plist.NewEncoder(buf)
	encoder.Indent("  ")
	encoder.Encode(xml)
	return buf.String()
}

func (p LaunchdPlist) Label() string {
	whitespace := regexp.MustCompile(`\s+`)
	return "launched." + whitespace.ReplaceAllString(strings.ToLower(p.Name), "_")
}

func (p LaunchdPlist) CronIntervals() []map[string]int {
	return GenerateCronIntervals(p.Minute, p.Hour, p.DayOfMonth, p.Month, p.Weekday)
}

func (p LaunchdPlist) EnvironmentMap() map[string]string {
	env := map[string]string{}
	for _, line := range strings.Split(p.Environment, "\r\n") {
		if left, right, found := strings.Cut(line, "="); found {
			env[left] = right
		}
	}
	return env
}

func (p LaunchdPlist) Validate() validator.ValidationErrorsTranslations {
	err := validate.Struct(p)
	if err == nil {
		return nil
	}
	validationErrors := err.(validator.ValidationErrors)
	return validationErrors.Translate(translate)
}

func validateMinuteCron(fl validator.FieldLevel) bool {
	input, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if input == "" {
		return true
	}
	return ValidateCronExpression(input, 0, 59)
}

func validateHourCron(fl validator.FieldLevel) bool {
	input, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if input == "" {
		return true
	}
	return ValidateCronExpression(input, 0, 23)
}

func validateDayOfMonthCron(fl validator.FieldLevel) bool {
	input, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if input == "" {
		return true
	}
	return ValidateCronExpression(input, 1, 31)
}

func validateMonthCron(fl validator.FieldLevel) bool {
	input, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if input == "" {
		return true
	}
	return ValidateCronExpression(input, 1, 12)
}

func validateWeekdayCron(fl validator.FieldLevel) bool {
	input, ok := fl.Field().Interface().(string)
	if !ok {
		return false
	}
	if input == "" {
		return true
	}
	return ValidateCronExpression(input, 0, 6)
}
