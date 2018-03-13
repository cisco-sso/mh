.PHONY: dockerBuild
dockerBuild: Dockerfile
	docker build -t ${DOCKER_IMAGE} .

.PHONY: dockerPush
dockerPush: dockerBuild
	docker push ${DOCKER_IMAGE}

.PHONY: dockerTest
dockerTest:
	docker run -it --rm ${DOCKER_IMAGE}

.PHONY: awsEcrCreateRepo
awsEcrCreateRepo:
	aws ecr create-repository --repository-name ${DOCKER_NAMESPACE}/${PROJECT}

