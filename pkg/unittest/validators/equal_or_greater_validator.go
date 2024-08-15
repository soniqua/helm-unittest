package validators

import (
	"fmt"
	"reflect"

	"github.com/helm-unittest/helm-unittest/pkg/unittest/valueutils"
	log "github.com/sirupsen/logrus"
)

// EqualOrGreaterValidator validate whether the value of Path is greater or equal to Value
type EqualOrGreaterValidator struct {
	Path  string
	Value interface{}
}

func (a EqualOrGreaterValidator) failInfo(msg string, index int, not bool) []string {
	return splitInfof(
		setFailFormat(not, true, false, false, " to be greater then or equal to, got"),
		index,
		a.Path,
		msg,
	)
}

// Validate implement Validatable
func (a EqualOrGreaterValidator) Validate(context *ValidateContext) (bool, []string) {
	log.WithField("validator", "ge").Debugln("expected content:", a.Value, "path:", a.Path)
	manifests, err := context.getManifests()
	if err != nil {
		return false, splitInfof(errorFormat, -1, err.Error())
	}

	validateSuccess := false
	validateErrors := make([]string, 0)

	for idx, manifest := range manifests {
		actual, err := valueutils.GetValueOfSetPath(manifest, a.Path)
		if err != nil {
			errorMessage := splitInfof(errorFormat, idx, err.Error())
			validateErrors = append(validateErrors, errorMessage...)
			continue
		}

		if len(actual) == 0 {
			errorMessage := splitInfof(errorFormat, idx, fmt.Sprintf("unknown path '%s'", a.Path))
			validateErrors = append(validateErrors, errorMessage...)
			continue
		}

		actType := reflect.TypeOf(actual[0])
		expType := reflect.TypeOf(a.Value)

		if actType != expType {
			errorMessage := splitInfof(errorFormat, idx, fmt.Sprintf("actual '%s' and expected '%s' types do not match", actType, expType))
			validateErrors = append(validateErrors, errorMessage...)
			continue
		}

		result, errors := compareValues(a.Value, actual[0], "greater", !context.Negative)
		if errors != nil {
			errorMessage := a.failInfo(errors[0], idx, context.Negative)
			validateErrors = append(validateErrors, errorMessage...)
		}

		validateSuccess = determineSuccess(idx, validateSuccess, result)
	}

	return validateSuccess, validateErrors
}
