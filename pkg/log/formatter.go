package log

import (
	"fmt"
	"sort"
	"strings"

	"github.com/sirupsen/logrus"
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	if entry == nil {
		return nil, nil
	}

	resourceType, ok := entry.Data["type"].(string)
	if !ok {
		return nil, nil
	}

	if _, ok := entry.Data["owner"]; !ok {
		return nil, nil
	}
	if _, ok := entry.Data["resource"]; !ok {
		return nil, nil
	}
	if _, ok := entry.Data["state"]; !ok {
		return nil, nil
	}

	owner := entry.Data["owner"].(string)
	resource := entry.Data["resource"].(string)
	state := entry.Data["state"].(int)

	var sortedFields = make([]string, 0)
	for k, v := range entry.Data {
		if strings.HasPrefix(k, "prop:") {
			sortedFields = append(sortedFields, fmt.Sprintf("%s: %q", k[5:], v))
		}
	}

	sort.Strings(sortedFields)

	msgColor := ReasonSuccess
	switch state {
	case 0, 1, 8:
		msgColor = ReasonSuccess
	case 2:
		msgColor = ReasonHold
	case 3:
		msgColor = ReasonRemoveTriggered
	case 4:
		msgColor = ReasonWaitDependency
	case 5:
		msgColor = ReasonWaitPending
	case 6:
		msgColor = ReasonError
	case 7:
		msgColor = ReasonSkip
	}

	msg := fmt.Sprintf("%s - %s - %s - %s - %s\n",
		ColorRegion.Sprint(owner),
		ColorResourceType.Sprint(resourceType),
		ColorResourceID.Sprint(resource),
		ColorResourceProperties.Sprintf("[%s]", strings.Join(sortedFields, ", ")),
		msgColor.Sprint(entry.Message))

	return []byte(msg), nil
}
