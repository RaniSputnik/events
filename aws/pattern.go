package aws

type Event struct {
	Source     string `json:"source"`
	DetailType string `json:"detail-type"`
}

type EventPattern struct {
	Source     []string `json:"source"`
	DetailType []string `json:"detail-type"`
}

func EventMatches(ev Event, pattern EventPattern) bool {
	if len(pattern.DetailType) > 0 {
		if !stringArrayContains(pattern.DetailType, ev.DetailType) {
			return false
		}
	}
	if len(pattern.Source) > 0 {
		if !stringArrayContains(pattern.Source, ev.Source) {
			return false
		}
	}
	return true
}

func stringArrayContains(a []string, search string) bool {
	for _, str := range a {
		if str == search {
			return true
		}
	}
	return false
}
