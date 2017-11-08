FROM ubuntu:17.04

LABEL author="Andrey Kuchin"

# Обновление списка пакетов
RUN apt-get -y update

#
# Установка postgresql
#
ENV PGVER 9.6
RUN apt-get install -y postgresql-$PGVER

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
RUN echo "host all  all    0.0.0.0/0  md5" >>\
    /etc/postgresql/$PGVER/main/pg_hba.conf

# And add ``listen_addresses`` to ``/etc/postgresql/$PGVER/main/postgresql.conf``
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

# Expose the PostgreSQL port
EXPOSE 5432

# Add VOLUMEs to allow backup of config, logs and databases
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Back to the root user
USER root

#
# Сборка проекта
#

# Установка golang
RUN apt-get install -y wget git && \
    wget https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf go1.9.1.linux-amd64.tar.gz && \
    mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg

# Выставляем переменную окружения для сборки проекта
ENV GOPATH $HOME/go

ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ADD ./ $GOPATH/src/github.com/nd-r/tech-db-forum/
RUN go install github.com/nd-r/tech-db-forum/

RUN echo "synchronous_commit='off'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "shared_buffers = 512MB" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "effective_cache_size = 1024MB" >> /etc/postgresql/$PGVER/main/postgresql.conf



WORKDIR ${GOPATH}/src/github.com/nd-r/tech-db-forum/

EXPOSE 5000

USER postgres
CMD  service postgresql start && tech-db-forum
# ["/usr/lib/postgresql/9.6/bin/postgres", "-D", "/var/lib/postgresql/9.6/main", "-c", "config_file=/etc/postgresql/9.6/main/postgresql.conf"] &&
