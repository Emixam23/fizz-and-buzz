Feature: health check should return 200 if the service is healthy
And return a status (body) as ok on "/" and "/health"


  Scenario: Endpoint "/" returns health
    When I send a "GET" request to "/"
    Then the response code should be 200
    And the response should match json:
      """
      {
          "status": "Ok"
      }
      """


  Scenario: Endpoint "/health" returns health
    When I send a "GET" request to "/health"
    Then the response code should be 200
    And the response should match json:
      """
      {
          "status": "Ok"
      }
      """
