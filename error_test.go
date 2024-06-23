package i18n

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestNewErr(t *testing.T) {
	tests := []struct {
		name    string
		section string
		key     string
		values  []M
		want    *I18nError
	}{
		{
			name:    "NoValues",
			section: "section1",
			key:     "key1",
			values:  nil,
			want:    &I18nError{baseError: &baseError{section: "section1", key: "key1"}},
		},
		{
			name:    "SingleValue",
			section: "section2",
			key:     "key2",
			values:  []M{{"foo": "bar"}},
			want:    &I18nError{baseError: &baseError{section: "section2", key: "key2", values: M{"foo": "bar"}}},
		},
		{
			name:    "MultipleValues",
			section: "section3",
			key:     "key3",
			values:  []M{{"foo": "bar"}, {"baz": "qux"}},
			want:    &I18nError{baseError: &baseError{section: "section3", key: "key3", values: M{"foo": "bar"}}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewErr(tt.section, tt.key, tt.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewErr() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewErrWithCode(t *testing.T) {
	type args struct {
		code    int
		section string
		key     string
		values  []M
	}
	tests := []struct {
		name string
		args args
		want *I18nError
	}{
		{
			"Empty values",
			args{
				code:    500,
				section: "test_section",
				key:     "test_key",
				values:  []M{},
			},
			&I18nError{
				code: 500,
				baseError: &baseError{
					section: "test_section",
					key:     "test_key",
				},
			},
		},
		{
			"Values provided",
			args{
				code:    400,
				section: "test_section",
				key:     "test_error_key",
				values:  []M{{"test value key": "test value"}},
			},
			&I18nError{
				code: 400,
				baseError: &baseError{
					section: "test_section",
					key:     "test_error_key",
					values:  M{"test value key": "test value"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewErrWithCode(tt.args.code, tt.args.section, tt.args.key, tt.args.values...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewErrWithCode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestI18nError_EmptyLocale_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		setErr  func() *I18nError
		want    string
		wantErr bool
	}{
		{
			name: "empty i18n error",
			setErr: func() *I18nError {
				return &I18nError{}
			},
			want: `null`,
		},
		{
			name: "filled i18n error",
			setErr: func() *I18nError {
				return &I18nError{
					baseError: &baseError{
						section: "user_section",
						key:     "error_key",
					},
				}
			},
			want: `{"user_section":"error_key"}`,
		},
		{
			name: "partial fill i18n error",
			setErr: func() *I18nError {
				return &I18nError{
					baseError: &baseError{
						section: "user_section",
					},
				}
			},
			want: `{"user_section":""}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.setErr()

			got, err := json.Marshal(e) // e.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("I18nError.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStr := string(got); gotStr != tt.want {
				t.Errorf("I18nError.MarshalJSON() = %v, want %v", gotStr, tt.want)
			}
		})
	}
}

type TestStruct struct {
	Data  string     `json:"data"`
	Error *I18nError `json:"error"`
}

func TestNestedI18nError_EmptyLocale_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		setErr  func() *TestStruct
		want    string
		wantErr bool
	}{
		{
			name: "empty nested i18n error",
			setErr: func() *TestStruct {
				return &TestStruct{
					Data:  "test_data",
					Error: nil,
				}
			},
			want: `{"data":"test_data","error":null}`,
		},
		{
			name: "filled i18n error",
			setErr: func() *TestStruct {
				return &TestStruct{
					Data: "test_data",
					Error: &I18nError{
						baseError: &baseError{
							section: "user_section",
							key:     "error_key",
						},
					},
				}
			},
			want: `{"data":"test_data","error":{"user_section":"error_key"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.setErr()

			got, err := json.Marshal(e)
			if (err != nil) != tt.wantErr {
				t.Errorf("I18nError.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStr := string(got); gotStr != tt.want {
				t.Errorf("I18nError.MarshalJSON() = %v, want %v", gotStr, tt.want)
			}
		})
	}
}
