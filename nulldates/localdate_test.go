package nulldates

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/koh789/go-local-date/dates"
	. "github.com/stretchr/testify/assert"
)

type NullLocalDateStruct struct {
	LocalDate LocalDate `json:"local_date,omitempty"`
}

func TestLocalDate_MarshalJSON(t *testing.T) {
	{
		jsonStruct := NullLocalDateStruct{
			LocalDate: NewLocalDate(2020, 3, 4),
		}
		expect := `{"local_date":"2020-03-04"}`
		jsonBytes, err := json.Marshal(jsonStruct)
		Nil(t, err)
		Equal(t, expect, string(jsonBytes))
	}
	{
		jsonStruct := NullLocalDateStruct{
			LocalDate: LocalDate{Valid: false},
		}
		expect := `{"local_date":null}`
		jsonBytes, err := json.Marshal(jsonStruct)
		Nil(t, err)
		Equal(t, expect, string(jsonBytes), "omitemptyは効かない")
	}
	{
		var jsonStruct NullLocalDateStruct
		expect := `{"local_date":null}`
		jsonBytes, err := json.Marshal(jsonStruct)
		Nil(t, err)
		Equal(t, expect, string(jsonBytes))
	}
}

func TestLocalDate_UnmarshalJSON(t *testing.T) {
	{
		var jsonStruct NullLocalDateStruct
		jsonString := `{
			"local_date": "2020-07-19"
		}`
		expect := NullLocalDateStruct{LocalDate: NewLocalDate(2020, 7, 19)}
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		Nil(t, err)
		Equal(t, expect, jsonStruct)
	}
	{
		var jsonStruct NullLocalDateStruct
		jsonString := `{
			"local_date": "2020/07/19"
		}`
		jsonString2 := `{
		"local_date": "2020-07-19 00:00:00"
		}`
		expect := NullLocalDateStruct{}
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		NotNil(t, err)
		Equal(t, expect, jsonStruct, "invalid format")

		err = json.Unmarshal([]byte(jsonString2), &jsonStruct)
		NotNil(t, err)
		Equal(t, expect, jsonStruct, "invalid format")
	}
	{
		var jsonStruct NullLocalDateStruct
		jsonString := `{"local_date": null}`
		expect := NullLocalDateStruct{LocalDate: LocalDate{}}
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		Nil(t, err)
		Equal(t, expect, jsonStruct)
	}
	{
		var jsonStruct NullLocalDateStruct
		jsonString := `{}`
		expect := NullLocalDateStruct{LocalDate: LocalDate{}}
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		Nil(t, err)
		Equal(t, expect, jsonStruct)
	}
	{
		var jsonStruct NullLocalDateStruct
		jsonString := `{
			"local_date": "unknownType"
		}`
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		NotNil(t, err)
	}
}

func TestLocalDate_UnmarshalFlag(t *testing.T) {
	for _, table := range []struct {
		title         string
		localDate     *LocalDate
		target        string
		expect        LocalDate
		errorOccurred bool
	}{
		{
			title:         "LocalDate.unmarshalFlag.localDate instance is nil,error発生",
			localDate:     nil,
			target:        "2022-04-01",
			expect:        LocalDate{LocalDate: dates.LocalDate{}, Valid: false},
			errorOccurred: true,
		},
		{
			title:         "LocalDate.unmarshalFlag.target string localDate is empty,emptyがセットされるが,errorは発生しない",
			localDate:     &LocalDate{},
			target:        "",
			expect:        LocalDate{LocalDate: dates.LocalDate{}, Valid: false},
			errorOccurred: false,
		},
		{
			title:         "LocalDate.unmarshalFlag.target string localDate is unexpected format, error occurred",
			localDate:     &LocalDate{},
			target:        "2022/04/01",
			expect:        LocalDate{LocalDate: dates.LocalDate{}, Valid: false},
			errorOccurred: true,
		},
		{
			title:         "LocalDate.unmarshalFlag.target string localDate's format: yyyy-MM-dd HH:mm:ss, parse error occurred",
			localDate:     &LocalDate{},
			target:        "2022-04-01 14:15:16",
			expect:        LocalDate{LocalDate: dates.LocalDate{}, Valid: false},
			errorOccurred: true,
		},
		{
			title:         "LocalDate.unmarshalFlag.target string localDate is valid format.",
			localDate:     &LocalDate{},
			target:        "2022-04-01",
			expect:        LocalDate{LocalDate: dates.LocalDate{Year: 2022, Month: 4, Day: 1}, Valid: true},
			errorOccurred: false,
		},
	} {
		t.Run(table.title, func(t *testing.T) {
			err := table.localDate.UnmarshalFlag(table.target)
			if table.errorOccurred {
				NotNil(t, err)
			} else {
				Nil(t, err)
			}
			if table.localDate != nil {
				Equal(t, table.expect, *table.localDate)
			}
		})
	}
}

func TestLocalDateFromPtr(t *testing.T) {
	for _, table := range []struct {
		title  string
		date   *dates.LocalDate
		expect LocalDate
	}{
		{
			title:  "nilを渡すと, Valid:falseのLocalDateが返却される",
			date:   nil,
			expect: LocalDate{Valid: false},
		},
		{
			title:  "empty localDateを渡すと, Valid:falseのLocalDateが返却される",
			date:   &dates.LocalDate{},
			expect: LocalDate{Valid: false},
		},
		{
			title: "localDateを渡すと, Valid:trueのLocalDateが返却される",
			date:  &dates.LocalDate{Year: 2022, Month: 4, Day: 5},
			expect: LocalDate{
				Valid:     true,
				LocalDate: dates.LocalDate{Year: 2022, Month: 4, Day: 5},
			},
		},
		{
			title:  "zero値のlocalDateを渡すと, Valid:falseのLocalDateが返却される",
			date:   &dates.LocalDate{Year: 0, Month: 0, Day: 0},
			expect: LocalDate{Valid: false},
		},
	} {
		t.Run(table.title, func(t *testing.T) {
			actual := LocalDateFromPtr(table.date)
			Equal(t, table.expect, actual)
		})
	}
}

type LocalDatetimeStruct struct {
	LocalDatetime LocalDatetime `json:"local_datetime,omitempty"`
}

func TestLocalDatetime_MarshalJSON(t *testing.T) {
	{
		jsonStruct := LocalDatetimeStruct{
			LocalDatetime: NewLocalDatetime(2020, 3, 4, 13, 13, 14),
		}
		expect := `{"local_datetime":"2020-03-04 13:13:14"}`
		jsonBytes, err := json.Marshal(jsonStruct)
		Nil(t, err)
		Equal(t, expect, string(jsonBytes))
	}
	{
		jsonStruct := LocalDatetimeStruct{
			LocalDatetime: LocalDatetime{Valid: false},
		}
		expect := `{"local_datetime":null}`
		jsonBytes, err := json.Marshal(jsonStruct)
		Nil(t, err)
		Equal(t, expect, string(jsonBytes), "omitemptyは効かない")
	}
	{
		var jsonStruct LocalDatetimeStruct
		expect := `{"local_datetime":null}`
		jsonBytes, err := json.Marshal(jsonStruct)
		Nil(t, err)
		Equal(t, expect, string(jsonBytes))
	}
}

func TestLocalDatetime_UnmarshalJSON(t *testing.T) {
	{
		var jsonStruct LocalDatetimeStruct
		jsonString := `{
			"local_datetime": "2020-07-19 11:14:15"
		}`
		expect := LocalDatetimeStruct{
			LocalDatetime: NewLocalDatetime(2020, 7, 19, 11, 14, 15),
		}
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		Nil(t, err)
		Equal(t, expect, jsonStruct)
	}
	{
		var jsonStruct LocalDatetimeStruct
		jsonString := `{
			"local_datetime": "2020/07/19 00:00:00"
		}`
		jsonString2 := `{
		"local_datetime": "2020-07-19T00:00:00"
		}`
		expect := LocalDatetimeStruct{}
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		NotNil(t, err)
		Equal(t, expect, jsonStruct, "invalid format")

		err = json.Unmarshal([]byte(jsonString2), &jsonStruct)
		NotNil(t, err)
		Equal(t, expect, jsonStruct, "invalid format")
	}
	{
		var jsonStruct LocalDatetimeStruct
		jsonString := `{"local_datetime": null}`
		expect := LocalDatetimeStruct{LocalDatetime: LocalDatetime{}}
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		Nil(t, err)
		Equal(t, expect, jsonStruct)
	}
	{
		var jsonStruct LocalDatetimeStruct
		jsonString := `{}`
		expect := LocalDatetimeStruct{LocalDatetime: LocalDatetime{}}
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		Nil(t, err)
		Equal(t, expect, jsonStruct)
	}
	{
		var jsonStruct LocalDatetimeStruct
		jsonString := `{
			"local_datetime": "unknownType"
		}`
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		NotNil(t, err)
	}
}

func TestLocalDatetime_UnmarshalFlag(t *testing.T) {
	for _, table := range []struct {
		title         string
		localDatetime *LocalDatetime
		target        string
		expect        LocalDatetime
		errorOccurred bool
	}{
		{
			title:         "LocalDatetime.unmarshalFlag.instance is nil,error発生",
			localDatetime: nil,
			target:        "2022-04-01 15:14:13",
			expect:        LocalDatetime{LocalDatetime: dates.LocalDatetime{}, Valid: false},
			errorOccurred: true,
		},
		{
			title:         "LocalDatetime.unmarshalFlag.target string localDate is empty,emptyがセットされるが,errorは発生しない",
			localDatetime: &LocalDatetime{},
			target:        "",
			expect:        LocalDatetime{LocalDatetime: dates.LocalDatetime{}, Valid: false},
			errorOccurred: false,
		},
		{
			title:         "LocalDatetime.unmarshalFlag.target string localDate is unexpected format, error発生",
			localDatetime: &LocalDatetime{},
			target:        "2022/04/01",
			expect:        LocalDatetime{LocalDatetime: dates.LocalDatetime{}, Valid: false},
			errorOccurred: true,
		},
		{
			title:         "LocalDatetime.unmarshalFlag.target string localDate's format: yyyy-MM-dd, parse error occurred",
			localDatetime: &LocalDatetime{},
			target:        "2022-04-01",
			expect:        LocalDatetime{LocalDatetime: dates.LocalDatetime{}, Valid: false},
			errorOccurred: true,
		},
		{
			title:         "LocalDatetime.unmarshalFlag.target string localDate is valid format.",
			localDatetime: &LocalDatetime{},
			target:        "2022-04-01 15:14:13",
			expect: LocalDatetime{LocalDatetime: dates.LocalDatetime{
				LocalDate: dates.LocalDate{Year: 2022, Month: 4, Day: 1},
				LocalTime: dates.LocalTime{Hour: 15, Minute: 14, Second: 13},
			}, Valid: true},
			errorOccurred: false,
		},
	} {
		t.Run(table.title, func(t *testing.T) {
			err := table.localDatetime.UnmarshalFlag(table.target)
			if table.errorOccurred {
				NotNil(t, err)
			} else {
				Nil(t, err)
			}
			if table.localDatetime != nil {
				Equal(t, table.expect, *table.localDatetime)
			}
		})
	}
}

func TestLocalDate_Scan(t *testing.T) {
	{
		value := "2020-02-01"
		expect := LocalDate{
			LocalDate: dates.LocalDate{Year: 2020, Month: 2, Day: 1},
			Valid:     true,
		}
		dt := new(LocalDate)
		err := dt.Scan(value)
		Nil(t, err, "valid value -> err is nil")
		Equal(t, expect.LocalDate, dt.LocalDate, "")
		Equal(t, true, dt.Valid, "")
	}
	{
		value := "2020/02/01"
		dt := new(LocalDate)
		err := dt.Scan(value)
		NotNil(t, err, "invalid format error occurred")
		Equal(t, dates.LocalDate{}, dt.LocalDate, "")
		Equal(t, false, dt.Valid, "")
	}
	{
		byteVal := []byte("2020-02-01")
		dt := new(LocalDate)
		err := dt.Scan(byteVal)
		NotNil(t, err, "invalid type value error occurred")
		Equal(t, dates.LocalDate{}, dt.LocalDate, "")
		Equal(t, false, dt.Valid, "")
	}
	{
		var value interface{}
		dtm := new(LocalDate)
		err := dtm.Scan(value)
		Nil(t, err, "nil value -> err is nil")
		Equal(t, dates.LocalDate{}, dtm.LocalDate, "")
		Equal(t, false, dtm.Valid, "")
	}
}

func TestLocalDate_Value(t *testing.T) {
	{
		emptyDT := LocalDate{
			LocalDate: dates.LocalDate{},
			Valid:     false,
		}
		value, err := emptyDT.Value()
		Nil(t, err, "empty -> err is  nil")
		Nil(t, value, "empty -> value is  nil")
	}
	{
		dt := LocalDate{
			LocalDate: dates.LocalDate{Year: 2020, Month: 1, Day: 1},
			Valid:     true,
		}
		value, err := dt.Value()
		Nil(t, err, "empty -> err is nil")
		Equal(t, "2020-01-01", value, "")
	}
}

func TestLocalDatetime_Scan(t *testing.T) {
	{
		value := "2020-02-01 15:10:15"
		expect := LocalDatetime{
			LocalDatetime: dates.LocalDatetime{
				LocalDate: dates.LocalDate{Year: 2020, Month: 2, Day: 1},
				LocalTime: dates.LocalTime{Hour: 15, Minute: 10, Second: 15},
			},
			Valid: true,
		}
		dtm := new(LocalDatetime)
		err := dtm.Scan(value)
		Nil(t, err, "valid value -> err is nil")
		Equal(t, expect.LocalDatetime, dtm.LocalDatetime, "")
		Equal(t, true, dtm.Valid, "")
	}
	{
		value := "2020/02/01 15:10:15"
		dtm := new(LocalDatetime)
		err := dtm.Scan(value)
		NotNil(t, err, "invalid format error occurred")
		Equal(t, dates.LocalDatetime{}, dtm.LocalDatetime, "")
		Equal(t, false, dtm.Valid, "")
	}
	{
		var value interface{}
		dtm := new(LocalDatetime)
		err := dtm.Scan(value)
		Nil(t, err, "nil value -> err is nil")
		Equal(t, dates.LocalDatetime{}, dtm.LocalDatetime, "")
		Equal(t, false, dtm.Valid, "")
	}
}

func TestLocalDatetime_Value(t *testing.T) {
	{
		emptyDTM := LocalDatetime{
			LocalDatetime: dates.LocalDatetime{},
			Valid:         false,
		}
		value, err := emptyDTM.Value()
		Nil(t, err, "empty -> err is nil")
		Nil(t, value, "empty -> value is nil")
	}

	{
		var emptyDTM LocalDatetime
		value, err := emptyDTM.Value()
		Nil(t, err, "not initialized NullLocalDateTime -> err is nil")
		Nil(t, value, "not initialized NullLocalDateTime -> value is nil")
	}
	{
		dtm := LocalDatetime{
			LocalDatetime: dates.LocalDatetime{
				LocalDate: dates.LocalDate{Year: 2020, Month: 1, Day: 1},
				LocalTime: dates.LocalTime{Hour: 13, Minute: 15, Second: 10},
			},
			Valid: true,
		}
		value, err := dtm.Value()
		Nil(t, err, "empty -> err is nil")
		Equal(t, "2020-01-01 13:15:10", value, "")
	}
}

func TestLocalDatetimeFromTime(t *testing.T) {
	{
		tm, _ := time.Parse(dates.RFC3339.String(), "2020-11-01T15:10:05+09:00")
		expect := dates.NewLocalDatetime(2020, 11, 1, 15, 10, 05)
		dtm := LocalDatetimeFromTime(tm)
		Equal(t, expect, dtm.LocalDatetime, "JST time -> expected")
		Equal(t, true, dtm.Valid, "")
	}
	{
		tm, _ := time.Parse(dates.RFC3339.String(), "2020-11-01T15:10:05Z")
		expect := dates.NewLocalDatetime(2020, 11, 1, 15, 10, 05)
		dtm := LocalDatetimeFromTime(tm)
		Equal(t, expect, dtm.LocalDatetime, "UTC time -> expected")
		Equal(t, true, dtm.Valid, "")
	}
	{
		tm, _ := time.Parse(dates.RFC3339.String(), "2020-11-01T15:10:05-10:00")
		expect := dates.NewLocalDatetime(2020, 11, 1, 15, 10, 05)
		dtm := LocalDatetimeFromTime(tm)
		Equal(t, expect, dtm.LocalDatetime, "Honolulu time -> expected")
		Equal(t, true, dtm.Valid, "")
	}
	{
		tm := time.Time{}
		expect := dates.LocalDatetime{}
		dtm := LocalDatetimeFromTime(tm)
		Equal(t, expect, dtm.LocalDatetime, "empty time -> expected")
		Equal(t, false, dtm.Valid, "")
	}
	{
		var tm time.Time
		expect := dates.LocalDatetime{}
		dtm := LocalDatetimeFromTime(tm)
		Equal(t, expect, dtm.LocalDatetime, "empty time -> expected")
		Equal(t, false, dtm.Valid, "")
	}
}

func TestLocalDatetimeFromTimePtr(t *testing.T) {

	time1 := time.Date(2022, 4, 1, 10, 0, 0, 0, dates.UTC.Location())
	time2 := time.Date(2022, 4, 1, 10, 0, 0, 0, dates.AsiaTokyo.Location())
	for _, table := range []struct {
		title  string
		time   *time.Time
		expect LocalDatetime
	}{
		{
			title:  "time nilの場合, Valid:false",
			time:   nil,
			expect: LocalDatetime{Valid: false},
		},
		{
			title:  "time emptyの場合, Valid:false",
			time:   &time.Time{},
			expect: LocalDatetime{Valid: false},
		},
		{
			title:  "not emptyなtimeの(UTC)場合, Valid:true",
			time:   &time1,
			expect: NewLocalDatetime(uint(time1.Year()), uint(time1.Month()), uint(time1.Day()), time1.Hour(), time1.Minute(), time1.Second()),
		},
		{
			title:  "not emptyなtime(JST)の場合, Valid:true",
			time:   &time2,
			expect: NewLocalDatetime(uint(time2.Year()), uint(time2.Month()), uint(time2.Day()), time2.Hour(), time2.Minute(), time2.Second()),
		},
	} {
		t.Run(table.title, func(t *testing.T) {
			actual := LocalDatetimeFromTimePtr(table.time)
			Equal(t, table.expect, actual)
		})
	}
}

func TestLocalDatetimeFromDate(t *testing.T) {
	{
		emptyDT := dates.LocalDate{}
		dtm := LocalDatetimeFromDate(emptyDT)
		expectDTM := emptyDT.LocalDatetime()
		Equal(t, false, dtm.Valid, "")
		Equal(t, expectDTM, dtm.LocalDatetime, "")
	}
	{
		var emptyDT dates.LocalDate
		dtm := LocalDatetimeFromDate(emptyDT)
		expectDTM := emptyDT.LocalDatetime()
		Equal(t, false, dtm.Valid, "")
		Equal(t, expectDTM, dtm.LocalDatetime, "")
	}
	{
		dt := dates.NewLocalDate(2020, 11, 1)
		expectDTM := dt.LocalDatetime()
		dtm := LocalDatetimeFromDate(dt)
		Equal(t, true, dtm.Valid, "")
		Equal(t, expectDTM, dtm.LocalDatetime, "")
	}
}

func TestLocalDatetimeFromDatetime(t *testing.T) {
	{
		emptyDTM := dates.LocalDatetime{}
		dtm := LocalDatetimeFromDatetime(emptyDTM)
		Equal(t, false, dtm.Valid, "")
		Equal(t, emptyDTM, dtm.LocalDatetime, "")
	}
	{
		var emptyDTM dates.LocalDatetime
		dtm := LocalDatetimeFromDatetime(emptyDTM)
		Equal(t, false, dtm.Valid, "")
		Equal(t, emptyDTM, dtm.LocalDatetime, "")
	}
	{
		localDatetime := dates.NewLocalDatetime(2020, 1, 1, 0, 0, 1)
		dtm := LocalDatetimeFromDatetime(localDatetime)
		Equal(t, true, dtm.Valid, "")
		Equal(t, localDatetime, dtm.LocalDatetime, "")
	}
}

func TestLocalDatetimeFromPtr(t *testing.T) {

	datetime := dates.NewLocalDatetime(2022, 4, 5, 5, 5, 5)
	zeroDatetetime := dates.NewLocalDatetime(0, 0, 0, 0, 0, 0)
	for _, table := range []struct {
		title    string
		datetime *dates.LocalDatetime
		expect   LocalDatetime
	}{
		{
			title:    "nilを渡すと, Valid:falseのLocalDatetimeが返却される",
			datetime: nil,
			expect:   LocalDatetime{Valid: false},
		},
		{
			title:    "empty localDateを渡すと, Valid:falseのLocalDatetimeが返却される",
			datetime: &dates.LocalDatetime{},
			expect:   LocalDatetime{Valid: false},
		},
		{
			title:    "localDateを渡すと, Valid:trueのLocalDatetimeが返却される",
			datetime: &datetime,
			expect: LocalDatetime{
				Valid:         true,
				LocalDatetime: datetime,
			},
		},
		{
			title:    "zero値のlocalDateを渡すと, Valid:falseのLocalDatetimeが返却される",
			datetime: &zeroDatetetime,
			expect:   LocalDatetime{Valid: false},
		},
	} {
		t.Run(table.title, func(t *testing.T) {
			actual := LocalDatetimeFromPtr(table.datetime)
			Equal(t, table.expect, actual)
		})
	}
}
