package models

import (
	"gorm.io/gorm"
	"log"
)

// SeedJobs - Seeds the Job model with metadata, description, and application links
func SeedJobs(db *gorm.DB) {
	// Check if jobs already seeded
	var count int64
	db.Model(&Job{}).Count(&count)
	if count > 0 {
		return // Already seeded
	}

	// Seed LinkedIn and Xing jobs into the main Job model
	jobs := []Job{
		// LinkedIn Jobs
		{JobID: "L001", Title: "Software Engineer", Company: "LinkedIn", Location: "Berlin", PostedDate: "2024-04-01", Link: "https://linkedin.com/jobs/1", Processed: true, Source: "LinkedIn", JobDescription: "We are looking for a skilled Software Engineer to build scalable systems.", JobType: "Full-time", Skills: "Go, REST, Microservices, Docker", JobLink: "https://apply.linkedin.com/job/1"},
		{JobID: "L002", Title: "DevOps Engineer", Company: "Google", Location: "Munich", PostedDate: "2024-04-02", Link: "https://linkedin.com/jobs/2", Processed: true, Source: "LinkedIn", JobDescription: "Join our DevOps team to manage CI/CD pipelines and cloud infrastructure.", JobType: "Full-time", Skills: "CI/CD, Jenkins, AWS, Docker, Kubernetes", JobLink: "https://apply.linkedin.com/job/2"},
		{JobID: "L003", Title: "Software Engineer", Company: "Meta", Location: "Hamburg", PostedDate: "2024-04-03", Link: "https://linkedin.com/jobs/3", Processed: true, Source: "LinkedIn", JobDescription: "Develop backend services with Go and microservices architecture.", JobType: "Remote", Skills: "Go, gRPC, PostgreSQL, Docker", JobLink: "https://apply.linkedin.com/job/3"},
		{JobID: "L004", Title: "DevOps Engineer", Company: "Amazon", Location: "Stuttgart", PostedDate: "2024-04-04", Link: "https://linkedin.com/jobs/4", Processed: true, Source: "LinkedIn", JobDescription: "Automate infrastructure with Terraform and Kubernetes.", JobType: "Contract", Skills: "Terraform, Kubernetes, Helm, AWS", JobLink: "https://apply.linkedin.com/job/4"},
		{JobID: "L005", Title: "Software Engineer", Company: "Apple", Location: "Frankfurt", PostedDate: "2024-04-05", Link: "https://linkedin.com/jobs/5", Processed: true, Source: "LinkedIn", JobDescription: "Contribute to the core platform used by millions of users.", JobType: "Full-time", Skills: "Go, Redis, Kafka, Prometheus", JobLink: "https://apply.linkedin.com/job/5"},
		{JobID: "L006", Title: "DevOps Engineer", Company: "IBM", Location: "Cologne", PostedDate: "2024-04-06", Link: "https://linkedin.com/jobs/6", Processed: true, Source: "LinkedIn", JobDescription: "Maintain and scale cloud-based infrastructure for large applications.", JobType: "Remote", Skills: "AWS, Terraform, Docker, Monitoring", JobLink: "https://apply.linkedin.com/job/6"},
		
		// Xing Jobs
		{JobID: "X001", Title: "Software Engineer", Company: "Xing", Location: "Hamburg", PostedDate: "2024-04-01", Link: "https://xing.com/jobs/1", Processed: true, Source: "Xing", JobDescription: "Join Xing as a Software Engineer working on scalable APIs.", JobType: "Full-time", Skills: "Go, REST APIs, MySQL, Docker", JobLink: "https://apply.xing.com/job/1"},
		{JobID: "X002", Title: "DevOps Engineer", Company: "SAP", Location: "Berlin", PostedDate: "2024-04-02", Link: "https://xing.com/jobs/2", Processed: true, Source: "Xing", JobDescription: "DevOps position focusing on automation and observability.", JobType: "Remote", Skills: "CI/CD, Grafana, Prometheus, Bash, Terraform", JobLink: "https://apply.xing.com/job/2"},
		{JobID: "X003", Title: "Software Engineer", Company: "Siemens", Location: "Munich", PostedDate: "2024-04-03", Link: "https://xing.com/jobs/3", Processed: true, Source: "Xing", JobDescription: "Help modernize our legacy systems into cloud-native services.", JobType: "Full-time", Skills: "AWS, Kubernetes, Go, Monolith Refactoring", JobLink: "https://apply.xing.com/job/3"},
		{JobID: "X004", Title: "DevOps Engineer", Company: "Allianz", Location: "Cologne", PostedDate: "2024-04-04", Link: "https://xing.com/jobs/4", Processed: true, Source: "Xing", JobDescription: "Implement CI/CD pipelines and improve deployment efficiency.", JobType: "Part-time", Skills: "GitLab CI, Docker, Helm, Kubernetes", JobLink: "https://apply.xing.com/job/4"},
	}
	

	// Insert the jobs into the database
	if err := db.Create(&jobs).Error; err != nil {
        log.Printf("Error seeding jobs: %v", err)
    }
	log.Printf("Successfully seeeded jobs")
}
