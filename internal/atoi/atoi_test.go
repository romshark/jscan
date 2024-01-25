package atoi_test

import (
	"math"
	"runtime"
	"strconv"
	"testing"

	"github.com/romshark/jscan/v2/internal/atoi"

	"github.com/stretchr/testify/require"
)

func TestU64(t *testing.T) {
	for _, td := range []struct {
		name     string
		input    string
		expect   uint64
		overflow bool
	}{
		{name: "0", input: `0`, expect: 0},
		{name: "1", input: `1`, expect: 1},
		{name: "12", input: `12`, expect: 12},
		{name: "123", input: `123`, expect: 123},
		{name: "1234", input: `1234`, expect: 1234},
		{name: "12345", input: `12345`, expect: 12345},
		{name: "123456", input: `123456`, expect: 123456},
		{name: "1234567", input: `1234567`, expect: 1234567},
		{name: "12345678", input: `12345678`, expect: 12345678},
		{name: "123456789", input: `123456789`, expect: 123456789},
		{name: "1234567891", input: `1234567891`, expect: 1234567891},
		{name: "12345678912", input: `12345678912`, expect: 12345678912},
		{name: "123456789123", input: `123456789123`, expect: 123456789123},
		{name: "1234567891234", input: `1234567891234`, expect: 1234567891234},
		{name: "12345678912345", input: `12345678912345`, expect: 12345678912345},
		{
			name:  "123456789123456",
			input: `123456789123456`, expect: 123456789123456,
		},
		{
			name:  "1234567891234567",
			input: `1234567891234567`, expect: 1234567891234567,
		},
		{
			name:  "12345678912345678",
			input: `12345678912345678`, expect: 12345678912345678,
		},
		{
			name:  "123456789123456789",
			input: `123456789123456789`, expect: 123456789123456789,
		},
		{
			name:  "1234567891234567891",
			input: `1234567891234567891`, expect: 1234567891234567891,
		},

		{name: "int32_max", input: `2147483647`, expect: math.MaxInt32},
		{name: "uint64_max", input: `18446744073709551615`, expect: math.MaxUint64},

		{name: "overflow_hi", input: `18446744073709551616`, overflow: true},
		{name: "overflow_hi2", input: `22222222222222222222`, overflow: true},
		{name: "overflow_l21", input: `222222222222222222222`, overflow: true},
	} {
		t.Run(td.name, func(t *testing.T) {
			a, overflow := atoi.U64(td.input)
			if td.overflow {
				require.True(t, overflow)
				require.Zero(t, a)
			} else {
				require.False(t, overflow)
				require.Equal(t, td.expect, a)
			}
		})
	}
}

func TestI64(t *testing.T) {
	for _, td := range []struct {
		name   string
		input  string
		expect int64
		err    bool
	}{
		{name: "0", input: `0`, expect: 0},
		{name: "1", input: `1`, expect: 1},
		{name: "12", input: `12`, expect: 12},
		{name: "123", input: `123`, expect: 123},
		{name: "1234", input: `1234`, expect: 1234},
		{name: "12345", input: `12345`, expect: 12345},
		{name: "123456", input: `123456`, expect: 123456},
		{name: "1234567", input: `1234567`, expect: 1234567},
		{name: "12345678", input: `12345678`, expect: 12345678},
		{name: "123456789", input: `123456789`, expect: 123456789},
		{name: "1234567891", input: `1234567891`, expect: 1234567891},
		{name: "12345678912", input: `12345678912`, expect: 12345678912},
		{name: "123456789123", input: `123456789123`, expect: 123456789123},
		{name: "1234567891234", input: `1234567891234`, expect: 1234567891234},
		{name: "12345678912345", input: `12345678912345`, expect: 12345678912345},
		{
			name:  "123456789123456",
			input: `123456789123456`, expect: 123456789123456,
		},
		{
			name:  "1234567891234567",
			input: `1234567891234567`, expect: 1234567891234567,
		},
		{
			name:  "12345678912345678",
			input: `12345678912345678`, expect: 12345678912345678,
		},
		{
			name:  "123456789123456789",
			input: `123456789123456789`, expect: 123456789123456789,
		},
		{
			name:  "1234567891234567891",
			input: `1234567891234567891`, expect: 1234567891234567891,
		},

		{name: "-1", input: `-1`, expect: -1},
		{name: "-12", input: `-12`, expect: -12},
		{name: "-123", input: `-123`, expect: -123},
		{name: "-1234", input: `-1234`, expect: -1234},
		{name: "-12345", input: `-12345`, expect: -12345},
		{name: "-123456", input: `-123456`, expect: -123456},
		{name: "-1234567", input: `-1234567`, expect: -1234567},
		{name: "-12345678", input: `-12345678`, expect: -12345678},
		{name: "-123456789", input: `-123456789`, expect: -123456789},
		{name: "-1234567891", input: `-1234567891`, expect: -1234567891},
		{name: "-12345678912", input: `-12345678912`, expect: -12345678912},
		{name: "-123456789123", input: `-123456789123`, expect: -123456789123},
		{name: "-1234567891234", input: `-1234567891234`, expect: -1234567891234},
		{
			name:  "-12345678912345",
			input: `-12345678912345`, expect: -12345678912345,
		},
		{
			name:  "-123456789123456",
			input: `-123456789123456`, expect: -123456789123456,
		},
		{
			name:  "-1234567891234567",
			input: `-1234567891234567`, expect: -1234567891234567,
		},
		{
			name:  "-12345678912345678",
			input: `-12345678912345678`, expect: -12345678912345678,
		},
		{
			name:  "-123456789123456789",
			input: `-123456789123456789`, expect: -123456789123456789,
		},
		{
			name:  "-1234567891234567891",
			input: `-1234567891234567891`, expect: -1234567891234567891,
		},

		{name: "int32_min", input: `-2147483648`, expect: math.MinInt32},
		{name: "int32_max", input: `2147483647`, expect: math.MaxInt32},
		{name: "int64_min", input: `-9223372036854775808`, expect: math.MinInt64},
		{name: "int64_max", input: `9223372036854775807`, expect: math.MaxInt64},

		{name: "overflow_lo", input: `-9223372036854775809`, err: true},
		{name: "overflow_hi", input: `9223372036854775808`, err: true},
		{name: "overflow_l20_neg", input: `-11111111111111111111`, err: true},
		{name: "overflow_l20_pos", input: `11111111111111111111`, err: true},
	} {
		t.Run(td.name, func(t *testing.T) {
			a, overflow := atoi.I64(td.input)
			if td.err {
				require.True(t, overflow)
				require.Zero(t, a)
			} else {
				require.False(t, overflow)
				require.Equal(t, td.expect, a)
			}
		})
	}
}

func TestU32(t *testing.T) {
	for _, td := range []struct {
		name   string
		input  string
		expect uint32
		err    bool
	}{
		{name: "0", input: `0`, expect: 0},
		{name: "1", input: `1`, expect: 1},
		{name: "12", input: `12`, expect: 12},
		{name: "123", input: `123`, expect: 123},
		{name: "1234", input: `1234`, expect: 1234},
		{name: "12345", input: `12345`, expect: 12345},
		{name: "123456", input: `123456`, expect: 123456},
		{name: "1234567", input: `1234567`, expect: 1234567},
		{name: "12345678", input: `12345678`, expect: 12345678},
		{name: "123456789", input: `123456789`, expect: 123456789},
		{name: "1234567891", input: `1234567891`, expect: 1234567891},
		{name: "uint32_max", input: `4294967295`, expect: math.MaxUint32},

		{name: "overflow_hi", input: `4294967296`, err: true},
		{name: "overflow_hi1", input: `4333333333`, err: true},
		{name: "overflow_hi2", input: `5555555555`, err: true},
		{name: "overflow_hi3", input: `6666666666`, err: true},
		{name: "overflow_hi4", input: `7777777777`, err: true},
		{name: "overflow_hi5", input: `8888888888`, err: true},
		{name: "overflow_hi6", input: `9999999999`, err: true},
		{name: "overflow_l11", input: `11111111111`, err: true},
	} {
		t.Run(td.name, func(t *testing.T) {
			a, overflow := atoi.U32(td.input)
			if td.err {
				require.True(t, overflow)
				require.Zero(t, a)
			} else {
				require.False(t, overflow)
				require.Equal(t, td.expect, a)
			}
		})
	}
}

func TestI32(t *testing.T) {
	for _, td := range []struct {
		name   string
		input  string
		expect int32
		err    bool
	}{
		{name: "0", input: `0`, expect: 0},
		{name: "1", input: `1`, expect: 1},
		{name: "12", input: `12`, expect: 12},
		{name: "123", input: `123`, expect: 123},
		{name: "1234", input: `1234`, expect: 1234},
		{name: "12345", input: `12345`, expect: 12345},
		{name: "123456", input: `123456`, expect: 123456},
		{name: "1234567", input: `1234567`, expect: 1234567},
		{name: "12345678", input: `12345678`, expect: 12345678},
		{name: "123456789", input: `123456789`, expect: 123456789},
		{name: "1234567891", input: `1234567891`, expect: 1234567891},

		{name: "-1", input: `-1`, expect: -1},
		{name: "-12", input: `-12`, expect: -12},
		{name: "-123", input: `-123`, expect: -123},
		{name: "-1234", input: `-1234`, expect: -1234},
		{name: "-12345", input: `-12345`, expect: -12345},
		{name: "-123456", input: `-123456`, expect: -123456},
		{name: "-1234567", input: `-1234567`, expect: -1234567},
		{name: "-12345678", input: `-12345678`, expect: -12345678},
		{name: "-123456789", input: `-123456789`, expect: -123456789},
		{name: "-1234567891", input: `-1234567891`, expect: -1234567891},

		{name: "int32_min", input: `-2147483648`, expect: math.MinInt32},
		{name: "int32_max", input: `2147483647`, expect: math.MaxInt32},

		{name: "overflow_int32_lo", input: `-2147483649`, err: true},
		{name: "overflow_int32_hi", input: `2147483648`, err: true},
		{name: "overflow_int64_min", input: `-9223372036854775808`, err: true},
		{name: "overflow_int64_max", input: `9223372036854775807`, err: true},
		{name: "overflow_int64_lo", input: `-9223372036854775809`, err: true},
		{name: "overflow_int64_hi", input: `9223372036854775808`, err: true},
		{name: "overflow_l20_neg", input: `-11111111111111111111`, err: true},
		{name: "overflow_l20_pos", input: `11111111111111111111`, err: true},
	} {
		t.Run(td.name, func(t *testing.T) {
			a, overflow := atoi.I32(td.input)
			if td.err {
				require.True(t, overflow)
				require.Zero(t, a)
			} else {
				require.False(t, overflow)
				require.Equal(t, td.expect, a)
			}
		})
	}
}

func TestI8(t *testing.T) {
	for _, td := range []struct {
		name   string
		input  string
		expect int8
		err    bool
	}{
		{name: "0", input: `0`, expect: 0},
		{name: "1", input: `1`, expect: 1},
		{name: "12", input: `12`, expect: 12},
		{name: "123", input: `123`, expect: 123},
		{name: "-1", input: `-1`, expect: -1},
		{name: "-12", input: `-12`, expect: -12},
		{name: "-123", input: `-123`, expect: -123},

		{name: "int8_min", input: `-128`, expect: math.MinInt8},
		{name: "int8_max", input: `127`, expect: math.MaxInt8},

		{name: "overflow_lo", input: `-129`, err: true},
		{name: "overflow_hi", input: `128`, err: true},
		{name: "overflow_int32_lo", input: `-2147483649`, err: true},
		{name: "overflow_int32_hi", input: `2147483648`, err: true},
		{name: "overflow_int64_min", input: `-9223372036854775808`, err: true},
		{name: "overflow_int64_max", input: `9223372036854775807`, err: true},
		{name: "overflow_int64_lo", input: `-9223372036854775809`, err: true},
		{name: "overflow_int64_hi", input: `9223372036854775808`, err: true},
	} {
		t.Run(td.name, func(t *testing.T) {
			a, overflow := atoi.I8(td.input)
			if td.err {
				require.True(t, overflow)
				require.Zero(t, a)
			} else {
				require.False(t, overflow)
				require.Equal(t, td.expect, a)
			}
		})
	}
}

func TestU8(t *testing.T) {
	for _, td := range []struct {
		name   string
		input  string
		expect uint8
		err    bool
	}{
		{name: "0", input: `0`, expect: 0},
		{name: "1", input: `1`, expect: 1},
		{name: "12", input: `12`, expect: 12},
		{name: "123", input: `123`, expect: 123},
		{name: "max_uint8", input: `255`, expect: 255},

		{name: "overflow_hi", input: `256`, err: true},
		{name: "overflow_hi1", input: `300`, err: true},
		{name: "overflow_hi1", input: `999`, err: true},
		{name: "overflow_l4", input: `1111`, err: true},
	} {
		t.Run(td.name, func(t *testing.T) {
			a, overflow := atoi.U8(td.input)
			if td.err {
				require.True(t, overflow)
				require.Zero(t, a)
			} else {
				require.False(t, overflow)
				require.Equal(t, td.expect, a)
			}
		})
	}
}

func TestI16(t *testing.T) {
	for _, td := range []struct {
		name   string
		input  string
		expect int16
		err    bool
	}{
		{name: "0", input: `0`, expect: 0},
		{name: "1", input: `1`, expect: 1},
		{name: "12", input: `12`, expect: 12},
		{name: "123", input: `123`, expect: 123},
		{name: "1234", input: `1234`, expect: 1234},
		{name: "12345", input: `12345`, expect: 12345},
		{name: "-1", input: `-1`, expect: -1},
		{name: "-12", input: `-12`, expect: -12},
		{name: "-123", input: `-123`, expect: -123},
		{name: "-1234", input: `-1234`, expect: -1234},
		{name: "-12345", input: `-12345`, expect: -12345},

		{name: "int16_min", input: `-32768`, expect: math.MinInt16},
		{name: "int16_max", input: `32767`, expect: math.MaxInt16},

		{name: "overflow_lo", input: `-32769`, err: true},
		{name: "overflow_hi", input: `32768`, err: true},
		{name: "overflow_int32_lo", input: `-2147483649`, err: true},
		{name: "overflow_int32_hi", input: `2147483648`, err: true},
		{name: "overflow_int64_min", input: `-9223372036854775808`, err: true},
		{name: "overflow_int64_max", input: `9223372036854775807`, err: true},
		{name: "overflow_int64_lo", input: `-9223372036854775809`, err: true},
		{name: "overflow_int64_hi", input: `9223372036854775808`, err: true},
	} {
		t.Run(td.name, func(t *testing.T) {
			a, overflow := atoi.I16(td.input)
			if td.err {
				require.True(t, overflow)
				require.Zero(t, a)
			} else {
				require.False(t, overflow)
				require.Equal(t, td.expect, a)
			}
		})
	}
}

func TestU16(t *testing.T) {
	for _, td := range []struct {
		name   string
		input  string
		expect uint16
		err    bool
	}{
		{name: "0", input: `0`, expect: 0},
		{name: "1", input: `1`, expect: 1},
		{name: "12", input: `12`, expect: 12},
		{name: "123", input: `123`, expect: 123},
		{name: "1234", input: `1234`, expect: 1234},
		{name: "12345", input: `12345`, expect: 12345},
		{name: "uint16_max", input: `65535`, expect: math.MaxUint16},

		{name: "overflow_hi", input: `65536`, err: true},
		{name: "overflow_hi1", input: `77777`, err: true},
		{name: "overflow_int32_hi", input: `2147483648`, err: true},
		{name: "overflow_int64_max", input: `9223372036854775807`, err: true},
	} {
		t.Run(td.name, func(t *testing.T) {
			a, overflow := atoi.U16(td.input)
			if td.err {
				require.True(t, overflow)
				require.Zero(t, a)
			} else {
				require.False(t, overflow)
				require.Equal(t, td.expect, a)
			}
		})
	}
}

func BenchmarkAtoi64(b *testing.B) {
	for _, bd := range []struct {
		name  string
		input string
	}{
		{name: "min_int64", input: `-9223372036854775808`},
		{name: "max_int64", input: `9223372036854775807`},
		{name: "1to9_____", input: `123456789`},
		{name: "neg_one__", input: `-1`},
		{name: "zero_____", input: `0`},
	} {
		b.Run(bd.name, func(b *testing.B) {
			var x int
			var x64 int64
			b.Run("strconv_atoi", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					var err error
					if x, err = strconv.Atoi(bd.input); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("strconv_parse", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					var err error
					if x64, err = strconv.ParseInt(bd.input, 10, 0); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("I64", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					var overflow bool
					if x64, overflow = atoi.I64(bd.input); overflow {
						b.Fatal("overflow")
					}
				}
			})
			runtime.KeepAlive(x)
			runtime.KeepAlive(x64)
		})
	}
}

func BenchmarkI8(b *testing.B) {
	for _, bd := range []struct {
		name  string
		input string
	}{
		{name: "min_int8", input: `-128`},
		{name: "max_int8", input: `127`},
		{name: "-1______", input: `-1`},
		{name: "zero____", input: `0`},
	} {
		b.Run(bd.name, func(b *testing.B) {
			var x int
			var x8 int8
			var x64 int64
			b.Run("strconv_atoi", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					var err error
					if x, err = strconv.Atoi(bd.input); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("strconv_parse", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					var err error
					if x64, err = strconv.ParseInt(bd.input, 10, 16); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("I16", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					var overflow bool
					if x8, overflow = atoi.I8(bd.input); overflow {
						b.Fatal("overflow")
					}
				}
			})
			runtime.KeepAlive(x)
			runtime.KeepAlive(x8)
			runtime.KeepAlive(x64)
		})
	}
}

func BenchmarkU8(b *testing.B) {
	for _, bd := range []struct {
		name  string
		input string
	}{
		{name: "max_uint8", input: `255`},
		{name: "123______", input: `123`},
		{name: "zero_____", input: `0`},
	} {
		b.Run(bd.name, func(b *testing.B) {
			var x int
			var x8 uint8
			var x64 uint64
			b.Run("strconv_atoi", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					var err error
					if x, err = strconv.Atoi(bd.input); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("strconv_parse", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					var err error
					if x64, err = strconv.ParseUint(bd.input, 10, 16); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("U8", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					var overflow bool
					if x8, overflow = atoi.U8(bd.input); overflow {
						b.Fatal("overflow")
					}
				}
			})
			runtime.KeepAlive(x)
			runtime.KeepAlive(x8)
			runtime.KeepAlive(x64)
		})
	}
}

func BenchmarkI16(b *testing.B) {
	for _, bd := range []struct {
		name  string
		input string
	}{
		{name: "min_int16", input: `-32768`},
		{name: "max_int16", input: `32767`},
		{name: "-1_______", input: `-1`},
		{name: "zero_____", input: `0`},
	} {
		b.Run(bd.name, func(b *testing.B) {
			var x int
			var x16 int16
			var x64 int64
			b.Run("strconv_atoi", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					var err error
					if x, err = strconv.Atoi(bd.input); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("strconv_parse", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					var err error
					if x64, err = strconv.ParseInt(bd.input, 10, 16); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("I16", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					var overflow bool
					if x16, overflow = atoi.I16(bd.input); overflow {
						b.Fatal("overflow")
					}
				}
			})
			runtime.KeepAlive(x)
			runtime.KeepAlive(x16)
			runtime.KeepAlive(x64)
		})
	}
}

func BenchmarkU16(b *testing.B) {
	for _, bd := range []struct {
		name  string
		input string
	}{
		{name: "max_int16", input: `65535`},
		{name: "32767____s", input: `32767`},
		{name: "zero_____", input: `0`},
	} {
		b.Run(bd.name, func(b *testing.B) {
			var x int
			var x16 uint16
			var x64 uint64
			b.Run("strconv_atoi", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					var err error
					if x, err = strconv.Atoi(bd.input); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("strconv_parse", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					var err error
					if x64, err = strconv.ParseUint(bd.input, 10, 16); err != nil {
						b.Fatal(err)
					}
				}
			})
			b.Run("U16", func(b *testing.B) {
				for n := 0; n < b.N; n++ {
					var overflow bool
					if x16, overflow = atoi.U16(bd.input); overflow {
						b.Fatal("overflow")
					}
				}
			})
			runtime.KeepAlive(x)
			runtime.KeepAlive(x16)
			runtime.KeepAlive(x64)
		})
	}
}

func BenchmarkU16Overflow(b *testing.B) {
	var x uint16
	var x64 uint64
	var err error
	inputOverflow := "65536"
	b.Run("U16", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			var overflow bool
			if x, overflow = atoi.U16(inputOverflow); !overflow {
				panic(err)
			}
		}
	})
	b.Run("strconv", func(b *testing.B) {
		for n := 0; n < b.N; n++ {
			if x64, err = strconv.ParseUint(inputOverflow, 10, 16); err == nil {
				panic(err)
			}
		}
	})
	runtime.KeepAlive(x)
	runtime.KeepAlive(x64)
	runtime.KeepAlive(err)
}
