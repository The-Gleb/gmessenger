CREATE TABLE "messages" (
  "id" bigserial PRIMARY KEY,
  "sender" varchar,
  "receiver" varchar,
  "text" text,
  "status" varchar,
  "created_at" timestamp
);

CREATE TABLE "users" (
  "login" varchar PRIMARY KEY,
  "username" varchar,
  "password" varchar
);

CREATE TABLE "sessions" (
  "token" varchar PRIMARY KEY,
  "user_login" varchar,
  "expiry" timestamp
);

ALTER TABLE "messages" ADD FOREIGN KEY ("sender") REFERENCES "users" ("login");

ALTER TABLE "messages" ADD FOREIGN KEY ("receiver") REFERENCES "users" ("login");

ALTER TABLE "sessions" ADD FOREIGN KEY ("user_login") REFERENCES "users" ("login");