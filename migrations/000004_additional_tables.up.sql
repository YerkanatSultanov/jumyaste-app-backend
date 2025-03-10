ALTER TABLE vacancies ADD COLUMN search_vector tsvector;

CREATE INDEX idx_vacancies_search ON vacancies USING GIN(search_vector);

CREATE FUNCTION update_search_vector() RETURNS trigger AS $$
BEGIN
    NEW.search_vector = to_tsvector('russian', NEW.title || ' ' || NEW.description || ' ' || array_to_string(NEW.skills, ' '));
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_update_search_vector
    BEFORE INSERT OR UPDATE ON vacancies
    FOR EACH ROW
EXECUTE FUNCTION update_search_vector();

