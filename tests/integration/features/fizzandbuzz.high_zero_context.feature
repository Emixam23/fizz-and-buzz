Feature: should return fizz and buzz request's result (context zero set to 1, so results starts from 1 to provided limit)
And return an array with the computed results with the request params, with a status code of 200
And return a 422 unprocessable error is the provided limit is not greater than "Zero" context starting point


  Scenario: Endpoint "/fizz-and-buzz" returns expected exercise results
    When Fizz and Buzz service has a zero context of "19"
    When I send "GET" requests to "/fizz-and-buzz" with the given params:
      | n1 | s1   | n2 | s2   | limit |
      | 3  | fizz | 5  | buzz | 31    |
    Then the response code should be 200
    And the response should match json:
      """
      {
        "request": {
          "n1": 3,
          "s1": "fizz",
          "n2": 5,
          "s2": "buzz",
          "limit": 31
        },
        "result": [
          "19",
          "buzz",
          "fizz",
          "22",
          "23",
          "fizz",
          "buzz",
          "26",
          "fizz",
          "28",
          "29",
          "fizzbuzz",
          "31"
        ]
      }
      """

  Scenario: Endpoint "/fizz-and-buzz" returns an error as limit is below zero from context
    When Fizz and Buzz service has a zero context of "19"
    When I send "GET" requests to "/fizz-and-buzz" with the given params:
      | n1 | s1   | n2 | s2   | limit |
      | 3  | fizz | 5  | buzz | 18    |
    Then the response code should be 422
    And the response should match json:
      """
      {
        "reason": "provided limit \"18\" must be greater than \"19\" (provided in configuration)"
      }
      """
