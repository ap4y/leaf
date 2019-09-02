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
        Historical: [{ interval: 0.3 }, { interval: 0.2 }]
      }
    },
    { card: "bar", stats: {} }
  ];
  statsList.deck = "Test";

  const el = statsList.element;
  expect(el.querySelector("#statsDeck").innerHTML).toEqual("Test");
  expect(el.querySelector("#statsList").children.length).toEqual(2);
  expect(el.querySelector("#statsCard").innerHTML).toEqual("foo");
  expect(el.querySelector("#reviewedAt").innerHTML).toEqual(
    "1/1/1970, 12:00:00 PM"
  );
  expect(el.querySelector("#interval").innerHTML).toEqual("5h");
  expect(el.querySelector("#difficulty").innerHTML).toEqual("1.3");
  expect(el.querySelector("#historical").innerHTML).toEqual("0.3, 0.2");
});
