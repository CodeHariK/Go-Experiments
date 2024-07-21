-- Add new schema named "public"
CREATE SCHEMA IF NOT EXISTS "public";
-- Set comment to schema: "public"
COMMENT ON SCHEMA "public" IS 'standard public schema';
-- Create "goose_db_version" table
CREATE TABLE "public"."goose_db_version" ("id" serial NOT NULL, "version_id" bigint NOT NULL, "is_applied" boolean NOT NULL, "tstamp" timestamp NULL DEFAULT now(), PRIMARY KEY ("id"));
-- Create "products" table
CREATE TABLE "public"."products" ("id" serial NOT NULL, "product_name" character varying(255) NOT NULL, "category" character varying(255) NOT NULL, "price" numeric(10,4) NOT NULL, "currency" character varying(12) NOT NULL DEFAULT 'USD', PRIMARY KEY ("id"), CONSTRAINT "products_product_name_key" UNIQUE ("product_name"), CONSTRAINT "products_currency_check" CHECK ((currency)::text = ANY ((ARRAY['USD'::character varying, 'INR'::character varying, 'EUR'::character varying])::text[])));
-- Create index "idx_products_category_id" to table: "products"
CREATE INDEX "idx_products_category_id" ON "public"."products" ("category");
-- Create "seller" table
CREATE TABLE "public"."seller" ("id" serial NOT NULL, "name" character varying(255) NOT NULL, "location" character varying(255) NOT NULL, PRIMARY KEY ("id"));
-- Create "inventory" table
CREATE TABLE "public"."inventory" ("id" serial NOT NULL, "product_id" integer NOT NULL, "seller_id" integer NOT NULL, "quantity" integer NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "inventory_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "inventory_seller_id_fkey" FOREIGN KEY ("seller_id") REFERENCES "public"."seller" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "idx_inventory_product_id" to table: "inventory"
CREATE INDEX "idx_inventory_product_id" ON "public"."inventory" ("product_id");
-- Create index "idx_inventory_seller_id_id" to table: "inventory"
CREATE INDEX "idx_inventory_seller_id_id" ON "public"."inventory" ("seller_id");
-- Create "locations" table
CREATE TABLE "public"."locations" ("id" serial NOT NULL, "name" character varying(255) NOT NULL, "address" character varying(255) NOT NULL, "city" character varying(100) NOT NULL, "state" character varying(100) NOT NULL, "country" character varying(100) NOT NULL, "postal_code" character varying(20) NOT NULL, "latitude" numeric(9,6) NOT NULL, "longitude" numeric(9,6) NOT NULL, PRIMARY KEY ("id"));
-- Create index "idx_locations_latitude_longitude" to table: "locations"
CREATE INDEX "idx_locations_latitude_longitude" ON "public"."locations" ("latitude", "longitude");
-- Create "users" table
CREATE TABLE "public"."users" ("id" serial NOT NULL, "username" character varying(50) NOT NULL, "email" character varying(50) NOT NULL, "is_admin" boolean NULL DEFAULT false, "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, "date_of_birth" date NULL, "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP, "phone_number" character varying(15) NOT NULL, "last_login" timestamp NULL, "location_id" integer NULL, PRIMARY KEY ("id"), CONSTRAINT "users_email_key" UNIQUE ("email"), CONSTRAINT "users_username_key" UNIQUE ("username"), CONSTRAINT "users_location_id_fkey" FOREIGN KEY ("location_id") REFERENCES "public"."locations" ("id") ON UPDATE NO ACTION ON DELETE SET NULL);
-- Create index "idx_users_email" to table: "users"
CREATE INDEX "idx_users_email" ON "public"."users" ("email");
-- Create "orders" table
CREATE TABLE "public"."orders" ("id" serial NOT NULL, "user_id" integer NOT NULL, "seller_id" integer NOT NULL, "status" character varying(12) NOT NULL DEFAULT 'PENDING', "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamp NULL DEFAULT CURRENT_TIMESTAMP, "total_amount" numeric(10,4) NOT NULL, "currency" character varying(12) NOT NULL DEFAULT 'USD', PRIMARY KEY ("id"), CONSTRAINT "orders_seller_id_fkey" FOREIGN KEY ("seller_id") REFERENCES "public"."seller" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "orders_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "orders_currency_check" CHECK ((currency)::text = ANY ((ARRAY['USD'::character varying, 'INR'::character varying, 'EUR'::character varying])::text[])), CONSTRAINT "orders_status_check" CHECK ((status)::text = ANY ((ARRAY['PENDING'::character varying, 'PROCESSING'::character varying, 'SHIPPED'::character varying, 'CANCELED'::character varying, 'REFUNDED'::character varying])::text[])));
-- Create index "idx_orders_fulfillment_center_id" to table: "orders"
CREATE INDEX "idx_orders_fulfillment_center_id" ON "public"."orders" ("seller_id");
-- Create index "idx_orders_user_id" to table: "orders"
CREATE INDEX "idx_orders_user_id" ON "public"."orders" ("user_id");
-- Create "order_items" table
CREATE TABLE "public"."order_items" ("id" serial NOT NULL, "order_id" integer NOT NULL, "product_id" integer NOT NULL, "quantity" integer NOT NULL, "price" numeric(10,4) NOT NULL, "currency" character varying(12) NOT NULL DEFAULT 'USD', PRIMARY KEY ("id"), CONSTRAINT "order_items_order_id_fkey" FOREIGN KEY ("order_id") REFERENCES "public"."orders" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "order_items_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "order_items_currency_check" CHECK ((currency)::text = ANY ((ARRAY['USD'::character varying, 'INR'::character varying, 'EUR'::character varying])::text[])));
-- Create index "idx_order_items_order_id" to table: "order_items"
CREATE INDEX "idx_order_items_order_id" ON "public"."order_items" ("order_id");
-- Create index "idx_order_items_product_id" to table: "order_items"
CREATE INDEX "idx_order_items_product_id" ON "public"."order_items" ("product_id");
-- Create "product_reviews" table
CREATE TABLE "public"."product_reviews" ("id" serial NOT NULL, "user_id" integer NOT NULL, "product_id" integer NOT NULL, "rating" integer NOT NULL, PRIMARY KEY ("id"), CONSTRAINT "product_reviews_product_id_fkey" FOREIGN KEY ("product_id") REFERENCES "public"."products" ("id") ON UPDATE NO ACTION ON DELETE CASCADE, CONSTRAINT "product_reviews_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "public"."users" ("id") ON UPDATE NO ACTION ON DELETE CASCADE);
-- Create index "idx_product_reviews_product_id" to table: "product_reviews"
CREATE INDEX "idx_product_reviews_product_id" ON "public"."product_reviews" ("product_id");
-- Create index "idx_product_reviews_user_id" to table: "product_reviews"
CREATE INDEX "idx_product_reviews_user_id" ON "public"."product_reviews" ("user_id");
