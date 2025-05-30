document.addEventListener("DOMContentLoaded", () => {
    const loginLink = document.getElementById("loginLink");
    const logoutLink = document.getElementById("logoutLink");

    if (loginLink && logoutLink) {
        if (localStorage.getItem("token")) {
            loginLink.style.display = "none";
            logoutLink.style.display = "inline";
        }

        logoutLink.addEventListener("click", () => {
            localStorage.removeItem("token");
            window.location.href = "index.html";
        });
    }

    const loginForm = document.getElementById("login-form");
    if (loginForm) {
        loginForm.addEventListener("submit", async (event) => {
            event.preventDefault();
            const email = document.getElementById("email").value;
            const password = document.getElementById("password").value;

            try {
                const response = await fetch("http://localhost:8080/login", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({ email, password })
                });

                const data = await response.json();
                if (response.ok) {
                    localStorage.setItem("token", data.token);
                    window.location.href = "index.html";
                } else {
                    alert("Login failed: " + (data.error || "Unknown error"));
                }
            } catch (error) {
                console.error("Login error:", error);
                alert("Login failed. Please try again.");
            }
        });
    }
});
