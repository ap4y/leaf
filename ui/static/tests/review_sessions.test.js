import ReviewSession from "../review_session.js";

describe("auto rater", () => {
  const session = {
    total: 10,
    left: 5,
    question: "foo"
  };

  test("render", () => {
    const reviewSession = new ReviewSession();
    expect(reviewSession.element).not.toBeNull();

    reviewSession.deck = "Test";
    reviewSession.session = session;

    const el = reviewSession.element;
    expect(el.querySelector("#deck").innerHTML).toEqual("Test");
    expect(el.querySelector("#progress").innerHTML).toEqual("5/10");
    expect(el.querySelector("#answer-state").innerHTML).toEqual("&nbsp;");
    expect(el.querySelector("#correct-answer").innerHTML).toEqual("&nbsp;");
  });

  test("submit incorrect", async () => {
    const reviewSession = new ReviewSession();
    reviewSession.session = session;
    reviewSession.resolveAnswer = () => ({ answer: "bar" });

    let rating = null;
    reviewSession.advanceSession = r => {
      rating = r;
    };

    const el = reviewSession.element;
    el.querySelector("#input-form").onsubmit({ preventDefault: () => {} });
    await new Promise(resolve => window.setTimeout(resolve, 100));
    expect(el.querySelector("#answer-state").innerHTML).toEqual("✕");
    expect(el.querySelector("#correct-answer").innerHTML).toEqual(
      '<span class="input-correct">bar</span>'
    );

    el.querySelector("#input-form").onsubmit({ preventDefault: () => {} });
    await new Promise(resolve => window.setTimeout(resolve, 100));
    expect(rating).toEqual(0);
  });

  test("mistake diff", async () => {
    const reviewSession = new ReviewSession();
    reviewSession.session = session;
    reviewSession.resolveAnswer = () => ({ answer: "にほんごのかくせい" });

    let rating = null;
    reviewSession.advanceSession = r => {
      rating = r;
    };

    const el = reviewSession.element;
    el.querySelector("#input").value = "にほんごこくせいい";
    el.querySelector("#input-form").onsubmit({ preventDefault: () => {} });
    await new Promise(resolve => window.setTimeout(resolve, 100));
    expect(el.querySelector("#correct-answer").innerHTML).toEqual(
      'にほんご<span class="input-correct">の</span>' +
        '<span class="input-mistake">こ</span><span class="input-correct">か</span>' +
        'くせい<span class="input-mistake">い</span>'
    );
  });

  test("submit correct", async () => {
    const reviewSession = new ReviewSession();
    reviewSession.session = session;
    reviewSession.resolveAnswer = () => ({ answer: "bar" });

    let rating = null;
    reviewSession.advanceSession = r => {
      rating = r;
    };

    const el = reviewSession.element;
    el.querySelector("#input").value = "bar";
    el.querySelector("#input-form").onsubmit({ preventDefault: () => {} });
    await new Promise(resolve => window.setTimeout(resolve, 100));
    expect(el.querySelector("#answer-state").innerHTML).toEqual("✓");
    expect(el.querySelector("#correct-answer").innerHTML).toEqual("&nbsp;");

    el.querySelector("#input-form").onsubmit({ preventDefault: () => {} });
    await new Promise(resolve => window.setTimeout(resolve, 100));
    expect(rating).toEqual(1);
  });

  test("submit correct with unicode separator", async () => {
    const reviewSession = new ReviewSession();
    reviewSession.session = session;
    reviewSession.resolveAnswer = () => ({ answer: "いち に" });

    let rating = null;
    reviewSession.advanceSession = r => {
      rating = r;
    };

    const el = reviewSession.element;
    el.querySelector("#input").value = "いち　に";
    el.querySelector("#input-form").onsubmit({ preventDefault: () => {} });
    await new Promise(resolve => window.setTimeout(resolve, 100));
    expect(el.querySelector("#answer-state").innerHTML).toEqual("✓");
    expect(el.querySelector("#correct-answer").innerHTML).toEqual("&nbsp;");

    el.querySelector("#input-form").onsubmit({ preventDefault: () => {} });
    await new Promise(resolve => window.setTimeout(resolve, 100));
    expect(rating).toEqual(1);
  });
});

describe("self rater", () => {
  const session = {
    total: 10,
    left: 5,
    question: "foo",
    rating_type: "self"
  };

  test("render", () => {
    const reviewSession = new ReviewSession();
    expect(reviewSession.element).not.toBeNull();

    reviewSession.deck = "Test";
    reviewSession.session = session;

    const el = reviewSession.element;
    expect(el.querySelector("#deck").innerHTML).toEqual("Test");
    expect(el.querySelector("#progress").innerHTML).toEqual("5/10");
    expect(el.querySelector("#self-answer").innerHTML).toEqual("&nbsp;");
    expect(el.querySelector("#advance").style.visibility).toEqual("visible");
    expect(el.querySelector("#rating").style.visibility).toEqual("hidden");
  });

  test("show answer", async () => {
    const reviewSession = new ReviewSession();
    reviewSession.session = session;
    reviewSession.resolveAnswer = () => ({ answer: "bar" });

    const el = reviewSession.element;
    el.querySelector("#review-form").onsubmit({ preventDefault: () => {} });
    await new Promise(resolve => window.setTimeout(resolve, 100));
    expect(el.querySelector("#self-answer").innerHTML).toEqual("bar");
    expect(el.querySelector("#advance").style.visibility).toEqual("hidden");
    expect(el.querySelector("#rating").style.visibility).toEqual("visible");
  });

  test("submit review", async () => {
    const reviewSession = new ReviewSession();
    reviewSession.session = session;
    reviewSession.resolveAnswer = () => ({ answer: "bar" });
    let rating = null;
    reviewSession.advanceSession = r => {
      rating = r;
    };

    const el = reviewSession.element;
    el.querySelector("#review-form").onsubmit({ preventDefault: () => {} });
    el.querySelector("#rating button").click({ preventDefault: () => {} });
    await new Promise(resolve => window.setTimeout(resolve, 100));
    expect(rating).toEqual(0);
  });
});
