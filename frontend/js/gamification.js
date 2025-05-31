document.addEventListener("DOMContentLoaded", () => {
    fetchLeaderboard();
    trackUserActions();
    checkDailyLogin();
});

let searchCount = 0;
let donationCount = 0;

async function fetchLeaderboard() {
    const response = await fetch('/gamification/leaderboard');
    const data = await response.json();
    const leaderboardTable = document.getElementById("leaderboard");
    leaderboardTable.innerHTML = "";

    data.forEach((user, index) => {
        const row = document.createElement("tr");
        let rankClass = "";
        if (index === 0) rankClass = "rank-1";
        if (index === 1) rankClass = "rank-2";
        if (index === 2) rankClass = "rank-3";
        
        row.innerHTML = `
            <td class="${rankClass}">${index + 1}</td>
            <td>${user.name}</td>
            <td>${user.achievements || "None"}</td>
            <td>${user.points}</td>
        `;
        leaderboardTable.appendChild(row);
    });
}

function trackUserActions() {
    document.getElementById("background-coins")?.addEventListener("click", () => {
        updatePoints("coin_click", "+1 Coin!");
    });

    document.getElementById("searchForm")?.addEventListener("submit", (event) => {
        event.preventDefault();
        searchCount++;
        if (searchCount === 5) {
            updatePoints("searcher", "🏆 Searcher Unlocked!");
            searchCount = 0;
        }
    });

    document.getElementById("donationForm")?.addEventListener("submit", (event) => {
        event.preventDefault();
        donationCount++;
        if (donationCount === 5) {
            updatePoints("active_donator", "🏆 Active Donator Unlocked!");
            donationCount = 0;
        }
    });
}

async function checkDailyLogin() {
    const lastLogin = localStorage.getItem("lastLogin");
    const today = new Date().toISOString().split("T")[0];
    if (lastLogin !== today) {
        await updatePoints("daily_login", "🏆 Daily Login Bonus!");
        localStorage.setItem("lastLogin", today);
    }
}

async function updatePoints(achievement, message = "") {
    const token = localStorage.getItem("token"); 
    if (!token) {
        console.error("No token found!");
        return;
    }

    const response = await fetch(`/gamification/update?achievement=${achievement}`, {
        method: "POST",
        headers: { "Authorization": `Bearer ${token}` }
    });

    const data = await response.json();
    if (data.message.includes("Achievement unlocked")) {
        showNotification(`🏆 ${data.achievement} Unlocked!`);
    } else if (message) {
        showNotification(message);
    }
}


function showNotification(message) {
    const notification = document.createElement("div");
    notification.className = "achievement-notification";
    notification.textContent = message;
    document.body.appendChild(notification);
    setTimeout(() => { notification.remove(); }, 3000);
}
