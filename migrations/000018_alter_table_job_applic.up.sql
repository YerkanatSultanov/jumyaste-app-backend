ALTER TABLE job_applications
    ADD COLUMN resume_id INT,
ADD CONSTRAINT fk_resume_id FOREIGN KEY (resume_id) REFERENCES resume(id) ON
DELETE
SET NULL;

ALTER TABLE job_applications
    ADD COLUMN ai_matching_score INTEGER;
