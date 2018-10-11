for d in `seq 1986 2018`
do
     echo "Re-Creating Table forest_practice_applications_wa_expired_$d"
     psql "host=127.0.0.1 sslmode=disable dbname=postgres user=postgres password=postgres" -c "DROP TABLE forest_practice_applications_wa_expired_$d"
     psql "host=127.0.0.1 sslmode=disable dbname=postgres user=postgres password=postgres" -c "CREATE TABLE forest_practice_applications_wa_expired_$d AS SELECT * FROM forest_practice_applications_wa WHERE EXTRACT(year from expiration) = $d AND decision != 'DISAPPROVED';"
done

for d in `seq 1986 2018`
do
     echo "Re-Creating Table forest_practice_applications_wa_evenage_expired_$d"
     psql "host=127.0.0.1 sslmode=disable dbname=postgres user=postgres password=postgres" -c "DROP TABLE forest_practice_applications_wa_evenage_expired_$d"
     psql "host=127.0.0.1 sslmode=disable dbname=postgres user=postgres password=postgres" -c "CREATE TABLE forest_practice_applications_wa_evenage_expired_$d AS SELECT * FROM forest_practice_applications_wa WHERE EXTRACT(year from expiration) = $d AND decision != 'DISAPPROVED' AND timharv_fp = 'EVEN-AGE';"
done