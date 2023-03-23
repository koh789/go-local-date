package nulldates

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/koh789/go-local-date/dates"
)

// LocalDate nullable localDate
type LocalDate struct {
	LocalDate dates.LocalDate
	Valid     bool
}

// Value for go-sql-driver
func (d LocalDate) Value() (driver.Value, error) {
	if d.Valid {
		return d.LocalDate.Value()
	}
	return nil, nil
}

// Scan for go-sql-driver
func (d *LocalDate) Scan(value interface{}) error {
	if d == nil || value == nil {
		d.LocalDate, d.Valid = dates.LocalDate{}, false
		return nil
	}
	scanErr := d.LocalDate.Scan(value)
	if scanErr != nil {
		d.Valid = false
	} else {
		d.Valid = true
	}
	return scanErr
}

// MarshalJSON for json return format: yyyy-MM-dd
func (d LocalDate) MarshalJSON() ([]byte, error) {
	if d.Valid {
		return d.LocalDate.MarshalJSON()
	}
	return json.Marshal(nil)
}

// UnmarshalJSON for json default format yyyy-MM-dd
func (d *LocalDate) UnmarshalJSON(data []byte) error {
	if d == nil {
		return fmt.Errorf("%w: nulldates.LocalDate. receiver is nil", dates.ErrUnmarshalJSON)
	}
	if len(data) == 0 || strings.EqualFold(string(data), "null") {
		d.LocalDate, d.Valid = dates.LocalDate{}, false
		return nil
	}
	if err := d.LocalDate.UnmarshalJSON(data); err != nil {
		d.LocalDate, d.Valid = dates.LocalDate{}, false
		return err
	}
	d.Valid = true
	return nil
}

func (d *LocalDate) UnmarshalFlag(s string) error {
	if d == nil {
		return fmt.Errorf("%w: nulldates.LocalDate. receiver is nil", dates.ErrUnmarshalFlag)
	}
	if len(s) == 0 {
		d.LocalDate, d.Valid = dates.LocalDate{}, false
		return nil
	}
	if err := d.LocalDate.UnmarshalFlag(s); err != nil {
		d.LocalDate, d.Valid = dates.LocalDate{}, false
		return err
	}
	d.Valid = true
	return nil
}

// NewLocalDate new LocalDate
func NewLocalDate(year, month, day int) LocalDate {
	localDate := dates.NewLocalDate(year, month, day)
	if localDate.IsZero() {
		return LocalDate{Valid: false}
	}
	return LocalDate{LocalDate: localDate, Valid: true}
}

func LocalDateFromDate(d dates.LocalDate) LocalDate {
	if d.IsZero() {
		return LocalDate{Valid: false}
	}
	return LocalDate{LocalDate: d, Valid: true}
}

func LocalDateFromPtr(d *dates.LocalDate) LocalDate {
	if d == nil {
		return LocalDate{Valid: false}
	}
	return LocalDateFromDate(*d)
}

// LocalDatetime nullable localDatetime
type LocalDatetime struct {
	LocalDatetime dates.LocalDatetime
	Valid         bool
}

// Value for go-sql-driver to value
func (dt LocalDatetime) Value() (driver.Value, error) {
	if dt.Valid {
		return dt.LocalDatetime.Value()
	}
	return nil, nil
}

// Scan for go-sql-driver
func (dt *LocalDatetime) Scan(value interface{}) error {
	if dt == nil || value == nil {
		dt.LocalDatetime, dt.Valid = dates.LocalDatetime{}, false
		return nil
	}
	scanErr := dt.LocalDatetime.Scan(value)
	if scanErr != nil {
		dt.Valid = false
	} else {
		dt.Valid = true
	}
	return scanErr
}

// MarshalJSON for json return format: yyyy-MM-dd hh:mm:ss
func (dt LocalDatetime) MarshalJSON() ([]byte, error) {
	if dt.Valid {
		return dt.LocalDatetime.MarshalJSON()
	}
	return json.Marshal(nil)
}

// UnmarshalJSON for json default format: yyyy-MM-dd hh:mm:ss
func (dt *LocalDatetime) UnmarshalJSON(data []byte) error {
	if dt == nil {
		return fmt.Errorf("%w: nulldates.LocalDatetime. receiver is nil", dates.ErrUnmarshalJSON)
	}
	if len(data) == 0 || strings.EqualFold(string(data), "null") {
		dt.LocalDatetime, dt.Valid = dates.LocalDatetime{}, false
		return nil
	}
	if err := dt.LocalDatetime.UnmarshalJSON(data); err != nil {
		dt.LocalDatetime, dt.Valid = dates.LocalDatetime{}, false
		return err
	}
	dt.Valid = true
	return nil
}

func (dt *LocalDatetime) UnmarshalFlag(s string) error {
	if dt == nil {
		return fmt.Errorf("%w: nulldates.LocalDatetime. receiver is nil", dates.ErrUnmarshalFlag)
	}
	if len(s) == 0 {
		dt.LocalDatetime, dt.Valid = dates.LocalDatetime{}, false
		return nil
	}
	if err := dt.LocalDatetime.UnmarshalFlag(s); err != nil {
		dt.LocalDatetime, dt.Valid = dates.LocalDatetime{}, false
		return err
	}
	dt.Valid = true
	return nil
}

// NewLocalDatetime new LocalDatetime
func NewLocalDatetime(year, month, day uint, hour, min, sec int) LocalDatetime {
	dtm := dates.NewLocalDatetime(year, month, day, hour, min, sec)
	if dtm.IsZero() {
		return LocalDatetime{Valid: false}
	}
	return LocalDatetime{LocalDatetime: dtm, Valid: true}
}

// LocalDatetimeFromTime time to LocalDatetime
func LocalDatetimeFromTime(t time.Time) LocalDatetime {
	if t.IsZero() {
		return LocalDatetime{Valid: false}
	}
	return LocalDatetime{
		LocalDatetime: dates.LocalDatetimeFromTime(t),
		Valid:         true,
	}
}

// LocalDatetimeFromTimePtr time pointer to LocalDatetime
func LocalDatetimeFromTimePtr(t *time.Time) LocalDatetime {
	if t == nil {
		return LocalDatetime{Valid: false}
	}
	return LocalDatetimeFromTime(*t)
}

// LocalDatetimeFromDate localDate to LocalDatetime
func LocalDatetimeFromDate(d dates.LocalDate) LocalDatetime {
	dtm := d.LocalDatetime()
	if dtm.IsZero() {
		return LocalDatetime{Valid: false}
	}
	return LocalDatetime{LocalDatetime: dtm, Valid: true}
}

// LocalDatetimeFromDatetime dtm to LocalDatetime
func LocalDatetimeFromDatetime(dt dates.LocalDatetime) LocalDatetime {
	if dt.IsZero() {
		return LocalDatetime{Valid: false}
	}
	return LocalDatetime{LocalDatetime: dt, Valid: true}
}

func LocalDatetimeFromPtr(dt *dates.LocalDatetime) LocalDatetime {
	if dt == nil {
		return LocalDatetime{Valid: false}
	}
	return LocalDatetimeFromDatetime(*dt)
}
