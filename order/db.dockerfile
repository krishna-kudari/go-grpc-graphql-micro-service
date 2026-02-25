FROM postgres:17

COPY order/up.sql /docker-entrypoint-initdb.d/1.sql

CMD ["postgres"]
