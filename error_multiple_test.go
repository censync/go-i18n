package i18n

import (
	"encoding/json"
	"testing"
)

func TestI18nMultipleError_EmptyLocale_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		setErr  func() *I18nMultipleError
		want    string
		wantErr bool
	}{
		{
			name: "empty i18n error",
			setErr: func() *I18nMultipleError {
				return &I18nMultipleError{}
			},
			want: `null`,
		},
		{
			name: "filled i18n multiple error",
			setErr: func() *I18nMultipleError {
				return &I18nMultipleError{
					errors: map[string]*BaseError{
						"field1": {
							section: "user_section1",
							key:     "error_key1",
						},
					},
				}
			},
			want: `{"field1":{"user_section1":"error_key1"}}`,
		},
		{
			name: "partial fill i18n error",
			setErr: func() *I18nMultipleError {
				return &I18nMultipleError{
					errors: map[string]*BaseError{
						"field1": {
							section: "user_section",
							key:     "error_key",
						},
					},
				}
			},
			want: `{"field1":{"user_section":"error_key"}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.setErr()

			got, err := json.Marshal(e) // e.MarshalJSON()
			if (err != nil) != tt.wantErr {
				t.Errorf("I18nMultipleError.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStr := string(got); gotStr != tt.want {
				t.Errorf("I18nMultipleError.MarshalJSON() = %v, want %v", gotStr, tt.want)
			}
		})
	}
}

type TestStructMultiple struct {
	Data   string             `json:"data"`
	Errors *I18nMultipleError `json:"errors"`
}

func TestNestedI18nMultipleError_EmptyLocale_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		setErr  func() *TestStructMultiple
		want    string
		wantErr bool
	}{
		{
			name: "empty nested i18n error",
			setErr: func() *TestStructMultiple {
				return &TestStructMultiple{
					Data:   "test_data",
					Errors: nil,
				}
			},
			want: `{"data":"test_data","errors":null}`,
		},
		{
			name: "filled i18n error",
			setErr: func() *TestStructMultiple {
				return &TestStructMultiple{
					Data: "test_data",
					Errors: &I18nMultipleError{
						errors: map[string]*BaseError{
							"field1": {
								section: "user_section1",
								key:     "error_key1",
							},
						},
					},
				}
			},
			want: `{"data":"test_data","errors":{"field1":{"user_section1":"error_key1"}}}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := tt.setErr()

			got, err := json.Marshal(e)
			if (err != nil) != tt.wantErr {
				t.Errorf("I18nMultipleError.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotStr := string(got); gotStr != tt.want {
				t.Errorf("I18nMultipleError.MarshalJSON() = %v, want %v", gotStr, tt.want)
			}
		})
	}
}
