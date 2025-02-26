document.addEventListener('DOMContentLoaded', () => {
    const newGameBtn = document.getElementById('new-game');
    const placeBetBtn = document.getElementById('place-bet');
    const hitBtn = document.getElementById('hit');
    const standBtn = document.getElementById('stand');
    const resetBtn = document.getElementById('reset');
    const messageEl = document.getElementById('message');
    const dealerCards = document.getElementById('dealer-cards');
    const playerCards = document.getElementById('player-cards');
    const dealerScore = document.getElementById('dealer-score');
    const playerScore = document.getElementById('player-score');
    const chipsEl = document.getElementById('chips');
    const currentBetEl = document.getElementById('current-bet');
    const betButtons = document.querySelectorAll('.bet-btn');
    const allInBtn = document.querySelector('.all-in');

    let currentBet = 0;

    // Event Listeners
    newGameBtn.addEventListener('click', startNewGame);
    placeBetBtn.addEventListener('click', placeBetAction);
    hitBtn.addEventListener('click', hit);
    standBtn.addEventListener('click', stand);
    resetBtn.addEventListener('click', resetGame);

    betButtons.forEach(btn => {
        if (!btn.classList.contains('all-in')) {
            btn.addEventListener('click', () => {
                const amount = parseInt(btn.dataset.amount);
                const currentChips = parseInt(chipsEl.textContent);
                if (currentChips >= amount) {
                    currentBet += amount;
                    updateBetDisplay(currentChips - amount);
                }
            });
        }
    });

    allInBtn.addEventListener('click', () => {
        const currentChips = parseInt(chipsEl.textContent);
        currentBet = currentChips;
        updateBetDisplay(0);
    });

    function updateBetDisplay(remainingChips) {
        currentBetEl.textContent = currentBet;
        chipsEl.textContent = remainingChips;
        
        // Zeige Reset-Button nur wenn Guthaben UND aktueller Einsatz 0 sind
        resetBtn.style.display = remainingChips <= 0 && currentBet <= 0 ? 'block' : 'none';
    }

    function animateCard(cardEl) {
        cardEl.style.opacity = '0';
        cardEl.style.transform = 'translateY(-20px)';
        setTimeout(() => {
            cardEl.style.transition = 'all 0.3s ease-out';
            cardEl.style.opacity = '1';
            cardEl.style.transform = 'translateY(0)';
        }, 50);
    }

    function startNewGame() {
        currentBet = 0;
        updateBetDisplay(parseInt(chipsEl.textContent));
        clearTable();
        fetch('/api/game/new', { method: 'POST' })
            .then(response => response.json())
            .then(updateGameState)
            .catch(handleError);
    }

    function clearTable() {
        dealerCards.innerHTML = '';
        playerCards.innerHTML = '';
        dealerScore.textContent = '';
        playerScore.textContent = '';
        messageEl.textContent = '';
    }

    function placeBetAction() {
        if (currentBet <= 0) {
            messageEl.textContent = 'Bitte setzen Sie einen Betrag!';
            messageEl.classList.add('pulse');
            setTimeout(() => messageEl.classList.remove('pulse'), 500);
            return;
        }

        fetch('/api/game/bet', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ amount: currentBet })
        })
            .then(response => response.json())
            .then(updateGameState)
            .catch(handleError);
    }

    function hit() {
        fetch('/api/game/hit', { method: 'POST' })
            .then(response => response.json())
            .then(updateGameState)
            .catch(handleError);
    }

    function stand() {
        fetch('/api/game/stand', { method: 'POST' })
            .then(response => response.json())
            .then(updateGameState)
            .catch(handleError);
    }

    function resetGame() {
        fetch('/api/game/new', { 
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ reset: true })
        })
            .then(response => response.json())
            .then(updateGameState)
            .catch(handleError);
    }

    function updateGameState(game) {
        // Karten aktualisieren
        dealerCards.innerHTML = '';
        playerCards.innerHTML = '';

        // Dealer-Karten mit versteckter erster Karte
        game.dealerHand.forEach((card, index) => {
            const cardEl = createCardElement(card);
            // Nur die erste Karte verdecken, wenn das Spiel noch läuft
            if (index === 0 && !game.isOver) {
                cardEl.className = 'card hidden';
            } else {
                cardEl.className = `card ${card.suit.toLowerCase()}`;
            }
            dealerCards.appendChild(cardEl);
            animateCard(cardEl);
        });

        game.playerHand.forEach(card => {
            const cardEl = createCardElement(card);
            playerCards.appendChild(cardEl);
            animateCard(cardEl);
        });

        // Punkte aktualisieren - verstecke Dealer-Punkte während des Spiels
        playerScore.textContent = calculateScore(game.playerHand);
        if (game.isOver) {
            dealerScore.textContent = calculateScore(game.dealerHand);
        } else {
            // Zeige nur die Punkte der sichtbaren Karten
            const visibleCards = game.dealerHand.slice(1);
            dealerScore.textContent = visibleCards.length > 0 ? calculateScore(visibleCards) : '?';
        }

        // Chips und Einsatz aktualisieren
        chipsEl.textContent = game.chips;
        currentBetEl.textContent = game.currentBet;
        currentBet = game.currentBet;

        // Reset-Button Anzeige aktualisieren - nur wenn beide Werte 0 sind
        resetBtn.style.display = game.chips <= 0 && game.currentBet <= 0 ? 'block' : 'none';

        // Nachricht aktualisieren
        messageEl.textContent = game.message;
        if (game.isOver) {
            messageEl.classList.add('pulse');
            setTimeout(() => messageEl.classList.remove('pulse'), 500);
        }

        // Buttons aktivieren/deaktivieren
        hitBtn.disabled = game.isOver || game.canBet;
        standBtn.disabled = game.isOver || game.canBet;
        placeBetBtn.disabled = !game.canBet;
        newGameBtn.disabled = !game.isOver;

        // Wett-Buttons aktivieren/deaktivieren
        betButtons.forEach(btn => {
            if (!btn.classList.contains('all-in')) {
                btn.disabled = !game.canBet || parseInt(btn.dataset.amount) > game.chips;
            }
        });
        allInBtn.disabled = !game.canBet || game.chips <= 0;
    }

    function createCardElement(card) {
        const cardEl = document.createElement('div');
        cardEl.className = `card ${card.suit.toLowerCase()}`;
        
        // Oberer Wert
        const valueTop = document.createElement('div');
        valueTop.className = 'value-top';
        valueTop.textContent = card.value;
        
        // Unterer Wert (gespiegelt)
        const valueBottom = document.createElement('div');
        valueBottom.className = 'value-bottom';
        valueBottom.textContent = card.value;
        
        // Kartensymbol in der Mitte
        const suitEl = document.createElement('div');
        suitEl.className = 'suit';
        suitEl.textContent = getSuitSymbol(card.suit);
        
        cardEl.appendChild(valueTop);
        cardEl.appendChild(suitEl);
        cardEl.appendChild(valueBottom);
        
        return cardEl;
    }

    function getSuitSymbol(suit) {
        switch (suit.toLowerCase()) {
            case 'herz': return '♥';
            case 'karo': return '♦';
            case 'pik': return '♠';
            case 'kreuz': return '♣';
            default: return suit;
        }
    }

    function calculateScore(hand) {
        let score = 0;
        let aces = 0;

        hand.forEach(card => {
            if (card.value === 'Ass') {
                aces++;
            } else {
                score += card.points;
            }
        });

        for (let i = 0; i < aces; i++) {
            if (score + 11 <= 21) {
                score += 11;
            } else {
                score += 1;
            }
        }

        return score;
    }

    function handleError(error) {
        console.error('Fehler:', error);
        messageEl.textContent = 'Ein Fehler ist aufgetreten. Bitte versuchen Sie es erneut.';
        messageEl.classList.add('pulse');
        setTimeout(() => messageEl.classList.remove('pulse'), 500);
    }

    function updateResetButtonVisibility() {
        const resetButton = document.getElementById('reset');
        const chips = parseInt(document.getElementById('chips').textContent);
        const currentBet = parseInt(document.getElementById('current-bet').textContent);
        
        if (chips === 0 && currentBet === 0) {
            resetButton.style.display = 'inline-block';
        } else {
            resetButton.style.display = 'none';
        }
    }

    function updateChips(amount) {
        const chipsElement = document.getElementById('chips');
        const currentChips = parseInt(chipsElement.textContent);
        chipsElement.textContent = currentChips + amount;
        updateResetButtonVisibility();
    }

    // Füge den initialen Aufruf hinzu
    updateResetButtonVisibility();

    // Casino effects - add money bills floating in the background
    function createMoneyBills() {
        const container = document.querySelector('body');
        const numBills = 15; // Number of money bills to create
        
        for (let i = 0; i < numBills; i++) {
            const bill = document.createElement('div');
            bill.classList.add('money-bill');
            
            // Random positioning and timing
            const leftPos = Math.random() * 100;
            const delay = Math.random() * 10;
            const duration = 10 + Math.random() * 20;
            const rotation = Math.random() * 360;
            const scale = 0.5 + Math.random() * 1;
            
            bill.style.left = `${leftPos}%`;
            bill.style.animationDelay = `${delay}s`;
            bill.style.animationDuration = `${duration}s`;
            bill.style.transform = `rotate(${rotation}deg) scale(${scale})`;
            
            container.appendChild(bill);
        }
    }

    // Call this when the page loads
    createMoneyBills();
}); 