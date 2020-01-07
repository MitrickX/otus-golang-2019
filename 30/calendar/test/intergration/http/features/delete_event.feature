Feature: Delete event
  As API client of calendar service
  I want to delete existing event

  Scenario: Delete event, 200 OK
    Given Clean DB
    Given Existing records:
      | id | name       | start_time        | end_time          | before_minutes  | notified_time |
      |  1 | test1      | 2019-12-21 14:00  | 2019-12-21 15:00  | 10              | nil           |
      |  2 | test2      | 2019-12-22 15:00  | 2019-12-22 17:00  | nil             | nil           |
    When I send "POST" request to "http://http:8888/delete_event" with "application/x-www-form-urlencoded" params:
    """
    id=2
    """
    Then The response code should be 200
    And The response contentType should be "application/json"
    And The response json should has field "result" with value match "deleted"
    And The records should match:
      | id | name       | start_time        | end_time          | before_minutes  | notified_time |
      | 1  | test1      | 2019-12-21 14:00  | 2019-12-21 15:00  | 10              | nil           |

  Scenario: Delete event, 400 invalid id
    When I send "POST" request to "http://http:8888/delete_event" with "application/x-www-form-urlencoded" params:
    """
    id=dwdfdf
    """
    Then The response code should be 400
    And The response contentType should be "application/json"
    And The response json should has field "error" with value match "^invalid id"
    And The records should match:
      | id | name       | start_time        | end_time          | before_minutes  | notified_time |
      | 1  | test1      | 2019-12-21 14:00  | 2019-12-21 15:00  | 10              | nil           |

  Scenario: Delete event, 200 OK
    Given Clean DB
    Given Existing records:
      | id | name       | start_time        | end_time          | before_minutes  | notified_time |
      |  1 | test1      | 2019-12-21 14:00  | 2019-12-21 15:00  | 10              | nil           |
    When I send "POST" request to "http://http:8888/delete_event" with "application/x-www-form-urlencoded" params:
    """
    id=2
    """
    Then The response code should be 200
    And The response contentType should be "application/json"
    And The response json should has field "error" with value match "^couldn't delete event"
    And The records should match:
      | id | name       | start_time        | end_time          | before_minutes  | notified_time |
      | 1  | test1      | 2019-12-21 14:00  | 2019-12-21 15:00  | 10              | nil           |