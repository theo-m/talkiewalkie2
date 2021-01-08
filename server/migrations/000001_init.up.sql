CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE "user"
(
    "id"           SERIAL PRIMARY KEY,
    "uuid"         UUID                     NOT NULL DEFAULT "uuid_generate_v4"(),
    "username"     TEXT                     NOT NULL UNIQUE,
    "email"        TEXT UNIQUE,
    "password"     TEXT,
    "token"        TEXT,
    "created_at"   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "updated_at"   TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "connected_at" TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TABLE "user_user"
(
    "id"             SERIAL PRIMARY KEY,
    "user_source_id" INTEGER REFERENCES "user" ("id") ON DELETE CASCADE NOT NULL,
    "user_target_id" INTEGER REFERENCES "user" ("id") ON DELETE CASCADE NOT NULL,
    UNIQUE ("user_source_id", "user_target_id")
);

CREATE TABLE "walk"
(
    "id"          SERIAL PRIMARY KEY,
    "uuid"        UUID                     NOT NULL DEFAULT "uuid_generate_v4"(),
    "title"       TEXT                     NOT NULL,
    "description" TEXT,
    "tags"        TEXT[]                   NOT NULL,
    "author_id"   INTEGER                  REFERENCES "user" ("id") ON DELETE SET NULL,
    "is_private"  BOOLEAN                  NOT NULL,
    "created_at"  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    "updated_at"  TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now()
);

CREATE TYPE USER_WALK_EDGE_TYPE AS ENUM ('LIKE', 'DISLIKE', 'SHARED');

CREATE TABLE "user_walk"
(
    "id"      SERIAL PRIMARY KEY,
    "user_id" INTEGER REFERENCES "user" ("id") ON DELETE CASCADE NOT NULL,
    "walk_id" INTEGER REFERENCES "walk" ("id") ON DELETE CASCADE NOT NULL,
    "type"    "user_walk_edge_type"                              NOT NULL,
    UNIQUE ("user_id", "walk_id", "type")
);

CREATE TABLE "walk_point"
(
    "id"         SERIAL PRIMARY KEY,
    "uuid"       UUID                                               NOT NULL DEFAULT "uuid_generate_v4"(),
    "text"       TEXT                                               NOT NULL,
    "walk_id"    INTEGER REFERENCES "walk" ("id") ON DELETE CASCADE NOT NULL,
    "created_at" TIMESTAMP WITH TIME ZONE                           NOT NULL DEFAULT now(),
    "updated_at" TIMESTAMP WITH TIME ZONE                           NOT NULL DEFAULT now()
);

CREATE TYPE ASSET_TYPE AS ENUM ('upload', 'text');

-- gps points:
-- https://stackoverflow.com/questions/8150721/which-DATA-TYPE-FOR-latitude-AND-longitude
CREATE TABLE "walk_point_asset"
(
    "id"            SERIAL PRIMARY KEY,
    "walk_point_id" INTEGER REFERENCES "walk_point" ("id") ON DELETE CASCADE NOT NULL,
    "author_id"     INTEGER                                                  REFERENCES "user" ("id") ON DELETE SET NULL,

    "asset_type"    ASSET_TYPE                                               NOT NULL,

    "bucket"        TEXT,
    "blob_path"     TEXT,
    "filename"      TEXT,

    "text"          TEXT,
    "created_at"    TIMESTAMP WITH TIME ZONE                                 NOT NULL DEFAULT now(),
    "updated_at"    TIMESTAMP WITH TIME ZONE                                 NOT NULL DEFAULT now(),

    CONSTRAINT "upload_asset_needs_blob_info" CHECK ( "asset_type" <> 'upload' OR
                                                      ("bucket" IS NOT NULL AND "blob_path" IS NOT NULL AND "filename" IS NOT NULL) ),
    CONSTRAINT "text_asset_needs_text_field" CHECK ( "asset_type" <> 'text' OR "text" IS NOT NULL)
);
