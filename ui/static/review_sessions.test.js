import ReviewSession from "./review_session.js";

test("render", () => {
  const reviewSession = new ReviewSession();
  expect(reviewSession.element).not.toBeNull();

  reviewSession.deck = "Test";
  reviewSession.session = {
    total: 10,
    left: 5,
    question: "foo"
  };

  const el = reviewSession.element;
  expect(el.querySelector("#deck").innerHTML).toEqual("Test");
  expect(el.querySelector("#progress").innerHTML).toEqual("5/10");
  expect(el.querySelector("#answerState").innerHTML).toEqual("&nbsp;");
  expect(el.querySelector("#correctAnswer").innerHTML).toEqual("&nbsp;");
});

test("submit incorrect", async () => {
  const reviewSession = new ReviewSession();
  reviewSession.session = {
    total: 10,
    left: 5,
    question: "foo"
  };
  reviewSession.submitAnswer = answer => ({
    is_correct: false,
    correct: "bar"
  });

  const el = reviewSession.element;
  el.querySelector("#inputForm").onsubmit({ preventDefault: () => {} });
  await new Promise(resolve => window.setTimeout(resolve, 100));
  expect(el.querySelector("#answerState").innerHTML).toEqual("✕");
  expect(el.querySelector("#correctAnswer").innerHTML).toEqual("bar");
});

test("submit correct", async () => {
  const reviewSession = new ReviewSession();
  reviewSession.session = {
    total: 10,
    left: 5,
    question: "foo"
  };
  reviewSession.submitAnswer = answer => ({ is_correct: true });

  const el = reviewSession.element;
  el.querySelector("#inputForm").onsubmit({ preventDefault: () => {} });
  await new Promise(resolve => window.setTimeout(resolve, 100));
  expect(el.querySelector("#answerState").innerHTML).toEqual("✓");
  expect(el.querySelector("#correctAnswer").innerHTML).toEqual("&nbsp;");
});

test("request new question", async () => {
  const reviewSession = new ReviewSession();
  reviewSession.session = {
    total: 10,
    left: 5,
    question: "foo"
  };
  reviewSession.submitAnswer = answer => ({ is_correct: true });
  let requested = false;
  reviewSession.getNextQuestion = () => {
    requested = true;
  };

  const el = reviewSession.element;
  el.querySelector("#inputForm").onsubmit({ preventDefault: () => {} });
  el.querySelector("#inputForm").onsubmit({ preventDefault: () => {} });
  await new Promise(resolve => window.setTimeout(resolve, 100));
  expect(requested).toBeTruthy();
});
