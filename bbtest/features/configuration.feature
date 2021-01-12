Feature: Service can be configured

  Scenario: configure log level to ERROR
    Given fio-bco is configured with
      | property  | value |
      | LOG_LEVEL | ERROR |
    And tenant CONFIGURATION_ERROR is onboarded

    Then journalctl of "fio-bco-import@CONFIGURATION_ERROR.service" contains following
    """
      Log level set to ERROR
    """
    And tenant CONFIGURATION_ERROR is offboarded

  Scenario: configure log level to INFO
    Given fio-bco is configured with
      | property  | value |
      | LOG_LEVEL | INFO  |
    And tenant CONFIGURATION_INFO is onboarded

    Then journalctl of "fio-bco-import@CONFIGURATION_INFO.service" contains following
    """
      Log level set to INFO
    """
    And tenant CONFIGURATION_INFO is offboarded

  Scenario: configure log level to INVALID
    Given fio-bco is configured with
      | property  | value   |
      | LOG_LEVEL | INVALID |
    And tenant CONFIGURATION_INVALID is onboarded

    Then journalctl of "fio-bco-import@CONFIGURATION_INVALID.service" contains following
    """
      Log level set to INFO
    """
    And tenant CONFIGURATION_INVALID is offboarded

  Scenario: configure log level to DEBUG
    Given fio-bco is configured with
      | property  | value |
      | LOG_LEVEL | DEBUG |
    And tenant CONFIGURATION_DEBUG is onboarded

    Then journalctl of "fio-bco-import@CONFIGURATION_DEBUG.service" contains following
    """
      Log level set to DEBUG
    """
    And tenant CONFIGURATION_DEBUG is offboarded
