CREATE TABLE IF NOT EXISTS "deductions" (
    id SERIAL PRIMARY KEY,
    "name" TEXT UNIQUE,
    maxAmount REAL
);
INSERT INTO "deductions" ("name", maxAmount)
VALUES ('personal', 60000.0),
('k-receipt', 50000.0);