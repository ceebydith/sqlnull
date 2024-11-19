package sqlnull_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/ceebydith/sqlnull"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
)

type OriginalSqlNullTest struct {
	FieldBool   sql.NullBool
	FieldByte   sql.NullByte
	FieldFloat  sql.NullFloat64
	FieldInt16  sql.NullInt16
	FieldInt32  sql.NullInt32
	FieldInt64  sql.NullInt64
	FieldString sql.NullString
	FieldTime   sql.NullTime
}

type BasicSqlNullTest struct {
	FieldBool   bool
	FieldByte   byte
	FieldFloat  float64
	FieldInt16  int16
	FieldInt32  int32
	FieldInt64  int64
	FieldString string
	FieldTime   time.Time
}

type NewSqlNullTest struct {
	FieldBool   *bool
	FieldByte   *byte
	FieldFloat  *float32
	FieldInt16  *int16
	FieldInt32  *int32
	FieldInt64  *int64
	FieldString *string
	FieldTime   *time.Time
}

type CustomBool bool
type CustomByte byte
type CustomFloat float64
type CustomInt16 int16
type CustomInt32 int32
type CustomInt64 int64
type CustomString string

type CustomSqlNullTest struct {
	FieldBool   *CustomBool
	FieldByte   *CustomByte
	FieldFloat  *CustomFloat
	FieldInt16  *CustomInt16
	FieldInt32  *CustomInt32
	FieldInt64  *CustomInt64
	FieldString *CustomString
	FieldTime   *time.Time
}

func makedatabase(values ...any) (*sql.Row, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `
		CREATE TABLE [test] (
			[field_bool] BOOLEAN,
			[field_int] INTEGER,
			[field_float] FLOAT,
			[field_time] DATETIME,
			[field_string] TEXT
		);
		INSERT INTO test (field_bool, field_int, field_float, field_time, field_string) VALUES (?, ?, ?, ?, ?);
	`
	_, err = db.Exec(query, values...)
	if err != nil {
		return nil, err
	}

	query = `SELECT field_bool, field_int, field_float, field_int, field_int, field_int, field_string, field_time, field_bool FROM test`
	row := db.QueryRow(query)
	return row, nil
}

func TestUnsupportedqlNullValued(t *testing.T) {
	row, err := makedatabase(true, 99, 88.89, time.Now(), "lorem ipsum")
	require.NoError(t, err)

	var test NewSqlNullTest
	alien := new(struct {
		Data bool
	})
	targets := sqlnull.Scanner(
		&test.FieldBool,
		&test.FieldByte,
		&test.FieldFloat,
		&test.FieldInt16,
		&test.FieldInt32,
		&test.FieldInt64,
		&test.FieldString,
		&test.FieldTime,
		&alien,
	)
	require.NoError(t, err)

	err = row.Scan(targets...)
	require.Error(t, err)
}

func TestNewSqlNullValued(t *testing.T) {
	row, err := makedatabase(true, 99, 88.89, time.Now(), "lorem ipsum")
	require.NoError(t, err)

	var test NewSqlNullTest
	targets := sqlnull.Scanner(
		&test.FieldBool,
		&test.FieldByte,
		&test.FieldFloat,
		&test.FieldInt16,
		&test.FieldInt32,
		&test.FieldInt64,
		&test.FieldString,
		&test.FieldTime,
		nil,
	)
	require.NoError(t, err)

	err = row.Scan(targets...)
	require.NoError(t, err)
	require.NotEmpty(t, test.FieldBool)
	require.Equal(t, true, *test.FieldBool)
	require.NotEmpty(t, test.FieldByte)
	require.Equal(t, byte(99), *test.FieldByte)
	require.NotEmpty(t, test.FieldFloat)
	require.Equal(t, float32(88.89), *test.FieldFloat)
	require.NotEmpty(t, test.FieldInt16)
	require.Equal(t, int16(99), *test.FieldInt16)
	require.NotEmpty(t, test.FieldInt32)
	require.Equal(t, int32(99), *test.FieldInt32)
	require.NotEmpty(t, test.FieldInt64)
	require.Equal(t, int64(99), *test.FieldInt64)
	require.NotEmpty(t, test.FieldString)
	require.Equal(t, "lorem ipsum", *test.FieldString)
	require.NotEmpty(t, test.FieldTime)
	require.Equal(t, true, !test.FieldTime.IsZero() && !test.FieldTime.After(time.Now()))
}

func TestNewSqlNullEmpty(t *testing.T) {
	row, err := makedatabase(nil, nil, nil, nil, nil)
	require.NoError(t, err)

	var test NewSqlNullTest = NewSqlNullTest{
		FieldBool:   new(bool),
		FieldFloat:  new(float32),
		FieldInt16:  new(int16),
		FieldInt32:  new(int32),
		FieldInt64:  new(int64),
		FieldString: new(string),
		FieldTime:   new(time.Time),
	}
	targets := sqlnull.Scanner(
		&test.FieldBool,
		&test.FieldByte,
		&test.FieldFloat,
		&test.FieldInt16,
		&test.FieldInt32,
		&test.FieldInt64,
		&test.FieldString,
		&test.FieldTime,
		nil,
	)
	require.NoError(t, err)

	err = row.Scan(targets...)
	require.NoError(t, err)
	require.Empty(t, test.FieldBool)
	require.Empty(t, test.FieldByte)
	require.Empty(t, test.FieldFloat)
	require.Empty(t, test.FieldInt16)
	require.Empty(t, test.FieldInt32)
	require.Empty(t, test.FieldInt64)
	require.Empty(t, test.FieldString)
	require.Empty(t, test.FieldTime)
}

func TestCustomSqlNullValued(t *testing.T) {
	row, err := makedatabase(true, 99, 88.89, time.Now(), "lorem ipsum")
	require.NoError(t, err)

	var test CustomSqlNullTest
	targets := sqlnull.Scanner(
		&test.FieldBool,
		&test.FieldByte,
		&test.FieldFloat,
		&test.FieldInt16,
		&test.FieldInt32,
		&test.FieldInt64,
		&test.FieldString,
		&test.FieldTime,
		nil,
	)
	require.NoError(t, err)

	err = row.Scan(targets...)
	require.NoError(t, err)
	require.NotEmpty(t, test.FieldBool)
	require.Equal(t, CustomBool(true), *test.FieldBool)
	require.NotEmpty(t, test.FieldByte)
	require.Equal(t, CustomByte(99), *test.FieldByte)
	require.NotEmpty(t, test.FieldFloat)
	require.Equal(t, CustomFloat(88.89), *test.FieldFloat)
	require.NotEmpty(t, test.FieldInt16)
	require.Equal(t, CustomInt16(99), *test.FieldInt16)
	require.NotEmpty(t, test.FieldInt32)
	require.Equal(t, CustomInt32(99), *test.FieldInt32)
	require.NotEmpty(t, test.FieldInt64)
	require.Equal(t, CustomInt64(99), *test.FieldInt64)
	require.NotEmpty(t, test.FieldString)
	require.Equal(t, CustomString("lorem ipsum"), *test.FieldString)
	require.NotEmpty(t, test.FieldTime)
	require.Equal(t, true, !test.FieldTime.IsZero() && !test.FieldTime.After(time.Now()))
}

func TestCustomSqlNullEmpty(t *testing.T) {
	row, err := makedatabase(nil, nil, nil, nil, nil)
	require.NoError(t, err)

	var test CustomSqlNullTest = CustomSqlNullTest{
		FieldBool:   new(CustomBool),
		FieldByte:   new(CustomByte),
		FieldFloat:  new(CustomFloat),
		FieldInt16:  new(CustomInt16),
		FieldInt32:  new(CustomInt32),
		FieldInt64:  new(CustomInt64),
		FieldString: new(CustomString),
		FieldTime:   new(time.Time),
	}
	targets := sqlnull.Scanner(
		&test.FieldBool,
		&test.FieldByte,
		&test.FieldFloat,
		&test.FieldInt16,
		&test.FieldInt32,
		&test.FieldInt64,
		&test.FieldString,
		&test.FieldTime,
		nil,
	)
	require.NoError(t, err)

	err = row.Scan(targets...)
	require.NoError(t, err)
	require.Empty(t, test.FieldBool)
	require.Empty(t, test.FieldByte)
	require.Empty(t, test.FieldFloat)
	require.Empty(t, test.FieldInt16)
	require.Empty(t, test.FieldInt32)
	require.Empty(t, test.FieldInt64)
	require.Empty(t, test.FieldString)
	require.Empty(t, test.FieldTime)
}

func TestOriginalSqlNullValues(t *testing.T) {
	row, err := makedatabase(true, 99, 88.89, time.Now(), "lorem ipsum")
	require.NoError(t, err)

	var test OriginalSqlNullTest
	targets := sqlnull.Scanner(
		&test.FieldBool,
		&test.FieldByte,
		&test.FieldFloat,
		&test.FieldInt16,
		&test.FieldInt32,
		&test.FieldInt64,
		&test.FieldString,
		&test.FieldTime,
		nil,
	)
	require.NoError(t, err)

	err = row.Scan(targets...)
	require.NoError(t, err)
	require.Equal(t, true, test.FieldBool.Valid)
	require.Equal(t, true, test.FieldBool.Bool)
	require.Equal(t, true, test.FieldByte.Valid)
	require.Equal(t, byte(99), test.FieldByte.Byte)
	require.Equal(t, true, test.FieldFloat.Valid)
	require.Equal(t, float64(88.89), test.FieldFloat.Float64)
	require.Equal(t, true, test.FieldInt16.Valid)
	require.Equal(t, int16(99), test.FieldInt16.Int16)
	require.Equal(t, true, test.FieldInt32.Valid)
	require.Equal(t, int32(99), test.FieldInt32.Int32)
	require.Equal(t, true, test.FieldInt64.Valid)
	require.Equal(t, int64(99), test.FieldInt64.Int64)
	require.Equal(t, true, test.FieldString.Valid)
	require.Equal(t, "lorem ipsum", test.FieldString.String)
	require.Equal(t, true, test.FieldTime.Valid)
	require.Equal(t, true, !test.FieldTime.Time.IsZero() && !test.FieldTime.Time.After(time.Now()))
}

func TestOriginalSqlNullEmpty(t *testing.T) {
	row, err := makedatabase(nil, nil, nil, nil, nil)
	require.NoError(t, err)

	var test OriginalSqlNullTest
	targets := sqlnull.Scanner(
		&test.FieldBool,
		&test.FieldByte,
		&test.FieldFloat,
		&test.FieldInt16,
		&test.FieldInt32,
		&test.FieldInt64,
		&test.FieldString,
		&test.FieldTime,
		nil,
	)
	require.NoError(t, err)

	err = row.Scan(targets...)
	require.NoError(t, err)
	require.Equal(t, false, test.FieldBool.Valid)
	require.Equal(t, false, test.FieldByte.Valid)
	require.Equal(t, false, test.FieldFloat.Valid)
	require.Equal(t, false, test.FieldInt16.Valid)
	require.Equal(t, false, test.FieldInt32.Valid)
	require.Equal(t, false, test.FieldInt64.Valid)
	require.Equal(t, false, test.FieldString.Valid)
	require.Equal(t, false, test.FieldTime.Valid)
}

func TestBasicSqlNullValued(t *testing.T) {
	row, err := makedatabase(true, 99, 88.89, time.Now(), "lorem ipsum")
	require.NoError(t, err)

	var test BasicSqlNullTest
	targets := sqlnull.Scanner(
		&test.FieldBool,
		&test.FieldByte,
		&test.FieldFloat,
		&test.FieldInt16,
		&test.FieldInt32,
		&test.FieldInt64,
		&test.FieldString,
		&test.FieldTime,
		nil,
	)
	require.NoError(t, err)

	err = row.Scan(targets...)
	require.NoError(t, err)
	require.Equal(t, true, test.FieldBool)
	require.Equal(t, byte(99), test.FieldByte)
	require.Equal(t, float64(88.89), test.FieldFloat)
	require.Equal(t, int16(99), test.FieldInt16)
	require.Equal(t, int32(99), test.FieldInt32)
	require.Equal(t, int64(99), test.FieldInt64)
	require.Equal(t, "lorem ipsum", test.FieldString)
	require.Equal(t, true, !test.FieldTime.IsZero() && !test.FieldTime.After(time.Now()))
}

func TestUnsupportedType(t *testing.T) {
	var some int
	null := sqlnull.New(some)
	err := null.Scan(99)
	require.Error(t, err)
}
