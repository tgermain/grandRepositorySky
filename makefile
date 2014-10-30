
buildDocker:
	docker build -t tgermain/repo_sky .

runDocker:
	docker run --publish 4444:4321 --publish 4444:4321/udp --name draluca --rm tgermain/repo_sky:latest -s /static/

pushNewVersion:
	docker push tgermain/repo_sky:latest 