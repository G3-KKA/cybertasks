CREATE INDEX task_timestamp_idx 
    ON tasks
    USING brin (created_at)
;