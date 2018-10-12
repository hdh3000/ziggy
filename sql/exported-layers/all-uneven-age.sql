SELECT gid, expiration, extract(year from expiration) as year, geom
FROM forest_practice_applications_wa a
WHERE a.timharv_fp = 'UNEVEN-AGE' AND decision != 'DISAPROVED'