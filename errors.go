package hyperion

import "fmt"

func appManagerTypeUnknownError(appManagerType string) error {
	return fmt.Errorf("Unknown app manager type: %q. Valid options are: %+v",
		appManagerType, gManagerTypes)
}

func appManagerTypeAlreadyRegisteredError(appManagerType string) error {
	return fmt.Errorf("Already registered app manager type: %q.",
		appManagerType)
}
