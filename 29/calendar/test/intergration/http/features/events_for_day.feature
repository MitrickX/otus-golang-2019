Feature: Get events for today
  As API client of calendar service
  I want get events that started today

  Scenario: Get from empty DB
    Given Clean DB

    When I send "GET" request to "http://localhost:8888/events_for_day"

    Then The response code should be 200
    And The response contentType should be "application/json"
    And The response json should match:
    """
      {"result": []}
    """

  Scenario: Get empty list, cause there is no any today
    Given Clean DB
    Given Existing records:
      | id | name       | start_time        | end_time          | before_minutes  | notified_time |
      |  1 | test1      | 2019-12-21 14:00  | 2019-12-21 15:00  | 10              | nil           |
      |  2 | test2      | 2021-12-22 15:00  | 2021-12-22 17:00  | nil             | nil           |

    When I send "GET" request to "http://localhost:8888/events_for_day"

    Then The response code should be 200
    And The response contentType should be "application/json"
    And The response json should match:
    """
      {"result": []}
    """