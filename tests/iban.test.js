
test("calculate IBAN - Czech Republic", async () => {
  const iban = require("../modules/iban.js")
  const nothing = undefined

  expect(iban.calculateCzech("1111", "2222222222")).toBe("CZ4911110000002222222222")
  expect(iban.calculateCzech(nothing, "2222222222")).toBe("2222222222")
  expect(iban.calculateCzech(nothing, nothing)).toBe(undefined)
})
