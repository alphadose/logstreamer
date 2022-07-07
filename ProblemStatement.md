# Tyk Platform squad - Go developer Take Home Exercise

Engineers at Tyk build highly concurrent systems in a multiple databases enviroments.
This take home assignment is designed to evaluate similar to real world tasks that are involved with this role. 

## Instructions
The goal of this excercise is to write an application that will read data from a file, parse it and write it in different data storages and services.

The application should be able to read the data from the `data.txt` file. Each line of the file contains a new JSON payload.

After reading and parsing the data, the aplication must write it to the database of your election (MongoDB or any SQL database) and also to a simple gRPC service that will receive the records and save them in memory. The aplication should log if the data writing was successful.

The error handling strategy is on the candidate’s decision. It also must contain a mechanism to control failures while writting data into the db or the grpc service (in case any of them is down).

## Submission Guidance

The finished solution **should:**
- Be written in Go.
- Contain clear and concise documented code.
- Have clear instructions on how to run it.
- Be well tested to the level you would expect in a working environment. Note that you can use Mocks for the storages tests or/and integrate docker-images on CI level.
- Be simple and concise.

It's okay to use third party libraries.


## How to submit your exercise

- Include your name in the README.
- Create a public [GitHub](https://help.github.com/en/articles/create-a-repo) repository with the code.
- Let us know you've completed the exercise sending the repo link via email.
