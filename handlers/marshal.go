package handlers

import (
	"encoding/json"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"RAAS/models"
	"RAAS/dto"
)

// MarshalStructToBson marshals any struct to BSON format (compatible with MongoDB)
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

// UnmarshalBsonToStruct unmarshals BSON into a provided struct reference
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

// General function to get embedded data from BSON (for both PersonalInfo and ProfessionalSummary)
func GetFieldFromBson(bsonData bson.M, output interface{}) error {
	if err := UnmarshalBsonToStruct(bsonData, output); err != nil {
		return err
	}
	return nil
}

// General function to set embedded data in BSON (for both PersonalInfo and ProfessionalSummary)
func SetFieldToBson(input interface{}, bsonData *bson.M) error {
	data, err := MarshalStructToBson(input)
	if err != nil {
		return err
	}
	*bsonData = data
	return nil
}

// GetPersonalInfo uses the general function to retrieve personal info from BSON
func GetPersonalInfo(seeker *models.Seeker) (*dto.PersonalInfoRequest, error) {
	var personalInfo dto.PersonalInfoRequest
	if err := GetFieldFromBson(seeker.PersonalInfo, &personalInfo); err != nil {
		return nil, err
	}
	return &personalInfo, nil
}

// SetPersonalInfo uses the general function to set personal info in BSON format
func SetPersonalInfo(seeker *models.Seeker, personalInfo *dto.PersonalInfoRequest) error {
	return SetFieldToBson(personalInfo, &seeker.PersonalInfo)
}

// GetProfessionalSummary uses the general function to retrieve professional summary
func GetProfessionalSummary(seeker *models.Seeker) (*dto.ProfessionalSummaryRequest, error) {
	var professionalSummary dto.ProfessionalSummaryRequest
	if err := GetFieldFromBson(seeker.ProfessionalSummary, &professionalSummary); err != nil {
		return nil, err
	}
	return &professionalSummary, nil
}

// SetProfessionalSummary uses the general function to set professional summary in BSON format
func SetProfessionalSummary(seeker *models.Seeker, professionalSummary *dto.ProfessionalSummaryRequest) error {
	return SetFieldToBson(professionalSummary, &seeker.ProfessionalSummary)
}
