tag := "1.0.0"

build repo:
    export KO_DOCKER_REPO={{repo}} && \
    ko build . --tags {{tag}}
