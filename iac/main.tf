resource "google_compute_instance" "go_app_instance" {
  name         = var.instance_name
  machine_type = "e2-medium"
  zone         = var.zone

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

    sudo apt-get update
    sudo apt-get install -y curl git postgresql postgresql-contrib

    GO_VERSION="1.22.3"
    GO_TAR="go$${GO_VERSION}.linux-amd64.tar.gz"
    wget https://go.dev/dl/$${GO_TAR}
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf $${GO_TAR}

    echo 'export PATH=$PATH:/usr/local/go/bin' | sudo tee -a /etc/profile
    export PATH=$PATH:/usr/local/go/bin

    sudo -u postgres psql -c "ALTER USER postgres PASSWORD '${var.db_password}';"
    sudo -u postgres createdb appdb

    git clone https://github.com/tamaqazaq/SREFinal.git /opt/app
    cd /opt/app
    /usr/local/go/bin/go mod tidy
    /usr/local/go/bin/go build -o app main.go

    nohup ./app > app.log 2>&1 &
  EOT



  tags = ["http-server"]
}
