Feature: Update event
  As API client of calendar service
  I want to update existing event

  Scenario: Update event, 200 OK
    Given Clean DB
    Given Existing records:
      | id | name       | start_time        | end_time          | before_minutes  | notified_time |
      |  1 | test1      | 2019-12-21 14:00  | 2019-12-21 15:00  | 10              | nil           |
      |  2 | test2      | 2019-12-22 15:00  | 2019-12-22 17:00  | nil             | nil           |
    When I send "POST" request to "http://http:8888/update_event" with "application/x-www-form-urlencoded" params:
    """
    id=2&name=updated&start=2019-12-22 15:15&end=2019-12-22 18:00&beforeMinutes=5
    """
    Then The response code should be 200
    And The response contentType should be "application/json"
    And The response json should has field "result" with value match "updated"
    And The records should match:
      | id | name       | start_time        | end_time          | before_minutes  | notified_time |
      | 1  | test1      | 2019-12-21 14:00  | 2019-12-21 15:00  | 10              | nil           |
      | 2  | updated    | 2019-12-22 15:15  | 2019-12-22 18:00  | 5               | nil           |

  Scenario: Update event, 400 invalid start date
    Given Clean DB
    When I send "POST" request to "http://http:8888/update_event" with "application/x-www-form-urlencoded" params:
    """
    id=2&name=updated&start=dfsdfsd&end=2019-12-22 18:00&beforeMinutes=5
    """
    Then The response code should be 400
    And The response contentType should be "application/json"
    And The response json should has field "error" with value match "^invalid format of datetime"
    And The DB should be clean

  Scenario: Update event, 400 invalid end date
    When I send "POST" request to "http://http:8888/update_event" with "application/x-www-form-urlencoded" params:
    """
    id=2&name=updated&start=2019-12-22 15:15&end=dfwefdfwe&beforeMinutes=5
    """
    Then The response code should be 400
    And The response contentType should be "application/json"
    And The response json should has field "error" with value match "^invalid format of datetime"
    And The DB should be clean

  Scenario: Update event, 400 invalid id
    When I send "POST" request to "http://http:8888/update_event" with "application/x-www-form-urlencoded" params:
    """
    id=dwdfdf&name=updated&start=2019-12-22 15:15&end=2019-12-22 18:00&beforeMinutes=5
    """
    Then The response code should be 400
    And The response contentType should be "application/json"
    And The response json should has field "error" with value match "^invalid id"
    And The DB should be clean

  Scenario: Update event, 200 OK
    Given Clean DB
    When I send "POST" request to "http://http:8888/update_event" with "application/x-www-form-urlencoded" params:
    """
    id=2&name=updated&start=2019-12-22 15:15&end=2019-12-22 18:00&beforeMinutes=5
    """
    Then The response code should be 200
    And The response contentType should be "application/json"
    And The response json should has field "error" with value match "^couldn't update event"
    And The DB should be clean