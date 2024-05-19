CREATE TABLE "messages" (
  "id" bigserial PRIMARY KEY,
  "sender_id" varchar,
  "group_id" bigserial,
  "text" text,
  "status" varchar,
  "created_at" timestamp
);

CREATE TABLE "group_user" (
  "group_id" bigserial,
  "user_id" varchar
);

CREATE TABLE "groups" (
  "id" bigserial PRIMARY KEY,
  "name" varchar,
  "created_at" timestamp
);

ALTER TABLE "messages" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id");

ALTER TABLE "group_user" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id");

