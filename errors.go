package hyperion

import "fmt"

func unknownAppManagerTypeError(appManagerType ManagerType) error {
	return fmt.Errorf("Unknown app manager type: %q. Valid options are: %+v", appManagerType, ManagerTypes)
}
