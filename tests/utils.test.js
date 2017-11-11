test("getMax should function same like Math.max", () => {
  const utils = require("../modules/utils.js")

  let set = [...new Array(10)].map(() => Math.round(Math.random() * 1000))

  expect(set.reduce(utils.getMax)).toBe(Math.max.apply(null, set))
})

