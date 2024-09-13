package pkg

type EventLogOption func(*OpEventLog)

func WithOpAction(opAction string) EventLogOption {
	return func(oel *OpEventLog) {
		oel.OPAction = opAction
	}
}

func WithDeviceMobile(deviceMobileNo string) EventLogOption {
	return func(oel *OpEventLog) {
		oel.DeviceMobileNo = deviceMobileNo
	}
}

func WithAndroidId(androidId string) EventLogOption {
	return func(oel *OpEventLog) {
		oel.AndroidId = androidId
	}
}
