Feature: should return fizz and buzz request's result (context zero set to 1, so results starts from 1 to provided limit)
And return an array with the computed results with the request params, with a status code of 200


  Scenario: Endpoint "/fizz-and-buzz" returns expected exercise results
    When I send "GET" requests to "/fizz-and-buzz" with the given params:
      | n1 | s1   | n2 | s2   | limit |
      | 3  | fizz | 5  | buzz | 19    |
    Then the response code should be 200
    And the response should match json:
      """
      {
        "request": {
          "n1": 3,
          "s1": "fizz",
          "n2": 5,
          "s2": "buzz",
          "limit": 19
        },
        "result": [
          "1",
          "2",
          "fizz",
          "4",
          "buzz",
          "fizz",
          "7",
          "8",
          "fizz",
          "buzz",
          "11",
          "fizz",
          "13",
          "14",
          "fizzbuzz",
          "16",
          "17",
          "fizz",
          "19"
        ]
      }
      """

  Scenario: Endpoint "/fizz-and-buzz" returns an error has n1 parameter is missing
    When I send "GET" requests to "/fizz-and-buzz" with the given params:
      | s1   | n2 | s2   | limit |
      | fizz | 5  | buzz | 19    |
    Then the response code should be 400
    And the response should match json:
      """
      {
        "reason": "Key: 'fizzAndBuzzRequestParams.N1' Error:Field validation for 'N1' failed on the 'required' tag"
      }
      """

  Scenario: Endpoint "/fizz-and-buzz" returns an error has s1 parameter is missing
    When I send "GET" requests to "/fizz-and-buzz" with the given params:
      | n1 | n2 | s2   | limit |
      | 3  | 5  | buzz | 19    |
    Then the response code should be 400
    And the response should match json:
      """
      {
        "reason": "Key: 'fizzAndBuzzRequestParams.S1' Error:Field validation for 'S1' failed on the 'required' tag"
      }
      """

  Scenario: Endpoint "/fizz-and-buzz" returns an error has n2 parameter is missing
    When I send "GET" requests to "/fizz-and-buzz" with the given params:
      | n1 | s1   | s2   | limit |
      | 3  | fizz | buzz | 19    |
    Then the response code should be 400
    And the response should match json:
      """
      {
        "reason": "Key: 'fizzAndBuzzRequestParams.N2' Error:Field validation for 'N2' failed on the 'required' tag"
      }
      """

  Scenario: Endpoint "/fizz-and-buzz" returns an error has s2 parameter is missing
    When I send "GET" requests to "/fizz-and-buzz" with the given params:
      | n1 | s1   | n2 | limit |
      | 3  | fizz | 5  | 19    |
    Then the response code should be 400
    And the response should match json:
      """
      {
        "reason": "Key: 'fizzAndBuzzRequestParams.S2' Error:Field validation for 'S2' failed on the 'required' tag"
      }
      """

  Scenario: Endpoint "/fizz-and-buzz" returns an error has n1 parameter is missing
    When I send "GET" requests to "/fizz-and-buzz" with the given params:
      | n1 | s1   | n2 | s2   |
      | 3  | fizz | 5  | buzz |
    Then the response code should be 400
    And the response should match json:
      """
      {
        "reason": "Key: 'fizzAndBuzzRequestParams.Limit' Error:Field validation for 'Limit' failed on the 'required' tag"
      }
      """

