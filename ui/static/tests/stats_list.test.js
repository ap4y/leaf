import StatsList from "../stats_list.js";

test("render", () => {
  const statsList = new StatsList();
  expect(statsList.element).not.toBeNull();

  statsList.stats = [
    {
      card: "foo",
      stats: {
        LastReviewedAt: new Date(0).toString(),
        Interval: 0.2,
        Difficulty: 1.3,
        Historical: [
          { interval: 0.3, factor: 0.3 },
          { interval: 0.2, factor: 0.3 }
        ]
      }
    },
    { card: "bar", stats: {} }
  ];
  statsList.deck = "Test";

  const el = statsList.element;
  expect(el.querySelector("#stats-deck").innerHTML).toEqual("Test");
  expect(el.querySelector("#stats-list").children.length).toEqual(2);
  expect(el.querySelector("#stats-card").innerHTML).toEqual("foo");
  expect(el.querySelector("#reviewed-at").innerHTML).toEqual(
    new Date(0).toLocaleString()
  );
  expect(el.querySelector("#interval").innerHTML).toEqual("5h");
  expect(el.querySelector("#difficulty").innerHTML).toEqual("1.3");
});
