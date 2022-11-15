package syncs

import (
	"fmt"
	"strings"
	"sync"
)

type multierror struct {
	errs []error

	lock sync.RWMutex
}

func (m *multierror) errorOrNil() error {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if len(m.errs) == 0 {
		return nil
	}

	return m
}

func (m *multierror) Add(err error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.errs = append(m.errs, err)
}

func (m *multierror) Error() string {
	m.lock.RLock()
	defer m.lock.RUnlock()

	if len(m.errs) == 0 {
		return ""
	}

	errText := make([]string, len(m.errs))
	for i, err := range m.errs {
		errText[i] = fmt.Sprintf("[%d] %s", i, err.Error())
	}

	return fmt.Sprintf("multierror: %d errors: %s", len(m.errs), strings.Join(errText, "; "))
}
