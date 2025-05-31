document.addEventListener("DOMContentLoaded", async () => {
    const donationForm = document.getElementById("donation-form");

    let stripe = Stripe("pk_test_51QqqjJFKXCRtVeXQtrqiOl61kWJBfdVBhHtJTEcB3JFjI6q1oil2O9vRQwduVIljGTJmD9N6sxOr1D4QrtgKfnIK00elxI0Oc5");
    let elements = stripe.elements();
    let cardElement = elements.create("card");
    cardElement.mount("#card-element"); // Mount Stripe Card UI

    if (donationForm) {
        donationForm.addEventListener("submit", async (event) => {
            event.preventDefault();
            const { paymentMethod, error } = await stripe.createPaymentMethod({
                type: "card",
                card: cardElement
            });
            // Get Campaign ID from URL
            const urlParams = new URLSearchParams(window.location.search);
            const campaignId = urlParams.get("id");
            if (!campaignId) {
                alert("Campaign ID is missing in the URL");
                return;
            }

            // Get donation amount (convert to cents for Stripe)
            const amount = parseFloat(document.getElementById("amount").value) * 100;
            if (!amount || amount <= 0) {
                alert("Invalid donation amount.");
                return;
            }

            // Get User Token
            const token = localStorage.getItem("token");
            if (!token) {
                alert("You must be logged in to donate.");
                window.location.href = "login.html";
                return;
            }

            const decoded = jwt_decode(token);
            const userID = decoded.UserID;
            console.log("User ID:", userID);

            try {

                const paymentIntentResponse = await fetch("http://localhost:8080/create-payment-intent", {
                    method: "POST",
                    headers: { "Content-Type": "application/json" },
                    body: JSON.stringify({   paymentMethodId: paymentMethod.id,
                        amount: parseFloat(document.getElementById("amount").value) * 100 })
                });

                const paymentIntentData = await paymentIntentResponse.json();
                if (!paymentIntentResponse.ok) {
                    throw new Error(paymentIntentData.error || "Failed to create payment intent");
                }

                const clientSecret = paymentIntentData.client_secret;

                // Step 2: Confirm Payment using Stripe
                const { paymentIntent, error } = await stripe.confirmCardPayment(clientSecret, {
                    payment_method: { card: cardElement }
                });

                if (error) {
                    console.error("Payment failed:", error);
                    alert("Payment failed. Try again.");
                    return;
                }

                console.log("Payment successful:", paymentIntent.id);

                // Step 3: Send Payment ID to backend to record the donation
                const donationResponse = await fetch(`http://localhost:8080/campaigns/${campaignId}/donations/`, {
                    method: "POST",
                    headers: {
                        "Content-Type": "application/json",
                        "Authorization": `Bearer ${token}`
                    },
                    body: JSON.stringify({
                        amount: amount / 100,
                        stripe_payment_id: paymentIntent.id,
                        user_id: userID
                    })
                });

                if (donationResponse.ok) {
                    alert("Donation successful!");
                    window.location.href = "/static/myDonations.html";
                } else {
                    alert("Donation failed.");
                }

            } catch (err) {
                console.error(err);
                alert("An error occurred while processing your donation.");
            }
        });
    }
});
