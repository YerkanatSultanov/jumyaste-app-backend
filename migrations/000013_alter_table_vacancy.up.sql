ALTER TABLE vacancies ADD COLUMN company_id INTEGER;

UPDATE vacancies SET company_id = 1;

ALTER TABLE vacancies ALTER COLUMN company_id SET NOT NULL;
ALTER TABLE vacancies ADD CONSTRAINT fk_company FOREIGN KEY (company_id) REFERENCES companies(id) ON DELETE CASCADE;
