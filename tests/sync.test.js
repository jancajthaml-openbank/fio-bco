jest.mock("jsonfile-promised")

beforeEach(() => {
  jest.clearAllMocks()
})

test("get transaction checkpoint, existing tenant and existing account number", async () => {
  const sync = require("../modules/sync.js")
  const jsonfile = require("jsonfile-promised")

  jsonfile.readFile = jest.fn(() => {
    return {
      "johny": {
        "CZ7920100000002400222233": {"idTransactionTo": 2151261787},
        "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
      },
      "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665}},
      "johny3": {"CZ7120100000002700968855": {"idTransactionTo": 14438087888}}
    }
  })

  const checkpoint = await sync.getTransactionCheckpoint("testdb.json", "johny", "CZ7120100000002700968855")
  expect(checkpoint).toBe(14434862430)
  expect(jsonfile.readFile.mock.calls[0][0]).toBe("testdb.json")
})

test("get transaction checkpoint, existing tenant and none-existing account number", async () => {
  const sync = require("../modules/sync.js")
  const jsonfile = require("jsonfile-promised")

  jsonfile.readFile = jest.fn(() => {
    return {
      "johny": {
        "CZ7920100000002400222233": {"idTransactionTo": 2151261787},
        "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
      },
      "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665}},
      "johny3": {"CZ7120100000002700968855": {"idTransactionTo": 14438087888}}
    }
  })

  const checkpoint = await sync.getTransactionCheckpoint("testdb.json", "johny", "CZ7120100000002700961234")
  expect(checkpoint).toBeNull()
  expect(jsonfile.readFile.mock.calls[0][0]).toBe("testdb.json")
})

test("get transaction checkpoint, none-existing tenant", async () => {
  const sync = require("../modules/sync.js")
  const jsonfile = require("jsonfile-promised")

  jsonfile.readFile = jest.fn(() => {
    return {
      "johny": {
        "CZ7920100000002400222233": {"idTransactionTo": 2151261787},
        "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
      },
      "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665}},
      "johny3": {"CZ7120100000002700968855": {"idTransactionTo": 14438087888}}
    }
  })

  const checkpoint = await sync.getTransactionCheckpoint("testdb.json", "the_tenant", "CZ7120100000002700968855")
  expect(checkpoint).toBeNull()
  expect(jsonfile.readFile.mock.calls[0][0]).toBe("testdb.json")
})

test("get transaction checkpoint, none-existing db file", async () => {
  const sync = require("../modules/sync.js")
  const jsonfile = require("jsonfile-promised")
  const mockedError = new Error()
  mockedError.code = "ENOENT"

  jsonfile.readFile = jest.fn(() => {
    throw mockedError
  })

  const checkpoint = await sync.getTransactionCheckpoint("testdb.json", "the_tenant", "CZ7120100000002700968855")

  expect(checkpoint).toBeNull()
  expect(jsonfile.readFile.mock.calls[0][0]).toBe("testdb.json")
})

test("get transaction checkpoint, rethrowing unknown error", async () => {
  const sync = require("../modules/sync.js")
  const jsonfile = require("jsonfile-promised")
  const VError = require('verror')
  const error = new Error('some error')

  jsonfile.readFile = jest.fn(() => {
    throw error
  })

  try {
    await sync.getTransactionCheckpoint("testdb.json", "the_tenant", "CZ7120100000002700968855")
  } catch (err) {
    expect(err).toEqual(new VError(error, `Error when reading DB file testdb.json`))
  }

  expect(jsonfile.readFile.mock.calls[0][0]).toBe("testdb.json")
})

test("set transaction checkpoint, none-existing tenant", async () => {
  const sync = require("../modules/sync.js")
  const jsonfile = require("jsonfile-promised")

  jsonfile.readFile = jest.fn(() => {
    return {
      "johny": {
        "CZ7920100000002400222233": {"idTransactionTo": 2151261787},
        "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
      },
      "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665}},
      "johny3": {"CZ7120100000002700968855": {"idTransactionTo": 14438087888}}
    }
  })
  jsonfile.writeFile = jest.fn()

  await sync.setTransactionCheckpoint("testdb.json", "newTenant", "accountNumber", "test_token", 12345)

  expect(jsonfile.writeFile.mock.calls[0][1]).toEqual({
    "johny": {
      "CZ7920100000002400222233": {"idTransactionTo": 2151261787},
      "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
    },
    "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665}},
    "johny3": {"CZ7120100000002700968855": {"idTransactionTo": 14438087888}},
    "newTenant": {"accountNumber": {"idTransactionTo": 12345, "token": "test_token"}}
  })
})

test("set transaction checkpoint, existing tenant, none-existing account number", async () => {
  const sync = require("../modules/sync.js")
  const jsonfile = require("jsonfile-promised")

  jsonfile.readFile = jest.fn(() => {
    return {
      "johny": {
        "CZ7920100000002400222233": {"idTransactionTo": 2151261787},
        "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
      },
      "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665}},
      "johny3": {"CZ7120100000002700968855": {"idTransactionTo": 14438087888}}
    }
  })
  jsonfile.writeFile = jest.fn()

  await sync.setTransactionCheckpoint("testdb.json", "johny3", "accountNumber", "test_token", 12345)

  expect(jsonfile.writeFile.mock.calls[0][1]).toEqual({
    "johny": {
      "CZ7920100000002400222233": {"idTransactionTo": 2151261787},
      "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
    },
    "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665}},
    "johny3": {
      "CZ7120100000002700968855": {"idTransactionTo": 14438087888},
      "accountNumber": {"idTransactionTo": 12345, "token": "test_token"}
    }
  })
})

test("set transaction checkpoint, none existing db file", async () => {
  const sync = require("../modules/sync.js")
  const jsonfile = require("jsonfile-promised")
  const mockedError = new Error()
  mockedError.code = "ENOENT"

  jsonfile.readFile = jest.fn(() => {
    throw mockedError
  })
  jsonfile.writeFile = jest.fn()

  await sync.setTransactionCheckpoint("testdb.json", "johny3", "accountNumber", "test_token", 12345)

  expect(jsonfile.writeFile.mock.calls[0][1]).toEqual( {"johny3": {"accountNumber": {"idTransactionTo": 12345, "token": "test_token" } } } )
})

test("set transaction checkpoint, rethrow unknown error", async () => {
  const sync = require("../modules/sync.js")
  const jsonfile = require("jsonfile-promised")
  const mockedError = new Error()
  mockedError.code = "UNKNOWN"

  jsonfile.readFile = jest.fn(() => {
    throw mockedError
  })

  let error
  try {
    await sync.setTransactionCheckpoint("testdb.json", "johny3", "accountNumber", "test_token", 12345)
  } catch (e) {
    error = e
  }

  expect(error).toEqual(mockedError)
})

test("getTransactionCheckpointByToken - token exists", async () => {
  const sync = require("../modules/sync.js")
  const jsonfile = require("jsonfile-promised")

  const testDb = "db"
  const testTenantName = "test_tenantName"
  const testToken = "test_token"
  const testIdTransactionTo = 14438087888

  jsonfile.readFile = jest.fn()
    .mockImplementationOnce(() => {
       return {
        "johny": {
          "CZ7920100000002400222233": {"idTransactionTo": 2151261787, "token": "another_token"},
          "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
        },
        "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665, "token": "another_token"}},
         [testTenantName]: {"CZ7120100000002700968855": {"idTransactionTo": testIdTransactionTo, "token": testToken}}
      }
    })

  const checkpoint = await sync.getTransactionCheckpointByToken(testDb, testTenantName, testToken)

  expect(jsonfile.readFile).toHaveBeenCalledTimes(1)
  expect(jsonfile.readFile.mock.calls[0][0]).toBe(testDb)
  expect(checkpoint).toBe(testIdTransactionTo)
})

test("getTransactionCheckpointByToken - tenant doesn't exist", async () => {
  const sync = require("../modules/sync.js")
  const jsonfile = require("jsonfile-promised")

  const testDb = "db"
  const testTenantName = "testTenant"
  const testToken = "test_token"

  jsonfile.readFile = jest.fn()
    .mockImplementationOnce(() => {
      return {
        "johny": {
          "CZ7920100000002400222233": {"idTransactionTo": 2151261787, "token": "another_token"},
          "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
        },
        "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665, "token": "another_token"}}
      }
    })

  const checkpoint = await sync.getTransactionCheckpointByToken(testDb, testTenantName, testToken)

  expect(checkpoint).toBeNull()
})