package validators

import (
	"context"
	"fmt"
	"testing"

	"github.com/aliparlakci/country-roads/models"
)

func TestUserValidator(t *testing.T) {
	tests := []struct {
		Case     models.NewUserForm
		Expected bool
	}{
		{
			Case:     models.NewUserForm{DisplayName: "Ali Parlakçı", Email: "parlakciali@sabanciuniv.edu", Phone: "+905420000000"},
			Expected: true,
		},
		{
			Case:     models.NewUserForm{DisplayName: "Ali Parlakçı", Email: "parlakciali@sabanciuniv.edu", Phone: "5420000000"},
			Expected: false,
		},
		{
			Case:     models.NewUserForm{DisplayName: "", Email: "aliparlakci@sabanciuniv.edu", Phone: "+905420000000"},
			Expected: false,
		},
		{
			Case:     models.NewUserForm{DisplayName: "Ali Parlakçı", Email: "aliparlakci@sabanciuniv", Phone: "+905420000000"},
			Expected: false,
		},
		{
			Case:     models.NewUserForm{DisplayName: "Ali Parlakçı", Email: "@sabanciuniv.edu", Phone: "+905420000000"},
			Expected: false,
		},
	}

	validatorFactory := ValidatorFactory{LocationFinder: nil}

	for _, tt := range tests {
		testName := fmt.Sprintf("%s__%s__%s", tt.Case.DisplayName, tt.Case.Email, tt.Case.Phone)
		t.Run(testName, func(t *testing.T) {
			userValidator, err := validatorFactory.GetValidator("users")
			if err != nil {
				t.Fatal(err)
			}

			if err := userValidator.SetDto(&tt.Case); err != nil {
				t.Fatal(err)
			}

			result, err := userValidator.Validate(context.TODO())
			if result != tt.Expected {
				t.Errorf("expected %v, got %v", tt.Expected, err)
			} else if tt.Expected == true && err != nil {
				t.Errorf("expected %v, got %v", tt.Expected, result)
			}
		})
	}
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		Case     string
		Expected bool
	}{
		{Case: "+1325420000000", Expected: true},
		{Case: "+905420000000", Expected: true},
		{Case: "05423530000", Expected: true},
		{Case: "asd\n+90 5423530000\nasd", Expected: false},
		{Case: "+5423530000", Expected: false},
		{Case: "12345", Expected: false},
		{Case: "aaaa", Expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			if result := ValidatePhone(tt.Case); result != tt.Expected {
				t.Errorf("expected %v, got %v", tt.Expected, result)
			}
		})
	}
}

func TestValidateDisplayName(t *testing.T) {
	tests := []struct {
		Case     string
		Expected bool
	}{
		{Case: "ali parlakci", Expected: true},
		{Case: "", Expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			if result := ValidateDisplayName(tt.Case); result != tt.Expected {
				t.Errorf("expected %v, got %v", tt.Expected, result)
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		Case          string
		ExpectedBool  bool
		ExpectedEmail string
	}{
		{Case: "aliparlakci@sabanciuniv.edu", ExpectedEmail: "aliparlakci@sabanciuniv.edu", ExpectedBool: true},
		{Case: "parlakciali@sabanciuniv.edu", ExpectedEmail: "parlakciali@sabanciuniv.edu", ExpectedBool: true},
		{Case: " aliparlakci@sabanciuniv.edu ", ExpectedEmail: "aliparlakci@sabanciuniv.edu", ExpectedBool: true},
		{Case: "parlakciali@gmail.com", ExpectedBool: false},
		{Case: "aliparlakci@sabanciuniv", ExpectedBool: false},
		{Case: "aliparlakcisabanciuniv", ExpectedBool: false},
		{Case: "aliparlakci@sabanciuniv.edu.tr", ExpectedBool: false},
	}

	for _, tt := range tests {
		t.Run(tt.Case, func(t *testing.T) {
			if email, success := ValidateEmail(tt.Case); success != tt.ExpectedBool {
				t.Errorf("expected %v, got %v", tt.ExpectedBool, success)
			} else {
				if email != tt.ExpectedEmail {
					t.Errorf("expected %v, got %v", tt.ExpectedEmail, email)
				}
			}
		})
	}
}
