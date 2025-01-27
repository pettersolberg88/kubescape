package policyhandler

import (
	"fmt"

	"github.com/armosec/kubescape/cautils/armotypes"
	"github.com/armosec/kubescape/cautils/opapolicy"
)

func (policyHandler *PolicyHandler) GetPoliciesFromBackend(notification *opapolicy.PolicyNotification) ([]opapolicy.Framework, []armotypes.PostureExceptionPolicy, error) {
	var errs error
	frameworks := []opapolicy.Framework{}
	exceptionPolicies := []armotypes.PostureExceptionPolicy{}

	// Get - cacli opa get
	for _, rule := range notification.Rules {
		switch rule.Kind {
		case opapolicy.KindFramework:
			receivedFramework, recExceptionPolicies, err := policyHandler.getFrameworkPolicies(rule.Name)
			if err != nil {
				errs = fmt.Errorf("%v\nKind: %v, Name: %s, error: %s", errs, rule.Kind, rule.Name, err.Error())
			}
			if receivedFramework != nil {
				frameworks = append(frameworks, *receivedFramework)
				if recExceptionPolicies != nil {
					exceptionPolicies = append(exceptionPolicies, recExceptionPolicies...)
				}
			}

		default:
			err := fmt.Errorf("missing rule kind, expected: %s", opapolicy.KindFramework)
			errs = fmt.Errorf("%s", err.Error())
		}
	}
	return frameworks, exceptionPolicies, errs
}

func (policyHandler *PolicyHandler) getFrameworkPolicies(policyName string) (*opapolicy.Framework, []armotypes.PostureExceptionPolicy, error) {
	receivedFramework, err := policyHandler.getters.PolicyGetter.GetFramework(policyName)
	if err != nil {
		return nil, nil, err
	}

	receivedException, err := policyHandler.getters.ExceptionsGetter.GetExceptions("", "")
	if err != nil {
		return receivedFramework, nil, err
	}

	return receivedFramework, receivedException, nil
}
