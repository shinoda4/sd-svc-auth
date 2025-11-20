
#!/bin/sh

if command -v migrate >/dev/null 2>&1; then
    echo "migrate is installed"
else
    echo "migrate is NOT installed"

    echo "Installing golang-migrate..."
    
    curl -L https://github.com/golang-migrate/migrate/releases/download/v4.19.0/migrate.linux-amd64.tar.gz | tar xvz -C /usr/local/bin
fi

echo "Database ready, prepare migrate..."

migrate -source file://db/migrations -database "$DATABASE_DSN" up





