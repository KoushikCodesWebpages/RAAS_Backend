#!/bin/bash

curl -X POST https://api-inference.huggingface.co/models/nouamanetazi/cover-letter-t5-base \
  -H "Authorization: Bearer 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" \
  -H "Content-Type: application/json" \
  -d "{\"inputs\": \"Write the body of a professional cover letter for a candidate named John Michael applying for the Software Engineer position at Xing. Education: Degree: Bachelor of Science in Computer Science from XYZ University in Computer Science. Achievements: Graduated with honors, President of Coding Club. Experience: Job Title: Software Engineer at TechNova Inc.. Responsibilities: Developed backend services, collaborated with frontend team, maintained CI/CD pipelines. Skills: [Go React PostgreSQL Docker AWS]. Key Achievements: Degree: Bachelor of Science in Computer Science from XYZ University in Computer Science. Achievements: Graduated with honors, President of Coding Club. The body should highlight the candidate's qualifications, experience, skills, and how they align with the company and the role. Avoid the introduction and closing.\"}"
