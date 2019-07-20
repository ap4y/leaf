const template = `
<header>
  <h3>Deck: <span id="deck"></span></h3>
  <h5>Progress: <span id="progress"></span></h5>
</header>

<main class="session container">
  <h1 id="question" class="title"></h1>

  <form id="inputForm">
    <input id="input" autofocus autocomplete="off"/>
    <input type="submit" value="⏎" />
  </form>

  <p id="result" class="result">
    <span id="answerState">&nbsp</span>
    <span id="correctAnswer">&nbsp</span>
  </p>
</main>
`;

export default class ReviewSession {
  constructor(session) {
    this._el = document.createElement("div");
    this._el.innerHTML = template;
    this._el.querySelector("#inputForm").onsubmit = e => {
      e.preventDefault();
      if (this.isAnswering) {
        this._resolveAnswer();
      } else {
        this._getNextQuestion();
      }
    };
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

  set submitAnswer(callback) {
    this._submitAnswer = callback;
  }

  set getNextQuestion(callback) {
    this._getNextQuestion = callback;
  }

  _render() {
    this.isAnswering = true;
    this._updateState();
  }

  _updateState() {
    const { question, total, left, answerLen } = this._session;

    if (left === 0) {
      window.history.back();
      return;
    }

    this._el.querySelector("#progress").innerHTML = `${total - left}/${total}`;
    this._el.querySelector("#question").innerHTML = question;
    this._el.querySelector("#answerState").innerHTML = "&nbsp";
    this._el.querySelector("#correctAnswer").innerHTML = "&nbsp";
    this._el.querySelector("#input").style.width = `${answerLen}ch`;
    this._el.querySelector("#input").value = "";
    this._el.querySelector("#input").focus();
  }

  async _resolveAnswer() {
    const answer = this._el.querySelector("#input").value;
    const answerState = this._el.querySelector("#answerState");
    const correctAnswer = this._el.querySelector("#correctAnswer");

    this.isAnswering = false;
    const result = await this._submitAnswer(answer);
    if (result.is_correct) {
      answerState.innerHTML = "✓";
      answerState.style.color = "green";
      correctAnswer.innerHTML = "&nbsp";
    } else {
      answerState.innerHTML = "✕";
      answerState.style.color = "red";
      correctAnswer.innerHTML = result.correct;
    }
  }
}
