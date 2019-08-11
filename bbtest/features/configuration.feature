Feature: Service can be configured

  Scenario: configure log level to DEBUG
    Given tenant CONFIGURATION_DEBUG is onboarded
    And   fio-bco is configured with
      | property  | value |
      | LOG_LEVEL | DEBUG |

    Then journalctl of "fio-bco-import@CONFIGURATION_DEBUG.service" contains following
    """
      Log level set to DEBUG
    """

  Scenario: configure log level to ERROR
    Given tenant CONFIGURATION_ERROR is onboarded
    And   fio-bco is configured with
      | property  | value |
      | LOG_LEVEL | ERROR |

    Then journalctl of "fio-bco-import@CONFIGURATION_ERROR.service" contains following
    """
      Log level set to ERROR
    """

  Scenario: configure log level to INFO
    Given tenant CONFIGURATION_INFO is onboarded
    And   fio-bco is configured with
      | property  | value |
      | LOG_LEVEL | INFO  |

    Then journalctl of "fio-bco-import@CONFIGURATION_INFO.service" contains following
    """
      Log level set to INFO
    """
