package models

import (
	//"gorm.io/driver/sqlserver"
	"gorm.io/gorm"

)
// DEVELOPMENT
func SeedJobs(db *gorm.DB) {
	// Check if jobs already seeded
	var count int64
	db.Model(&LinkedInJobMetaData{}).Count(&count)
	if count > 0 {
		return // Already seeded
	}

	// Seed LinkedIn jobs: only "Software Engineer" and "DevOps Engineer"
	linkedinJobs := []LinkedInJobMetaData{
		{ID: "lnk1", JobID: "L001", Title: "Software Engineer", Company: "LinkedIn", Location: "Berlin", PostedDate: "2024-04-01", Link: "https://linkedin.com/jobs/1", Processed: true},
		{ID: "lnk2", JobID: "L002", Title: "DevOps Engineer", Company: "Google", Location: "Munich", PostedDate: "2024-04-02", Link: "https://linkedin.com/jobs/2", Processed: true},
		{ID: "lnk3", JobID: "L003", Title: "Software Engineer", Company: "Meta", Location: "Hamburg", PostedDate: "2024-04-03", Link: "https://linkedin.com/jobs/3", Processed: true},
		{ID: "lnk4", JobID: "L004", Title: "DevOps Engineer", Company: "Amazon", Location: "Stuttgart", PostedDate: "2024-04-04", Link: "https://linkedin.com/jobs/4", Processed: true},
		{ID: "lnk5", JobID: "L005", Title: "Software Engineer", Company: "Apple", Location: "Frankfurt", PostedDate: "2024-04-05", Link: "https://linkedin.com/jobs/5", Processed: true},
		{ID: "lnk6", JobID: "L006", Title: "DevOps Engineer", Company: "IBM", Location: "Cologne", PostedDate: "2024-04-06", Link: "https://linkedin.com/jobs/6", Processed: true},
		{ID: "lnk7", JobID: "L007", Title: "Software Engineer", Company: "Netflix", Location: "Leipzig", PostedDate: "2024-04-07", Link: "https://linkedin.com/jobs/7", Processed: true},
		{ID: "lnk8", JobID: "L008", Title: "DevOps Engineer", Company: "Spotify", Location: "Dresden", PostedDate: "2024-04-08", Link: "https://linkedin.com/jobs/8", Processed: true},
		{ID: "lnk9", JobID: "L009", Title: "Software Engineer", Company: "Tesla", Location: "Bonn", PostedDate: "2024-04-09", Link: "https://linkedin.com/jobs/9", Processed: true},
		{ID: "lnk10", JobID: "L010", Title: "DevOps Engineer", Company: "Intel", Location: "Nuremberg", PostedDate: "2024-04-10", Link: "https://linkedin.com/jobs/10", Processed: true},
	}
	db.Create(&linkedinJobs)

	// Seed Xing jobs: same titles
	xingJobs := []XingJobMetaData{
		{ID: "xg1", JobID: "X001", Title: "Software Engineer", Company: "Xing", Location: "Hamburg", PostedDate: "2024-04-01", Link: "https://xing.com/jobs/1", Processed: true},
		{ID: "xg2", JobID: "X002", Title: "DevOps Engineer", Company: "SAP", Location: "Berlin", PostedDate: "2024-04-02", Link: "https://xing.com/jobs/2", Processed: true},
		{ID: "xg3", JobID: "X003", Title: "Software Engineer", Company: "Siemens", Location: "Munich", PostedDate: "2024-04-03", Link: "https://xing.com/jobs/3", Processed: true},
		{ID: "xg4", JobID: "X004", Title: "DevOps Engineer", Company: "Allianz", Location: "Cologne", PostedDate: "2024-04-04", Link: "https://xing.com/jobs/4", Processed: true},
		{ID: "xg5", JobID: "X005", Title: "Software Engineer", Company: "Bosch", Location: "Stuttgart", PostedDate: "2024-04-05", Link: "https://xing.com/jobs/5", Processed: true},
		{ID: "xg6", JobID: "X006", Title: "DevOps Engineer", Company: "Volkswagen", Location: "Wolfsburg", PostedDate: "2024-04-06", Link: "https://xing.com/jobs/6", Processed: true},
		{ID: "xg7", JobID: "X007", Title: "Software Engineer", Company: "BMW", Location: "Leipzig", PostedDate: "2024-04-07", Link: "https://xing.com/jobs/7", Processed: true},
		{ID: "xg8", JobID: "X008", Title: "DevOps Engineer", Company: "Deutsche Bank", Location: "Frankfurt", PostedDate: "2024-04-08", Link: "https://xing.com/jobs/8", Processed: true},
		{ID: "xg9", JobID: "X009", Title: "Software Engineer", Company: "Telekom", Location: "Bonn", PostedDate: "2024-04-09", Link: "https://xing.com/jobs/9", Processed: true},
		{ID: "xg10", JobID: "X010", Title: "DevOps Engineer", Company: "ZF Group", Location: "Friedrichshafen", PostedDate: "2024-04-10", Link: "https://xing.com/jobs/10", Processed: true},
	}
	db.Create(&xingJobs)

	// Seed LinkedIn application links
	linkedinAppLinks := []LinkedInJobApplicationLink{
		{JobID: "lnk1", JobLink: "https://apply.linkedin.com/job/1"},
		{JobID: "lnk2", JobLink: "https://apply.linkedin.com/job/2"},
		{JobID: "lnk3", JobLink: "https://apply.linkedin.com/job/3"},
		{JobID: "lnk4", JobLink: "https://apply.linkedin.com/job/4"},
		{JobID: "lnk5", JobLink: "https://apply.linkedin.com/job/5"},
		{JobID: "lnk6", JobLink: "https://apply.linkedin.com/job/6"},
		{JobID: "lnk7", JobLink: "https://apply.linkedin.com/job/7"},
		{JobID: "lnk8", JobLink: "https://apply.linkedin.com/job/8"},
		{JobID: "lnk9", JobLink: "https://apply.linkedin.com/job/9"},
		{JobID: "lnk10", JobLink: "https://apply.linkedin.com/job/10"},
	}
	db.Create(&linkedinAppLinks)

	// Seed Xing application links
	xingAppLinks := []XingJobApplicationLink{
		{JobID: "xg1", JobLink: "https://apply.xing.com/job/1"},
		{JobID: "xg2", JobLink: "https://apply.xing.com/job/2"},
		{JobID: "xg3", JobLink: "https://apply.xing.com/job/3"},
		{JobID: "xg4", JobLink: "https://apply.xing.com/job/4"},
		{JobID: "xg5", JobLink: "https://apply.xing.com/job/5"},
		{JobID: "xg6", JobLink: "https://apply.xing.com/job/6"},
		{JobID: "xg7", JobLink: "https://apply.xing.com/job/7"},
		{JobID: "xg8", JobLink: "https://apply.xing.com/job/8"},
		{JobID: "xg9", JobLink: "https://apply.xing.com/job/9"},
		{JobID: "xg10", JobLink: "https://apply.xing.com/job/10"},
	}
	db.Create(&xingAppLinks)

	// Seed LinkedIn job descriptions
	linkedinDescriptions := []LinkedInJobDescription{
		{JobID: "lnk1", JobLink: "https://apply.linkedin.com/job/1", JobDescription: "We are looking for a skilled Software Engineer to build scalable systems."},
		{JobID: "lnk2", JobLink: "https://apply.linkedin.com/job/2", JobDescription: "Join our DevOps team to manage CI/CD pipelines and cloud infrastructure."},
		{JobID: "lnk3", JobLink: "https://apply.linkedin.com/job/3", JobDescription: "Develop backend services with Go and microservices architecture."},
		{JobID: "lnk4", JobLink: "https://apply.linkedin.com/job/4", JobDescription: "Automate infrastructure with Terraform and Kubernetes."},
		{JobID: "lnk5", JobLink: "https://apply.linkedin.com/job/5", JobDescription: "Contribute to the core platform used by millions of users."},
		{JobID: "lnk6", JobLink: "https://apply.linkedin.com/job/6", JobDescription: "Maintain and scale cloud-based infrastructure for large applications."},
		{JobID: "lnk7", JobLink: "https://apply.linkedin.com/job/7", JobDescription: "Collaborate with cross-functional teams to build software solutions."},
		{JobID: "lnk8", JobLink: "https://apply.linkedin.com/job/8", JobDescription: "Improve system reliability and monitoring in a fast-paced DevOps team."},
		{JobID: "lnk9", JobLink: "https://apply.linkedin.com/job/9", JobDescription: "Participate in the development of highly available distributed systems."},
		{JobID: "lnk10", JobLink: "https://apply.linkedin.com/job/10", JobDescription: "Support infrastructure as code initiatives across all teams."},
	}

	// Seed Xing job descriptions
	xingDescriptions := []XingJobDescription{
		{JobID: "xg1", JobLink: "https://apply.xing.com/job/1", JobDescription: "Join Xing as a Software Engineer working on scalable APIs."},
		{JobID: "xg2", JobLink: "https://apply.xing.com/job/2", JobDescription: "DevOps position focusing on automation and observability."},
		{JobID: "xg3", JobLink: "https://apply.xing.com/job/3", JobDescription: "Help modernize our legacy systems into cloud-native services."},
		{JobID: "xg4", JobLink: "https://apply.xing.com/job/4", JobDescription: "Implement CI/CD pipelines and improve deployment efficiency."},
		{JobID: "xg5", JobLink: "https://apply.xing.com/job/5", JobDescription: "Work closely with data teams to deploy scalable services."},
		{JobID: "xg6", JobLink: "https://apply.xing.com/job/6", JobDescription: "Manage cloud resources and write automation scripts for operations."},
		{JobID: "xg7", JobLink: "https://apply.xing.com/job/7", JobDescription: "Design backend solutions to handle high-volume user traffic."},
		{JobID: "xg8", JobLink: "https://apply.xing.com/job/8", JobDescription: "Ensure uptime and performance of mission-critical applications."},
		{JobID: "xg9", JobLink: "https://apply.xing.com/job/9", JobDescription: "Develop high-performance APIs using Go and modern tech stack."},
		{JobID: "xg10", JobLink: "https://apply.xing.com/job/10", JobDescription: "Drive DevOps practices across teams and improve deployment workflows."},
	}

	db.Create(&linkedinDescriptions)
	db.Create(&xingDescriptions)

	}
