const urlParams = new URLSearchParams(window.location.search);
const campaignId = urlParams.get("id");
document.addEventListener("DOMContentLoaded", () => {
    loadCampaignDetails();
    document.getElementById("toggle-donations-btn").addEventListener("click", toggleDonations);
});
async function loadCampaignDetails() {
    try {
        const response = await fetch(`/campaigns/${campaignId}`); // Adjust API URL
        const data = await response.json();
        const campaign = data.campaign;
        const creator = data.creator;

        document.getElementById("campaign-title").innerText = campaign.title;
        document.getElementById("campaign-description").innerText = campaign.description;
        document.getElementById("campaign-goal").innerText = campaign.target_amount;
        document.getElementById("campaign-raised").innerText = campaign.amount_raised;

        const mediaContainer = document.getElementById("campaign-media");
        mediaContainer.innerHTML = `<img src="${campaign.media_path}" alt="Campaign Image" class="campaign-img">`;
        const goal = Number(campaign.target_amount) || 1;
        const raised = Number(campaign.amount_raised) || 0;
        const progressPercentage = Math.min((raised / goal) * 100, 100);
        document.getElementById("progress-bar").style.width = progressPercentage + "%";
        document.getElementById("progress-percent").innerText = progressPercentage.toFixed(1) + "%";

        const detailsContainer = document.querySelector(".campaign-details");
        detailsContainer.innerHTML = `
            <h3>Contact Information</h3>
            <p>
                <strong>Creator:</strong> ${creator.name}<br>
                <strong>Email:</strong> <a href="mailto:${creator.email}">${creator.email}</a>
            </p>
      
        `;


        document.getElementById("donate-btn").addEventListener("click", () => {
            window.location.href = `/static/donation.html?id=${campaignId}`;
        });

    } catch (error) {
        console.error("Error loading campaign:", error);
    }
}

async function loadDonations() {
    try {
        const token = localStorage.getItem('token');

        const options = {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}` // Adjust based on your backend requirements
            }
        };
        const donationResponse = await fetch(`/campaigns/${campaignId}/donations`, options);
        const donations = await donationResponse.json();

        const donationList = document.getElementById("donation-list");
        donationList.innerHTML = '';

        if (donations.length === 0) {
            donationList.innerHTML = '<li class="donation-item">No donations yet.</li>';
        } else {
            donations.forEach(item  => {
                const donation = item.donation;
                const user = item.user;
                const listItem = document.createElement("li");
                listItem.classList.add("donation-item");
                listItem.innerHTML = `
          <div>
            <strong>${user.name}</strong> (<a href="mailto:${user.email}">${user.email}</a>) donated
          </div>
          <div>
            $${Number(donation.amount).toFixed(2)} <span>on ${new Date(donation.donation_date).toLocaleDateString()}</span>
          </div>
        `;
                donationList.appendChild(listItem);
            });
        }
    } catch (error) {
        console.error("Error loading donations:", error);
    }
}

function toggleDonations() {
    const donationsContainer = document.getElementById("donations-container");
    if (donationsContainer.style.display === "none" || donationsContainer.style.display === "") {
        loadDonations();
        donationsContainer.style.display = "block";
        document.getElementById("toggle-donations-btn").innerText = "Hide Donations";
    } else {
        donationsContainer.style.display = "none";
        document.getElementById("toggle-donations-btn").innerText = "View Donations";
    }
}

