ALTER TABLE job_applications
DROP
CONSTRAINT fk_resume_id;

ALTER TABLE job_applications
DROP
COLUMN resume_id;

ALTER TABLE job_applications
DROP
COLUMN ai_matching_score;
