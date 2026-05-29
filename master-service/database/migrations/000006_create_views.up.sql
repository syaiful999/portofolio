CREATE OR REPLACE VIEW master.v_master_user AS
SELECT
    u.id,
    u.is_active,
    u.created_by,
    u.created_date,
    u.modified_by,
    u.modified_date,
    u.user_name,
    u.name,
    u.email,
    u.role_id,
    role_enum.enum_value AS role_description,
    role_enum.enum_code AS role_code,
    u.department_id,
    dept_enum.enum_value AS department_name,
    u.outsource_id,
    outs_enum.enum_value AS outsource_name,
    u.location_id,
    loc_enum.enum_value AS location_name,
    u.picture,
    u.is_status
FROM master.master_user u
LEFT JOIN master.master_enum role_enum ON u.role_id = role_enum.id
LEFT JOIN master.master_enum dept_enum ON u.department_id = dept_enum.id
LEFT JOIN master.master_enum outs_enum ON u.outsource_id = outs_enum.id
LEFT JOIN master.master_enum loc_enum ON u.location_id = loc_enum.id
WHERE u.is_active = TRUE;

CREATE OR REPLACE VIEW master.v_master_user_groupby_role AS
SELECT
    u.role_id,
    role_enum.enum_code AS role_code,
    role_enum.enum_value AS role_description,
    COUNT(u.id)::INT AS "count"
FROM master.master_user u
INNER JOIN master.master_enum role_enum ON u.role_id = role_enum.id
WHERE u.is_active = TRUE
GROUP BY u.role_id, role_enum.enum_code, role_enum.enum_value;
