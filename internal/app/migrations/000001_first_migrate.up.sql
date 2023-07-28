CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "users" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "first_name" varchar(255),
    "last_name" varchar(255),
    "email" varchar(255),
    "username" varchar(50) UNIQUE,
    "password" varchar(255),
    "role" varchar(50) DEFAULT 'User',
    "avatar" varchar(255),
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP
);


CREATE TABLE "edu_centers" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" varchar(255) UNIQUE,
    "html_description" text,
    "address" varchar(255),
    "location" POINT,
    "owner_id" uuid REFERENCES "users" ("id"),
    "cover_image" varchar(250),
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" TIMESTAMP
);


CREATE TABLE "contacts" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "instagram" varchar(255),
    "telegram" varchar(255),
    "website" varchar(255),
    "phone_number" varchar(50),
    "edu_center_id" uuid REFERENCES "edu_centers" ("id"),
    "user_id" uuid REFERENCES "users" ("id"),
    CHECK (
        (edu_center_id IS NOT NULL AND user_id IS NULL) OR
        (edu_center_id IS NULL AND user_id IS NOT NULL)
    )
);


CREATE TABLE "courses" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "name" varchar(255),
    "description" text,
    "teacher" varchar(255),
    "edu_center_id" uuid REFERENCES "edu_centers" ("id"),
    "created_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "updated_at" TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    "deleted_at" timestamp
);

CREATE TABLE "edu_center_images" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "edu_center_id" uuid REFERENCES "edu_centers" ("id"),
    "image_link" varchar(255)
);

CREATE TABLE "ratings" (
    "id" uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
    "score" int,
    "owner_id" uuid REFERENCES "users" ("id"),
    "edu_center_id" uuid REFERENCES "edu_centers" ("id"),
    "course_id" uuid REFERENCES "courses" ("id"),
    CHECK (
        (edu_center_id IS NOT NULL AND course_id IS NULL) OR
        (edu_center_id IS NULL AND course_id IS NOT NULL)
    )
);
