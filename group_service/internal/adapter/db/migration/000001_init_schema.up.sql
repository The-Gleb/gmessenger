CREATE TABLE "messages" (
  "id" bigserial PRIMARY KEY,
  "sender" varchar,
  "group_id" bigserial,
  "text" text,
  "status" varchar,
  "created_at" timestamp
);

CREATE TABLE "members" (
  "group_id" bigserial,
  "member_login" varchar
);

CREATE TABLE "groups" (
  "id" bigserial PRIMARY KEY,
  "name" varchar,
  "created_at" timestamp
);

ALTER TABLE "messages" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id");

ALTER TABLE "members" ADD FOREIGN KEY ("group_id") REFERENCES "groups" ("id");

