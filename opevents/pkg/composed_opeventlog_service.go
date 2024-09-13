package pkg

import (
	"fmt"
	"reflect"

	"go.uber.org/multierr"
	"go.uber.org/zap/zapcore"
)

type composedOpEventLogService struct {
	registedOpEventLogService []IOpEventLogService
}

var _ IOpEventLogService = (*composedOpEventLogService)(nil)

func newComposedOpEventLogService(underlyingService ...IOpEventLogService) *composedOpEventLogService {
	s := &composedOpEventLogService{
		registedOpEventLogService: make([]IOpEventLogService, 0),
	}

	s.registedOpEventLogService = append(s.registedOpEventLogService, underlyingService...)
	return s
}

func (s *composedOpEventLogService) registEventLogService(underlyingService ...IOpEventLogService) {
	s.registedOpEventLogService = append(s.registedOpEventLogService, underlyingService...)
}

func (s *composedOpEventLogService) Save(item *OpEventLog) error {
	if s.registedOpEventLogService == nil || len(s.registedOpEventLogService) <= 0 {
		return nil
	}
	var err error
	for _, eachService := range s.registedOpEventLogService {
		innerErr := eachService.Save(item)
		if innerErr != nil {
			innerErr = fmt.Errorf("err occur in {%v}.Save method,err:%v", reflect.TypeOf(eachService), innerErr)
			err = multierr.Append(err, innerErr)
		}
	}
	return err
}

func (s *composedOpEventLogService) SaveDebugOpEventLog(message string, opts ...EventLogOption) (err error) {
	if s.registedOpEventLogService == nil || len(s.registedOpEventLogService) <= 0 {
		return nil
	}
	for _, eachService := range s.registedOpEventLogService {
		innerErr := eachService.SaveDebugOpEventLog(message, opts...)
		if innerErr != nil {
			innerErr = fmt.Errorf("err occur in {%v}.SaveDebugOpEventLog method,err:%v", reflect.TypeOf(eachService), innerErr)
			err = multierr.Append(err, innerErr)
		}
	}
	return err
}

func (s *composedOpEventLogService) SaveWarnOpEventLog(message string, opts ...EventLogOption) (err error) {
	if s.registedOpEventLogService == nil || len(s.registedOpEventLogService) <= 0 {
		return nil
	}
	for _, eachService := range s.registedOpEventLogService {
		innerErr := eachService.SaveWarnOpEventLog(message, opts...)
		if innerErr != nil {
			innerErr = fmt.Errorf("err occur in {%v}.SaveWarnOpEventLog method,err:%v", reflect.TypeOf(eachService), innerErr)
			err = multierr.Append(err, innerErr)
		}
	}
	return err
}

func (s *composedOpEventLogService) SaveErrorOpEventLog(message string, opts ...EventLogOption) (err error) {
	if s.registedOpEventLogService == nil || len(s.registedOpEventLogService) <= 0 {
		return nil
	}
	for _, eachService := range s.registedOpEventLogService {
		innerErr := eachService.SaveErrorOpEventLog(message, opts...)
		if innerErr != nil {
			innerErr = fmt.Errorf("err occur in {%v}.SaveErrorOpEventLog method,err:%v", reflect.TypeOf(eachService), innerErr)
			err = multierr.Append(err, innerErr)
		}
	}
	return err
}

func (s *composedOpEventLogService) SaveOpEventLog(message string, logLevel zapcore.Level, opts ...EventLogOption) (err error) {
	if s.registedOpEventLogService == nil || len(s.registedOpEventLogService) <= 0 {
		return nil
	}
	for _, eachService := range s.registedOpEventLogService {
		innerErr := eachService.SaveOpEventLog(message, logLevel, opts...)
		if innerErr != nil {
			innerErr = fmt.Errorf("err occur in {%v}.SaveErrorOpEventLog method,err:%v", reflect.TypeOf(eachService), innerErr)
			err = multierr.Append(err, innerErr)
		}
	}
	return err
}
