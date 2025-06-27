Feature: should return stats regarding fizz and buzz "played" requests, with a count for each combination
And return an array with the computed stats on "/stats" (sorted or not, based on given params), with a status code of 200
And return empty an array if no results are present to be computed, with a status code of 404

  Background:
    Given now is the "2021-11-03 19:00:00+02"


  Scenario: Endpoint "/stats" returns nothing/404 because no fizz and buzz requests were recorded
    When I send a "GET" request to "/stats"
    Then the response code should be 404
    And the response should match json:
      """
      []
      """

  Scenario: Endpoint "/stats" returns the stats for the recorded fizz and buzz requests
    When I send "GET" requests to "/fizz-and-buzz" with the given params:
      | n1 | s1   | n2 | s2         | limit |
      | 3  | fizz | 5  | buzz       | 19    |
      | 5  | neko | 8  | demon      | 41    |
      | 4  | bmw  | 5  | koeingsegg | 21    |
      | 3  | fizz | 5  | buzz       | 19    |
    Then the response codes should be 200
    When I send a "GET" request to "/stats"
    Then the response code should be 200
    And the response should match json:
      """
      [
        {
          "n1": 4,
          "s1": "bmw",
          "n2": 5,
          "s2": "koeingsegg",
          "limit": 21,
          "count": 1
        },
        {
          "n1": 5,
          "s1": "neko",
          "n2": 8,
          "s2": "demon",
          "limit": 41,
          "count": 1
        },
        {
          "n1": 3,
          "s1": "fizz",
          "n2": 5,
          "s2": "buzz",
          "limit": 19,
          "count": 2
        }
      ]
      """

  Scenario: Endpoint "/stats" returns the stats for the recorded fizz and buzz requests (with limit parameter given)
    When I send "GET" requests to "/fizz-and-buzz" with the given params:
      | n1 | s1   | n2 | s2         | limit |
      | 3  | fizz | 5  | buzz       | 19    |
      | 5  | neko | 8  | demon      | 41    |
      | 4  | bmw  | 5  | koeingsegg | 21    |
      | 3  | fizz | 5  | buzz       | 19    |
    Then the response codes should be 200
    When I send a "GET" request to "/stats" with the given params:
      | sorted |
      | true   |
    Then the response code should be 200
    And the response should match json:
      """
      [
        {
          "n1": 3,
          "s1": "fizz",
          "n2": 5,
          "s2": "buzz",
          "limit": 19,
          "count": 2
        },
        {
          "n1": 4,
          "s1": "bmw",
          "n2": 5,
          "s2": "koeingsegg",
          "limit": 21,
          "count": 1
        },
        {
          "n1": 5,
          "s1": "neko",
          "n2": 8,
          "s2": "demon",
          "limit": 41,
          "count": 1
        }
      ]
      """

  Scenario: Endpoint "/stats" returns an error has provided parameter is not valid
    When I send a "GET" request to "/stats?sorted=test"
    Then the response code should be 400
    And the response should match json:
      """
      {
        "reason":"strconv.ParseBool: parsing \"test\": invalid syntax"
      }
      """

