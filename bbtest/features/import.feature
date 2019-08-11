Feature: Fio Bank import

  Scenario: import from gateway token
    Given fio gateway contains following statements
      | key       | value |

    And tenant IMPORT is onboarded
    And fio-bco is configured with
      | property  | value |
      | SYNC_RATE |    1s |
    And token IMPORT/importToken is created
