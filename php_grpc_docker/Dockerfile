FROM php:8.1-fpm-alpine

RUN apk update && apk add  --no-cache \
    grpc \
    protobuf \
    grpc-plugins
    #grpc-plugins \
    #protoc

COPY docker/php/php.ini-development "$PHP_INI_DIR/php.ini"

RUN pecl config-set php_ini "$PHP_INI_DIR/php.ini"

# Set working directory
WORKDIR /var/www

# Copy existing application directory permissions
RUN chown -R www-data:www-data /var/www/

# Setup supervisor as root user
USER root
RUN apk add --no-progress --quiet --no-cache supervisor
COPY docker/supervisord.conf /etc/supervisord.conf

# Switch to www-user
USER www-data
EXPOSE 8080

CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]
