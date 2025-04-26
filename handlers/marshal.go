package handlers

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"RAAS/models"
	"RAAS/dto"
)

func MarshalStructToBson(input interface{}) (bson.M, error) {
	if input == nil {
		return nil, errors.New("cannot marshal a nil input")
	}
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	var bsonData bson.M
	if err := json.Unmarshal(data, &bsonData); err != nil {
		return nil, err
	}
	return bsonData, nil
}
func UnmarshalBsonToStruct(bsonData bson.M, output interface{}) error {
	if bsonData == nil {
		return errors.New("empty or nil BSON data")
	}
	data, err := json.Marshal(bsonData)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, output)
}

func MarshalArrayToBson(input interface{}) ([]byte, error) {
	if input == nil {
		return nil, errors.New("cannot marshal a nil input")
	}

	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}

	var bsonData []bson.M
	if err := json.Unmarshal(data, &bsonData); err != nil {
		return nil, err
	}
	return bson.Marshal(bsonData)
}

func UnmarshalBsonToArray(bsonData []byte, output interface{}) error {
	if bsonData == nil {
		return errors.New("empty or nil BSON data")
	}
	var bsonArray []bson.M
	if err := bson.Unmarshal(bsonData, &bsonArray); err != nil {
		return err
	}
	return json.Unmarshal(bsonData, output)
}

func GetFieldFromBson(bsonData bson.M, output interface{}) error {
	if err := UnmarshalBsonToStruct(bsonData, output); err != nil {
		return err
	}
	return nil
}

func SetFieldToBson(input interface{}, bsonData *bson.M) error {
	data, err := MarshalStructToBson(input)
	if err != nil {
		return err
	}
	*bsonData = data
	return nil
}

func GetEmbeddedData(seeker *models.Seeker, fieldName string, output interface{}) error {
	var field bson.M
	switch fieldName {
	case "personal_info":
		field = seeker.PersonalInfo
	case "professional_summary":
		field = seeker.ProfessionalSummary
	case "work_experiences":
		field = seeker.WorkExperiences
	case "education":
		field = seeker.Educations
	case "certificates":
		field = seeker.Certificates
	case "languages":
		field = seeker.Languages
	default:
		return errors.New("invalid field name")
	}
	return GetFieldFromBson(field, output)
}

func SetEmbeddedData(seeker *models.Seeker, fieldName string, input interface{}) error {
	var field *bson.M
	switch fieldName {
	case "personal_info":
		field = &seeker.PersonalInfo
	case "professional_summary":
		field = &seeker.ProfessionalSummary
	case "work_experiences":
		field = &seeker.WorkExperiences
	case "education":
		field = &seeker.Educations
	case "certificates":
		field = &seeker.Certificates
	case "languages":
		field = &seeker.Languages
	default:
		return errors.New("invalid field name")
	}

	return SetFieldToBson(input, field)
}

func GetPersonalInfo(seeker *models.Seeker) (*dto.PersonalInfoRequest, error) {
	var personalInfo dto.PersonalInfoRequest
	if err := GetEmbeddedData(seeker, "personal_info", &personalInfo); err != nil {
		return nil, err
	}
	return &personalInfo, nil
}

func SetPersonalInfo(seeker *models.Seeker, personalInfo *dto.PersonalInfoRequest) error {
	return SetEmbeddedData(seeker, "personal_info", personalInfo)
}

func GetProfessionalSummary(seeker *models.Seeker) (*dto.ProfessionalSummaryRequest, error) {
	var professionalSummary dto.ProfessionalSummaryRequest
	if err := GetEmbeddedData(seeker, "professional_summary", &professionalSummary); err != nil {
		return nil, err
	}
	return &professionalSummary, nil
}

func SetProfessionalSummary(seeker *models.Seeker, professionalSummary *dto.ProfessionalSummaryRequest) error {
	return SetEmbeddedData(seeker, "professional_summary", professionalSummary)
}

func GetWorkExperience(seeker *models.Seeker) ([]dto.WorkExperienceRequest, error) {
	var workExperiences []dto.WorkExperienceRequest
	if err := GetEmbeddedData(seeker, "work_experiences", &workExperiences); err != nil {
		return nil, err
	}
	return workExperiences, nil
}

func SetWorkExperience(seeker *models.Seeker, workExperiences []dto.WorkExperienceRequest) error {
	return SetEmbeddedData(seeker, "work_experiences", workExperiences)
}
