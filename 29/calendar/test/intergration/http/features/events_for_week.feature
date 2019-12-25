Feature: Get events for current week
  As API client of calendar service
  I want get events that started in this week

  Scenario: Get from empty DB
    Given Clean DB

    When I send "GET" request to "http://localhost:8888/events_for_week"

    Then The response code should be 200
    And The response contentType should be "application/json"
    And The response json should match:
    """
      {"result": []}
    """

  Scenario: Get empty list, cause there is no any this week
    Given Clean DB
    Given Existing records:
      | id | name       | start_time        | end_time          | before_minutes  | notified_time |
      |  1 | test1      | 2010-12-21 14:00  | 2010-12-21 15:00  | 10              | nil           |
      |  2 | test2      | 2026-12-22 15:00  | 2026-12-22 17:00  | nil             | nil           |

    When I send "GET" request to "http://localhost:8888/events_for_week"

    Then The response code should be 200
    And The response contentType should be "application/json"
    And The response json should match:
    """
      {"result": []}
    """

  Scenario: Get not empty list
    Given Clean DB
    Given Existing records:
    # Special notations for substitute current week day
    #	Mon
    #	Tue
    #	Wed
    #	Thu
    #	Fri
    #	Sat
    # 	Sun
    #   H - 24 format hours with leading zero: 00, 01, 02, ... 22, 23, 24
    #   i - minutes with leading zero: 00, 01, ... , 58, 59
    #   s - seconds with leading zero: 00, 01, ... , 58, 59
      | id | name       | start_time          | end_time            | before_minutes  | notified_time |
      |  1 | test1      | 2019-12-21 14:00    | 2019-12-21 15:00    |  10             | nil           |
      |  2 | test2      | 2021-12-22 15:00    | 2021-12-22 17:00    | nil             | nil           |
      |  3 | test3      | Mon 10:00:00        | Mon 15:00:00        |  4              | nil           |
      |  4 | test4      | 2018-10-12 17:00    | 2018-10-13 11:00    | nil             | nil           |
      |  5 | test5      | Tue 14:00:00        | Tue 16:00:00        | nil             | nil           |
      |  6 | test6      | Wed 14:00:00        | Wed 16:00:00        | nil             | nil           |
      |  7 | test7      | Thu 14:00:00        | Thu 16:00:00        | nil             | nil           |
      |  8 | test8      | Fri 14:00:00        | Fri 17:00:00        | 100             | nil           |
      |  9 | test9      | 2022-01-12 01:00:00 | 2022-01-12 02:00:00 |  10             | nil           |
      | 10 | test10     | Sat H:i:s           | Sat H:i:s           |   1             | nil           |
      | 11 | test11     | 2001-01-01 12:00:00 | 2001-01-12 H:i:s    |   2             | nil           |
      | 12 | test12     | Sun H:i:s           | Sun H:i:s           | nil             | nil           |

    When I send "GET" request to "http://localhost:8888/events_for_week"

    Then The response code should be 200
    And The response contentType should be "application/json"
    And The response json is EventListResponse filled with events of ids "3,5,6,7,8,10,12"