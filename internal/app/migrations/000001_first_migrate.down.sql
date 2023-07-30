SELECT
    id,
    name,
    html_description,
    address,
    location,
    owner_id,
    6371 * ACOS(
        SIN(RADIANS($ 1)) * SIN(RADIANS(location [0])) + COS(RADIANS($ 1)) * COS(RADIANS(location [0])) * COS(RADIANS($ 2 - location [1]))
    ) AS distance,
    created_at,
    updated_at
FROM
    edu_centers
WHERE
    CASE
        WHEN $ 5 <> 0 THEN 
        6371 * ACOS(
            SIN(RADIANS($ 1)) * SIN(RADIANS(location [0])) + COS(RADIANS($ 1)) * COS(RADIANS(location [0])) * COS(RADIANS($ 2 - location [1]))
        ) <= $ 5 
        ELSE TRUE 
    END
ORDER BY
    distance
LIMIT
    $ 3 OFFSET $ 4;