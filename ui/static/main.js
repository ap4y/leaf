import DeckList from "./deck_list.js";
import ReviewSession from "./review_session.js";
import StatsList from "./stats_list.js";

class App {
  render() {
    window.onpopstate = () => {
      this.showDecks();
    };

    this.showDecks();
    const hash = window.location.hash;
    window.history.replaceState({}, "DeckApp", "/");
    if (hash.startsWith("#stats")) {
      this.showStats(hash.replace("#stats-", ""));
      return;
    }
    if (hash) this.startSession(hash.replace("#", ""));
  }

  showDecks() {
    document.getElementById("decks").style.display = null;
    document.getElementById("session").style.display = "none";
    document.getElementById("stats").style.display = "none";
    const deckList = new DeckList();
    deckList.render();
  }

  startSession(deck, cardsReady) {
    if (cardsReady === 0) return;

    window.history.pushState({ deck }, `Review: ${deck}`, `#${deck}`);
    document.getElementById("decks").style.display = "none";
    document.getElementById("session").style.display = null;
    const session = new ReviewSession();
    session.render();
  }

  showStats(deck) {
    window.history.pushState({ deck }, `Stats: ${deck}`, `#stats-${deck}`);
    document.getElementById("decks").style.display = "none";
    document.getElementById("stats").style.display = null;
    const stats = new StatsList();
    stats.render();
  }
}

window.app = new App();
window.app.render();
