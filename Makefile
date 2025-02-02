test:
	go test ./...

# Install ebiten for ubuntu
# https://ebitengine.org/en/documents/install.html#Debian_/_Ubuntu
install:
	sudo apt install libc6-dev libgl1-mesa-dev libxcursor-dev libxi-dev libxinerama-dev libxrandr-dev libxxf86vm-dev libasound2-dev pkg-config
	sudo apt install libglfw3 libglfw3-dev libglew-dev

example1:
	go run example/simple_string_matcher/main.go

example2:
	go run example/multi_string_matcher/main.go

example3:
	go run example/traveling_salesman/*.go
