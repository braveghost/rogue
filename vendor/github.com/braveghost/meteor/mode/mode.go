package mode

import (
	"github.com/pkg/errors"
)

var (
	ModeSetNameError = errors.New("mode setting error, ['pro','dev','local','sta','test']")
	ModeSetNumError  = errors.New("mode setting error, [1,2,3,3,5]")
	ModeNameError    = errors.New("mode name error, ['pro','dev','local','sta','test']")
)

var (
	modeName = map[ModeType]string{
		ModeLocal: "local",
		ModeDev:   "dev",
		ModeTest:  "test",
		ModeStage: "stage",
		ModePro:   "pro",
	}

	modeValue = map[string]ModeType{
		"local": ModeLocal,
		"dev":   ModeDev,
		"test":  ModeTest,
		"stage": ModeStage,
		"pro":   ModePro,
	}
)

type ModeType int

func (mt ModeType) String() string {
	return modeName[mt]
}

func (mt ModeType) Int() int {
	return int(mt)
}

func (mt ModeType) Int32() int32 {
	return int32(mt)
}

func (mt ModeType) Int64() int64 {
	return int64(mt)
}

const (
	ModeLocal ModeType = iota + 1 // localhost
	ModeDev                       // development
	ModeTest                      // test
	ModeStage                     // stage
	ModePro                       // production
)

var (
	defaultMode = ModeLocal
)

func IsOnline(md ModeType) bool {
	return md != ModeLocal
}

func SetModeByNum(m int) error {
	tmp := ModeType(m)

	if _, ok := modeName[tmp]; ! ok {
		return ModeSetNumError
	}
	defaultMode = tmp
	return nil
}

func SetModeByName(name string) error {
	m, ok := modeValue[name]
	if ! ok {
		return ModeSetNameError
	}
	defaultMode = m
	return nil
}

func GetDefaultMode() ModeType {
	return defaultMode
}

func GetModeByName(name string) (ModeType, error) {
	if m, ok := modeValue[name]; ! ok {
		return m, nil
	}
	return -1, ModeNameError
}

func GetModeByNum(m int) (ModeType, error) {
	tmp := ModeType(m)
	if _, ok := modeName[tmp]; ! ok {
		return tmp, nil
	}
	return -1, ModeNameError
}
