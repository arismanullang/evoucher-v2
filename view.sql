-- m_blasts

 SELECT b.id,
    b.subject,
    b.sender,
    b.company_id,
    b.program_id,
    b.image_header,
    b.image_footer,
    b.email_content,
    b.template,
    b.created_at,
    b.created_by,
    b.updated_at,
    b.updated_by,
    b.status,
    ( SELECT array_to_json(array_agg(row_to_json(p.*))) AS array_to_json
           FROM programs p
          WHERE p.id::text = b.program_id::text AND p.status = 'created'::status) AS program,
    ( SELECT array_to_json(array_agg(row_to_json(br.*))) AS array_to_json
           FROM blast_recipients br
          WHERE b.id::text = br.blast_id::text AND br.status = 'created'::status) AS recipients
   FROM blasts b;

--    


-- m_cashouts

 SELECT c.id,
    c.company_id,
    c.code,
    p.name AS outlet_name,
    c.outlet_id,
    c.bank_name,
    c.bank_company_name,
    c.bank_account,
    c.reference_no,
    c.amount,
    c.payment_method,
    c.created_at,
    c.created_by,
    c.updated_at,
    c.updated_by,
    c.status
   FROM cashouts c
     JOIN outlets p ON p.id::text = c.outlet_id::text;

--


-- m_channels

 SELECT c.id,
    c.name,
    c.company_id,
    c.description,
    c.is_super,
    c.created_at,
    c.created_by,
    c.updated_at,
    c.updated_by,
    c.status,
    ( SELECT array_to_json(array_agg(row_to_json(t.*))) AS array_to_json
           FROM tags t,
            object_tags pt
          WHERE t.id::text = pt.tag_id::text AND pt.object_id::text = c.id::text AND pt.object_category::text = 'channels'::text AND pt.status = 'created'::status) AS channel_tags
   FROM channels c;

--    


-- m_outlets

 SELECT o.id,
    o.name,
    o.description::json AS description,
    o.emails,
    o.company_id,
    o.created_by,
    o.created_at,
    o.updated_by,
    o.updated_at,
    o.status,
    o.tags,
    o.bank,
    ( SELECT array_to_json(array_agg(row_to_json(pb.*))) AS array_to_json
           FROM outlet_banks pb
          WHERE pb.outlet_id::text = o.id::text) AS outlet_banks,
    ( SELECT array_to_json(array_agg(row_to_json(t.*))) AS array_to_json
           FROM tags t,
            object_tags pt
          WHERE t.id::text = pt.tag_id::text AND pt.object_id::text = o.id::text AND pt.object_category::text = 'outlets'::text AND pt.status = 'created'::status) AS outlet_tags
   FROM outlets o;

--    

-- m_programs

 SELECT p.id,
    p.company_id,
    p.name,
    p.type,
    p.value,
    p.max_value,
    p.stock,
    p.image_url,
    p.description,
    p.template,
    p.channel_id,
    p.is_reimburse,
    p.voucher_format,
    p.price,
    p.rule,
    p.state,
    p.start_date,
    p.end_date,
    p.created_by,
    p.created_at,
    p.updated_by,
    p.updated_at,
    p.status,
    count(v.id) AS claimed,
    COALESCE(sum(
        CASE
            WHEN v.state = 'used'::voucher_state THEN 1
            WHEN v.state = 'paid'::voucher_state THEN 1
            ELSE 0
        END), 0::bigint) AS used,
    COALESCE(sum(
        CASE
            WHEN v.state = 'paid'::voucher_state THEN 1
            ELSE 0
        END), 0::bigint) AS paid,
    ( SELECT array_to_json(array_agg(row_to_json(t.*))) AS array_to_json
           FROM tags t,
            object_tags pt
          WHERE t.id::text = pt.tag_id::text AND pt.object_id::text = p.id::text AND pt.object_category::text = 'programs'::text AND pt.status = 'created'::status) AS program_tags,
    ( SELECT array_to_json(array_agg(row_to_json(c.*))) AS array_to_json
           FROM channels c
          WHERE c.id::text = p.channel_id::text AND c.status = 'created'::status) AS program_channels
   FROM programs p
     LEFT JOIN vouchers v ON p.id::text = v.program_id::text
  GROUP BY p.id;

--   


-- m_transactions

 SELECT t.id,
    t.company_id,
    t.transaction_code,
    p.id AS outlet_id,
    p.name AS outlet_name,
    p.description AS outlet_description,
    t.total_amount,
    t.holder,
    t.created_by,
    t.created_at,
    t.updated_by,
    t.updated_at,
    t.status
   FROM transactions t
     JOIN outlets p ON t.outlet_id::text = p.id::text;

--  

-- m_vouchers

 SELECT v.id,
    v.code,
    v.reference_no,
    v.holder,
    v.holder_detail,
    v.program_id,
    p.image_url AS program_img_url,
    p.name AS program_name,
    p.value AS program_value,
    p.max_value AS program_max_value,
    p.company_id,
    v.valid_at,
    v.expired_at,
    v.state,
    v.created_at,
    v.created_by,
    v.updated_by,
    v.updated_at,
    v.status
   FROM vouchers v
     JOIN programs p ON v.program_id::text = p.id::text;

--  
