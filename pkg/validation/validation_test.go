package validation

import (
	"testing"

	ut "github.com/go-playground/universal-translator"
	valid "github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
)

type testStruct struct {
	Field1 int    `json:"field1,omitempty" validate:"eq=1"`
	Field2 string `json:"field2,omitempty" validate:"oneof=value1 value2 value3"`
	Field3 string `json:"field3,omitempty" validate:"required_if=Field1 1 Field4 test"`
	Field4 string `json:"field4,omitempty"`
}

func TestValidateStruct(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		validation, err := New(WithPredefinedErrorMessages(), WithJSONNamesForStructFields())
		assert.NoError(t, err)
		tStruct := testStruct{
			Field1: 1,
			Field2: "value1",
			Field3: "field3",
			Field4: "test",
		}
		err = validation.ValidateStruct(tStruct)
		assert.NoError(t, err)
	})

	t.Run("error case - default error messages", func(t *testing.T) {
		validation, err := New()
		assert.NoError(t, err)
		testStruct := testStruct{
			Field1: 1,
			Field2: "value4",
			Field4: "test",
		}
		err = validation.ValidateStruct(testStruct)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Key: 'testStruct.Field2' Error:Field validation for 'Field2' failed on the 'oneof'")
		assert.Contains(t, err.Error(), "tagKey: 'testStruct.Field3' Error:Field validation for 'Field3' failed on the 'required_if' tag")
	})

	t.Run("error case - predefined error messages without naming fields from json tags", func(t *testing.T) {
		validation, err := New(WithPredefinedErrorMessages())
		assert.NoError(t, err)
		testStruct := testStruct{
			Field1: 1,
			Field2: "value4",
			Field4: "test",
		}
		err = validation.ValidateStruct(testStruct)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "The Field2 must use the following values: value1 value2 value3.")
		assert.Contains(t, err.Error(), "The Field3 is required since the Field1 is 1.")
	})

	t.Run("failed case - predefined error messages with naming fields from json tags", func(t *testing.T) {
		validation, err := New(WithPredefinedErrorMessages(), WithJSONNamesForStructFields())
		assert.NoError(t, err)
		testStruct := testStruct{
			Field1: 1,
			Field2: "value4",
			Field4: "test",
		}
		err = validation.ValidateStruct(testStruct)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "The field2 must use the following values: value1 value2 value3.")
		assert.Contains(t, err.Error(), "The field3 is required since the Field1 is 1.")
	})

	t.Run("failed case - custom error messages", func(t *testing.T) {
		translations := []Translation{
			{
				Tag: "oneof",
				RegisterFn: func(ut ut.Translator) error {
					return ut.Add("oneof", "Custom error message for {0}. ", true)
				},
				TranslationFn: func(ut ut.Translator, fe valid.FieldError) string {
					t, _ := ut.T("oneof", fe.Field())
					return t
				},
			}, {
				Tag: "required_if",
				RegisterFn: func(ut ut.Translator) error {
					return ut.Add("required_if", "Custom error message for {0}. ", true)
				},
				TranslationFn: func(ut ut.Translator, fe valid.FieldError) string {
					t, _ := ut.T("required_if", fe.Field())
					return t
				},
			},
		}
		validation, err := New(WithCustomErrorMessages(translations))
		assert.NoError(t, err)
		testStruct := testStruct{
			Field1: 1,
			Field2: "value4",
			Field4: "test",
		}
		err = validation.ValidateStruct(testStruct)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "Custom error message for Field2")
		assert.Contains(t, err.Error(), "Custom error message for Field3")
	})
}
