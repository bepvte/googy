CREATE TABLE IF NOT EXISTS "perms" (
  "priority" int NOT NULL,
  "state" boolean NOT NULL,
  "where" text not null,
  "what" text not null,
  "type" int not null,
  "comment" text not null,
  "guild" text not null
);

CREATE TABLE IF NOT EXISTS store (myid text PRIMARY KEY);