Feature: Get events for today
  As API client of calendar service
  I want get events that started today

  Scenario: Get from empty DB
    Given Clean DB

    When I send "GET" request to "http://http:8888/events_for_day"

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

    When I send "GET" request to "http://http:8888/events_for_day"

    Then The response code should be 200
    And The response contentType should be "application/json"
    And The response json should match:
    """
      {"result": []}
    """

  Scenario: Get not empty list
    Given Clean DB
    Given Existing records:
    # Special notation Y-m-d H:i:s will substitute components or current datetime
    #   Y - full year
    #   m - month with leading zero: 01, 02, ... 12
    #   d - day of month with leading zero: 01, 02, ... 29, 30, 31
    #   H - 24 format hours with leading zero: 00, 01, 02, ... 22, 23, 24
    #   i - minutes with leading zero: 00, 01, ... , 58, 59
    #   s - seconds with leading zero: 00, 01, ... , 58, 59
      | id | name       | start_time          | end_time            | before_minutes  | notified_time |
      |  1 | test1      | 2019-12-21 14:00    | 2019-12-21 15:00    | 10              | nil           |
      |  2 | test2      | 2021-12-22 15:00    | 2021-12-22 17:00    | nil             | nil           |
      |  3 | test3      | Y-m-d 10:00:00      | Y-m-d 15:00:00      | 4               | nil           |
      |  4 | test4      | 2018-10-12 17:00    | 2018-10-13 11:00    | nil             | nil           |
      |  5 | test5      | Y-m-d 14:00:00      | Y-m-d 16:00:00      | nil             | nil           |
    When I send "GET" request to "http://http:8888/events_for_day"

    Then The response code should be 200
    And The response contentType should be "application/json"
    And The response json is EventListResponse filled with events of ids "3,5"