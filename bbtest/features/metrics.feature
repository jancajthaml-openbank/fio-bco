Feature: Metrics test

  Scenario: metrics have expected keys
    Given tenant M1 is onboarded
    And   fio-bco is configured with
      | property            | value |
      | METRICS_REFRESHRATE |    1s |

    Then metrics file reports/blackbox-tests/metrics/metrics.M1.json should have following keys:
      | key                      |
      | createdTokens            |
      | deletedTokens            |
      | importedTransactions     |
      | importedTransfers        |
      | syncLatency              |
    And metrics file reports/blackbox-tests/metrics/metrics.M1.json has permissions -rw-r--r--

    And metrics file reports/blackbox-tests/metrics/metrics.json should have following keys:
      | key                      |
      | createTokenLatency       |
      | deleteTokenLatency       |
      | getTokenLatency          |
      | memoryAllocated          |
    And metrics file reports/blackbox-tests/metrics/metrics.json has permissions -rw-r--r--

  Scenario: metrics can remembers previous values after reboot
    Given tenant M2 is onboarded
    And   fio-bco is configured with
      | property            | value |
      | METRICS_REFRESHRATE |    1s |

    Then metrics file reports/blackbox-tests/metrics/metrics.M2.json reports:
      | key                      | value |
      | createdTokens            |     0 |
      | deletedTokens            |     0 |
      | importedTransactions     |     0 |
      | importedTransfers        |     0 |

    When token M2/A is created
    Then metrics file reports/blackbox-tests/metrics/metrics.M2.json reports:
      | key                      | value |
      | createdTokens            |     1 |
      | deletedTokens            |     0 |
      | importedTransactions     |     0 |
      | importedTransfers        |     0 |

    When restart unit "fio-bco-import@M2.service"
    Then metrics file reports/blackbox-tests/metrics/metrics.M2.json reports:
      | key                      | value |
      | createdTokens            |     1 |
      | deletedTokens            |     0 |
      | importedTransactions     |     0 |
      | importedTransfers        |     0 |
