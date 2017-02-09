SET statement_timeout = 0;
SET lock_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;
COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';

SET search_path = public, pg_catalog;
CREATE TYPE account_role AS ENUM (
    'suadmin',
    'admin',
    'operator',
    'user'
);
CREATE TYPE payment_type AS ENUM (
    'cash',
    'credit',
    'debit'
);
CREATE TYPE redeemtion_method AS ENUM (
    'qr',
    'token'
);
CREATE TYPE status AS ENUM (
    'created',
    'deleted'
);
CREATE TYPE variant_type AS ENUM (
    'bulk',
    'on-demand'
);
CREATE TYPE voucher_state AS ENUM (
    'created',
    'active',
    'used',
    'paid',
    'deleted'
);
CREATE TYPE voucher_type AS ENUM (
    'cash',
    'discount',
    'promo'
);
CREATE FUNCTION new_id() RETURNS text
    LANGUAGE plpgsql
    AS $$
DECLARE
  v_chars TEXT[] := '{0,1,2,3,4,5,6,7,8,9,A,B,C,D,E,F,G,H,I,J,K,L,M,N,O,P,Q,R,S,T,U,V,W,X,Y,Z,a,b,c,d,e,f,g,h,i,j,k,l,m,n,o,p,q,r,s,t,u,v,w,x,y,z,_,-}';
  out_result TEXT := '';
  i INTEGER := 0;
BEGIN
  FOR i IN 1..8 LOOP
    out_result := out_result || v_chars[1+RANDOM()*(ARRAY_LENGTH(v_chars, 1)-1)];
  END LOOP;
  RETURN out_result;
END;
$$;
SET default_tablespace = '';
SET default_with_oids = false;
CREATE TABLE account_roles (
    id character varying(8) DEFAULT new_id() NOT NULL,
    role_detail character varying(32) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    status status DEFAULT 'created'::status NOT NULL
);
CREATE TABLE account_rules (
    id character varying(8) DEFAULT new_id() NOT NULL,
    role_id character varying(8) NOT NULL,
    feature_id character varying(8) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    status status DEFAULT 'created'::status NOT NULL
);
CREATE TABLE accounts (
    id character varying(8) DEFAULT new_id() NOT NULL,
    company_id character varying(8) NOT NULL,
    user_id character varying(8) NOT NULL,
    account_role account_role DEFAULT 'user'::account_role NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL,
    assign_by character varying(8),
    serial_number character varying(32)
);
CREATE TABLE broadcast_users (
    id character varying(8) DEFAULT new_id() NOT NULL,
    company_id character varying(8) NOT NULL,
    variant_id character varying(8) NOT NULL,
    account_id character varying(8) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL
);
CREATE TABLE features (
    id character varying(8) DEFAULT new_id() NOT NULL,
    feature_detail character varying(32) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    status status DEFAULT 'created'::status NOT NULL
);
CREATE TABLE transaction_details (
    id character varying(8) DEFAULT new_id() NOT NULL,
    transaction_id character varying(8) NOT NULL,
    voucher_id character varying(8) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL,
    company_id character varying(8) NOT NULL
);
CREATE TABLE transactions (
    id character varying(8) DEFAULT new_id() NOT NULL,
    pic_merchant character varying(8) NOT NULL,
    company_id character varying(8) NOT NULL,
    transaction_code character varying(16) NOT NULL,
    total_transaction numeric(24,2) NOT NULL,
    discount_value numeric(24,2) NOT NULL,
    payment_type character varying(16) DEFAULT 'cash'::payment_type NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL
);
CREATE TABLE variant_users (
    id character varying(8) DEFAULT new_id() NOT NULL,
    company_id character varying(8) NOT NULL,
    variant_id character varying(8) NOT NULL,
    account_id character varying(8) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL
);
CREATE TABLE variants (
    id character varying(8) DEFAULT new_id() NOT NULL,
    company_id character varying(8) NOT NULL,
    variant_name character varying(64) NOT NULL,
    variant_type character varying(64) DEFAULT 'on-demand'::variant_type NOT NULL,
    point_needed numeric NOT NULL,
    allow_accumulative character varying(8) NOT NULL,
    start_date timestamp with time zone DEFAULT now() NOT NULL,
    end_date timestamp with time zone DEFAULT now() NOT NULL,
    img_url character varying(256) NOT NULL,
    variant_tnc character varying(256) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    deleted_by character varying(8),
    deleted_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL,
    discount_value numeric(24,2),
    max_usage_voucher numeric(24,2),
    max_quantity_voucher numeric(24,2),
    redeemtion_method character varying(16) DEFAULT 'qr'::redeemtion_method,
    voucher_type character varying(16) DEFAULT 'cash'::voucher_type
);
CREATE TABLE vouchers (
    id character varying(8) DEFAULT new_id() NOT NULL,
    voucher_code character varying(16) NOT NULL,
    reference_no character varying(64) NOT NULL,
    account_id character varying(8) NOT NULL,
    variant_id character varying(8) NOT NULL,
    valid_at timestamp with time zone DEFAULT now() NOT NULL,
    expired_at timestamp with time zone DEFAULT now() NOT NULL,
    discount_value numeric(24,2) NOT NULL,
    state voucher_state DEFAULT 'created'::voucher_state NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    deleted_by character varying(8),
    deleted_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL,
    token character varying(32)
);
