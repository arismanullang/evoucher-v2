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
CREATE TYPE redemption_method AS ENUM (
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
    name character varying(16) NOT NULL,
    billing character varying(16) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL,
    alias CHARACTER VARYING(8),
    CONSTRAINT accounts_pkey PRIMARY KEY (id)
);

CREATE SEQUENCE broadcast_users_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

CREATE TABLE broadcast_users (
    id serial NOT NULL ,
    state CHARACTER VARYING(8) NOT NULL,
    program_id CHARACTER VARYING(8) NOT NULL,
    target CHARACTER VARYING(256) NOT NULL,
    created_by CHARACTER VARYING(8) DEFAULT 'unknown'::CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_by CHARACTER VARYING(8),
    updated_at TIMESTAMP WITH TIME ZONE,
    status status DEFAULT 'created'::status NOT NULL,
    description CHARACTER VARYING(64),
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
    id INTEGER DEFAULT nextval('features_id_seq'::regclass) NOT NULL,
    detail CHARACTER VARYING(32) NOT NULL,
    created_by CHARACTER VARYING(8) DEFAULT 'unknown'::CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    status status DEFAULT 'created'::status NOT NULL ,
    category CHARACTER VARYING(20) NOT NULL,
    CONSTRAINT features_pkey PRIMARY KEY (id)
);
ALTER SEQUENCE features_id_seq OWNED BY features.id;


CREATE TABLE partners (
    id CHARACTER VARYING(8) DEFAULT new_id() NOT NULL,
    name CHARACTER VARYING(32) NOT NULL,
    serial_number CHARACTER VARYING(32),
    tag TEXT,
    description TEXT,
    account_id CHARACTER VARYING(8) NOT NULL DEFAULT 'unknown'::CHARACTER VARYING,
    created_by CHARACTER VARYING(8) DEFAULT 'unknown'::CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_by CHARACTER VARYING(8),
    updated_at TIMESTAMP WITH TIME ZONE,
    status status DEFAULT 'created'::status NOT NULL,
    CONSTRAINT partners_pkey PRIMARY KEY (id)
);

CREATE SEQUENCE program_partners_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

CREATE TABLE program_partners (
    id INTEGER DEFAULT nextval('valid_partners_id_seq'::regclass) NOT NULL,
    program_id CHARACTER VARYING(8) NOT NULL,
    partner_id CHARACTER VARYING(8) NOT NULL,
    created_by CHARACTER VARYING(8) DEFAULT 'unknown'::CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_by CHARACTER VARYING(8),
    updated_at TIMESTAMP WITH TIME ZONE ,
    status status DEFAULT 'created'::status NOT NULL ,
    CONSTRAINT valid_partners_pkey PRIMARY KEY (id)
);
ALTER SEQUENCE program_partners_id_seq OWNED BY program_partners.id;


CREATE TABLE programs (
    id character varying(8) DEFAULT new_id() NOT NULL,
    account_id character varying(8) NOT NULL,
    name character varying(64) NOT NULL,
    type character varying(64) DEFAULT 'on-demand'::variant_type NOT NULL,
    voucher_format_id integer NOT NULL default 0,
    voucher_type character varying(16) DEFAULT 'cash'::voucher_type,
    voucher_price numeric NOT NULL,
    allow_accumulative character varying(8) NOT NULL,
    start_date timestamp with time zone DEFAULT now() NOT NULL,
    end_date timestamp with time zone DEFAULT now() NOT NULL,
    start_hour character varying(8) NOT NULL,
    end_hour character varying(8) NOT NULL,
    voucher_value numeric(24,2),
    max_generate_voucher numeric(24,2),
    max_quantity_voucher numeric(24,2),
    redemption_method character varying(16) DEFAULT 'qr'::redemption_method,
    img_url character varying(256) NOT NULL,
    tnc text NOT NULL,
    description character varying(256) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    deleted_by character varying(8),
    deleted_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL,
    valid_voucher_start TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
    valid_voucher_end TIMESTAMP WITHOUT TIME ZONE DEFAULT now(),
    voucher_lifetime NUMERIC,
    validity_days TEXT,
    max_redeem_voucher NUMERIC(24,2) NOT NULL DEFAULT 1,
    CONSTRAINT programs_pkey PRIMARY KEY (id)
);


CREATE TABLE public.role_features
(
  id integer NOT NULL DEFAULT nextval('rules_id_seq'::regclass),
  role_id character varying(8) NOT NULL,
  feature_id character varying(8) NOT NULL,
  created_by character varying(8) NOT NULL DEFAULT 'unknown'::character varying,
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  status status NOT NULL DEFAULT 'created'::status,
  CONSTRAINT rules_pkey PRIMARY KEY (id)
)


CREATE TABLE roles (
    id CHARACTER VARYING(8) DEFAULT new_id() NOT NULL,
    detail character varying(32) NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    status status DEFAULT 'created'::status NOT NULL,
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
    status status DEFAULT 'created'::status NOT NULL,
    CONSTRAINT rules_pkey PRIMARY KEY (id)
);
ALTER SEQUENCE rules_id_seq OWNED BY rules.id;


CREATE TABLE public.tags
(
  id character varying(8) NOT NULL DEFAULT new_id(),
  value character varying(16),
  created_at timestamp with time zone NOT NULL DEFAULT now(),
  status status NOT NULL DEFAULT 'created'::status,
  created_by character varying(8) NOT NULL DEFAULT 'unknown'::character varying,
  updated_at timestamp with time zone,
  updated_by character varying(8) DEFAULT 'unknown'::character varying,
  CONSTRAINT tags_pkey PRIMARY KEY (id)
)


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
    discount_value numeric(24,2) NOT NULL,
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
    status status DEFAULT 'created'::status NOT NULL,
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
    status status DEFAULT 'created'::status NOT NULL,
    CONSTRAINT users_pkey PRIMARY KEY (id)
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
    holder character varying(64) NOT NULL,
    holder_email character varying(32),
    holder_phone character varying(16),
    holder_description character varying(64),
    program_id character varying(8) NOT NULL,
    valid_at timestamp with time zone DEFAULT now() NOT NULL,
    expired_at timestamp with time zone DEFAULT now() NOT NULL,
    voucher_value numeric(24,2) NOT NULL,
    state voucher_state DEFAULT 'created'::voucher_state NOT NULL,
    created_by character varying(8) DEFAULT 'unknown'::character varying NOT NULL,
    created_at timestamp with time zone DEFAULT now() NOT NULL,
    updated_by character varying(8),
    updated_at timestamp with time zone,
    deleted_by character varying(8),
    deleted_at timestamp with time zone,
    status status DEFAULT 'created'::status NOT NULL,
    CONSTRAINT vouchers_pkey PRIMARY KEY (id)
);
