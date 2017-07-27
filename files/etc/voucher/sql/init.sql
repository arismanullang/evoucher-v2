CREATE TYPE voucher_format AS ENUM (
    'Alphabet',
    'Numerals',
    'Alphanumeric'
);

CREATE TYPE broadcast_state AS ENUM (
    'created',
    'broadcast',
    'void'
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


CREATE TABLE accounts (
    id character varying(8) DEFAULT new_id() NOT NULL,
    account_name character varying(16) NOT NULL,
    billing character varying(16) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL ,
    CONSTRAINT accounts_pkey PRIMARY KEY (id)
);

CREATE SEQUENCE broadcast_users_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

CREATE TABLE broadcast_users (
    id integer DEFAULT nextval('broadcast_users_id_seq'::regclass) NOT NULL ,
    state character varying(8) NOT NULL,
    variant_id character varying(8) NOT NULL,
    target character varying(256) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL,

    CONSTRAINT broadcast_users_pkey PRIMARY KEY (id)
);
ALTER SEQUENCE broadcast_users_id_seq OWNED BY broadcast_users.id;


CREATE SEQUENCE features_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;
CREATE TABLE features (
    id integer DEFAULT nextval('features_id_seq'::regclass) NOT NULL,
    feature_detail character varying(32) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    status status DEFAULT 'created'::status NOT NULL ,
    CONSTRAINT features_pkey PRIMARY KEY (id)
);
ALTER SEQUENCE features_id_seq OWNED BY features.id;


CREATE TABLE partners (
    id character varying(8) DEFAULT new_id() NOT NULL,
    partner_name character varying(32) NOT NULL,
    serial_number character varying(32),
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    status status DEFAULT 'created'::status NOT NULL ,
    CONSTRAINT partners_pkey PRIMARY KEY (id)
);

CREATE TABLE roles (
    id character varying(8) DEFAULT new_id() NOT NULL,
    role_detail character varying(32) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    status status DEFAULT 'created'::status NOT NULL ,
    CONSTRAINT roles_pkey PRIMARY KEY (id)
);


CREATE SEQUENCE rules_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

CREATE TABLE rules (
    id integer DEFAULT nextval('rules_id_seq'::regclass) NOT NULL,
    role_id character varying(8) NOT NULL,
    feature_id character varying(8) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    status status DEFAULT 'created'::status NOT NULL ,
    CONSTRAINT rules_pkey PRIMARY KEY (id)
);
ALTER SEQUENCE rules_id_seq OWNED BY rules.id;


CREATE SEQUENCE transaction_details_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

CREATE TABLE transaction_details (
    id integer DEFAULT nextval('transaction_details_id_seq'::regclass) NOT NULL,
    transaction_id character varying(8) NOT NULL,
    voucher_id character varying(8) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL ,
    CONSTRAINT transaction_details_pkey PRIMARY KEY (id)
);
ALTER SEQUENCE transaction_details_id_seq OWNED BY transaction_details.id;


CREATE TABLE transactions (
    id character varying(8) DEFAULT new_id() NOT NULL,
    account_id character varying(8) NOT NULL,
    partner_id character varying(8) NOT NULL,
    transaction_code character varying(16) NOT NULL,
    total_transaction numeric(24,2) NOT NULL,
    discount_value numeric(24,2) NOT NULL,
    payment_type character varying(16) DEFAULT 'cash'::payment_type NOT NULL,
    token character varying(32),
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL,
    CONSTRAINT transactions_pkey PRIMARY KEY (id)
);


CREATE SEQUENCE user_accounts_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

CREATE TABLE user_accounts (
    id integer DEFAULT nextval('user_accounts_id_seq'::regclass) NOT NULL,
    user_id character varying(8) NOT NULL,
    account_id character varying(8) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL ,
    CONSTRAINT user_accounts_pkey PRIMARY KEY (id)
);
ALTER SEQUENCE user_accounts_id_seq OWNED BY user_accounts.id;


CREATE SEQUENCE user_roles_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

CREATE TABLE user_roles (
    id integer DEFAULT nextval('user_roles_id_seq'::regclass) NOT NULL,
    user_id character varying(8) NOT NULL,
    role_id character varying(8) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL ,
    CONSTRAINT user_roles_pkey PRIMARY KEY (id)
);
ALTER SEQUENCE user_roles_id_seq OWNED BY user_roles.id;


CREATE TABLE users (
    id character varying(8) DEFAULT new_id() NOT NULL,
    username character varying(32) NOT NULL,
    password character varying(64) NOT NULL,
    email character varying(64) NOT NULL,
    phone character varying(32) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL ,
    CONSTRAINT users_pkey PRIMARY KEY (id)
);


CREATE SEQUENCE variant_partners_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

CREATE TABLE variant_partners (
    id integer DEFAULT nextval('valid_partners_id_seq'::regclass) NOT NULL,
    variant_id character varying(8) NOT NULL,
    partner_id character varying(8) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL ,
    CONSTRAINT valid_partners_pkey PRIMARY KEY (id)
);
ALTER SEQUENCE variant_partners_id_seq OWNED BY variant_partners.id;


CREATE TABLE variants (
    id character varying(8) DEFAULT new_id() NOT NULL,
    account_id character varying(8) NOT NULL,
    variant_name character varying(64) NOT NULL,
    variant_type character varying(64) DEFAULT 'on-demand'::variant_type NOT NULL,
    voucher_format_id integer NOT NULL default 0,
    voucher_type character varying(16) DEFAULT 'cash'::voucher_type,
    voucher_price numeric NOT NULL,
    allow_accumulative character varying(8) NOT NULL,
    start_date timestamp with time zone DEFAULT now() NOT NULL,
    end_date timestamp with time zone DEFAULT now() NOT NULL,
    discount_value numeric(24,2),
    max_generate_voucher numeric(24,2),
    max_quantity_voucher numeric(24,2),
    redeemtion_method character varying(16) DEFAULT 'qr'::redeemtion_method,
    img_url character varying(256) NOT NULL,
    variant_tnc character varying(256) NOT NULL,
    variant_description character varying(256) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    deleted_by character varying(8),
    deleted_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL,
    CONSTRAINT variants_pkey PRIMARY KEY (id)
);


CREATE SEQUENCE voucher_formats_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

CREATE TABLE voucher_formats (
    id integer DEFAULT nextval('voucher_formats_id_seq'::regclass) NOT NULL,
    prefix character varying(8),
    postfix character varying(8),
    body character varying(8),
    format_type voucher_format default 'Numerals' NOT NULL,
    length numeric NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    status status DEFAULT 'created'::status NOT NULL ,
    CONSTRAINT voucher_formats_pkey PRIMARY KEY (id)
);
ALTER SEQUENCE voucher_formats_id_seq OWNED BY voucher_formats.id;


CREATE TABLE vouchers (
    id character varying(8) DEFAULT new_id() NOT NULL,
    voucher_code character varying(16) NOT NULL,
    reference_no character varying(64) NOT NULL,
    holder character varying(8) NOT NULL,
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
    status status DEFAULT 'created'::status NOT NULL ,
    CONSTRAINT vouchers_pkey PRIMARY KEY (id)
);
