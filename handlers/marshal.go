package handlers

import (
	"encoding/json"
	"errors"

	"RAAS/dto"
	"RAAS/models"
	"go.mongodb.org/mongo-driver/bson"
)

// --- General marshal/unmarshal helpers ---

func MarshalStructToBson(input interface{}) (bson.M, error) {
	if input == nil {
		return nil, errors.New("cannot marshal a nil input")
	}
	data, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	var bsonData bson.M
	err = json.Unmarshal(data, &bsonData)
	return bsonData, err
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
	return bson.Marshal(input)
}

func UnmarshalBsonToArray(bsonData []byte, output interface{}) error {
	return bson.Unmarshal(bsonData, output)
}

func GetPersonalInfo(seeker *models.Seeker) (*dto.PersonalInfoRequest, error) {
	var personalInfo dto.PersonalInfoRequest
	if seeker.PersonalInfo == nil {
		return nil, errors.New("personal info is nil")
	}
	if err := UnmarshalBsonToStruct(seeker.PersonalInfo, &personalInfo); err != nil {
		return nil, err
	}
	return &personalInfo, nil
}

func SetPersonalInfo(seeker *models.Seeker, personalInfo *dto.PersonalInfoRequest) error {
	bsonData, err := MarshalStructToBson(personalInfo)
	if err != nil {
		return err
	}
	seeker.PersonalInfo = bsonData
	return nil
}

func GetProfessionalSummary(seeker *models.Seeker) (*dto.ProfessionalSummaryRequest, error) {
	var professionalSummary dto.ProfessionalSummaryRequest
	if seeker.ProfessionalSummary == nil {
		return nil, errors.New("professional summary is nil")
	}
	if err := UnmarshalBsonToStruct(seeker.ProfessionalSummary, &professionalSummary); err != nil {
		return nil, err
	}
	return &professionalSummary, nil
}

func SetProfessionalSummary(seeker *models.Seeker, professionalSummary *dto.ProfessionalSummaryRequest) error {
	bsonData, err := MarshalStructToBson(professionalSummary)
	if err != nil {
		return err
	}
	seeker.ProfessionalSummary = bsonData
	return nil
}
func GetWorkExperience(seeker *models.Seeker) ([]bson.M, error) {
    if len(seeker.WorkExperiences) == 0 {
        return []bson.M{}, nil
    }
    return seeker.WorkExperiences, nil
}

// SetWorkExperience sets the work experiences for a Seeker using an array of bson.M.
func SetWorkExperience(seeker *models.Seeker, workExperiences []bson.M) error {
    seeker.WorkExperiences = workExperiences
    return nil
}


func AppendToWorkExperience(seeker *models.Seeker, newWorkExperience dto.WorkExperienceRequest) error {
    // Check if the WorkExperiences array is nil or empty, if so, initialize it
    if seeker.WorkExperiences == nil {
        seeker.WorkExperiences = []bson.M{}
    }

    // Append the new work experience as a bson.M document
    workExperienceBson := bson.M{
        "job_title":           newWorkExperience.JobTitle,
        "company_name":        newWorkExperience.CompanyName,
        "employment_type":     newWorkExperience.EmploymentType,
        "start_date":          newWorkExperience.StartDate,
        "end_date":            newWorkExperience.EndDate,
        "key_responsibilities": newWorkExperience.KeyResponsibilities,
    }

    // Append the work experience to the array
    seeker.WorkExperiences = append(seeker.WorkExperiences, workExperienceBson)

    return nil
}



// func GetLanguages(seeker *models.Seeker) ([]dto.LanguageRequest, error) {
// 	var languages []dto.LanguageRequest
// 	if seeker.Languages == nil {
// 		return nil, errors.New("languages are nil")
// 	}
// 	languagesData, err := bson.Marshal(seeker.Languages)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := UnmarshalBsonToArray(languagesData, &languages); err != nil {
// 		return nil, err
// 	}
// 	return languages, nil
// }




// func SetLanguages(seeker *models.Seeker, languages []dto.LanguageRequest) error {
// 	data, err := MarshalArrayToBson(languages)
// 	if err != nil {
// 		return err
// 	}
// 	var bsonData bson.M
// 	if err := bson.Unmarshal(data, &bsonData); err != nil {
// 		return err
// 	}
// 	seeker.Languages = bsonData
// 	return nil
// }

// func GetCertificates(seeker *models.Seeker) ([]dto.CertificateRequest, error) {
// 	var certificates []dto.CertificateRequest
// 	if seeker.Certificates == nil {
// 		return nil, errors.New("certificates are nil")
// 	}
// 	certificatesData, err := bson.Marshal(seeker.Certificates)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := UnmarshalBsonToArray(certificatesData, &certificates); err != nil {
// 		return nil, err
// 	}
// 	return certificates, nil
// }


// func SetCertificates(seeker *models.Seeker, certificates []dto.CertificateRequest) error {
// 	data, err := MarshalArrayToBson(certificates)
// 	if err != nil {
// 		return err
// 	}
// 	var bsonData bson.M
// 	if err := bson.Unmarshal(data, &bsonData); err != nil {
// 		return err
// 	}

// 	seeker.Certificates = bsonData
// 	return nil
// }

// func GetEducations(seeker *models.Seeker) ([]dto.EducationRequest, error) {
// 	var educations []dto.EducationRequest
// 	if seeker.Educations == nil {
// 		return nil, errors.New("educations are nil")
// 	}
// 	educationsData, err := bson.Marshal(seeker.Educations)
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := UnmarshalBsonToArray(educationsData, &educations); err != nil {
// 		return nil, err
// 	}
// 	return educations, nil
// }

// func SetEducations(seeker *models.Seeker, educations []dto.EducationRequest) error {
// 	data, err := MarshalArrayToBson(educations)
// 	if err != nil {
// 		return err
// 	}
// 	var bsonData bson.M
// 	if err := bson.Unmarshal(data, &bsonData); err != nil {
// 		return err
// 	}
// 	seeker.Educations = bsonData
// 	return nil
// }
