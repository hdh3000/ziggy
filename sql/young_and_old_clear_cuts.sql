-- 8-12 yo clear cuts bordering or inside of 18-22 yo clearcuts


DROP TABLE	young_and_old_clear_cuts_wa;

CREATE TABLE 
	young_and_old_clear_cuts_wa
AS 
WITH cuts as (
	SELECT 
	a.gid as youngCut,
	b.gid as oldCut
	FROM
		forest_practice_applications_wa a
	CROSS JOIN (SELECT 
	*
	FROM
		forest_practice_applications_wa
	WHERE
		decision != 'DISAPPROVED'
	AND
		timharv_fp = 'EVEN-AGE'
	AND 
		EXTRACT(year from expiration) > 1995 
		AND 
		EXTRACT(year from expiration) < 2001
	) b 
	WHERE
		a.decision != 'DISAPPROVED'
	AND
		a.timharv_fp = 'EVEN-AGE'
	AND 
		EXTRACT(year from a.expiration) > 2005 
	AND 
		EXTRACT(year from a.expiration) < 2011
	AND	
		( 
			st_touches(b.geom, a.geom) 
		OR 
			st_intersects(b.geom, a.geom)
		)
)

SELECT 
	l.* 
FROM 
	forest_practice_applications_wa l
INNER JOIN 
	(select distinct(youngCut) as cut from cuts UNION select distinct(oldCut) as cut from cuts) y on y.cut = l.gid order by l.expiration desc;
	