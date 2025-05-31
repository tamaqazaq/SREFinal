resource "google_sql_database_instance" "postgres" {
  name             = "sre-postgres"
  database_version = "POSTGRES_13"
  region           = var.region

  settings {
    tier = "db-f1-micro"

    ip_configuration {
      ipv4_enabled = true
      authorized_networks {
        name  = "allow-all"
        value = "0.0.0.0/0"
      }
    }
  }
}

resource "google_sql_user" "postgres_user" {
  name     = "postgres"
  instance = google_sql_database_instance.postgres.name
  password = var.db_password
}

resource "google_sql_database" "appdb" {
  name     = "appdb"
  instance = google_sql_database_instance.postgres.name
}

resource "google_compute_instance" "go_app_instance" {
  name         = var.instance_name
  machine_type = "e2-medium"
  zone         = var.zone

  depends_on = [google_sql_database_instance.postgres]

  boot_disk {
    initialize_params {
      image = "debian-cloud/debian-11"
    }
  }

  network_interface {
    network = "default"
    access_config {}
  }

  metadata_startup_script = <<-EOT
    #!/bin/bash
    set -e

    sudo apt-get update
    sudo apt-get install -y curl git wget

    # Install Go
    GO_VERSION="1.22.3"
    GO_TAR="go$${GO_VERSION}.linux-amd64.tar.gz"
    wget https://go.dev/dl/$${GO_TAR}
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf $${GO_TAR}
    echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
    export PATH=$PATH:/usr/local/go/bin

    # Clone repo
    git clone https://github.com/tamaqazaq/SREFinal.git /opt/app
    cd /opt/app

    # Set up .env file
    cat <<EOF > .env
DB_HOST=34.41.114.8
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=${var.db_password}
DB_NAME=appdb
STRIPE_KEY=sk_test_...
SMTP_HOST=smtp.gmail.com
SMTP_PORT=587
SMTP_USER=yermukhanovdaulet@gmail.com
SMTP_PASSWORD=...
EOF

    # Export env vars explicitly
    export DB_HOST=34.41.114.8
    export DB_PORT=5432
    export DB_USER=postgres
    export DB_PASSWORD=${var.db_password}
    export DB_NAME=appdb
    export STRIPE_KEY=sk_test_...
    export SMTP_HOST=smtp.gmail.com
    export SMTP_PORT=587
    export SMTP_USER=yermukhanovdaulet@gmail.com
    export SMTP_PASSWORD=airasxxjbwzgbmjj

    # Build and run
    go mod tidy
    go build -o app main.go
    nohup ./app > /opt/app/app.log 2>&1 &
  EOT

  tags = ["http-server"]
}

