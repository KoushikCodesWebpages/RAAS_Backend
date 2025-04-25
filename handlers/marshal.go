package handlers

// import (
// 	"encoding/json"
// 	"errors"
// 	"gorm.io/datatypes"


// 	"RAAS/models"
// 	"RAAS/dto"
// )

// // MarshalStructToJSON marshals any struct to JSON and returns datatypes.JSON
// func MarshalStructToJSON(input interface{}) (datatypes.JSON, error) {
// 	if input == nil {
// 		return nil, errors.New("cannot marshal a nil input")
// 	}
// 	data, err := json.Marshal(input)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return datatypes.JSON(data), nil
// }

// // UnmarshalJSONToStruct unmarshals datatypes.JSON into a provided struct reference
// func UnmarshalJSONToStruct(jsonData datatypes.JSON, output interface{}) error {
// 	if len(jsonData) == 0 {
// 		return errors.New("empty or nil JSON data")
// 	}
// 	return json.Unmarshal(jsonData, output)
// }

// // GetPersonalInfo uses the global unmarshal utility
// func GetPersonalInfo(seeker *models.Seeker) (*dto.PersonalInfoRequest, error) {
// 	var personalInfo dto.PersonalInfoRequest
// 	if err := UnmarshalJSONToStruct(seeker.PersonalInfo, &personalInfo); err != nil {
// 		return nil, err
// 	}
// 	return &personalInfo, nil
// }

// // SetPersonalInfo uses the global marshal utility
// func SetPersonalInfo(seeker *models.Seeker, personalInfo *dto.PersonalInfoRequest) error {
// 	jsonData, err := MarshalStructToJSON(personalInfo)
// 	if err != nil {
// 		return err
// 	}
// 	seeker.PersonalInfo = jsonData
// 	return nil
// }


// // GetProfessionalSummary uses the global unmarshal utility
// func GetProfessionalSummary(seeker *models.Seeker) (*dto.ProfessionalSummaryRequest, error) {
// 	var professionalSummary dto.ProfessionalSummaryRequest
// 	if err := UnmarshalJSONToStruct(seeker.ProfessionalSummary, &professionalSummary); err != nil {
// 		return nil, err
// 	}
// 	return &professionalSummary, nil
// }


// // SetProfessionalSummary uses the global marshal utility
// func SetProfessionalSummary(seeker *models.Seeker, professionalSummary *dto.ProfessionalSummaryRequest) error {
// 	jsonData, err := MarshalStructToJSON(professionalSummary)
// 	if err != nil {
// 		return err
// 	}
// 	seeker.ProfessionalSummary = jsonData
// 	return nil
// }
