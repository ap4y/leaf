import DeckList from "../deck_list.js";

test("render", () => {
  const deckList = new DeckList();
  expect(deckList.element).not.toBeNull();

  deckList.decks = [
    { name: "foo", cards_ready: 10, next_review_at: new Date().toString() },
    { name: "bar", cards_ready: 10, next_review_at: new Date().toString() }
  ];

  const el = deckList.element;
  expect(el.children.length).toEqual(2);

  const child = el.children[0];
  expect(child.querySelector("a").text).toEqual("bar");
  expect(child.querySelector("code").innerHTML).toEqual("10");
});

test("render not ready", () => {
  const deckList = new DeckList();
  deckList.decks = [
    { name: "foo", cards_ready: 0, next_review_at: new Date(0).toString() }
  ];

  const child = deckList.element.children[0];
  expect(child.querySelector("code").innerHTML).toEqual("available now");
});

test("deck click", () => {
  const deckList = new DeckList();
  deckList.decks = [
    { name: "foo", cards_ready: 10, next_review_at: new Date(0) }
  ];

  const el = deckList.element;
  let event = {};
  window.app = {
    startSession: (deck, count) => {
      event = { deck, count };
    }
  };
  el.querySelector("a").click();

  expect(event.deck).toEqual("foo");
  expect(event.count).toEqual(10);
});

test("stats click", () => {
  const deckList = new DeckList();
  deckList.decks = [
    { name: "foo", cards_ready: 10, next_review_at: new Date(0) }
  ];

  const el = deckList.element;
  let event = {};
  window.app = {
    showStats: deck => {
      event = { deck };
    }
  };
  el.querySelector(".stats-link").click();

  expect(event.deck).toEqual("foo");
});
