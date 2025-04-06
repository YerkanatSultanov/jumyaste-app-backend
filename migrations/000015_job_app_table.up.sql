CREATE TABLE job_applications
(
    id         SERIAL PRIMARY KEY,
    user_id    INT          NOT NULL,
    vacancy_id INT          NOT NULL,
    first_name VARCHAR(255) NOT NULL,
    last_name  VARCHAR(255) NOT NULL,
    email      VARCHAR(255) NOT NULL,
    status     VARCHAR(50) CHECK (status IN ('new', 'invited', 'interview', 'accepted', 'rejected')) DEFAULT 'new',
    applied_at TIMESTAMP                                                                             DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
    FOREIGN KEY (vacancy_id) REFERENCES vacancies (id) ON DELETE CASCADE
);
