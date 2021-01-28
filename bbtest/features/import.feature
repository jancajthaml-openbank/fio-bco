Feature: Fio Bank import

  Scenario: import from gateway token
    Given fio gateway contains following statements
      | key       | value |
    And fio-bco is configured with
      | property  | value |
      | SYNC_RATE |    1s |
    And tenant IMPORT is onboarded
    And token IMPORT/importToken is created
