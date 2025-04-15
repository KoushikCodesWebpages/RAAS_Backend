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
    {JobID: "L001", JobLink: "https://apply.linkedin.com/job/1"},
    {JobID: "L002", JobLink: "https://apply.linkedin.com/job/2"},
    {JobID: "L003", JobLink: "https://apply.linkedin.com/job/3"},
    {JobID: "L004", JobLink: "https://apply.linkedin.com/job/4"},
    {JobID: "L005", JobLink: "https://apply.linkedin.com/job/5"},
    {JobID: "L006", JobLink: "https://apply.linkedin.com/job/6"},
    {JobID: "L007", JobLink: "https://apply.linkedin.com/job/7"},
    {JobID: "L008", JobLink: "https://apply.linkedin.com/job/8"},
    {JobID: "L009", JobLink: "https://apply.linkedin.com/job/9"},
    {JobID: "L010", JobLink: "https://apply.linkedin.com/job/10"},
}
db.Create(&linkedinAppLinks)

// Seed Xing application links
xingAppLinks := []XingJobApplicationLink{
    {JobID: "X001", JobLink: "https://apply.xing.com/job/1"},
    {JobID: "X002", JobLink: "https://apply.xing.com/job/2"},
    {JobID: "X003", JobLink: "https://apply.xing.com/job/3"},
    {JobID: "X004", JobLink: "https://apply.xing.com/job/4"},
    {JobID: "X005", JobLink: "https://apply.xing.com/job/5"},
    {JobID: "X006", JobLink: "https://apply.xing.com/job/6"},
    {JobID: "X007", JobLink: "https://apply.xing.com/job/7"},
    {JobID: "X008", JobLink: "https://apply.xing.com/job/8"},
    {JobID: "X009", JobLink: "https://apply.xing.com/job/9"},
    {JobID: "X010", JobLink: "https://apply.xing.com/job/10"},
}
db.Create(&xingAppLinks)

	// Seed LinkedIn job descriptions
// Seed LinkedIn job descriptions with additional fields
linkedinDescriptions := []LinkedInJobDescription{
	{JobID: "L001", JobLink: "https://apply.linkedin.com/job/1", JobDescription: "We are looking for a skilled Software Engineer to build scalable systems.", JobType: "Full-time", Skills: "Go, REST, Microservices, Docker"},
	{JobID: "L002", JobLink: "https://apply.linkedin.com/job/2", JobDescription: "Join our DevOps team to manage CI/CD pipelines and cloud infrastructure.", JobType: "Full-time", Skills: "CI/CD, Jenkins, AWS, Docker, Kubernetes"},
	{JobID: "L003", JobLink: "https://apply.linkedin.com/job/3", JobDescription: "Develop backend services with Go and microservices architecture.", JobType: "Remote", Skills: "Go, gRPC, PostgreSQL, Docker"},
	{JobID: "L004", JobLink: "https://apply.linkedin.com/job/4", JobDescription: "Automate infrastructure with Terraform and Kubernetes.", JobType: "Contract", Skills: "Terraform, Kubernetes, Helm, AWS"},
	{JobID: "L005", JobLink: "https://apply.linkedin.com/job/5", JobDescription: "Contribute to the core platform used by millions of users.", JobType: "Full-time", Skills: "Go, Redis, Kafka, Prometheus"},
	{JobID: "L006", JobLink: "https://apply.linkedin.com/job/6", JobDescription: "Maintain and scale cloud-based infrastructure for large applications.", JobType: "Remote", Skills: "AWS, Terraform, Docker, Monitoring"},
	{JobID: "L007", JobLink: "https://apply.linkedin.com/job/7", JobDescription: "Collaborate with cross-functional teams to build software solutions.", JobType: "Full-time", Skills: "Go, Teamwork, APIs, SQL"},
	{JobID: "L008", JobLink: "https://apply.linkedin.com/job/8", JobDescription: "Improve system reliability and monitoring in a fast-paced DevOps team.", JobType: "Part-time", Skills: "Prometheus, Grafana, Alertmanager, On-call"},
	{JobID: "L009", JobLink: "https://apply.linkedin.com/job/9", JobDescription: "Participate in the development of highly available distributed systems.", JobType: "Full-time", Skills: "Go, Load Balancing, Kafka, Redis"},
	{JobID: "L010", JobLink: "https://apply.linkedin.com/job/10", JobDescription: "Support infrastructure as code initiatives across all teams.", JobType: "Contract", Skills: "Terraform, GitOps, Kubernetes, CI/CD"},
}


// Seed Xing job descriptions with additional fields
xingDescriptions := []XingJobDescription{
	{JobID: "X001", JobLink: "https://apply.xing.com/job/1", JobDescription: "Join Xing as a Software Engineer working on scalable APIs.", JobType: "Full-time", Skills: "Go, REST APIs, MySQL, Docker"},
	{JobID: "X002", JobLink: "https://apply.xing.com/job/2", JobDescription: "DevOps position focusing on automation and observability.", JobType: "Remote", Skills: "CI/CD, Grafana, Prometheus, Bash, Terraform"},
	{JobID: "X003", JobLink: "https://apply.xing.com/job/3", JobDescription: "Help modernize our legacy systems into cloud-native services.", JobType: "Full-time", Skills: "AWS, Kubernetes, Go, Monolith Refactoring"},
	{JobID: "X004", JobLink: "https://apply.xing.com/job/4", JobDescription: "Implement CI/CD pipelines and improve deployment efficiency.", JobType: "Part-time", Skills: "GitLab CI, Docker, Helm, Kubernetes"},
	{JobID: "X005", JobLink: "https://apply.xing.com/job/5", JobDescription: "Work closely with data teams to deploy scalable services.", JobType: "Full-time", Skills: "Go, Kafka, PostgreSQL, BigQuery"},
	{JobID: "X006", JobLink: "https://apply.xing.com/job/6", JobDescription: "Manage cloud resources and write automation scripts for operations.", JobType: "Full-time", Skills: "Terraform, AWS, Bash, Scripting"},
	{JobID: "X007", JobLink: "https://apply.xing.com/job/7", JobDescription: "Design backend solutions to handle high-volume user traffic.", JobType: "Remote", Skills: "Go, NATS, MongoDB, Scalability"},
	{JobID: "X008", JobLink: "https://apply.xing.com/job/8", JobDescription: "Ensure uptime and performance of mission-critical applications.", JobType: "Full-time", Skills: "Kubernetes, On-call, SRE, Logging"},
	{JobID: "X009", JobLink: "https://apply.xing.com/job/9", JobDescription: "Develop high-performance APIs using Go and modern tech stack.", JobType: "Contract", Skills: "Go, Echo, GORM, Redis"},
	{JobID: "X010", JobLink: "https://apply.xing.com/job/10", JobDescription: "Drive DevOps practices across teams and improve deployment workflows.", JobType: "Full-time", Skills: "DevOps, CI/CD, Kubernetes, Culture"},
}


db.Create(&linkedinDescriptions)
db.Create(&xingDescriptions)

}
