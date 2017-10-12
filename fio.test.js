let rewire = require("rewire");
let fio = rewire("./fio");

test("get account number from fio transaction", () => {
  let getAccount = fio.__get__("getAccount");
  expect(getAccount({
    "column2": {
      "value": "CZ7920100000002400222345"
    }
  })).toBe("CZ7920100000002400222345");
});

test("get default account number from fio transaction", () => {
  let getAccount = fio.__get__("getAccount");
  expect(getAccount({})).toBe("FIO");
});