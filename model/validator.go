package model

import "github.com/go-playground/validator"

var validate *validator.Validate

//RegisterValidator :
func RegisterValidator() {
	validate = validator.New()
	validate.RegisterStructValidation(userStructLevelValidation, ObjectTag{})
}

func userStructLevelValidation(sl validator.StructLevel) {

	tag := sl.Current().Interface().(ObjectTag)

	if len(tag.TagID) == 0 && len(tag.ObjectID) == 0 {
		sl.ReportError(tag.TagID, "TagID", "tag_id", "tag_id", "")
		sl.ReportError(tag.ObjectID, "ObjectID", "object_id", "object_id", "")
	}

	// plus can do more, even with different tag than "fnameorlname"
}
