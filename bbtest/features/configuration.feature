Feature: Service can be configured

  Scenario: configure log level to ERROR
    Given fio-bco is configured with
      | property  | value |
      | LOG_LEVEL | ERROR |

    Then journalctl of "fio-bco-rest.service" contains following
    """
      Log level set to ERROR
    """

  Scenario: configure log level to INFO
    Given fio-bco is configured with
      | property  | value |
      | LOG_LEVEL | INFO  |

    Then journalctl of "fio-bco-rest.service" contains following
    """
      Log level set to INFO
    """

  Scenario: configure log level to INVALID
    Given fio-bco is configured with
      | property  | value   |
      | LOG_LEVEL | INVALID |

    Then journalctl of "fio-bco-rest.service" contains following
    """
      Log level set to INFO
    """

  Scenario: configure log level to DEBUG
    Given fio-bco is configured with
      | property  | value |
      | LOG_LEVEL | DEBUG |

    Then journalctl of "fio-bco-rest.service" contains following
    """
      Log level set to DEBUG
    """
