data "abrha_ssh_key" "ssh_key" {
  name = "terraform-vm"
}

resource "abrha_vm" "web" {
  image  = "ubuntu24-cloudinit-qcow2"
  name   = data.external.vm_name.result.name
  region = "frankfurt"
  size   = "deLinuxVPS4"
  ssh_keys = [
    data.abrha_ssh_key.ssh_key.id
  ]

  connection {
    host        = self.ipv4_address
    user        = "root"
    type        = "ssh"
    private_key = file(var.ssh_private_key_path)
    timeout     = "2m"
  }

  provisioner "remote-exec" {
    inline = [
      "export PATH=$PATH:/usr/bin",
      # Install Apache
      "apt update",
      "apt -y install apache2"
    ]
  }
}