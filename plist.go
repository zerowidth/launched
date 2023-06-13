package main

import (
	"encoding/base64"
	"encoding/json"
	"net/url"

	"github.com/go-playground/locales/en"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
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
	Name              string `json:"name,omitempty" form:"name" validate:"required"`
	Command           string `json:"command,omitempty" form:"command" validate:"required"`
	StartInterval     string `json:"start_interval,omitempty" form:"start_interval" validate:"omitempty,number"`
	Minute            string `json:"minute,omitempty" form:"minute" validate:"minute-cron"`
	Hour              string `json:"hour,omitempty" form:"hour" validate:"hour-cron"`
	DayOfMonth        string `json:"day_of_month,omitempty" form:"day_of_month" validate:"day-of-month-cron"`
	Month             string `json:"month,omitempty" form:"month" validate:"month-cron"`
	Weekday           string `json:"weekday,omitempty" form:"weekday" validate:"weekday-cron"`
	RunAtLoad         string `json:"run_at_load" form:"run_at_load"`
	RestartOnCrash    string `json:"restart_on_crash" form:"restart_on_crash"`
	StartOnMount      string `json:"start_on_mount" form:"start_on_mount"`
	QueueDirectories  string `json:"queue_directories,omitempty" form:"queue_directories"`
	Environment       string `json:"environment,omitempty" form:"environment"`
	User              string `json:"user,omitempty" form:"user"`
	Group             string `json:"group,omitempty" form:"group"`
	WorkingDirectory  string `json:"working_directory,omitempty" form:"working_directory"`
	StandardOutPath   string `json:"standard_out_path,omitempty" form:"standard_out_path"`
	StandardErrorPath string `json:"standard_error_path,omitempty" form:"standard_error_path"`
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

func (p LaunchdPlist) Encode() string {
	encoded, _ := json.Marshal(p)
	return base64.RawURLEncoding.EncodeToString(encoded)
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
