import DeckList from "./deck_list.js";
import ReviewSession from "./review_session.js";
import StatsList from "./stats_list.js";

class App {
  constructor() {
    this.deckList = new DeckList();
    this._decks = document.getElementById("decks");
    this._decks.appendChild(this.deckList.element);

    this.statsList = new StatsList();
    this._stats = document.getElementById("stats");
    this._stats.appendChild(this.statsList.element);

    this.reviewSession = new ReviewSession();
    this.reviewSession.resolveAnswer = () => this._resolveAnswer();
    this.reviewSession.advanceSession = async score => {
      const session = await this._advanceSession(score);
      this.reviewSession.session = session;
    };
    this._session = document.getElementById("session");
    this._session.appendChild(this.reviewSession.element);
  }

  render() {
    window.onpopstate = () => {
      this.showDecks();
    };

    const hash = window.location.hash;
    window.history.replaceState({}, "DeckApp", "/");
    if (!hash) {
      this.showDecks();
      return;
    }

    if (hash.startsWith("#stats")) {
      this.showStats(hash.replace("#stats-", ""));
    } else {
      this.startSession(hash.replace("#", ""));
    }
  }

  async showDecks() {
    this._decks.style.display = null;
    this._session.style.display = "none";
    this._stats.style.display = "none";

    this.deckList.decks = await this._fetchDecks();
  }

  async startSession(deck, cardsReady) {
    if (cardsReady === 0) return;

    window.history.pushState({ deck }, `Review: ${deck}`, `#${deck}`);
    this._decks.style.display = "none";
    this._session.style.display = null;
    this._stats.style.display = "none";

    this.reviewSession.deck = deck;
    this.reviewSession.session = await this._startSession(deck);
  }

  async showStats(deck) {
    window.history.pushState({ deck }, `Stats: ${deck}`, `#stats-${deck}`);

    this._decks.style.display = "none";
    this._session.style.display = "none";
    this._stats.style.display = null;

    this.statsList.deck = deck;
    this.statsList.stats = await this._fetchStats(deck);
  }

  async _request(path, options = {}) {
    const res = await window.fetch(path, options);
    if (res.ok) return await res.json();

    alert(`Failed to fetch: ${await res.text()}`);
    return null;
  }

  _fetchDecks() {
    return this._request("decks");
  }

  _fetchStats(deck) {
    return this._request(`stats/${deck}`);
  }

  _startSession(deck) {
    return this._request(`start/${deck}`, {
      method: "POST"
    });
  }

  _resolveAnswer() {
    return this._request("resolve");
  }

  _advanceSession(score) {
    return this._request("advance", {
      method: "POST",
      body: JSON.stringify({ score })
    });
  }
}

window.app = new App();
window.app.render();
