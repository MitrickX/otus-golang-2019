Feature: Get notification
  As API client of calendar service
  I want receive notifications about events

  Scenario: Get notification
    Given Clean DB
    Given Existing records:
    # Special notation +someDuration says someDuration from now
    # Special notation Y-m-d H:i:s will substitute components or current datetime
    #   Y - full year
    #   m - month with leading zero: 01, 02, ... 12
    #   d - day of month with leading zero: 01, 02, ... 29, 30, 31
    #   H - 24 format hours with leading zero: 00, 01, 02, ... 22, 23, 24
    #   i - minutes with leading zero: 00, 01, ... , 58, 59
    #   s - seconds with leading zero: 00, 01, ... , 58, 59
      | id | name       | start_time        | end_time          | before_minutes  | notified_time |
      |  1 | test1      | +11m              | +100m             | 10              | nil           |
      |  2 | test2      | +11m              | +100m             | 10              | Y-m-d H:i:s   |
      |  3 | test3      | +11m              | +100m             | nil             | nil           |
      |  4 | test4      | +100m             | +200m             | 99              | nil           |
    When I after wait "90s" should receive notification about events of ids "1,4"
    And Field "notified_time" of records "1,4" must be not nil