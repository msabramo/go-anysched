package hyperion

import "fmt"

func appManagerTypeUnknownError(appManagerType string) error {
	return fmt.Errorf("unknown app manager type: %q. Valid options are: %+v",
		appManagerType, gManagerTypes)
}

func appManagerTypeAlreadyRegisteredError(appManagerType string) error {
	return fmt.Errorf("already registered app manager type: %q",
		appManagerType)
}
