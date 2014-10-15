
buildDocker:
	docker build -t testator .

runDocker:
	docker run --publish 4444:4321 --publish 4444:4321/udp --name test --rm testator:latest