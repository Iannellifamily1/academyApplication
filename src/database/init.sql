CREATE TABLE "Staff" (
    "ID" SERIAL PRIMARY KEY,
    "Email" VARCHAR(255) NOT NULL UNIQUE,
    "Name" VARCHAR(255) NOT NULL,
    "Password" VARCHAR(255) NOT NULL,
    "Salt" VARCHAR(255),
    "Token" VARCHAR(255)
);

CREATE TABLE "Role" (
    "ID" SERIAL PRIMARY KEY,
    "Name" VARCHAR(255) NOT NULL UNIQUE
);

CREATE TABLE "StaffWithRoles" (
    "staff_id" INT REFERENCES "Staff"("ID") ON DELETE CASCADE,
    "role_id" INT REFERENCES "Role"("ID") ON DELETE CASCADE,
    PRIMARY KEY (staff_id, role_id)
);
