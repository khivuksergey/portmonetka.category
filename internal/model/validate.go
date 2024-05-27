package model

import "github.com/go-playground/validator/v10"

func GetCategoryValidator() *validator.Validate {
	v := validator.New(validator.WithRequiredStructEnabled())
	//TODO Add custom validation messages
	//v.RegisterStructValidation(validateCategoryUpdate, CategoryUpdateDTO{})

	//// Add custom message for the validation error
	//v.RegisterTranslation("atleastonefieldrequired", en.Translations, func(ut validator.Translator) error {
	//	return ut.Add("atleastonefieldrequired", "At least one field must be provided", true)
	//}, func(ut validator.Translator, fe validator.FieldError) string {
	//	t, _ := ut.T("atleastonefieldrequired", fe.Field())
	//	return t
	//})

	return v
}

//func validateCategoryUpdate(sl validator.StructLevel) {
//	category := sl.Current().Interface().(CategoryUpdateDTO)
//
//	if category.Name == nil && category.Description == nil && category.Currency == nil && category.InitialAmount == nil {
//		sl.ReportError(category, "CategoryUpdateDTO", "", "atleastonefieldrequired", "")
//	}
//}
