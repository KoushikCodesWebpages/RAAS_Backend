package repository

import (
	"encoding/json"
	"errors"

	"RAAS/internal/dto"
	"RAAS/internal/models"
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


// GetEducation retrieves the education information of the seeker
func GetEducation(seeker *models.Seeker) ([]bson.M, error) {
    if len(seeker.Education) == 0 {
        return []bson.M{}, nil
    }
    return seeker.Education, nil
}

// SetEducation sets the education information for a Seeker using an array of bson.M.
func SetEducation(seeker *models.Seeker, educations []bson.M) error {
    seeker.Education = educations
    return nil
}

// AppendToEducation adds a new education entry to the Seeker's education list
func AppendToEducation(seeker *models.Seeker, newEducation dto.EducationRequest) error {
    // Check if the Educations array is nil or empty, if so, initialize it
    if seeker.Education == nil {
        seeker.Education = []bson.M{}
    }

    // Create a new education entry as a bson.M document
    educationBson := bson.M{
        "degree":        newEducation.Degree,
        "institution":   newEducation.Institution,
        "field_of_study": newEducation.FieldOfStudy,
        "start_date":    newEducation.StartDate,
        "end_date":      newEducation.EndDate,
        "achievements":  newEducation.Achievements,
    }

    // Append the new education entry to the Educations array
    seeker.Education = append(seeker.Education, educationBson)

    return nil
}

// GetCertificates retrieves the certificate information of the seeker
func GetCertificates(seeker *models.Seeker) ([]bson.M, error) {
	if len(seeker.Certificates) == 0 {
		return []bson.M{}, nil
	}
	return seeker.Certificates, nil
}

// SetCertificates sets the certificate information for a Seeker using an array of bson.M
func SetCertificates(seeker *models.Seeker, certificates []bson.M) error {
	seeker.Certificates = certificates
	return nil
}

// AppendToCertificates adds a new certificate entry to the Seeker's certificates list
func AppendToCertificates(seeker *models.Seeker, newCertificate dto.CertificateRequest, certificateFile string) error {
	// Check if the Certificates array is nil or empty, if so, initialize it
	if seeker.Certificates == nil {
		seeker.Certificates = []bson.M{}
	}

	// Create a new certificate entry as a bson.M document
	certificateBson := bson.M{
		"certificate_name":   newCertificate.CertificateName,
		"certificate_file":   certificateFile,
	}

	// Only add certificate_number if it's not nil
	if newCertificate.CertificateNumber != nil {
		certificateBson["certificate_number"] = *newCertificate.CertificateNumber
	}

	// Append the new certificate entry to the Certificates array
	seeker.Certificates = append(seeker.Certificates, certificateBson)

	return nil
}

// GetLanguages retrieves the language information of the seeker
func GetLanguages(seeker *models.Seeker) ([]bson.M, error) {
    if len(seeker.Languages) == 0 {
        return []bson.M{}, nil
    }
    return seeker.Languages, nil
}

// SetLanguages sets the language information for a Seeker using an array of bson.M
func SetLanguages(seeker *models.Seeker, languages []bson.M) error {
    seeker.Languages = languages
    return nil
}
// AppendToLanguages adds a new language entry to the Seeker's languages list
func AppendToLanguages(seeker *models.Seeker, newLanguage dto.LanguageRequest, languageFile string) error {
    // Check if the Languages array is nil or empty, if so, initialize it
    if seeker.Languages == nil {
        seeker.Languages = []bson.M{}
    }

    // Create a new language entry as a bson.M document
    languageBson := bson.M{
        "language":         newLanguage.LanguageName,
        "proficiency":      newLanguage.ProficiencyLevel,
        "certificate_file": languageFile,
    }

    // Append the new language entry to the Languages array
    seeker.Languages = append(seeker.Languages, languageBson)

    return nil
}
