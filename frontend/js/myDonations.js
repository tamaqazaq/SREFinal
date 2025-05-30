
document.addEventListener('DOMContentLoaded', fetchDonations);


async function fetchDonations() {
    const token = localStorage.getItem("token");
    const decoded = jwt_decode(token);
    const userId = decoded.UserID;
    try {


        const response = await fetch(`http://localhost:8080/donations/my/${userId}`, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${token}`
            }
        });
        const donations = await response.json();
        const combinedDonations = combineDonations(donations);
        displayDonations(combinedDonations);
    } catch (error) {
        console.error('Error fetching donations:', error);
        document.getElementById('donations-container').innerHTML =
            '<p class="error-message">Error loading donations. Please try again later.</p>';
    }
}


function combineDonations(donations) {

    if (!donations || !Array.isArray(donations)) {
        return [];
    }

    const combined = donations.reduce((acc, item) => {
        // Ensure that the item has a donation property.
        if (!item.donation) {
            console.warn("Skipping item with missing donation:", item);
            return acc;
        }
        const campaignId = item.campaign.campaign_id;
        const donationAmount = item.donation.amount || 0;
        const donationDate = item.donation.donation_date || "";

        if (!acc[campaignId]) {
            acc[campaignId] = {
                campaign: item.campaign,
                totalAmount: donationAmount,
                donationDates: [donationDate]
            };
        } else {
            acc[campaignId].totalAmount += donationAmount;
            acc[campaignId].donationDates.push(donationDate);
        }
        return acc;
    }, {});

    return Object.values(combined);
}

function displayDonations(combinedDonations) {
    const container = document.getElementById('donations-container');
    container.innerHTML = '';

    if (!combinedDonations.length) {
        container.innerHTML = '<p class="no-donations">You haven’t made any donations yet.</p>';
        return;
    }

    combinedDonations.forEach(item => {
        const { campaign, totalAmount, donationDates} = item;
        const card = document.createElement('div');
        card.className = 'donation-card';
        const sortedDates = donationDates.slice().sort((a, b) => new Date(a) - new Date(b));
        const latestDonationDate = sortedDates[sortedDates.length - 1];
        const imgSrc = campaign.media_path ? campaign.media_path : 'placeholder.jpg';
        card.innerHTML = `
      <div class = "campaign-media"> <img src="${imgSrc}" alt="${campaign.title}"> </div>
      <div class="donation-details">
        <h3>${campaign.title}</h3>
        <p>${campaign.description}</p>
        <p class="donation-amount">Total Donated: $${totalAmount.toFixed(2)}</p>
         <p class="donation-date">Last donated on: ${new Date(latestDonationDate).toLocaleDateString()}</p>
      </div>
    `;
        card.addEventListener("click", (e) => {
               window.location.href = `/static/campaign-details.html?id=${campaign.campaign_id}`;
        });

        container.appendChild(card);
    });
}