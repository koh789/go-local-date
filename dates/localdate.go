package dates

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"
)

// time const
const (
	LocalDateRegex                   = "^(\\d{4})-(\\d{1,2})-(\\d{1,2})"
	LocalDateTimeRegex               = "^(\\d{4})-(\\d{1,2})-(\\d{1,2})\\ (\\d{1,2}):(\\d{1,2}):(\\d{1,2})"
	MaxYear            uint          = 999999999 //mysqlの最大値とは異なるため注意.
	MinYear            uint          = 0
	MaxMonthOfYear     uint          = 12
	MinMonthOfYear     uint          = 1
	MaxDayOfMonth      uint          = 31
	MinDayOfMonth      uint          = 1
	MaxHourOfDay       uint          = 23
	MinHourOfDay       uint          = 1
	MaxMinuteOfHour    uint          = 59
	MinMinuteOfHour    uint          = 0
	MaxSecOfMinute     uint          = 59
	MinSecOfMinute     uint          = 0
	MinDuration        time.Duration = -1 << 63
	MaxDuration        time.Duration = 1<<63 - 1
	FirstUnixInAD      int64         = -62135596800
)

// format list
const (
	Month           = Format("200601")
	MonthHyphen     = Format("2006-01")
	DateAbbreviated = Format("060102")
	Date            = Format("20060102")
	DateHyphen      = Format("2006-01-02")
	DateSlash       = Format("2006/01/02")
	DateHour        = Format("2006010215")
	DateHourHyphen  = Format("2006-01-02 15")
	DateHourSlash   = Format("2006/01/02 15")
	DateRFC3339     = Format("2006-01-02 Z07:00")
	DateTime        = Format("20060102150405")
	DateTimeHyphen  = Format("2006-01-02 15:04:05")
	DateTimeSlash   = Format("2006/01/02 15:04:05")
	ANSIC           = Format(time.ANSIC)
	UnixDate        = Format(time.UnixDate)
	RubyDate        = Format(time.RubyDate)
	RFC822          = Format(time.RFC822)
	RFC822Z         = Format(time.RFC822Z)
	RFC850          = Format(time.RFC850)
	RFC1123         = Format(time.RFC1123)
	RFC1123Z        = Format(time.RFC1123Z)
	RFC3339         = Format(time.RFC3339)
	RFC3339Nano     = Format(time.RFC3339Nano)
)

// Format type alias format
type Format string

// String to string
func (f Format) String() string {
	return string(f)
}

var (
	ErrFutureThanEndDate     = errors.New("start is a date future than end")
	ErrFutureOrSameDateAsEnd = errors.New("start is in the future or on the same date as end")
	ErrEmptyDate             = errors.New("date is empty")
	ErrLengthOfPeriod        = errors.New("too long or too short a period")
	ErrIncorrectDivisionDays = errors.New("the number of days specified is less than or equal to 0")
	ErrOutOfRangeDate        = errors.New("date out of range")
	ErrScan                  = errors.New("failed to scan")
	ErrMarshalJSON           = errors.New("failed to marshal json")
	ErrUnmarshalJSON         = errors.New("failed to unmarshal json")
	ErrUnmarshalFlag         = errors.New("failed to unmarshal flag")
	ErrParse                 = errors.New("failed to parse")
)

// Timezone timezone type
type Timezone int

// timezone enums
const (
	UTC Timezone = iota
	AsiaTokyo
)

var _TimezoneZoneIDMap = map[Timezone]string{
	UTC:       "UTC",
	AsiaTokyo: "Asia/Tokyo",
}
var _TimezoneOffsetMap = map[Timezone]int{
	UTC:       0 * 60 * 60,
	AsiaTokyo: +9 * 60 * 60,
}

// ZoneID to zoneID
func (tz Timezone) ZoneID() string {
	return _TimezoneZoneIDMap[tz]
}

// Offset to offset
func (tz Timezone) Offset() int {
	return _TimezoneOffsetMap[tz]
}

// Location to location
func (tz Timezone) Location() *time.Location {
	return time.FixedZone(tz.ZoneID(), tz.Offset())
}

func MarshalJSON(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, ErrMarshalJSON
	}
	return b, nil
}

// LocalDate represents a date without a timezone
type LocalDate struct {
	Year  uint
	Month uint
	Day   uint
}

// Valid valid localDate
func (d LocalDate) Valid() (LocalDate, error) {
	if d.Year < MinYear || MaxYear < d.Year {
		return d, fmt.Errorf("%w: localDate out of range ! year :%d", ErrOutOfRangeDate, d.Year)
	}
	if d.Month < MinMonthOfYear || MaxMonthOfYear < d.Month {
		return d, fmt.Errorf("%w: localDate out of range ! month: %d", ErrOutOfRangeDate, d.Month)
	}
	if d.Day < MinDayOfMonth || MaxDayOfMonth < d.Day {
		return d, fmt.Errorf("%w: localDate out of range ! day: %d", ErrOutOfRangeDate, d.Day)
	}
	return d, nil
}

// Value for go-sql-driver
func (d LocalDate) Value() (driver.Value, error) {
	year, month, day := d.SplitString()
	return year + "-" + month + "-" + day, nil
}

// Scan for go-sql-driver
func (d *LocalDate) Scan(value interface{}) error {
	if d == nil || value == nil {
		return fmt.Errorf("%w: nil value %v", ErrScan, value)
	}
	if sv, ce := driver.String.ConvertValue(value); ce == nil {
		if v, ok := sv.(string); ok {
			groups, ge := groupSubMatch(v, LocalDateRegex)
			if ge != nil {
				return fmt.Errorf("%w: failed to convert localDate! %v", ErrScan, ge.Error())
			} else if len(groups) < 4 {
				return fmt.Errorf("%w: failed to convert localDate! ( in grouping ) len: %d", ErrScan, len(groups))
			}
			year, ye := strconv.Atoi(groups[1])
			month, me := strconv.Atoi(groups[2])
			day, de := strconv.Atoi(groups[3])
			if ye != nil || me != nil || de != nil {
				return fmt.Errorf("%w: failed to convert localDate! groups [ %s, %s, %s ]", ErrScan, groups[1], groups[2], groups[3])
			}
			*d = LocalDate{Year: uint(year), Month: uint(month), Day: uint(day)}
			return nil
		}
	}
	return fmt.Errorf("%w: failed to scan localDate", ErrScan)
}

func groupSubMatch(target, regex string) ([]string, error) {
	reg, err := regexp.Compile(regex)
	if err != nil {
		return make([]string, 0), err
	}
	return reg.FindStringSubmatch(target), nil
}

// String to string
func (d LocalDate) String() string {
	val, _ := d.Value()
	return val.(string)
}

// SplitString returns year, month and day with the same number of digits.
//
//	(year,month,day)= (2000,1,3) ---> "2000", "01", "03"
func (d LocalDate) SplitString() (yearStr, monthStr, dayStr string) {
	var year, month, day string
	switch {
	case d.Year < 10:
		year = "000" + strconv.Itoa(int(d.Year))
	case d.Year < 100:
		year = "00" + strconv.Itoa(int(d.Year))
	case d.Year < 1000:
		year = "0" + strconv.Itoa(int(d.Year))
	default:
		year = strconv.Itoa(int(d.Year))
	}
	if d.Month < 10 {
		month = "0" + strconv.Itoa(int(d.Month))
	} else {
		month = strconv.Itoa(int(d.Month))
	}

	if d.Day < 10 {
		day = "0" + strconv.Itoa(int(d.Day))
	} else {
		day = strconv.Itoa(int(d.Day))
	}
	return year, month, day
}

// ToTime convert to Time type based on location
func (d LocalDate) ToTime(loc *time.Location) time.Time {
	return time.Date(int(d.Year), time.Month(int(d.Month)), int(d.Day), 0, 0, 0, 0, loc)
}

// ToTimeUtc convert to Time type based on UTC.
func (d LocalDate) ToTimeUtc() time.Time {
	loc := UTC.Location()
	return d.ToTime(loc)
}

// Before reports whether the localDate instant t is before u.
func (d LocalDate) Before(targetDate LocalDate) bool {
	firstDate := d.ToTimeUtc()
	secondDate := targetDate.ToTimeUtc()

	return firstDate.Before(secondDate)
}

// After reports whether the localDate instant t is after u.
func (d LocalDate) After(targetDate LocalDate) bool {
	firstDate := d.ToTimeUtc()
	secondDate := targetDate.ToTimeUtc()
	return firstDate.After(secondDate)
}

// Equal localDate Equal?
func (d LocalDate) Equal(targetDate LocalDate) bool {
	return d.ToTimeUtc().Equal(targetDate.ToTimeUtc())
}

// Between localDate between ?
func (d LocalDate) Between(start, end LocalDate) bool {
	return (d.After(start) || d.Equal(start)) && (d.Equal(end) || d.Before(end))
}

// IsZero localDate is zero?
func (d LocalDate) IsZero() bool {
	return d.Year == 0 && d.Month == 0 && d.Day == 0
}

// LocalDatetime date to LocalDatetime
func (d LocalDate) LocalDatetime() LocalDatetime {
	return LocalDatetime{
		LocalDate: d,
		LocalTime: LocalTime{Hour: 0, Minute: 0, Second: 0},
	}
}

// Sub dtm - target 290年以上の期間は扱えません
func (d LocalDate) Sub(target LocalDate) (time.Duration, bool) {
	duration := d.ToTimeUtc().Sub(target.ToTimeUtc())
	if duration == MaxDuration || duration == MinDuration {
		return duration, false
	}
	return duration, true
}

func (d LocalDate) Add(duration time.Duration) LocalDate {
	time := d.ToTimeUtc()
	return LocalDateFromTime(time.Add(duration))
}

// MarshalJSON for json return format: yyyy-MM-dd
func (d LocalDate) MarshalJSON() ([]byte, error) {
	if d.IsZero() {
		return MarshalJSON(nil)
	}
	return MarshalJSON(d.String())
}

// UnmarshalJSON for json default format yyyy-MM-dd
func (d *LocalDate) UnmarshalJSON(data []byte) error {
	if d == nil || len(data) == 0 {
		return fmt.Errorf("%w: localDate, receiver is nil or data len is 0", ErrUnmarshalJSON)
	}
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("%w: failed to unmarshal localDate, err: %v", ErrUnmarshalJSON, err)
	}
	date, err := ParseLocalDate(DateHyphen, str)
	if err != nil {
		return fmt.Errorf("%w: failed to parse localDate, err: %v", ErrUnmarshalJSON, err)
	}
	*d = date
	return nil
}

func (d *LocalDate) UnmarshalFlag(s string) error {
	if d == nil || len(s) == 0 {
		return fmt.Errorf("%w: localDate. receiver is nil or data len is 0", ErrUnmarshalFlag)
	}
	date, err := ParseLocalDate(DateHyphen, s)
	if err != nil {
		return err
	}
	*d = date
	return nil
}

// NewLocalDate new localDate
func NewLocalDate(year, month, day int) LocalDate {
	tm := time.Date(year, time.Month(month), day, 0, 0, 0, 0, UTC.Location())
	if tm.Unix() < FirstUnixInAD {
		return LocalDate{}
	}
	return LocalDate{Year: uint(tm.Year()), Month: uint(tm.Month()), Day: uint(tm.Day())}
}

// ParseLocalDate parse localDate by string
func ParseLocalDate(f Format, t string) (LocalDate, error) {
	loc := UTC.Location() //localdateのため、このtimezoneは使用しない

	tm, err := time.ParseInLocation(f.String(), t, loc)
	if err != nil {
		return LocalDate{}, fmt.Errorf("%w: %v", ErrParse, err)
	}
	return LocalDateFromTime(tm), nil
}

type (
	LocalDatePeriod struct {
		Start LocalDate
		End   LocalDate
	}
	LocalDatePeriods []LocalDatePeriod
)

func DivideDatePeriod(start, end LocalDate, day int) (LocalDatePeriods, error) {
	if start.IsZero() || end.IsZero() {
		return nil, fmt.Errorf("%w: start:%v,end:%v", ErrEmptyDate, start, end)
	}
	if day < 1 {
		return nil, ErrIncorrectDivisionDays
	}
	duration, ok := end.Sub(start)
	if !ok {
		return nil, fmt.Errorf("%w: start:%v,end:%v", ErrLengthOfPeriod, start, end)
	}
	if duration <= time.Duration(day)*24*time.Hour {
		return LocalDatePeriods{{Start: start, End: end}}, nil
	}
	periods := make([]LocalDatePeriod, 0)
	var tmpStart LocalDate
	dayCounter, oneDay := 0, time.Duration(24)*time.Hour
	for d := start; ; d = d.Add(oneDay) {
		dayCounter++
		if dayCounter == 1 {
			tmpStart = d
		}
		if dayCounter == day {
			periods = append(periods, LocalDatePeriod{Start: tmpStart, End: d})
			// 日数, tmpStartをリセット
			dayCounter, tmpStart = 0, LocalDate{}
		}
		if d.After(end) || d.Equal(end) {
			periods = append(periods, LocalDatePeriod{Start: tmpStart, End: d})
			break
		}
	}
	return periods, nil
}

// LocalTime
type LocalTime struct {
	Hour   uint
	Minute uint
	Second uint
}

// Valid validate localTime
func (t LocalTime) Valid() (LocalTime, error) {
	if t.Hour < MinHourOfDay || MaxHourOfDay < t.Hour {
		return t, fmt.Errorf("%w: localTime hour :%d", ErrOutOfRangeDate, t.Hour)
	}
	if t.Minute < MinMinuteOfHour || MaxMinuteOfHour < t.Minute {
		return t, fmt.Errorf("%w: localTime minute: %d", ErrOutOfRangeDate, t.Minute)
	}
	if t.Second < MinSecOfMinute || MaxSecOfMinute < t.Second {
		return t, fmt.Errorf("%w: localTime second: %d", ErrOutOfRangeDate, t.Second)
	}
	return t, nil
}

// SplitString (hour, min, sec)= (12,1,3) ---> "12", "01", "03"
func (t LocalTime) SplitString() (hour, min, sec string) {
	var h, m, d string
	if t.Hour < 10 {
		h = "0" + strconv.Itoa(int(t.Hour))
	} else {
		h = strconv.Itoa(int(t.Hour))
	}
	if t.Minute < 10 {
		m = "0" + strconv.Itoa(int(t.Minute))
	} else {
		m = strconv.Itoa(int(t.Minute))
	}

	if t.Second < 10 {
		d = "0" + strconv.Itoa(int(t.Second))
	} else {
		d = strconv.Itoa(int(t.Second))
	}
	return h, m, d
}

// IsZero localTime is Zero
func (t LocalTime) IsZero() bool {
	return t.Hour == 0 && t.Minute == 0 && t.Second == 0
}

func LocalDateFromTime(tm time.Time) LocalDate {
	return LocalDate{
		Year:  uint(tm.Year()),
		Month: uint(tm.Month()),
		Day:   uint(tm.Day()),
	}
}

// NewLocalTime new localTime
func NewLocalTime(hour, min, sec uint) (LocalTime, error) {
	return LocalTime{Hour: hour, Minute: min, Second: sec}.Valid()
}

// LocalDatetime local datetime
type LocalDatetime struct {
	LocalDate LocalDate
	LocalTime LocalTime
}

// Value for go-sql-driver
func (dt LocalDatetime) Value() (driver.Value, error) {
	y, m, d := dt.LocalDate.SplitString()
	h, min, sec := dt.LocalTime.SplitString()
	return y + "-" + m + "-" + d + " " + h + ":" + min + ":" + sec, nil
}

// String localDatetime to string
func (dt LocalDatetime) String() string {
	val, _ := dt.Value()
	return val.(string)
}

// Scan for go-sql-driver
func (dt *LocalDatetime) Scan(value interface{}) error {
	if dt == nil || value == nil {
		return fmt.Errorf("%w: nil value %v", ErrScan, value)
	}
	if sv, ce := driver.String.ConvertValue(value); ce == nil {
		if v, ok := sv.(string); ok {
			groups, ge := groupSubMatch(v, LocalDateTimeRegex)
			if ge != nil {
				return fmt.Errorf("%w: localDatetime %v", ErrScan, ge)
			} else if len(groups) < 7 {
				return fmt.Errorf("%w: localDatetime (in grouping) len: %v", ErrScan, strconv.Itoa(len(groups)))
			}
			y, ye := strconv.Atoi(groups[1])
			m, me := strconv.Atoi(groups[2])
			d, de := strconv.Atoi(groups[3])
			h, he := strconv.Atoi(groups[4])
			min, minErr := strconv.Atoi(groups[5])
			sec, se := strconv.Atoi(groups[6])

			if ye != nil || me != nil || de != nil || he != nil || minErr != nil || se != nil {
				return fmt.Errorf("%w: localDatetime groups: %v", ErrScan, groups)
			}
			*dt = LocalDatetime{
				LocalDate: LocalDate{Year: uint(y), Month: uint(m), Day: uint(d)},
				LocalTime: LocalTime{Hour: uint(h), Minute: uint(min), Second: uint(sec)},
			}
			return nil
		}
	}
	return fmt.Errorf("%w: localDatetime", ErrScan)
}

// MarshalJSON for json return format: yyyy-MM-dd hh:mm:ss
func (dt LocalDatetime) MarshalJSON() ([]byte, error) {
	if dt.IsZero() {
		return MarshalJSON(nil)
	}
	return MarshalJSON(dt.String())
}

// UnmarshalJSON for json default format: yyyy-MM-dd hh:mm:ss
//
//	If you want to specify FMTs individually, add a custom Unmarshal receiver on the caller side as follows
//	type FmtLocalDatetime LocalDatetime
//	func (dt *LocalDatetime) UnmarshalJSON(data []byte) error {
//	   custom unmarshal
//	}
func (dt *LocalDatetime) UnmarshalJSON(data []byte) error {
	if dt == nil || len(data) == 0 {
		return fmt.Errorf("%w: localDatetime. receiver is nil or data len is 0", ErrUnmarshalJSON)
	}
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("%w: failed to unmarshal localDatetime. err: %v", ErrUnmarshalJSON, err)
	}
	time, err := time.ParseInLocation(DateTimeHyphen.String(), str, UTC.Location())
	if err != nil {
		return fmt.Errorf("%w: failed to parse localDatetime, err: %v", ErrUnmarshalJSON, err)
	}
	*dt = LocalDatetimeFromTime(time)
	return nil
}

func (dt *LocalDatetime) UnmarshalFlag(s string) error {
	datetime, err := ParseLocalDatetime(DateTimeHyphen, s)
	if err != nil {
		return err
	}
	*dt = datetime
	return nil
}

// ToTime convert to Time type based on Location.
func (dt LocalDatetime) ToTime(loc *time.Location) time.Time {
	return time.Date(int(dt.LocalDate.Year), time.Month(int(dt.LocalDate.Month)), int(dt.LocalDate.Day),
		int(dt.LocalTime.Hour), int(dt.LocalTime.Minute), int(dt.LocalTime.Second), 0, loc)
}

// ToTimeUtc UTCベースでTime型へ変換します.
func (dt LocalDatetime) ToTimeUtc() time.Time {
	loc := UTC.Location()
	return dt.ToTime(loc)
}

// Before localDatetime is before
func (dt LocalDatetime) Before(targetDateTime LocalDatetime) bool {
	firstDateTime := dt.ToTimeUtc()
	secondDateTime := targetDateTime.ToTimeUtc()
	return firstDateTime.Before(secondDateTime)
}

// After localDatetime is after
func (dt LocalDatetime) After(targetDateTime LocalDatetime) bool {
	firstDateTime := dt.ToTimeUtc()
	secondDateTime := targetDateTime.ToTimeUtc()

	return firstDateTime.After(secondDateTime)
}

// BeforeEqual localDatetime is before or equal
func (dt LocalDatetime) BeforeEqual(targetDateTime LocalDatetime) bool {
	return dt.Before(targetDateTime) || dt.Equal(targetDateTime)
}

// AfterEqual localDatetime is after or equal
func (dt LocalDatetime) AfterEqual(targetDateTime LocalDatetime) bool {
	return dt.After(targetDateTime) || dt.Equal(targetDateTime)
}

// Equal localDateme is equal
func (dt LocalDatetime) Equal(targetDateTime LocalDatetime) bool {
	return dt.ToTimeUtc().Equal(targetDateTime.ToTimeUtc())
}

// Between localDateime is between
func (dt LocalDatetime) Between(start, end LocalDatetime) bool {
	return (dt.After(start) || dt.Equal(start)) && (dt.Equal(end) || dt.Before(end))
}

// Sub dtm - target 290年以上の期間は扱えません
func (dt LocalDatetime) Sub(target LocalDatetime) (time.Duration, bool) {
	duration := dt.ToTimeUtc().Sub(target.ToTimeUtc())
	if duration == MaxDuration || duration == MinDuration {
		return duration, false
	}
	return duration, true
}

// Add localDateime add duration
func (dt LocalDatetime) Add(d time.Duration) LocalDatetime {
	loc := UTC.Location()
	addedTime := dt.ToTime(loc).Add(d)
	return LocalDatetimeFromTime(addedTime)
}

// AddDate localDateime add date
func (dt LocalDatetime) AddDate(year, month, day int) LocalDatetime {
	loc := UTC.Location()
	addedTime := dt.ToTime(loc).AddDate(year, month, day)
	return LocalDatetimeFromTime(addedTime)
}

// IsZero localDatetime is zero?
func (dt LocalDatetime) IsZero() bool {
	return dt.LocalDate.IsZero() && dt.LocalTime.IsZero()
}

// IsNotZero localDatetime is not zero?
func (dt LocalDatetime) IsNotZero() bool {
	return !dt.IsZero()
}

// NewLocalDatetime new localDatetime.
// If a number exceeding the maximum value is given, such as day: 40, a calendar calculation is performed and initialized.
// If the result of the calendar calculation is in BC, it will return empty.
func NewLocalDatetime(year, month, day uint, hour, min, sec int) LocalDatetime {
	tm := time.Date(int(year), time.Month(month), int(day), hour, min, sec, 0, UTC.Location())
	if tm.Unix() < FirstUnixInAD {
		return LocalDatetime{}
	}
	localTime := LocalTime{Hour: uint(tm.Hour()), Minute: uint(tm.Minute()), Second: uint(tm.Second())}
	localDate := LocalDate{Year: uint(tm.Year()), Month: uint(tm.Month()), Day: uint(tm.Day())}
	return LocalDatetime{LocalDate: localDate, LocalTime: localTime}
}

// LocalDatetimeFromTime converts from time to LocalDateTime. (tz is ignored)
func LocalDatetimeFromTime(tm time.Time) LocalDatetime {
	// conversion from time is not necessary because the error that occurs is due to the validity of the date.
	return NewLocalDatetime(uint(tm.Year()), uint(tm.Month()), uint(tm.Day()), tm.Hour(), tm.Minute(), tm.Second())
}

// NowLocalDatetimeJst now localDatetime jst
func NowLocalDatetimeJst() LocalDatetime {
	loc := AsiaTokyo.Location()
	return LocalDatetimeFromTime(time.Now().In(loc))
}

// NowLocalDatetimeUtc now localDateime utc
func NowLocalDatetimeUtc() LocalDatetime {
	return LocalDatetimeFromTime(time.Now().In(time.UTC))
}

// ParseLocalDatetime parse localDatetime by string
func ParseLocalDatetime(f Format, t string) (LocalDatetime, error) {
	loc := UTC.Location() //localdatetimeのため、このtimezoneは使用しない

	tm, err := time.ParseInLocation(f.String(), t, loc)
	if err != nil {
		return LocalDatetime{}, fmt.Errorf("%w: err: %v", ErrParse, err)
	}
	return LocalDatetimeFromTime(tm), nil
}
