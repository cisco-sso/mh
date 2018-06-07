.PHONY: dockerBuild
dockerBuild: Dockerfile
	docker build -t ${DOCKER_IMAGE} .

## We use Docker Hub's Automated Build, so this target isn't commonly used.
.PHONY: dockerPush
dockerPush: dockerBuild
	docker push ${DOCKER_IMAGE}

.PHONY: dockerTest
dockerTest:
	docker run -it --rm ${DOCKER_IMAGE}
