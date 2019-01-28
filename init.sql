CREATE TYPE payment_type AS ENUM (
    'cash',
    'credit',
    'debit'
);
CREATE TYPE status AS ENUM (
    'created',
    'deleted'
);
CREATE TYPE voucher_state AS ENUM (
    'created',
    'used',
    'paid',
    'deleted'
);
CREATE TYPE voucher_type AS ENUM (
    'amount',
    'discount',
    'item'
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

--Companies : bisa di ganti client juno
CREATE TABLE companies (
    id CHARACTER VARYING(8) DEFAULT new_id() NOT NULL,
    client_key CHARACTER VARYING(64) NOT NULL,
    client_secret CHARACTER VARYING(64) NOT NULL,
    created_by CHARACTER VARYING(8) DEFAULT 'system'::CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_by CHARACTER VARYING(8),
    updated_at TIMESTAMP WITH TIME ZONE,
    status status DEFAULT 'created'::status NOT NULL,

    CONSTRAINT accounts_pkey PRIMARY KEY (id)
);

--customers (additional obj)
CREATE TABLE customers (
    id CHARACTER VARYING(8) DEFAULT new_id() NOT NULL ,
    name CHARACTER VARYING(128) NOT NULL,
    mobile_phone CHARACTER VARYING(16) NOT NULL,
    email CHARACTER VARYING(256) NOT NULL,
    ref_id CHARACTER VARYING(64),
    company_id CHARACTER VARYING(8),
    created_by CHARACTER VARYING(8) DEFAULT 'system'::CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_by CHARACTER VARYING(8),
    updated_at TIMESTAMP WITH TIME ZONE ,
    deleted_by CHARACTER VARYING(8),
    deleted_at TIMESTAMP WITH TIME ZONE,
    status status DEFAULT 'created'::status NOT NULL,    
    
    CONSTRAINT customers_pkey PRIMARY KEY (id)
);

--Partner Tags
CREATE SEQUENCE customer_tag_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;
CREATE TABLE customer_tags
(
    id serial NOT NULL,
    customer_id CHARACTER VARYING(8) NOT NULL,
    tag_id CHARACTER VARYING(8) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    created_by CHARACTER VARYING(8) NOT NULL DEFAULT 'unknown'::CHARACTER VARYING,
    status status NOT NULL DEFAULT 'created'::status,
)
ALTER SEQUENCE customer_tag_id_seq OWNED BY customer_tags.id;


--partners
CREATE TABLE partners (
    id CHARACTER VARYING(8) DEFAULT new_id() NOT NULL,
    name CHARACTER VARYING(32) NOT NULL,
    description TEXT,
    company_id CHARACTER VARYING(8) NOT NULL DEFAULT 'system'::CHARACTER VARYING,
    created_by CHARACTER VARYING(8) DEFAULT 'system'::CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_by CHARACTER VARYING(8),
    updated_at TIMESTAMP WITH TIME ZONE,
    status status DEFAULT 'created'::status NOT NULL,
    CONSTRAINT partners_pkey PRIMARY KEY (id)
);

--tags (additional obj)
CREATE TABLE tags
(
    id CHARACTER VARYING(8) NOT NULL DEFAULT new_id(),
    name CHARACTER VARYING(32),
    company_id CHARACTER VARYING (16)
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    created_by CHARACTER VARYING(8) NOT NULL DEFAULT 'unknown'::CHARACTER VARYING,
    updated_at TIMESTAMP WITH TIME ZONE,
    updated_by CHARACTER VARYING(8) DEFAULT 'unknown'::CHARACTER VARYING,
    status status NOT NULL DEFAULT 'created'::status,

    CONSTRAINT tags_pkey PRIMARY KEY (id)
)

--Partner Tags
CREATE SEQUENCE partner_tag_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;
CREATE TABLE partner_tags
(
    id serial NOT NULL,
    partner_id CHARACTER VARYING(8) NOT NULL,
    tag_id CHARACTER VARYING(8) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT now(),
    created_by CHARACTER VARYING(8) NOT NULL DEFAULT 'unknown'::CHARACTER VARYING,
    status status NOT NULL DEFAULT 'created'::status,
)
ALTER SEQUENCE partner_tag_id_seq OWNED BY partner_tags.id;


--Programs
CREATE TABLE programs (
    id CHARACTER VARYING(8) DEFAULT new_id() NOT NULL,
    company_id CHARACTER VARYING(8) NOT NULL,
    name CHARACTER VARYING(64) NOT NULL,
    type CHARACTER VARYING(64) DEFAULT NOT NULL, --bulk / fix 
    voucher_type CHARACTER VARYING(16) DEFAULT 'amount'::voucher_type,
    value numeric(24.4) NOT NULL,
    price numeric(24.4) NOT NULL,
    price_type CHARACTER VARYING(32),
    stock numeric,
    img_url CHARACTER VARYING(8) NOT NULL,
    description text,        
    template text,
    voucher_format text,
    rule_id CHARACTER VARYING (8)
    created_by CHARACTER VARYING(8) DEFAULT 'unknown'::CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_by CHARACTER VARYING(8),
    updated_at TIMESTAMP WITH TIME ZONE,
    deleted_by CHARACTER VARYING(8),
    deleted_at TIMESTAMP WITH TIME ZONE,
    status status DEFAULT 'created'::status NOT NULL,
    CONSTRAINT programs_pkey PRIMARY KEY (id)
);

--Program Partners
CREATE SEQUENCE program_partner_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

CREATE TABLE program_partners (
    id SERIAL NOT NULL,
    program_id CHARACTER VARYING(8) NOT NULL,
    partner_id CHARACTER VARYING(8) NOT NULL,
    created_by CHARACTER VARYING(8) DEFAULT 'system'::CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_by CHARACTER VARYING(8),
    updated_at TIMESTAMP WITH TIME ZONE ,
    status status DEFAULT 'created'::status NOT NULL ,

    CONSTRAINT partners_pkey PRIMARY KEY (id)
);
ALTER SEQUENCE program_partner_id_seq OWNED BY program_partners.id;



--Voucher
CREATE TABLE vouchers (
    id CHARACTER VARYING(8) DEFAULT new_id() NOT NULL,
    voucher_code CHARACTER VARYING(16) NOT NULL,
    reference_no CHARACTER VARYING(64) NOT NULL,
    holder text, 
    program_id CHARACTER VARYING(8) NOT NULL,
    valid_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    expired_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    state voucher_state DEFAULT 'created'::voucher_state NOT NULL,
    created_by CHARACTER VARYING(8) DEFAULT 'unknown'::CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_by CHARACTER VARYING(8),
    updated_at TIMESTAMP WITH TIME ZONE,
    status status DEFAULT 'created'::status NOT NULL,
    CONSTRAINT vouchers_pkey PRIMARY KEY (id)
);

--Campaign
CREATE TABLE campaigns (
    id CHARACTER VARYING(8) DEFAULT new_id() NOT NULL,
    program_id CHARACTER VARYING(8),
    voucher_id CHARACTER VARYING(8),
    customer_id CHARACTER VARYING(8),
    created_by CHARACTER VARYING(8) DEFAULT 'unknown'::CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_by CHARACTER VARYING(8),
    updated_at TIMESTAMP WITH TIME ZONE,
    status status DEFAULT 'created'::status NOT NULL,
)

CREATE TABLE transactions (
    id CHARACTER VARYING(8) DEFAULT new_id() NOT NULL,
    company_id CHARACTER VARYING(8) NOT NULL,
    transaction_code CHARACTER VARYING(16) NOT NULL,
    total_amount numeric(24,2) NOT NULL,
    holder text ,
    partner_id CHARACTER VARYING(8),
    created_by CHARACTER VARYING(8) DEFAULT 'unknown'::CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_by CHARACTER VARYING(8),
    updated_at TIMESTAMP WITH TIME ZONE,
    status status DEFAULT 'created'::status NOT NULL,
    CONSTRAINT transactions_pkey PRIMARY KEY (id)
);


CREATE SEQUENCE transaction_details_id_seq
START WITH 1
INCREMENT BY 1
NO MINVALUE
NO MAXVALUE
CACHE 1;

-- ux mobile scan dl bru pilih voucher
CREATE TABLE transaction_details (
    id integer DEFAULT nextval('transaction_details_id_seq'::regclass) NOT NULL,
    transaction_id CHARACTER VARYING(8) NOT NULL,
    program_id CHARACTER VARYING(8),
    voucher_id CHARACTER VARYING(8) NOT NULL,
    created_by CHARACTER VARYING(8) DEFAULT 'unknown'::CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_by CHARACTER VARYING(8),
    updated_at TIMESTAMP WITH TIME ZONE,
    status status DEFAULT 'created'::status NOT NULL ,
    CONSTRAINT transaction_details_pkey PRIMARY KEY (id)
    -- totamount
    -- sku
    -- user
);
ALTER SEQUENCE transaction_details_id_seq OWNED BY transaction_details.id;

-- ###### RULE #####
    allow_cross_program bool -- 1 x 100k + 2 x 50k voucher , di perlukan taging cross program , berkorelasi dengan "allow_accumulative"    
    allow_accumulative bool, -- akumulasi penggunaan multyple voucher per transaksi
        -> max_redeem_voucher -- max usange "accumulative" voucher per transaksi
    max_generate_by_program -- maksimum redeem vocher dalam 1 program 
        -> max_generate_by_day -- redeem perhari nya

    --DISTRIBUTION
    start_date -- start program (01 Nov) batas awal untuk distribusi / pengambilan voucher
    end_date -- end program (31 - Dec)  batas awal ahir distribusi / pengambilan voucher
    --USAGE
        start_hour -- spesifik jam penggunaan vocher ( 10 AM)
        end_hour  -- spesifik jam penggunaan vocher (10 PM)
        validity_days -- spesifik hari yg di perbolehkan dalam penggunaan  ( sen , kamis , sabtu)

        --Voucher Lifetime
        voucher_lifetime -- umur voucher (3 bln/ 90 hari)  , > start_date 
        --## OR ## --
        valid_voucher_start -- tanggal berlakunya voucher (15 )
        valid_voucher_end -- tanggal berlakunya vocher 

        max_usage_by_program -- usage voucher dalam 1 program 
        -> max_usage_by_day -- usage perhari nya
    
    


-- #### CASHOUT ####--  
create table cashouts as (
    id  CHARACTER VARYING(8) ,
    code CHARACTER VARYING (32) ,
    partner_id CHARACTER VARYING(8),
    payment_method CHARACTER VARYING(32),
    amount numeric(24,4),
    create_date TIMESTAMPTZ ,
    state state,
    created_by CHARACTER VARYING(8) DEFAULT 'unknown'::CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_by CHARACTER VARYING(8),
    updated_at TIMESTAMP WITH TIME ZONE,
    status status DEFAULT 'created'::status NOT NULL ,
)
create table cashout_details as (
    id  CHARACTER VARYING(8) ,
    cashout_id CHARACTER VARYING (32) ,
    transaction_id CHARACTER VARYING(8)
    created_by CHARACTER VARYING(8) DEFAULT 'unknown'::CHARACTER VARYING NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL,
    updated_by CHARACTER VARYING(8),
    updated_at TIMESTAMP WITH TIME ZONE,
)

create table users as (
    id  CHARACTER VARYING(8) ,
    
)

-- partner  , company(user) 
-- transactions , INVOICE (report) ,

-- dasboard outlet/partner ???
-- upload bukti transfer untuk validasi cashout
