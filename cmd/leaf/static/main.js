class App {
  constructor() {
    this.isAnswering = true;
  }

  render() {
    document.getElementById("inputForm").onsubmit = e => {
      e.preventDefault();
      if (this.isAnswering) {
        this._resolveAnswer();
      } else {
        this._nextQuestion();
      }
    };

    this._nextQuestion();
  }

  async _nextQuestion() {
    const state = await this._fetchNext();
    this.isAnswering = true;
    document.getElementById("input").value = "";
    answerState.innerHTML = "&nbsp";
    correctAnswer.innerHTML = "&nbsp";

    document.getElementById("deck").innerHTML = state.deck;
    document.getElementById("progress").innerHTML = `${state.total -
      state.left}/${state.total}`;

    if (state.left === 0) {
      document.getElementById("session").style.display = "none";
      document.getElementById("finished").style.display = "block";
      return;
    } else {
      document.getElementById("session").style.display = "flex";
      document.getElementById("finished").style.display = "none";
    }

    document.getElementById("question").innerHTML = state.question;
    document.getElementById("input").style.width = `${state.answerLen}ch`;
  }

  async _resolveAnswer() {
    const answer = document.getElementById("input").value;
    const answerState = document.getElementById("answerState");
    const correctAnswer = document.getElementById("correctAnswer");

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

  async _fetchNext() {
    return fetch("next").then(res => {
      if (res.ok) return res.json();

      res.text().then(text => alert(`Failed to fetch next question: ${text}`));
      return nil;
    });
  }

  async _submitAnswer(answer) {
    return fetch("resolve", {
      method: "POST",
      body: JSON.stringify({ answer })
    }).then(res => {
      if (res.ok) return res.json();

      res.text().then(text => alert(`Failed to fetch next question: ${text}`));
      return nil;
    });
  }
}

const app = new App();
app.render();
