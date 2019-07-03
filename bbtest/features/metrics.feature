@metrics
Feature: Metrics test

  Scenario: metrics have expected keys
    And   tenant M2 is onbdoarded
    And   fio-bco is reconfigured with
    """
      METRICS_REFRESHRATE=1s
    """

    Then metrics file /reports/metrics.M2.json should have following keys:
    """
      createdTokens
      deletedTokens
      exportAccountLatency
      exportTransactionLatency
      importTransactionLatency
      importAccountLatency
      importedAccounts
      importedTransfers
      exportedAccounts
      exportedTransfers
      syncLatency
    """
    And metrics file /reports/metrics.M2.json has permissions -rw-r--r--
    And metrics file /reports/metrics.json should have following keys:
    """
      createTokenLatency
      deleteTokenLatency
      getTokenLatency
    """
    And metrics file /reports/metrics.json has permissions -rw-r--r--
