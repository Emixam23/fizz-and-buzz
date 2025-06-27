Feature: should return, based on the stats, regarding fizz and buzz requests, the most "played" params with the relative count of it
And return a json payload containing the most "played" request params combination, with a status code of 200
And return empty an array if no results are present to be computed, with a status code of 404

  Background:
    Given now is the "2021-11-03 19:00:00+02"


  Scenario: Endpoint "/stats/most_used" returns nothing/404 because no fizz and buzz requests were recorded
    When I send a "GET" request to "/stats/most_used"
    Then the response code should be 404
    And the response should match json:
      """
      {}
      """

  Scenario: Endpoint "/stats/most_used" returns the most "played" request params combination
    When I send "GET" requests to "/fizz-and-buzz" with the given params:
      | n1 | s1   | n2 | s2         | limit |
      | 3  | fizz | 5  | buzz       | 19    |
      | 5  | neko | 8  | demon      | 41    |
      | 4  | bmw  | 5  | koeingsegg | 21    |
      | 3  | fizz | 5  | buzz       | 19    |
    Then the response codes should be 200
    When I send a "GET" request to "/stats/most_used"
    Then the response code should be 200
    And the response should match json:
      """
      {
        "n1": 3,
        "s1": "fizz",
        "n2": 5,
        "s2": "buzz",
        "limit": 19,
        "count": 2
      }
      """

