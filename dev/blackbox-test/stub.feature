Feature: API test
  Background: Basic orchestration
    Given container server should be running
    And   container queue should be running
    And   container vault should be running
    And   server is listening on 8080
    And   server is healthy
