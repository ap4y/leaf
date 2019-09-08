import { AutoRater, SelfRater } from "./rater.js";

const template = `
<header>
  <h3>Deck: <span id="deck"></span></h3>
  <h5>Progress: <span id="progress"></span></h5>
</header>

<main id="session" class="session container">
  <h1 id="question" class="title"></h1>
</main>
`;

export default class ReviewSession {
  constructor() {
    this._el = document.createElement("div");
    this._el.innerHTML = template;

    this.selfRater = new SelfRater();
    this.selfRater.onSubmit = answer =>
      this._handleRater(this.selfRater, answer);

    this.autoRater = new AutoRater();
    this.autoRater.onSubmit = answer =>
      this._handleRater(this.autoRater, answer);
  }

  get element() {
    return this._el;
  }

  set deck(deck) {
    this._el.querySelector("#deck").innerHTML = deck;
  }

  set session(session) {
    this._session = session;
    this._render();
  }

  set resolveAnswer(callback) {
    this._resolveAnswer = callback;
  }

  set advanceSession(callback) {
    this._advanceSession = callback;
  }

  _render() {
    this.isAnswering = true;
    this._updateState();
  }

  _updateState() {
    const { question, total, left, rating_type } = this._session;

    if (left === 0) {
      window.history.back();
      return;
    }

    this._el.querySelector("#progress").innerHTML = `${total - left}/${total}`;
    this._el.querySelector("#question").innerHTML = question;

    const rater = rating_type === "self" ? this.selfRater : this.autoRater;
    const session = this._el.querySelector("#session");
    if (session.children.length === 1) {
      session.appendChild(rater.element);
    } else {
      session.replaceChild(rater.element, session.lastChild);
    }
    rater.showQuestion(this._session);
  }

  async _handleRater(rater, score) {
    if (this.isAnswering) {
      this.isAnswering = false;
      const { answer } = await this._resolveAnswer();
      rater.showResult(answer);
    } else {
      if (typeof score !== "number") return;
      this._advanceSession(score);
    }
  }
}
