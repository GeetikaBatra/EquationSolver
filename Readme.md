# Equation Solver

## To build and push the docker image to docker hub
`make build-push`

## Run service in local
`make run-local`

## Run tests in local
`make run-test`

## POST REQUEST

`curl -X POST https://equationsolver.azurewebsites.net/equate -d "first_equation=x+y=3&second_equation=x-y=5"`
