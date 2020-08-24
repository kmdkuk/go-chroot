# -*- mode: ruby -*-
# vi: set ft=ruby :

# All Vagrant configuration is done below. The "2" in Vagrant.configure
# configures the configuration version (we support older styles for
# backwards compatibility). Please don't change it unless you know what
# you're doing.
Vagrant.configure("2") do |config|
  # The most common configuration options are documented and commented below.
  # For a complete reference, please see the online documentation at
  # https://docs.vagrantup.com.

  # Every Vagrant development environment requires a box. You can search for
  # boxes at https://vagrantcloud.com/search.
  config.vm.box = "bento/ubuntu-20.04"

  # Share an additional folder to the guest VM. The first argument is
  # the path on the host to the actual folder. The second argument is
  # the path on the guest to mount the folder. And the optional third
  # argument is a set of non-required options.
  config.vm.synced_folder "./", "/home/vagrant/go-chroot", create: "true"

  # Provider-specific configuration so you can fine-tune various
  # backing providers for Vagrant. These expose provider-specific options.
  # Example for VirtualBox:
  #
  config.vm.provider "virtualbox" do |vb|
    # Display the VirtualBox GUI when booting the machine
    # vb.gui = true
  
    # Customize the amount of memory on the VM:
    vb.memory = "4096"
    vb.cpus = 2
  end
  #
  # View the documentation for the provider you are using for more
  # information on available options.

  # Enable provisioning with a shell script. Additional provisioners such as
  # Ansible, Chef, Docker, Puppet and Salt are also available. Please see the
  # documentation for more information about their specific syntax and use.
  config.vm.provision "shell", inline: <<-SHELL
    apt-get update
    apt-get -y upgrade
    apt-get install -y git wget vim build-essential
  SHELL

  config.vm.provision :vagrant_user, type: "shell", privileged: false, inline: <<-SHELL
    echo $HOME

    if [ ! -d ~/.goenv ]; then
      git clone https://github.com/syndbg/goenv.git ~/.goenv
      echo 'export GOENV_ROOT="$HOME/.goenv"' >> ~/.bash_profile
      echo 'export PATH="$GOENV_ROOT/bin:$PATH"' >> ~/.bash_profile
      echo 'eval "$(goenv init -)"' >> ~/.bash_profile
      echo 'export PATH="$GOROOT/bin:$PATH"' >> ~/.bash_profile
      echo 'export PATH="$PATH:$GOPATH/bin"' >> ~/.bash_profile
    fi

    if [ ! -e .vimrc ]; then
      wget https://raw.githubusercontent.com/kmdkuk/MyDotFiles/master/.vimrc
    fi
  SHELL
end
