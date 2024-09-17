CREATE TABLE "contacts"(
    "contact_id" SERIAL PRIMARY KEY,
    "first_name" VARCHAR NOT NULL,
    "last_name" VARCHAR NOT NULL,
    "phone_number" VARCHAR NOT NULL,
    "street" VARCHAR NOT NULL,
    "created_at" TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP NOT NULL
);

-- CREATE UNIQUE INDEX ON "contacts" ("phone_number");
CREATE UNIQUE INDEX "contacts_phone_number_key" ON "contacts"("phone_number");