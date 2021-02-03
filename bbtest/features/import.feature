Feature: Fio Bank import

  Scenario: import from gateway token
    Given fio gateway contains following statements
      | key       | value |
    And fio-bco is configured with
      | property  | value |
      | SYNC_RATE |    8h |
    And tenant IMPORT is onboarded
    And token IMPORT/importToken is created
    And token IMPORT/importToken is ordered to synchronize
