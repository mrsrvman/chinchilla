
[![Circle CI](https://circleci.com/gh/benschw/chinchilla.svg?style=svg)](https://circleci.com/gh/benschw/chinchilla)

[Latest build](https://circle-artifacts.com/gh/benschw/chinchilla/latest/artifacts/0/home/ubuntu/.go_project/src/github.com/benschw/chinchilla/dist/chinchilla_linux_amd64.gz)

[Latest Release](https://github.com/benschw/chinchilla/releases/latest)

# Chinchilla

A service which connects to Rabbitmq queues and delivers messages to REST endpoints.

Chinchilla can be configured either with a yaml config, or with consul. In either case, the
config backend is watched for changes and will live-update your running chinchilla daemon to
reflect the new configuration.

## local demo

	wget https://dl.bintray.com/mitchellh/consul/0.5.2_linux_amd64.zip
	unzip 0.5.2_linux_amd64.zip
	wget https://dl.bintray.com/mitchellh/consul/0.5.2_web_ui.zip
	unzip 0.5.2_web_ui.zip
	mv dist /tmp/web-ui

### terminal 1

	./consul.sh

### terminal 2

	# install some endpoints
	./fixture-data.sh 
	
	# build and start the service
	go build
	SRVLB_HOST=127.0.0.1:8600 ./chinchilla -secret-keyring ./test-keys/.secring.gpg


### terminal 3

	# run a mock rest service
	go run ./example/cmd/server/serve.go


### terminal 4

	# publish some messages to flow through the system
	go run ./example/cmd/publisher/publish.go -queue demo.foo
	go run ./example/cmd/publisher/publish.go -queue demo.bar
	go run ./example/cmd/publisher/publish.go -queue demo.bar -body "hello galaxy"




## testing

### Install and configure Rabbitmq

	sudo aptitude

	# allow guest login other than loopback (for consul discovery)
	cat << EOF > /etc/rabbitmq/rabbitmq.config
	[{rabbit,[
		{loopback_users, []}
	]}].
	EOF

### Install and configure Consul

	# get the consul binary
	wget https://dl.bintray.com/mitchellh/consul/0.5.2_linux_amd64.zip
	unzip 0.5.2_linux_amd64.zip

	# get the web ui and drop it in your /tmp dir
	wget https://dl.bintray.com/mitchellh/consul/0.5.2_web_ui.zip
	unzip 0.5.2_web_ui.zip
	mv dist /tmp/web-ui


	# run wrapper script that will configure some demo info
	./consul.sh

### Encrypt Rabbit Credentials in your Config
see [https://github.com/xordataexchange/crypt] for details


Update `app.batch` with your info, and run the following to generate your keys:

	gpg2 --batch --armor --gen-key app.batch

This will generate `.pubring.gpg` and `.secring.gpg` for encrypting and
decrypting rabbitmq credentials in your configuration backend.

#### app.batch

	%echo Generating a configuration OpenPGP key
	Key-Type: default
	Subkey-Type: default
	Name-Real: app
	Name-Comment: app configuration key
	Name-Email: app@example.com
	Expire-Date: 0
	%pubring .pubring.gpg
	%secring .secring.gpg
	%commit
	%echo done



