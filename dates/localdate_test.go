package dates

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
	"time"

	. "github.com/stretchr/testify/assert"
)

type LocalDateStruct struct {
	LocalDate LocalDate `json:"local_date,omitempty"`
}
type LocalDateStruct2 struct {
	LocalDate LocalDate `json:"local_date"`
}

func TestLocalDate_MarshalJSON(t *testing.T) {
	{
		jsonStruct := LocalDateStruct{LocalDate: NewLocalDate(1999, 12, 31)}
		jsonBytes, err := json.Marshal(jsonStruct)
		expectJSON := `{"local_date":"1999-12-31"}`
		Nil(t, err)
		Equal(t, expectJSON, string(jsonBytes))
	}
	{
		var emptyStruct LocalDateStruct
		jsonBytes, err := json.Marshal(emptyStruct)
		expectJSON := `{"local_date":null}`
		Nil(t, err)
		Equal(t, expectJSON, string(jsonBytes), "omitemptyは効かない")
	}
	{
		var emptyStruct LocalDateStruct2
		jsonBytes, err := json.Marshal(emptyStruct)
		expectJSON := `{"local_date":null}`
		Nil(t, err)
		Equal(t, expectJSON, string(jsonBytes))
	}
}

func TestLocalDate_UnmarshalJSON(t *testing.T) {
	{
		jsonByte := []byte("")
		var localDate *LocalDate
		err := localDate.UnmarshalJSON(jsonByte)
		NotNil(t, err, "localDate is nil")
	}
	{
		var jsonByte []byte
		cpiType := new(LocalDate)
		err := cpiType.UnmarshalJSON(jsonByte)
		NotNil(t, err, "json is nil")
	}
	{
		jsonByte1 := []byte(`"2020-11-05"`)
		jsonByte2 := []byte(`"1980-07-11"`)
		expect1 := NewLocalDate(2020, 11, 5)
		expect2 := NewLocalDate(1980, 7, 11)
		localDate := new(LocalDate)
		err1 := localDate.UnmarshalJSON(jsonByte1)
		Nil(t, err1)
		Equal(t, expect1, *localDate)
		err2 := localDate.UnmarshalJSON(jsonByte2)
		Nil(t, err2)
		Equal(t, expect2, *localDate)
	}
	{
		jsonString := `{
						"local_date": "2020-07-03"
		}`
		expect := LocalDateStruct{LocalDate: NewLocalDate(2020, 7, 3)}
		var jsonStruct LocalDateStruct
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		Nil(t, err)
		Equal(t, expect, jsonStruct)
	}
	{
		jsonString := `{
						"local_date": "2020/07/03"
		}`
		expect := LocalDateStruct{}
		var jsonStruct LocalDateStruct
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		NotNil(t, err, "invalid format")
		Equal(t, expect, jsonStruct)
	}
	{
		jsonString := `{
						"local_date": "2020-07-03 10:14:12"
		}`
		expect := LocalDateStruct{}
		var jsonStruct LocalDateStruct
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		NotNil(t, err, "invalid format")
		Equal(t, expect, jsonStruct)
	}
	{
		emptyJSONString := `{}`
		var expect LocalDateStruct
		var jsonStruct LocalDateStruct
		err := json.Unmarshal([]byte(emptyJSONString), &jsonStruct)
		Nil(t, err)
		Equal(t, expect, jsonStruct)
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
			title:         "LocalDate.unmarshalFlag.localDate instance is nil. error occurred",
			localDate:     nil,
			target:        "2022-04-01",
			expect:        LocalDate{},
			errorOccurred: true,
		},
		{
			title:         "LocalDate.unmarshalFlag.target string localDate is empty, error occurred",
			localDate:     &LocalDate{},
			target:        "",
			expect:        LocalDate{},
			errorOccurred: true,
		},
		{
			title:         "LocalDate.unmarshalFlag.target string localDate is unexpected format, error occurred",
			localDate:     &LocalDate{},
			target:        "2020/04/01",
			expect:        LocalDate{},
			errorOccurred: true,
		},
		{
			title:         "LocalDate.unmarshalFlag.target string localDate's format: yyyy-MM-dd HH:mm:ss, parse error occurred",
			localDate:     &LocalDate{},
			target:        "2020-04-01 20:10:10",
			expect:        LocalDate{},
			errorOccurred: true,
		},
		{
			title:     "unmarshalFlag.target string localDate is valid format",
			localDate: &LocalDate{},
			target:    "2020-04-01",
			expect:    LocalDate{Year: 2020, Month: 4, Day: 1},
		},
	} {
		t.Run(table.title, func(t *testing.T) {
			err := table.localDate.UnmarshalFlag(table.target)
			if table.errorOccurred {
				NotNil(t, err)
			} else {
				Nil(t, err)
				Equal(t, table.expect, *table.localDate)
			}
		})
	}
}

func TestLocalDate_NewLocalDate(t *testing.T) {
	{
		localDate := NewLocalDate(0, 0, 0)
		Equal(t, LocalDate{}, localDate, "紀元前の場合空を返却する")
	}
	{
		localDate := NewLocalDate(0, -1, 0)
		Equal(t, LocalDate{}, localDate, "紀元前の場合空を返却する")
	}
	{
		localDate := NewLocalDate(1, 1, 1)
		expect := LocalDate{Year: 1, Month: 1, Day: 1}
		Equal(t, expect, localDate, "紀元前の場合空を返却する")
	}
	{
		localDate := NewLocalDate(2018, 10, 21)
		const year = "2018"
		const month = "10"
		const day = "21"
		y, m, d := localDate.SplitString()

		Equal(t, year, y, " year")
		Equal(t, month, m, " month")
		Equal(t, day, d, " day")
	}
	{
		localDate := NewLocalDate(9, 2, 1)
		const year = "0009"
		const month = "02"
		const day = "01"
		y, m, d := localDate.SplitString()
		Equal(t, year, y, " year")
		Equal(t, month, m, " month")
		Equal(t, day, d, " day")
	}
	{
		localDate := NewLocalDate(18, 2, 1)
		const year = "0018"
		const month = "02"
		const day = "01"
		y, m, d := localDate.SplitString()
		Equal(t, year, y, " year")
		Equal(t, month, m, " month")
		Equal(t, day, d, " day")
	}
	{
		dt := NewLocalDate(549, 14, 1)
		expect := LocalDate{Year: 550, Month: 2, Day: 1}
		Equal(t, expect, dt, "12を超える月を指定した場合, 計算後の結果が算出される")
	}
	{
		dt := NewLocalDate(549, 4, 33)
		expect := LocalDate{Year: 549, Month: 5, Day: 3}
		Equal(t, expect, dt, "31を超える日を指定した場合, 計算後の結果が算出される")
	}
}

func TestParseLocalDate(t *testing.T) {
	for _, table := range []struct {
		title         string
		format        Format
		target        string
		expect        LocalDate
		errorOccurred bool
	}{
		{
			title:         "ParseLocalDate.yyyy-MM-dd.想定通り取得できる",
			format:        DateHyphen,
			target:        "2020-04-05",
			expect:        NewLocalDate(2020, 4, 5),
			errorOccurred: false,
		},
		{
			title:         "ParseLocalDate.yyyy/MM/ddで渡す.想定通り取得できる",
			format:        DateSlash,
			target:        "2020/04/05",
			expect:        NewLocalDate(2020, 4, 5),
			errorOccurred: false,
		},
		{
			title:         "ParseLocalDate.yyyy-MM-dd HH:mm:ssでFMT:yyyy-MM-dd HH:mm:ssを渡すと,戻り値はyyyy-MM-dd",
			format:        DateTimeHyphen,
			target:        "2020-04-05 14:15:18",
			expect:        NewLocalDate(2020, 4, 5),
			errorOccurred: false,
		},
		{
			title:         "ParseLocalDate.yyyy-MM-ddTHH:mm:ssZで異なるFMT:yyyy-MM-dd HH:mm:ss渡すと,error発生",
			format:        DateTimeHyphen,
			target:        "2020-04-05T16:00:00+06:00",
			expect:        LocalDate{},
			errorOccurred: true,
		},
		{
			title:         "ParseLocalDate.yyyy-MM-ddTHH:mm:ssZ渡しFMT:RFC3339を指定しても,timezoneは無視される",
			format:        RFC3339,
			target:        "2020-04-05T01:00:00+09:00",
			expect:        NewLocalDate(2020, 4, 5),
			errorOccurred: false,
		},
		{
			title:         "ParseLocalDate.emptyを渡すと,エラー発生",
			format:        DateHyphen,
			target:        "",
			expect:        LocalDate{},
			errorOccurred: true,
		},
	} {
		t.Run(table.title, func(t *testing.T) {
			actual, err := ParseLocalDate(table.format, table.target)
			Equal(t, table.expect, actual)
			if table.errorOccurred {
				NotNil(t, err)
			} else {
				Nil(t, err)
			}
		})
	}
}

func TestLocalDate_LocalDatetime(t *testing.T) {
	{
		dt := NewLocalDate(2020, 11, 1)
		dtm := dt.LocalDatetime()
		expect := NewLocalDatetime(2020, 11, 1, 0, 0, 0)
		Equal(t, expect, dtm, "")
	}
	{
		dt := LocalDate{}
		dtm := dt.LocalDatetime()
		expect := LocalDatetime{}
		Equal(t, expect, dtm, "")
	}
	{
		var dt LocalDate
		dtm := dt.LocalDatetime()
		expect := LocalDatetime{}
		Equal(t, expect, dtm, "")
	}
}

func TestLocalDate_Sub(t *testing.T) {
	{
		dt1 := NewLocalDate(2020, 7, 10)
		dt2 := NewLocalDate(2020, 7, 01)
		duration, _ := dt1.Sub(dt2)
		expect := time.Duration(9) * time.Hour * 24
		Equal(t, expect, duration, "")
	}
	{
		dt1 := NewLocalDate(2020, 7, 10)
		dt2 := NewLocalDate(2020, 7, 10)
		duration, _ := dt1.Sub(dt2)
		expect := time.Duration(0) * time.Hour * 24
		Equal(t, expect, duration, "")
	}
	{
		dt1, dt2 := LocalDate{}, LocalDate{}
		duration, _ := dt1.Sub(dt2)
		expect := time.Duration(0) * time.Hour * 24
		Equal(t, expect, duration, "")
	}
	{
		dt1 := NewLocalDate(2020, 7, 10)
		dt2 := NewLocalDate(2019, 9, 10)
		duration, _ := dt1.Sub(dt2)
		expect := time.Duration(304) * time.Hour * 24
		Equal(t, expect, duration, "")
	}
}

func TestLocalDatePeriod_Divide(t *testing.T) {
	for _, table := range []struct {
		title     string
		input     LocalDatePeriod
		chunkDay  int
		expect    LocalDatePeriods
		expectErr error
	}{
		{
			title: "指定した日数で分割できる",
			input: LocalDatePeriod{
				Start: NewLocalDate(2022, 4, 1),
				End:   NewLocalDate(2022, 4, 7),
			},
			chunkDay: 3,
			expect: LocalDatePeriods{
				{Start: NewLocalDate(2022, 4, 1), End: NewLocalDate(2022, 4, 3)},
				{Start: NewLocalDate(2022, 4, 4), End: NewLocalDate(2022, 4, 6)},
				{Start: NewLocalDate(2022, 4, 7), End: NewLocalDate(2022, 4, 7)},
			},
		},
		{
			title: "月をまたいでいても指定した日数で分割できる",
			input: LocalDatePeriod{
				Start: NewLocalDate(2022, 3, 30),
				End:   NewLocalDate(2022, 4, 6),
			},
			chunkDay: 3,
			expect: LocalDatePeriods{
				{Start: NewLocalDate(2022, 3, 30), End: NewLocalDate(2022, 4, 1)},
				{Start: NewLocalDate(2022, 4, 2), End: NewLocalDate(2022, 4, 4)},
				{Start: NewLocalDate(2022, 4, 5), End: NewLocalDate(2022, 4, 6)},
			},
		},
		{
			title: "Startがemptyの場合,エラー発生",
			input: LocalDatePeriod{
				End: NewLocalDate(2022, 4, 6),
			},
			chunkDay:  3,
			expect:    nil,
			expectErr: ErrEmptyDate,
		},
		{
			title: "Endがemptyの場合,エラー発生",
			input: LocalDatePeriod{
				Start: NewLocalDate(2022, 4, 6),
			},
			chunkDay:  3,
			expect:    nil,
			expectErr: ErrEmptyDate,
		},
		{
			title: "chunkDayが0の場合,エラー発生.",
			input: LocalDatePeriod{
				Start: NewLocalDate(2022, 4, 1),
				End:   NewLocalDate(2022, 4, 6),
			},
			chunkDay:  0,
			expect:    nil,
			expectErr: ErrIncorrectDivisionDays,
		},
		{
			title: "期間が290年以上の場合,エラー発生",
			input: LocalDatePeriod{
				Start: NewLocalDate(1500, 4, 1),
				End:   NewLocalDate(2022, 4, 6),
			},
			chunkDay:  3,
			expect:    nil,
			expectErr: ErrLengthOfPeriod,
		},
		{
			title: "start,endが同一の場合,そのまま返却される.",
			input: LocalDatePeriod{
				Start: NewLocalDate(2022, 4, 1),
				End:   NewLocalDate(2022, 4, 1),
			},
			chunkDay: 1,
			expect: LocalDatePeriods{
				{Start: NewLocalDate(2022, 4, 1), End: NewLocalDate(2022, 4, 1)},
			},
		},
		{
			title: "start,endの期間より分割日数のほうが大きい場合,そのまま返却される.",
			input: LocalDatePeriod{
				Start: NewLocalDate(2022, 4, 1),
				End:   NewLocalDate(2022, 4, 3),
			},
			chunkDay: 5,
			expect: LocalDatePeriods{
				{Start: NewLocalDate(2022, 4, 1), End: NewLocalDate(2022, 4, 3)},
			},
		},
	} {
		t.Run(table.title, func(t *testing.T) {
			actual, err := DivideDatePeriod(table.input.Start, table.input.End, table.chunkDay)
			if table.expectErr != nil {
				NotNil(t, err)
				True(t, errors.Is(err, table.expectErr))
			} else {
				ElementsMatch(t, table.expect, actual)
			}

		})
	}
}

type LocalDatetimeStruct struct {
	LocalDatetime LocalDatetime `json:"local_datetime,omitempty"`
}
type LocalDatetimeStruct2 struct {
	LocalDatetime LocalDatetime `json:"local_datetime"`
}

func TestLocalDatetime_MarshalJSON(t *testing.T) {
	{
		jsonStruct := LocalDatetimeStruct{
			LocalDatetime: NewLocalDatetime(1999, 12, 31, 10, 11, 12),
		}
		jsonBytes, err := json.Marshal(jsonStruct)
		expectJSON := `{"local_datetime":"1999-12-31 10:11:12"}`
		Nil(t, err)
		Equal(t, expectJSON, string(jsonBytes))
	}
	{
		var emptyStruct LocalDatetimeStruct
		jsonBytes, err := json.Marshal(emptyStruct)
		expectJSON := `{"local_datetime":null}`
		Nil(t, err)
		Equal(t, expectJSON, string(jsonBytes), "omitemptyは効かない")
	}
	{
		var emptyStruct LocalDatetimeStruct2
		jsonBytes, err := json.Marshal(emptyStruct)
		expectJSON := `{"local_datetime":null}`
		Nil(t, err)
		Equal(t, expectJSON, string(jsonBytes))
	}
}

func TestLocalDatetime_UnmarshalJSON(t *testing.T) {
	{
		jsonByte := []byte("")
		var localDatetime *LocalDatetime
		err := localDatetime.UnmarshalJSON(jsonByte)
		NotNil(t, err, "localDatetime is nil")
	}
	{
		var jsonByte []byte
		localDatetime := new(LocalDatetime)
		err := localDatetime.UnmarshalJSON(jsonByte)
		NotNil(t, err, "json is nil")
	}
	{
		jsonByte1 := []byte(`"2020-11-05 13:14:15"`)
		jsonByte2 := []byte(`"1980-07-11 05:06:07"`)
		expect1 := NewLocalDatetime(2020, 11, 5, 13, 14, 15)
		expect2 := NewLocalDatetime(1980, 7, 11, 5, 6, 7)
		localDatetime := new(LocalDatetime)
		err1 := localDatetime.UnmarshalJSON(jsonByte1)
		Nil(t, err1)
		Equal(t, expect1, *localDatetime)
		err2 := localDatetime.UnmarshalJSON(jsonByte2)
		Nil(t, err2)
		Equal(t, expect2, *localDatetime)
	}
	{
		jsonString := `{
						"local_datetime": "2020-07-03 11:14:27"
		}`
		expect := LocalDatetimeStruct{
			LocalDatetime: NewLocalDatetime(2020, 7, 3, 11, 14, 27),
		}
		var jsonStruct LocalDatetimeStruct
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		Nil(t, err)
		Equal(t, expect, jsonStruct)
	}
	{
		jsonString := `{
						"local_datetime": "2020/07/03 11:14:15"
		}`
		expect := LocalDatetimeStruct{}
		var jsonStruct LocalDatetimeStruct
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		NotNil(t, err, "invalid format")
		Equal(t, expect, jsonStruct)
	}
	{
		jsonString := `{
						"local_datetime": "2020-07-03T10:14:12"
		}`
		expect := LocalDatetimeStruct{}
		var jsonStruct LocalDatetimeStruct
		err := json.Unmarshal([]byte(jsonString), &jsonStruct)
		NotNil(t, err, "invalid format")
		Equal(t, expect, jsonStruct)
	}
	{
		emptyJSONString := `{}`
		var expect LocalDateStruct
		var jsonStruct LocalDateStruct
		err := json.Unmarshal([]byte(emptyJSONString), &jsonStruct)
		Nil(t, err)
		Equal(t, expect, jsonStruct)
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
			title:         "LocalDatetime.unmarshalFlag.localDatetime instance is nil. error occurred",
			localDatetime: nil,
			target:        "2022-04-01",
			expect:        LocalDatetime{},
			errorOccurred: true,
		},
		{
			title:         "LocalDatetime.unmarshalFlag.target string localDatetime is empty, error occurred",
			localDatetime: &LocalDatetime{},
			target:        "",
			expect:        LocalDatetime{},
			errorOccurred: true,
		},
		{
			title:         "LocalDatetime.unmarshalFlag.target string localDatetime is unexpected format, error occurred",
			localDatetime: &LocalDatetime{},
			target:        "2020/04/01 15:15:15",
			expect:        LocalDatetime{},
			errorOccurred: true,
		},
		{
			title:         "LocalDatetime.unmarshalFlag.target string localDatetime's format: yyyy-MM-dd, parse error occurred",
			localDatetime: &LocalDatetime{},
			target:        "2020-04-01",
			expect:        LocalDatetime{},
			errorOccurred: true,
		},
		{
			title:         "LocalDatetime.unmarshalFlag.target string localDatetime is valid format",
			localDatetime: &LocalDatetime{},
			target:        "2020-04-01 15:16:17",
			expect: LocalDatetime{
				LocalDate: LocalDate{Year: 2020, Month: 4, Day: 1},
				LocalTime: LocalTime{Hour: 15, Minute: 16, Second: 17},
			},
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

func TestLocaltime_LocalTimeApply(t *testing.T) {
	d1 := &LocalDate{Year: 2000, Month: 12, Day: 1}
	d2 := &LocalDate{Year: 2000, Month: 12, Day: 1}
	fmt.Println(d1)
	fmt.Println(d2)
	{
		localTime, err := NewLocalTime(3, 4, 5)
		const (
			hour = "03"
			min  = "04"
			sec  = "05"
		)
		Nil(t, err, "")
		h, m, s := localTime.SplitString()
		Equal(t, hour, h, "")
		Equal(t, min, m, "")
		Equal(t, sec, s, "")
	}
	{
		localTime, err := NewLocalTime(23, 14, 51)
		const (
			hour = "23"
			min  = "14"
			sec  = "51"
		)
		Nil(t, err, "")
		h, m, s := localTime.SplitString()
		Equal(t, hour, h, "")
		Equal(t, min, m, "")
		Equal(t, sec, s, "")
	}
	{
		_, err := NewLocalTime(24, 14, 51)
		NotNil(t, err, " 時間の上限を超えています. ")
	}
	{
		_, err := NewLocalTime(12, 60, 51)
		NotNil(t, err, "分の上限を超えています. ")
	}
	{
		_, err := NewLocalTime(12, 40, 60)
		NotNil(t, err, "秒の上限を超えています. ")
	}
}

func TestNewLocalDatetime(t *testing.T) {
	{
		dtm := NewLocalDatetime(0, 0, 0, 0, 0, 0)
		Equal(t, LocalDatetime{}, dtm, "紀元前以前の日付を渡した場合, 空が返却される.")
	}
	{
		dtm := NewLocalDatetime(0, 0, 0, -1, 0, 0)
		Equal(t, LocalDatetime{}, dtm, "紀元前以前の日付を渡した場合, 空が返却される.")
	}
	{
		dtm := NewLocalDatetime(1, 1, 1, 0, 0, 0)
		expect := LocalDatetime{
			LocalDate: LocalDate{Year: 1, Month: 1, Day: 1},
			LocalTime: LocalTime{Hour: 0, Minute: 0, Second: 0},
		}
		Equal(t, expect, dtm, "西暦元年はOK")
	}
	{
		actual := NewLocalDatetime(2018, 10, 21, 16, 0, 0)
		expect := LocalDatetime{
			LocalDate: LocalDate{Year: 2018, Month: 10, Day: 21},
			LocalTime: LocalTime{Hour: 16, Minute: 0, Second: 0},
		}
		Equal(t, expect, actual, "")
	}
	{
		actual := NewLocalDatetime(9, 2, 1, 13, 0, 0)
		expect := LocalDatetime{
			LocalDate: LocalDate{Year: 9, Month: 2, Day: 1},
			LocalTime: LocalTime{Hour: 13, Minute: 0, Second: 0},
		}
		Equal(t, expect, actual, "")
	}
	{
		actual := NewLocalDatetime(549, 14, 1, 13, 0, 0)
		expect := LocalDatetime{
			LocalDate: LocalDate{Year: 550, Month: 2, Day: 1},
			LocalTime: LocalTime{Hour: 13, Minute: 0, Second: 0},
		}
		Equal(t, expect, actual, "12を超える月を指定した場合, 計算後の結果が算出される")
	}
	{
		actual := NewLocalDatetime(549, 4, 33, 13, 80, 80)
		expect := LocalDatetime{
			LocalDate: LocalDate{Year: 549, Month: 5, Day: 3},
			LocalTime: LocalTime{Hour: 14, Minute: 21, Second: 20},
		}
		Equal(t, expect, actual, "dayが31以上,min,secが60以上を指定した場合, 計算後の結果が算出される")
	}
	{
		actual1 := NewLocalDatetime(2018, 10, 21, 16, 0, 0)
		expect1 := LocalDatetime{
			LocalDate: LocalDate{Year: 2018, Month: 10, Day: 21},
			LocalTime: LocalTime{Hour: 16, Minute: 0, Second: 0},
		}
		Equal(t, expect1, actual1, "timezoneに依存しない. timezone指定なし")

		time.Local = time.FixedZone("Pacific/Honolulu", -10*60*60)
		actual2 := NewLocalDatetime(2018, 10, 21, 16, 0, 0)
		Equal(t, expect1, actual2, "timezoneに依存しない. timezone: Pacific/Honolulu")

		time.Local = time.FixedZone("Asia/Magadan", 12*60*60)
		actual3 := NewLocalDatetime(2018, 10, 21, 16, 0, 0)
		Equal(t, expect1, actual3, "timezoneに依存しない. timezone: Asia/Magadan")
	}
}

func TestLocalDatetime_After(t *testing.T) {
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 0)
		isAfter := dateTime.After(targetDateTime)
		Equal(t, isAfter, true, "日付が未来の場合が正しく判定される")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		isAfter := dateTime.After(targetDateTime)
		Equal(t, isAfter, false, "同時刻の場合false")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 2)
		isAfter := dateTime.After(targetDateTime)
		Equal(t, isAfter, false, "日付が過去の場合false")
	}
	{
		var dateTime LocalDatetime
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 0)
		isAfter := dateTime.After(targetDateTime)
		Equal(t, isAfter, false, "日付が空の場合は0000-00-00として判定される")
	}
}

func TestLocalDatetime_AfterEqual(t *testing.T) {
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 0)
		isAfter := dateTime.AfterEqual(targetDateTime)
		Equal(t, isAfter, true, "日付が未来の場合が正しく判定される")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		isAfter := dateTime.AfterEqual(targetDateTime)
		Equal(t, isAfter, true, "同時刻の場合true")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 2)
		isAfter := dateTime.AfterEqual(targetDateTime)
		Equal(t, isAfter, false, "日付が過去の場合false")
	}
	{
		var dateTime LocalDatetime
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 0)
		isAfter := dateTime.AfterEqual(targetDateTime)
		Equal(t, isAfter, false, "日付が空の場合は0000-00-00として判定される")
	}
}

func TestLocalDatetime_Between(t *testing.T) {
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		start := NewLocalDatetime(2000, 3, 1, 15, 1, 1)
		end := NewLocalDatetime(2020, 3, 1, 15, 1, 2)
		isBetween := dateTime.Between(start, end)
		Equal(t, isBetween, true, "日付が指定期間内の場合true")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		start := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		end := NewLocalDatetime(2020, 3, 1, 15, 1, 2)
		isBetween := dateTime.Between(start, end)
		Equal(t, isBetween, true, "日付がstartと同時刻, endより以前の場合true")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		start := NewLocalDatetime(2020, 3, 1, 15, 1, 0)
		end := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		isBetween := dateTime.Between(start, end)
		Equal(t, isBetween, true, "日付がstartより後以降, endと同時刻の場合true")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		start := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		end := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		isBetween := dateTime.Between(start, end)
		Equal(t, isBetween, true, "日付がstart, endと同時刻の場合true")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		start := NewLocalDatetime(2020, 3, 1, 15, 2, 15)
		end := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		isBetween := dateTime.Between(start, end)
		Equal(t, isBetween, false, "日付がstart以前, endと同時刻の場合false")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		start := NewLocalDatetime(2020, 3, 1, 15, 1, 0)
		end := NewLocalDatetime(2020, 3, 1, 15, 1, 0)
		isBetween := dateTime.Between(start, end)
		Equal(t, isBetween, false, "日付がstartより後, endより後の場合false")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		start := NewLocalDatetime(2020, 3, 1, 15, 1, 3)
		end := NewLocalDatetime(2020, 3, 1, 15, 1, 0)
		isBetween := dateTime.Between(start, end)
		Equal(t, isBetween, false, "日付がstartより前, endより後の場合false")
	}
	{
		var dateTime LocalDatetime
		start := NewLocalDatetime(2020, 3, 1, 15, 1, 3)
		end := NewLocalDatetime(2020, 3, 1, 15, 1, 0)
		isBetween := dateTime.Between(start, end)
		Equal(t, isBetween, false, "対象日付が空の場合, start, endに値セットの場合false")
	}
}

func TestLocalDatetime_Sub(t *testing.T) {
	{
		dtm := NewLocalDatetime(2020, 6, 10, 3, 15, 10)
		target := NewLocalDatetime(2020, 6, 9, 21, 15, 10)
		duration, ok := dtm.Sub(target)
		Equal(t, true, ok, "")
		Equal(t, float64(6), duration.Hours(), "")
		Equal(t, float64(6*60), duration.Minutes(), "")
		Equal(t, float64(6*60*60), duration.Seconds(), "")
	}
	{
		dtm := NewLocalDatetime(2020, 6, 9, 21, 15, 10)
		target := NewLocalDatetime(2020, 6, 10, 3, 15, 10)
		duration, ok := dtm.Sub(target)
		Equal(t, true, ok, "")
		Equal(t, float64(-6), duration.Hours(), "dtm < target sub hours -> minus")
		Equal(t, float64(-6*60), duration.Minutes(), "dtm < target sub hours -> minus")
		Equal(t, float64(-6*60*60), duration.Seconds(), "dtm < target sub hours -> minus")
	}
	{
		var dtm LocalDatetime
		var target LocalDatetime
		duration, ok := dtm.Sub(target)
		Equal(t, true, ok, "dtm, target both empty -> ok")
		Equal(t, float64(0), duration.Hours(), "dtm, target  both empty. sub hours -> 0")
	}
	{
		dtm := NewLocalDatetime(2020, 6, 9, 21, 15, 10)
		var target LocalDatetime
		duration, ok := dtm.Sub(target)

		Equal(t, false, ok, "only target is empty -> 290年以上のためfalse")
		Equal(t, MaxDuration, duration, "only target is empty -> 290年以上のためmaxDuration")
	}
	{
		dtm := NewLocalDatetime(2000, 6, 30, 0, 0, 0)
		target := NewLocalDatetime(1710, 7, 30, 0, 0, 0)
		duration, ok := dtm.Sub(target)
		Equal(t, true, ok, "duration less than 290 years, but year sub is 290year -> true")
		Equal(t, float64(2541384), duration.Hours(), "\"duration less than 290 years, but year sub is 290year -> expected")
	}
	{
		dtm := NewLocalDatetime(280, 6, 9, 21, 15, 10)
		var target LocalDatetime
		_, ok := dtm.Sub(target)
		Equal(t, true, ok, "dtm year 280,  target is empty -> 290年未満以上のためtrue")
	}
}

func TestLocalDate_IsZero(t *testing.T) {
	{
		var dt LocalDate
		Equal(t, true, dt.IsZero(), "")
	}
	{
		dt := LocalDate{}
		Equal(t, true, dt.IsZero(), "")
	}
	{
		dt := LocalDate{Year: 0, Month: 0, Day: 0}
		Equal(t, true, dt.IsZero(), "")
	}
	{
		dt := LocalDate{Year: 2020, Month: 4, Day: 1}
		Equal(t, false, dt.IsZero(), "")
	}
	{
		dt := NewLocalDate(2020, 1, 1)
		Equal(t, false, dt.IsZero(), "")
	}
}

func TestLocalDatetime_Before(t *testing.T) {
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 2)
		isAfter := dateTime.Before(targetDateTime)
		Equal(t, isAfter, true, "日付が過去の場合が正しく判定される")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		isAfter := dateTime.Before(targetDateTime)
		Equal(t, isAfter, false, "同時刻の場合が正しく判定される")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 0)
		isAfter := dateTime.Before(targetDateTime)
		Equal(t, isAfter, false, "日付が未来の場合false")
	}
	{
		var dateTime LocalDatetime
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 0)
		isAfter := dateTime.Before(targetDateTime)
		Equal(t, isAfter, true, "日付が空の場合は0000-00-00として判定される")
	}
}

func TestLocalDatetime_BeforeEqual(t *testing.T) {
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 2)
		isAfter := dateTime.BeforeEqual(targetDateTime)
		Equal(t, isAfter, true, "日付が過去の場合が正しく判定される")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		isAfter := dateTime.BeforeEqual(targetDateTime)
		Equal(t, isAfter, true, "同時刻の場合が正しく判定される")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 0)
		isAfter := dateTime.BeforeEqual(targetDateTime)
		Equal(t, isAfter, false, "日付が未来の場合false")
	}
	{
		var dateTime LocalDatetime
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 0)
		isAfter := dateTime.BeforeEqual(targetDateTime)
		Equal(t, isAfter, true, "日付が空の場合は0000-00-00として判定される")
	}
}

func TestLocalDatetime_Equal(t *testing.T) {
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		isEqual := dateTime.Equal(targetDateTime)
		Equal(t, isEqual, true, "同時刻の場合true")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 2)
		isEqual := dateTime.Equal(targetDateTime)
		Equal(t, isEqual, false, "異なる時刻の場合false")
	}
	{
		dateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		targetDateTime := LocalDatetime{}
		isEqual := dateTime.Equal(targetDateTime)
		Equal(t, isEqual, false, "対象時刻がemptyの場合false")
	}
	{
		dateTime := LocalDatetime{}
		targetDateTime := NewLocalDatetime(2020, 3, 1, 15, 1, 1)
		isEqual := dateTime.Equal(targetDateTime)
		Equal(t, isEqual, false, "自身の時刻がemptyの場合false")
	}
	{
		dateTime := LocalDatetime{}
		targetDateTime := LocalDatetime{}
		isEqual := dateTime.Equal(targetDateTime)
		Equal(t, isEqual, true, "自身,対象時刻がemptyの場合true")
	}
}

func TestApplyLocalDatetimeByTime(t *testing.T) {
	{
		loc := UTC.Location()
		year := 2020
		month := 3
		day := 3
		hour := 16
		min := 15
		sec := 20
		time := time.Date(year, time.Month(month), day, hour, min, sec, 0, loc)
		actual := LocalDatetimeFromTime(time)
		expect := NewLocalDatetime(uint(year), uint(month), uint(day), hour, min, sec)
		Equal(t, actual, expect, "time -> datetimeの変換 jst timezoneは無視される")
	}
	{
		loc := AsiaTokyo.Location()
		year := 2020
		month := 3
		day := 3
		hour := 16
		min := 15
		sec := 20
		time := time.Date(year, time.Month(month), day, hour, min, sec, 0, loc)
		actual := LocalDatetimeFromTime(time)
		expect := NewLocalDatetime(uint(year), uint(month), uint(day), hour, min, sec)
		Equal(t, actual, expect, "time -> datetimeの変換 utc timezoneは無視される")
	}
	{
		year := 2020
		month := 4
		day := 5
		hour := 16
		min := 0
		sec := 0
		tm, _ := time.Parse(RFC3339.String(), "2020-04-05T16:00:00+03:00")
		actual := LocalDatetimeFromTime(tm)
		expect := NewLocalDatetime(uint(year), uint(month), uint(day), hour, min, sec)
		Equal(t, actual, expect, "time -> datetimeの変換 timezoneは無視される")
	}
}

func TestParseLocalDatetimeInLoc(t *testing.T) {
	{
		actual, _ := ParseLocalDatetime(DateTimeHyphen, "2020-04-05 16:00:00")
		expect := NewLocalDatetime(2020, 4, 5, 16, 0, 0)
		Equal(t, actual, expect, "yyyy-MM-dd hh:mm:ssの変換. loc:JST指定. ")
	}
	{
		actual, _ := ParseLocalDatetime(DateTimeHyphen, "2020-04-05 16:00:00")
		expect := NewLocalDatetime(2020, 4, 5, 16, 0, 0)
		Equal(t, actual, expect, "yyyy-MM-dd hh:mm:ssの変換. loc:UTC指定. ")
	}
	{
		actual, _ := ParseLocalDatetime(RFC3339, "2020-04-05T16:00:00+06:00")
		expect := NewLocalDatetime(2020, 4, 5, 16, 0, 0)
		Equal(t, actual, expect, "tz:+06:00. loc:UTC指定. timezoneは無視される")
	}
	{
		time.Local = time.FixedZone("Local", 0)
		actual, _ := ParseLocalDatetime(DateTimeSlash, "2020/04/05 16:00:00")
		expect := NewLocalDatetime(2020, 4, 5, 16, 0, 0)
		Equal(t, actual, expect, "localのtimezoneに依存しない. local:UTC. ")
		time.Local = time.FixedZone("Local", -6*60*60)
		actual2, _ := ParseLocalDatetime(DateTimeSlash, "2020/04/05 16:00:00")
		expect2 := NewLocalDatetime(2020, 4, 5, 16, 0, 0)
		Equal(t, actual2, expect2, "localのtimezoneに依存しない. local:America/Chicago. ")
	}
}

func TestLocalDatetime_Add(t *testing.T) {
	{
		dtm := NewLocalDatetime(2020, 4, 15, 20, 20, 20)
		expect := NewLocalDatetime(2020, 4, 15, 23, 20, 20)
		actual := dtm.Add(time.Duration(3) * time.Hour)
		Equal(t, actual, expect, "3h加算")
	}
	{
		dtm := NewLocalDatetime(2020, 4, 15, 20, 20, 20)
		expect := NewLocalDatetime(2020, 4, 15, 21, 40, 20)
		actual := dtm.Add(time.Duration(80) * time.Minute)
		Equal(t, actual, expect, "80min加算")
	}
	{
		dtm := NewLocalDatetime(2020, 4, 15, 20, 20, 20)
		expect := NewLocalDatetime(2020, 4, 12, 20, 20, 20)
		actual := dtm.Add(time.Duration(-72) * time.Hour)
		Equal(t, actual, expect, "-72h減算")
	}
	{
		time.Local = time.FixedZone("Local", 0)
		dtm := NewLocalDatetime(2020, 4, 15, 20, 20, 20)
		expect := NewLocalDatetime(2020, 4, 12, 20, 20, 20)
		actual := dtm.Add(time.Duration(-72) * time.Hour)
		Equal(t, actual, expect, "-72h減算. tz:UTC. timezoneに依存しない")
	}
}

func TestLocalDatetime_AddDate(t *testing.T) {
	{
		dtm := NewLocalDatetime(2020, 4, 15, 20, 20, 20)
		expect := NewLocalDatetime(2020, 7, 15, 20, 20, 20)
		actual := dtm.AddDate(0, 3, 0)
		Equal(t, actual, expect, "3ヶ月加算")
	}
	{
		dtm := NewLocalDatetime(2020, 4, 15, 20, 20, 20)
		expect := NewLocalDatetime(2022, 5, 16, 20, 20, 20)
		actual := dtm.AddDate(2, 1, 1)
		Equal(t, actual, expect, "2年, 1ヶ月 ,1日加算")
	}
	{
		dtm := NewLocalDatetime(2020, 4, 15, 20, 20, 20)
		expect := NewLocalDatetime(2020, 4, 15, 20, 20, 20)
		actual := dtm.AddDate(0, 0, 0)
		Equal(t, actual, expect, "0加算")
	}
	{
		dtm := NewLocalDatetime(2020, 4, 15, 20, 20, 20)
		expect := NewLocalDatetime(2021, 2, 14, 20, 20, 20)
		actual := dtm.AddDate(0, 9, 30)
		Equal(t, actual, expect, "月,日を年度と月が変わるように加算")
	}
	{
		time.Local = time.FixedZone("Local", 0)
		dtm := NewLocalDatetime(2020, 4, 15, 20, 20, 20)
		expect := NewLocalDatetime(2020, 4, 12, 20, 20, 20)
		actual := dtm.Add(time.Duration(-72) * time.Hour)
		Equal(t, actual, expect, "-72h減算. tz:UTC. timezoneに依存しない")
	}
}

func TestLocalDatetime_IsZero(t *testing.T) {
	{
		var dtm LocalDatetime
		Equal(t, true, dtm.IsZero(), "")
	}
	{
		dtm := LocalDatetime{}
		Equal(t, true, dtm.IsZero(), "")
	}
	{
		dtm := LocalDatetime{LocalDate: LocalDate{}, LocalTime: LocalTime{}}
		Equal(t, true, dtm.IsZero(), "")
	}
	{
		dtm := LocalDatetime{LocalDate: LocalDate{Year: 0, Month: 0, Day: 0}, LocalTime: LocalTime{Hour: 0, Minute: 0, Second: 0}}
		Equal(t, true, dtm.IsZero(), "")
	}
	{
		dtm := NewLocalDatetime(2020, 1, 1, 14, 14, 14)
		Equal(t, false, dtm.IsZero(), "")
	}
}

func TestLocalDate_Scan(t *testing.T) {
	{
		value := "2020-05-01"
		expect := LocalDate{Year: 2020, Month: 5, Day: 1}
		dt := new(LocalDate)
		err := dt.Scan(value)
		Nil(t, err, "valid value -> err is nil")
		Equal(t, expect, *dt, "valid value -> date is expected")
	}
	{
		emptyValue := ""
		expect := LocalDate{}
		dt := new(LocalDate)
		err := dt.Scan(emptyValue)
		NotNil(t, err, "empty value -> err occurred!")
		Equal(t, expect, *dt, "empty value -> date is empty")
	}
	{
		var nilValue []byte
		expect := LocalDate{}
		dt := new(LocalDate)
		err := dt.Scan(nilValue)
		NotNil(t, err, "nil value -> err occurred!")
		Equal(t, expect, *dt, "nil value -> date is empty")
	}
	{
		invalidFormatValue := "2020/01/03"
		expect := LocalDate{}
		dt := new(LocalDate)
		err := dt.Scan(invalidFormatValue)
		NotNil(t, err, "invalid format value -> err is occurred!")
		Equal(t, expect, *dt, "invalid format value -> date is empty")
	}
	{
		validDatetimeValue := "2020-01-03 15:00:00"
		expect := LocalDate{Year: 2020, Month: 1, Day: 3}
		dt := new(LocalDate)
		err := dt.Scan(validDatetimeValue)
		Nil(t, err, "valid datetime format value -> err is nil")
		Equal(t, expect, *dt, "valid datetime format value -> date is expected")
	}
}

func TestLocalDate_Value(t *testing.T) {
	{
		var emptyDate LocalDate
		expected := "0000-00-00"
		value, err := emptyDate.Value()
		Nil(t, err, "empty value -> err is nil")
		Equal(t, expected, value, "empty value -> result is empty")
	}
	{
		date := LocalDate{Year: 2020, Month: 5, Day: 10}
		expected := "2020-05-10"
		value, err := date.Value()
		Nil(t, err, "valid value -> err is nil")
		Equal(t, expected, value, "valid value -> value is expected")
	}
}

func TestLocalDatetime_Scan(t *testing.T) {
	{
		value := "2020-05-01 10:15:35"
		expect := LocalDatetime{
			LocalDate: LocalDate{Year: 2020, Month: 5, Day: 1},
			LocalTime: LocalTime{Hour: 10, Minute: 15, Second: 35},
		}
		dtm := new(LocalDatetime)
		err := dtm.Scan(value)
		Nil(t, err, "valid value -> err is nil")
		Equal(t, expect, *dtm, "valid value -> datetime is expected")
	}
	{
		emptyValue := ""
		expect := LocalDatetime{}
		dt := new(LocalDatetime)
		err := dt.Scan(emptyValue)
		NotNil(t, err, "empty value -> err occurred!")
		Equal(t, expect, *dt, "empty value -> datetime is empty")
	}
	{
		var nilValue []byte
		expect := LocalDatetime{}
		dt := new(LocalDatetime)
		err := dt.Scan(nilValue)
		NotNil(t, err, "nil value -> err occurred!")
		Equal(t, expect, *dt, "nil value -> datetime is empty")
	}
	{
		invalidFormatValue := "2020/01/03 10:15:35"
		expect := LocalDatetime{}
		dt := new(LocalDatetime)
		err := dt.Scan(invalidFormatValue)
		NotNil(t, err, "invalid format value -> err is occurred!")
		Equal(t, expect, *dt, "invalid format value -> datetime is empty")
	}
	{
		validDatetimeValue := "2020-01-03 15:20:35"
		expect := LocalDatetime{
			LocalDate: LocalDate{Year: 2020, Month: 1, Day: 3},
			LocalTime: LocalTime{Hour: 15, Minute: 20, Second: 35},
		}
		dt := new(LocalDatetime)
		err := dt.Scan(validDatetimeValue)
		Nil(t, err, "valid datetime format value -> err is nil")
		Equal(t, expect, *dt, "valid datetime format value -> date is expected")
	}
}
