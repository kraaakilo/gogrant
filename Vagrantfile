Vagrant.configure("2") do |config|
  config.vm.box = "dazdazdaz"
  config.vm.hostname = "dazdazdaz"
  config.vm.network "private_network", ip: "192.168.56.200"

  # config.vm.synced_folder "gogrant", "/home/vagrant/gogrant"

  config.vm.provider "virtualbox" do |vb|
    vb.memory = "2048"
    vb.cpus = 4
    vb.name = "dazdazdaz"
  end

  config.vm.provision "shell" do |shell|
    shell.inline = <<-SHELL
      echo "Hello from gogrant!"
    SHELL
  end

end

