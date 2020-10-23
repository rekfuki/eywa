.PHONY: setup
setup:
	sudo apt install ansible
	ansible-galaxy install gantsign.golang