# go-test

## Introduction

The purpose of this test is to ensure you have some basic knowledge of the golang's standard http library, and that you are familiar with objects. Please do not fork the repository and do not push your solution back to github. Please compress your complete git solution and send it to us as zip package. PLease ensure you include .git in your folder so that view your commits.

## Deliverable
* Clone repository
* Commit your solution
* zip the folder
* Send it back to Elective

## Test
Create go application with:
* http server `GET http://localhost/countries` endpoint that invokes go function that will start printing a list of countries to the currently open web socket connections at `GET ws://localhost/events`. 1 country every 0.5 second at a time.
* The connection to `GET http://localhost/countries` should do the work asyncroniously and respond with 200 OK imediatally (not wait for data to be send to websocket).
* All open websocket connections should start printing a list of countries, one at a time in 0.5 second interval
* Use the interface provided in `index.html` to print countries and `countries.txt` to populate the list

![alt text](https://github.com/electivegroup/go-test/blob/master/go-test.png "Solution diagram")

## Evaluation
Evaluation criteria:

1. Validity
* The software must do what it is supposed to do.

2. Maintainability
* The source code should use interfaces.
* The source code should follow golang naming conventions.
* The source code must be readable.
* The source code must be tested.
* The source code must be commented.

3. Exploitability
* The support team must be able to troubleshoot with the logs.
* Any configuration must be easy.
* The software must be documented.
* The software must be easy to deploy.

## Help
Should you have any questions please ask the team. Good luck!
