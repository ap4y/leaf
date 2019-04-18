class DeckList {
  async render() {
    const decks = await this._fetchDecks();
    document.getElementById("deckList").innerHTML = decks
      .sort((a, b) => a.name > b.name)
      .map(
        ({ name, cards_ready, next_review_at }) =>
          `<li>
    <a href="#${name}" onclick="app.startSession('${name}'); return false;">${name}</a>
    <code>${this._reviewStats(cards_ready, new Date(next_review_at))}</code>
</li>`
      )
      .join("");
  }

  async _fetchDecks() {
    return window.fetch("decks").then(res => {
      if (res.ok) return res.json();

      return res.text().then(text => alert(`Failed to fetch decks: ${text}`));
    });
  }

  _reviewStats(ready, next) {
    if (ready > 0) return ready;

    const diff = next - new Date();
    if (diff < 1000) return "available now";

    let d = next;
    d = [
      "0" + d.getDate(),
      "0" + (d.getMonth() + 1),
      "" + d.getFullYear(),
      "0" + d.getHours(),
      "0" + d.getMinutes()
    ].map(component => component.slice(-2));

    return d.slice(0, 3).join(".") + " " + d.slice(3).join(":");
  }
}

class ReviewSession {
  constructor() {
    this.isAnswering = true;
    this.deck = null;
    this.session = null;
  }

  async render() {
    document.getElementById("inputForm").onsubmit = e => {
      e.preventDefault();
      if (this.isAnswering) {
        this._resolveAnswer();
      } else {
        this._nextQuestion();
      }
    };

    this.deck = window.history.state.deck;
    document.getElementById("deck").innerHTML = this.deck;
    this.session = await this._startSession(this.deck);
    this._updateState();
  }

  _updateState() {
    const { question, total, left, answerLen } = this.session;

    if (left === 0) {
      window.history.back();
      return;
    }

    document.getElementById("progress").innerHTML = `${total - left}/${total}`;
    document.getElementById("question").innerHTML = question;
    document.getElementById("input").style.width = `${answerLen}ch`;
    document.getElementById("answerState").innerHTML = "&nbsp";
    document.getElementById("correctAnswer").innerHTML = "&nbsp";
    document.getElementById("input").value = "";
    document.getElementById("input").focus();
  }

  async _nextQuestion() {
    this.session = await this._fetchNext();
    this.isAnswering = true;
    this._updateState();
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

  async _startSession(deck) {
    return window
      .fetch(`start/${deck}`, {
        method: "POST"
      })
      .then(res => {
        if (res.ok) return res.json();

        return res.text().then(text => alert(`Failed to start new session: ${text}`));
      });
  }

  async _fetchNext() {
    return window.fetch("next").then(res => {
      if (res.ok) return res.json();

      return res.text().then(text => alert(`Failed to fetch next question: ${text}`));
    });
  }

  async _submitAnswer(answer) {
    return window
      .fetch("resolve", {
        method: "POST",
        body: JSON.stringify({ answer })
      })
      .then(res => {
        if (res.ok) return res.json();

        return res.text().then(text => alert(`Failed to submit answer: ${text}`));
      });
  }
}

class App {
  render() {
    window.onpopstate = () => {
      this.showDecks();
    };

    this.showDecks();
    const hash = window.location.hash;
    window.history.replaceState({}, "DeckApp", "/");
    if (hash) this.startSession(hash.replace("#", ""));
  }

  showDecks() {
    document.getElementById("decks").style.display = null;
    document.getElementById("session").style.display = "none";
    const deckList = new DeckList();
    deckList.render();
  }

  startSession(deck) {
    window.history.pushState({ deck }, `Review: ${deck}`, `#${deck}`);
    document.getElementById("decks").style.display = "none";
    document.getElementById("session").style.display = null;
    const session = new ReviewSession();
    session.render();
  }
}

const app = new App();
app.render();
