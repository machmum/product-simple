-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE "public"."product_source"
(
    "id"          UUID           NOT NULL DEFAULT uuid_generate_v4(),
    "name"        VARCHAR(255)   NOT NULL,
    "sku"         VARCHAR(255)            DEFAULT NULL,
    "description" TEXT           DEFAULT NULL,
    "created_by"  UUID                    DEFAULT NULL,
    "created_at"  TIMESTAMPTZ(6) NOT NULL DEFAULT NOW(),
    "updated_by"  UUID                    DEFAULT NULL,
    "updated_at"  TIMESTAMPTZ(6)          DEFAULT NULL,
    PRIMARY KEY ("id")
);

INSERT INTO "public"."product_source" ("id", "name", "sku", "description", "created_by", "created_at")
VALUES (DEFAULT, 'product-temp', 'product-temp', 'product description', '5afff7e2-301e-4161-a424-066a93a75b88', NOW());

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE IF EXISTS "public"."product_source" CASCADE;
