document.addEventListener("DOMContentLoaded", () => {
    const registrationForm = document.getElementById("registration-form");

    if (registrationForm) {
        registrationForm.addEventListener("submit", async (event) => {
            event.preventDefault();

            const user = {
                name: document.getElementById("name").value,
                email: document.getElementById("email").value,
                password: document.getElementById("password").value,
                role: document.getElementById("role").value
            };

            try {
                const response = await fetch("/register", {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                    },
                    body: JSON.stringify(user)
                });

                const data = await response.json();
                const messageDiv = document.getElementById("message");

                if (response.ok) {
                    messageDiv.innerHTML = `<div class="alert alert-success">Registration successful! Redirecting...</div>`;

                    // Если сервер возвращает токен, сразу логиним пользователя
                    if (data.token) {
                        localStorage.setItem("token", data.token);
                        window.location.href = "index.html";
                    } else {
                        setTimeout(() => {
                            window.location.href = "login.html";
                        }, 2000);
                    }
                } else {
                    messageDiv.innerHTML = `<div class="alert alert-danger">${data.error || "Registration failed"}</div>`;
                }
            } catch (error) {
                console.error("Registration error:", error);
                document.getElementById("message").innerHTML = `<div class="alert alert-danger">Error during registration. Please try again.</div>`;
            }
        });
    }
});
