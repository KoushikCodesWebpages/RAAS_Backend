curl -s -X POST https://cvncover-production.up.railway.app/generate/resume \
-H "Authorization: Bearer 9b1bf2a8ff2de605dce2131859d5cf0925a4fe9189c3dccfcdcfa7c933937664" \
-H "Content-Type: application/json" \
-d @- --output resume.pdf <<EOF
{
 "user_details": {
  "name": "Alex Johnson",
  "designation": "Software Engineer",
  "address": "Prenzlauer Allee 172, Berlin 10409",
  "contact": "+49 17624931591",
  "email": "Ramani.mallempuri@gmail.com",
  "portfolio": "www.reallygreatsite.com",
  "linkedin": "www.linkedin.com/alex",
  "tools": ["Golang", "React", "PostgreSQL", "Docker", "Kubernetes", "Python", "AI/ML"],
  "skills": ["Project Management", "Public Relations", "Teamwork", "Time Management", "Leadership", "Effective Communication", "Critical Thinking"],
  "education": [
        "Master of Computer Science Engineering, BORCELLE UNIVERSITY, 2029 - 2030",
        "Bachelor of Computer Science Engineering, BORCELLE UNIVERSITY, 2025 - 2029, GPA: 3.8 / 4.0"
    ],
  "experience_summary": [
        "Software Engineer at TechCorp (2019-2023): Led full-stack development for SaaS platform.",
        "Backend Developer at StartUpXYZ (2018-2019): Built REST APIs and microservices."
    ],
  "certifications": [
        "AWS Certified Developer",
        "ML Certified Developer"
    ],
  "languages": ["English: Fluent", "French: Fluent", "German: Basics", "Spanish: Intermediate"]
 },
  "job_description": {
    "job_title": "Product Manager",
    "company": "GimmeMore",
    "location": "Berlin, DE",
    "job_type": "Full-Time",
    "responsibilities": [
      "Support the development and optimization of the personality assessment platform TerraYou",
      "Shape and execute product vision for TerraYou",
      "Define strategy and roadmap",
      "Translate product vision into execution-level requirements including measurable KPIs",
      "Improve retention, conversion, and monetization by formulating hypotheses, conducting tests, and iteratively optimizing strategies",
      "Identify cost-effective and high-impact opportunities that expand TerraYou’s product offering and elevate its position in the market"
    ],
    "qualifications": [
      "At least 3 years experience as a product manager, ideally with a focus on the performance of an international website and/or mobile application",
      "Strong understanding of technical principles, data analysis tools, and A/B testing",
      "An affinity for project management and experience working in an agile environment",
      "Extremely talented communicator and enjoy working with a team and with business stakeholders",
      "Comfortable taking ownership of tasks, being a source of knowledge for your product, and taking the lead", 
      "Fluent in written and spoken English"
    ],
    "skills": [
      "Technical principles",
      "Data analysis tools",
      "A/B testing",
      "Project management",
      "Communication",
      "Leadership"
    ],
    "benefits": [
      "Shape the company’s future and work on expanding TerraYou brand within a diverse, talented, and international team",
      "Generous educational budget, sponsored gym membership through Urban Sports, and regular team sports activities",
      "Performance-based bonus scheme",
      "Various team events throughout the year",
      "Home-office friendly employer"
    ]
  }
}
EOF