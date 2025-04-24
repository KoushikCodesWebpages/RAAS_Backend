package config

import "github.com/spf13/viper"

type HuggingFaceConfig struct {
	    // New fields for Hugging Face models and API key
		HFModelForMS1               string
		HFModelForMS2               string
		HFModelForMS3               string
		HFModelForMS4               string
		HFModelForMS5               string
		HFModelForMS6               string
		HFModelForMS7               string
		HFModelForMS8               string
		HFModelForMS9               string
		HFModelForMS10              string
	
		// CV MODELS
		HFModelForCL1               string
		HFModelForCL2               string
		HFModelForCL3               string
		HFModelForCL4               string
		HFModelForCL5               string
		HFModelForCL6               string
		HFModelForCL7               string
		HFModelForCL8               string
		HFModelForCL9               string
		HFModelForCL10              string
	
		HFBaseAPIUrl                string
	
		HFAPIKey                    string

}

func LoadHuggingFaceConfig() (*HuggingFaceConfig, error) {
	HuggingFaceConfig := &HuggingFaceConfig{

	// Load Hugging Face model and API key values from environment variables
	HFModelForMS1:               viper.GetString("HF_MODEL_FOR_MS_1"),
	HFModelForMS2:               viper.GetString("HF_MODEL_FOR_MS_2"),
	HFModelForMS3:               viper.GetString("HF_MODEL_FOR_MS_3"),
	HFModelForMS4:               viper.GetString("HF_MODEL_FOR_MS_4"),
	HFModelForMS5:               viper.GetString("HF_MODEL_FOR_MS_5"),
	HFModelForMS6:               viper.GetString("HF_MODEL_FOR_MS_6"),
	HFModelForMS7:               viper.GetString("HF_MODEL_FOR_MS_7"),
	HFModelForMS8:               viper.GetString("HF_MODEL_FOR_MS_8"),
	HFModelForMS9:               viper.GetString("HF_MODEL_FOR_MS_9"),
	HFModelForMS10:              viper.GetString("HF_MODEL_FOR_MS_10"),
	
	//CV GEN

	HFModelForCL1: 				viper.GetString("HF_MODEL_FOR_CL_1"),
	HFModelForCL2: 				viper.GetString("HF_MODEL_FOR_CL_2"),
	HFModelForCL3: 				viper.GetString("HF_MODEL_FOR_CL_3"),
	HFModelForCL4: 				viper.GetString("HF_MODEL_FOR_CL_4"),
	HFModelForCL5: 				viper.GetString("HF_MODEL_FOR_CL_5"),
	HFModelForCL6: 				viper.GetString("HF_MODEL_FOR_CL_6"),
	HFModelForCL7: 				viper.GetString("HF_MODEL_FOR_CL_7"),
	HFModelForCL8: 				viper.GetString("HF_MODEL_FOR_CL_8"),
	HFModelForCL9: 				viper.GetString("HF_MODEL_FOR_CL_9"),
	HFModelForCL10: 			viper.GetString("HF_MODEL_FOR_CL_10"),


	HFAPIKey:              		viper.GetString("HF_API_KEY"),

	HFBaseAPIUrl:     			viper.GetString("HF_BASE_API_URL"),
		

    }

    return HuggingFaceConfig, nil
}
