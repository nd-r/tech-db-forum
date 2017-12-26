FROM ubuntu:17.04

LABEL author="Andrey Kuchin"

# Обновление списка пакетов
RUN apt-get -y update

#
# Установка postgresql
#
ENV PGVER 9.6
RUN apt-get install -y postgresql-$PGVER wget git

RUN wget https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf go1.9.1.linux-amd64.tar.gz && \
    mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg

# Run the rest of the commands as the ``postgres``
# user created by the ``postgres-$PGVER`` package 
# when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN rm -rf /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "local   all             postgres                                peer\n\
local   all             docker                                md5\n\
host    all             all             127.0.0.1/32            md5\n\
host all  all    0.0.0.0/0  md5" >>\
    /etc/postgresql/$PGVER/main/pg_hba.conf

RUN echo "unix_socket_directories = '/var/run/postgresql/'\n\
synchronous_commit='off'\n\
shared_buffers = 512MB\n\
effective_cache_size = 1024MB\n\
wal_writer_delay = 2000ms\n" >> /etc/postgresql/$PGVER/main/postgresql.conf

#RUN echo "log_duration = on" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "log_min_duration_statement = 10" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "log_checkpoints = on" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "log_filename = 'postgresql-%Y-%m-%d_%H%M%S'" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "log_directory = '/var/log/postgresql'" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "log_destination = 'csvlog'" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "log_destination = 'csvlog'" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "logging_collector = on" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "log_disconnections = on" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "log_lock_waits = on" >> /etc/postgresql/$PGVER/main/postgresql.conf
#RUN echo "log_temp_files = 0 " >> /etc/postgresql/$PGVER/main/postgresql.conf

EXPOSE 5432
# EXPOSE 1111

VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
USER root

# Установка golang

# Выставляем переменную окружения для сборки проекта
ENV GOPATH $HOME/go

ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ADD ./ $GOPATH/src/github.com/nd-r/tech-db-forum/
RUN go install github.com/nd-r/tech-db-forum/

WORKDIR ${GOPATH}/src/github.com/nd-r/tech-db-forum/

EXPOSE 5000

USER postgres
CMD  service postgresql start && tech-db-forum
# CMD ["/usr/lib/postgresql/9.6/bin/postgres", "-D", "/var/lib/postgresql/9.6/main", "-c", "config_file=/etc/postgresql/9.6/main/postgresql.conf"] &&
