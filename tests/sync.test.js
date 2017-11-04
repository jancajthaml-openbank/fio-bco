jest.mock("jsonfile-promised");

beforeEach(() => {
  jest.clearAllMocks();
});

test("get transaction checkpoint, existing tenant and existing account number", async () => {
  const sync = require("../modules/sync.js");
  const jsonfile = require("jsonfile-promised");

  jsonfile.readFile = jest.fn(() => {
    return {
      "johny": {
        "CZ7920100000002400222233": {"idTransactionTo": 2151261787},
        "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
      },
      "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665}},
      "johny3": {"CZ7120100000002700968855": {"idTransactionTo": 14438087888}}
    }
  });

  const checkpoint = await sync.getTransactionCheckpoint("testdb.json", "johny", "CZ7120100000002700968855");
  expect(checkpoint).toBe(14434862430);
  expect(jsonfile.readFile.mock.calls[0][0]).toBe("testdb.json");
});

test("get transaction checkpoint, existing tenant and none-existing account number", async () => {
  const sync = require("../modules/sync.js");
  const jsonfile = require("jsonfile-promised");

  jsonfile.readFile = jest.fn(() => {
    return {
      "johny": {
        "CZ7920100000002400222233": {"idTransactionTo": 2151261787},
        "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
      },
      "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665}},
      "johny3": {"CZ7120100000002700968855": {"idTransactionTo": 14438087888}}
    }
  });

  const checkpoint = await sync.getTransactionCheckpoint("testdb.json", "johny", "CZ7120100000002700961234");
  expect(checkpoint).toBeNull();
  expect(jsonfile.readFile.mock.calls[0][0]).toBe("testdb.json");
});

test("get transaction checkpoint, none-existing tenant", async () => {
  const sync = require("../modules/sync.js");
  const jsonfile = require("jsonfile-promised");

  jsonfile.readFile = jest.fn(() => {
    return {
      "johny": {
        "CZ7920100000002400222233": {"idTransactionTo": 2151261787},
        "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
      },
      "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665}},
      "johny3": {"CZ7120100000002700968855": {"idTransactionTo": 14438087888}}
    }
  });

  const checkpoint = await sync.getTransactionCheckpoint("testdb.json", "the_tenant", "CZ7120100000002700968855");
  expect(checkpoint).toBeNull();
  expect(jsonfile.readFile.mock.calls[0][0]).toBe("testdb.json");
});

test("get transaction checkpoint, none-existing db file", async () => {
  const sync = require("../modules/sync.js");
  const jsonfile = require("jsonfile-promised");
  const mockedError = new Error();
  mockedError.code = "ENOENT";

  jsonfile.readFile = jest.fn(() => {
    throw mockedError;
  });

  const checkpoint = await sync.getTransactionCheckpoint("testdb.json", "the_tenant", "CZ7120100000002700968855");

  expect(checkpoint).toBeNull();
  expect(jsonfile.readFile.mock.calls[0][0]).toBe("testdb.json");
});

test("get transaction checkpoint, rethrowing unknown error", async () => {
  const sync = require("../modules/sync.js");
  const jsonfile = require("jsonfile-promised");
  const mockedError = new Error();
  mockedError.code = "UNKNOWN";

  jsonfile.readFile = jest.fn(() => {
    throw mockedError;
  });

  let error;
  try {
    const checkpoint = await sync.getTransactionCheckpoint("testdb.json", "the_tenant", "CZ7120100000002700968855");
  } catch (e) {
    error = e;
  }
  expect(error).toBe(mockedError);
  expect(jsonfile.readFile.mock.calls[0][0]).toBe("testdb.json");
});

test("set transaction checkpoint, none-existing tenant", async () => {
  const sync = require("../modules/sync.js");
  const jsonfile = require("jsonfile-promised");

  jsonfile.readFile = jest.fn(() => {
    return {
      "johny": {
        "CZ7920100000002400222233": {"idTransactionTo": 2151261787},
        "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
      },
      "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665}},
      "johny3": {"CZ7120100000002700968855": {"idTransactionTo": 14438087888}}
    };
  });
  jsonfile.writeFile = jest.fn();

  await sync.setTransactionCheckpoint("testdb.json", "newTenant", "accountNumber", 12345);

  expect(jsonfile.writeFile.mock.calls[0][1]).toEqual({
    "johny": {
      "CZ7920100000002400222233": {"idTransactionTo": 2151261787},
      "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
    },
    "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665}},
    "johny3": {"CZ7120100000002700968855": {"idTransactionTo": 14438087888}},
    "newTenant": {"accountNumber": {"idTransactionTo": 12345}}
  });
});

test("set transaction checkpoint, existing tenant, none-existing account number", async () => {
  const sync = require("../modules/sync.js");
  const jsonfile = require("jsonfile-promised");

  jsonfile.readFile = jest.fn(() => {
    return {
      "johny": {
        "CZ7920100000002400222233": {"idTransactionTo": 2151261787},
        "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
      },
      "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665}},
      "johny3": {"CZ7120100000002700968855": {"idTransactionTo": 14438087888}}
    };
  });
  jsonfile.writeFile = jest.fn();

  await sync.setTransactionCheckpoint("testdb.json", "johny3", "accountNumber", 12345);

  expect(jsonfile.writeFile.mock.calls[0][1]).toEqual({
    "johny": {
      "CZ7920100000002400222233": {"idTransactionTo": 2151261787},
      "CZ7120100000002700968855": {"idTransactionTo": 14434862430}
    },
    "johny2": {"CZ7120100000002700968855": {"idTransactionTo": 14434862665}},
    "johny3": {
      "CZ7120100000002700968855": {"idTransactionTo": 14438087888},
      "accountNumber": {"idTransactionTo": 12345}
    }
  });
});

test("set transaction checkpoint, none existing db file", async () => {
  const sync = require("../modules/sync.js");
  const jsonfile = require("jsonfile-promised");
  const mockedError = new Error();
  mockedError.code = "ENOENT";

  jsonfile.readFile = jest.fn(() => {
    throw mockedError;
  });
  jsonfile.writeFile = jest.fn();

  await sync.setTransactionCheckpoint("testdb.json", "johny3", "accountNumber", 12345);

  expect(jsonfile.writeFile.mock.calls[0][1]).toEqual( {"johny3": {"accountNumber": {"idTransactionTo": 12345 } } } );
});

test("set transaction checkpoint, rethrow unknown error", async () => {
  const sync = require("../modules/sync.js");
  const jsonfile = require("jsonfile-promised");
  const mockedError = new Error();
  mockedError.code = "UNKNOWN";

  jsonfile.readFile = jest.fn(() => {
    throw mockedError;
  });

  let error;
  try {
    await sync.setTransactionCheckpoint("testdb.json", "johny3", "accountNumber", 12345);
  } catch (e) {
    error = e;
  }

  expect(error).toEqual(mockedError);
});