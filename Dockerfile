FROM ubuntu:17.04

LABEL author="Andrey Kuchin"

# Обновление списка пакетов
RUN apt-get -y update

#
# Установка postgresql, wget, git, golang
#
ENV PGVER 9.6

RUN apt-get install -y postgresql-$PGVER wget git && \
    wget https://storage.googleapis.com/golang/go1.9.1.linux-amd64.tar.gz

RUN tar -C /usr/local -xzf go1.9.1.linux-amd64.tar.gz && \
    mkdir go && mkdir go/src && mkdir go/bin && mkdir go/pkg

# Выставляем переменную окружения для сборки проекта

ENV GOPATH $HOME/go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

ADD ./ $GOPATH/src/github.com/nd-r/tech-db-forum/
# RUN go get github.com/nd-r/tech-db-forum/
RUN go install github.com/nd-r/tech-db-forum/

# Run the rest of the commands as the ``postgres``
# user created by the ``postgres-$PGVER`` package 
# when it was ``apt-get installed``
USER postgres

# Create a PostgreSQL role named ``docker`` with ``docker`` as the password and
# then create a database `docker` owned by the ``docker`` role.
RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
    createdb -O docker docker &&\
     &&\
    /etc/init.d/postgresql stop

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >>\
 /etc/postgresql/$PGVER/main/pg_hba.conf

# And add ``listen_addresses`` to ``/etc/postgresql/$PGVER/main/postgresql.conf``
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

# Expose the PostgreSQL port
EXPOSE 5432
EXPOSE 5000

# Back to the root user
USER root

CMD service postgresql start && tech-db-forum
