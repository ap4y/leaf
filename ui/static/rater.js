const autoRated = `
<form id="input-form" class="input-form">
  <input id="input" autofocus autocomplete="off"/>
  <input type="submit" class="submit-button" value="⏎" />
</form>

<p id="result" class="result">
  <span id="answer-state">&nbsp</span>
  <span id="correct-answer">&nbsp</span>
</p>
`;

export class AutoRater {
  constructor() {
    this.score = 0;
    this._el = document.createElement("div");
    this._el.innerHTML = autoRated;
    this._el.querySelector("#input-form").onsubmit = e => {
      e.preventDefault();
      this._onSubmit(this.score);
    };
  }

  get element() {
    return this._el;
  }

  set onSubmit(callback) {
    this._onSubmit = callback;
  }

  showQuestion({ answer_length }) {
    this._el.querySelector("#answer-state").innerHTML = "&nbsp";
    this._el.querySelector("#correct-answer").innerHTML = "&nbsp";

    const input = this._el.querySelector("#input");
    input.style.width = `${2 * answer_length}ch`;
    input.value = "";
    input.focus();
  }

  showResult(answer) {
    const userInput = this._el.querySelector("#input").value;
    const answerState = this._el.querySelector("#answer-state");
    const correctAnswer = this._el.querySelector("#correct-answer");

    const pattern = answer.split(/\s/).join("\\s");
    if (userInput.match(new RegExp(pattern))) {
      answerState.innerHTML = "✓";
      answerState.style.color = "green";
      correctAnswer.innerHTML = "&nbsp";
      this.score = 1;
    } else {
      answerState.innerHTML = "✕";
      answerState.style.color = "red";
      correctAnswer.innerHTML = answer;
      this.score = 0;
    }
  }
}

const selfRated = `
<p id="result" class="result">
  <span id="self-answer">&nbsp</span>
</p>

<form id="review-form" class="rating-form">
  <input id="advance" type="submit" value="Show Answer" accesskey="space"/>

  <div id="rating">
    <button value="0" accesskey="1">Again</button>
    <button value="1" accesskey="2">Hard</button>
    <button value="2" accesskey="3">Good</button>
    <button value="3" accesskey="4">Easy</button>
  </div>
</form>
`;

export class SelfRater {
  constructor() {
    this._el = document.createElement("div");
    this._el.innerHTML = selfRated;

    document.addEventListener("keydown", e => {
      if (!this._el.parentNode) return null;

      switch (e.code) {
        case "Space":
          return this._onSubmit();
        case "Digit1":
          return this._onSubmit(0);
        case "Digit2":
          return this._onSubmit(1);
        case "Digit3":
          return this._onSubmit(2);
        case "Digit4":
          return this._onSubmit(3);
      }
      return null;
    });
    this._el.querySelector("#review-form").onsubmit = e => {
      e.preventDefault();
      this._onSubmit();
    };
    this._el.querySelectorAll("button").forEach(input =>
      input.addEventListener("click", e => {
        e.preventDefault();
        this._onSubmit(Number.parseInt(e.target.value));
      })
    );
  }

  get element() {
    return this._el;
  }

  set onSubmit(callback) {
    this._onSubmit = callback;
  }

  showQuestion() {
    this._el.querySelector("#self-answer").innerHTML = "&nbsp";
    this._el.querySelector("#rating").style.visibility = "hidden";
    this._el.querySelector("#advance").style.visibility = "visible";
  }

  showResult(answer) {
    const correctAnswer = this._el.querySelector("#self-answer");
    correctAnswer.innerHTML = answer;

    this._el.querySelector("#rating").style.visibility = "visible";
    this._el.querySelector("#advance").style.visibility = "hidden";
  }
}
