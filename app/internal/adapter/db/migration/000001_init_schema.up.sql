CREATE TABLE "users" (
  "id" bigserial PRIMARY KEY,
  "email" varchar UNIQUE,
  "username" varchar
);

CREATE TABLE "messages" (
    "id" bigserial PRIMARY KEY,
    "sender_id" bigserial REFERENCES users ("id"),
    "receiver_id" bigserial REFERENCES users ("id"),
    "text" text,
    "status" varchar,
    "created_at" timestamp
);

CREATE TABLE  "user_password" (
    "user_id" bigserial REFERENCES users ("id"),
    "password" varchar
);

CREATE TABLE "oauth2_info" (
  "user_id" bigserial REFERENCES users ("id"),
  "provider" varchar,
  "provider_id" bigserial,
  UNIQUE ("user_id", "provider")
);

CREATE TABLE "sessions" (
  "session_id" bigserial PRIMARY KEY,
  "session_token" varchar UNIQUE,
  "user_id" bigserial REFERENCES users ("id"),
  "expiry" timestamp
);