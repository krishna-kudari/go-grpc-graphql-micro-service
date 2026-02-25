FROM postgres:17

COPY account/up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]
