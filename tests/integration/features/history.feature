Feature: should return fizz and buzz request's history
And return an array with the history of "played" requests among their params, with a status code of 200
And return an empty array if no history is available, with a status code of 404

  Background:
    Given now is the "2021-11-03 19:00:00+02"


  Scenario: Endpoint "/history" returns nothing/404 because no fizz and buzz requests were recorded
    When I send a "GET" request to "/history"
    Then the response code should be 404
    And the response should match json:
      """
      []
      """


  Scenario: Endpoint "/history" returns the history of fizz and buzz requests
    When I send "GET" requests to "/fizz-and-buzz" with the given params:
      | n1 | s1   | n2 | s2         | limit |
      | 3  | fizz | 5  | buzz       | 19    |
      | 5  | neko | 8  | demon      | 41    |
      | 4  | bmw  | 5  | koeingsegg | 21    |
    Then the response codes should be 200
    When I send a "GET" request to "/history"
    Then the response code should be 200
    And the response should match json:
      """
      [
        {
          "id": 3,
          "request_date": "2021-11-02T19:00:00Z",
          "n1": 4,
          "s1": "bmw",
          "n2": 5,
          "s2": "koeingsegg",
          "limit": 21
        },
        {
          "id": 2,
          "request_date": "2021-11-02T19:00:00Z",
          "n1": 5,
          "s1": "neko",
          "n2": 8,
          "s2": "demon",
          "limit": 41
        },
        {
          "id": 1,
          "request_date": "2021-11-02T19:00:00Z",
          "n1": 3,
          "s1": "fizz",
          "n2": 5,
          "s2": "buzz",
          "limit": 19
        }
      ]
      """


  Scenario: Endpoint "/history" returns the history of fizz and buzz requests (with limit parameter given)
    When I send "GET" requests to "/fizz-and-buzz" with the given params:
      | n1 | s1   | n2 | s2         | limit |
      | 3  | fizz | 5  | buzz       | 19    |
      | 5  | neko | 8  | demon      | 41    |
      | 4  | bmw  | 5  | koeingsegg | 21    |
    Then the response codes should be 200
    When I send a "GET" request to "/history" with the given params:
      | limit |
      | 2     |
    Then the response code should be 200
    And the response should match json:
      """
      [
        {
          "id": 3,
          "request_date": "2021-11-02T19:00:00Z",
          "n1": 4,
          "s1": "bmw",
          "n2": 5,
          "s2": "koeingsegg",
          "limit": 21
        },
        {
          "id": 2,
          "request_date": "2021-11-02T19:00:00Z",
          "n1": 5,
          "s1": "neko",
          "n2": 8,
          "s2": "demon",
          "limit": 41
        }
      ]
      """

  Scenario: Endpoint "/history" returns an error has provided parameter is not valid
    When I send a "GET" request to "/history?limit=test"
    Then the response code should be 400
    And the response should match json:
      """
      {
        "reason":"strconv.ParseUint: parsing \"test\": invalid syntax"
      }
      """

