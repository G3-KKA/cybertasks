CREATE ROLE bookuser PASSWORD NULL;

GRANT SELECT , INSERT , UPDATE , DELETE
    ON tasks
    TO bookuser 
    GRANTED BY postgres
;