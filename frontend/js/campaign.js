let campaignList;

document.addEventListener("DOMContentLoaded", async () => {
    campaignList = document.getElementById("campaign-list");
    const searchForm = document.getElementById("searchForm");
    const searchInput = document.getElementById("searchInput");

    // Ensure required DOM elements are present
    if (!searchForm || !searchInput || !campaignList) {
        console.error("Required DOM elements are missing.");
        return;
    }

    const urlParams = new URLSearchParams(window.location.search);
    const category = urlParams.get("category") || "";
    const query = urlParams.get("search") || "";
    const targetAmount = urlParams.get("target_amount") || "";
    const amountRaised = urlParams.get("amount_raised") || "";

    fetchCampaigns(query, category, targetAmount, amountRaised);

    // Set the initial value and listener on the category dropdown
    categoryFilter();

    document.getElementById("target-amount").addEventListener("change", updateFilters);
    document.getElementById("amount-raised").addEventListener("change", updateFilters);
    document.getElementById("category").addEventListener("change", updateFilters);
    // Event listener for search form submission
    searchForm.addEventListener("submit", (e) => {
        e.preventDefault();
        const query = searchInput.value.trim();

        // Instead of replacing the URL, merge the search query with existing filter parameters.
        const currentParams = new URLSearchParams(window.location.search);
        if (query) {
            currentParams.set("search", query);
        } else {
            currentParams.delete("search");
        }
        // Preserve existing filters (category, target_amount, amount_raised) in currentParams
        history.pushState({}, "", `?${currentParams.toString()}`);

        // Retrieve the up-to-date filter values from the merged URL parameters
        const newCategory = currentParams.get("category") || "";
        const newTargetAmount = currentParams.get("target_amount") || "";
        const newAmountRaised = currentParams.get("amount_raised") || "";

        // Fetch campaigns with search AND the active filters.
        fetchCampaigns(query, newCategory, newTargetAmount, newAmountRaised);
    });

    // Handle back/forward navigation to maintain filter state
    window.addEventListener('popstate', () => {
        const queryParams = new URLSearchParams(window.location.search);
        const query = queryParams.get('search') || '';
        const category = queryParams.get('category') || '';
        const targetAmount = queryParams.get("target_amount") || "";
        const amountRaised = queryParams.get("amount_raised") || "";
        searchInput.value = query; // Update the search input field
        fetchCampaigns(query, category, targetAmount, amountRaised);
    });
});

async function fetchCampaigns(query = "", category = "", targetAmount = "", amountRaised = "") {
    const params = new URLSearchParams();
    if (query) params.set("search", query);
    if (category) params.set("category", category);
    if (targetAmount) params.set("target_amount", targetAmount);
    if (amountRaised) params.set("amount_raised", amountRaised);

    const url = `http://localhost:8080/campaigns?${params.toString()}`;
    console.log("Fetching campaigns with URL:", url);

    // If a category is selected, update the UI to show the filter form
    if (category !== "") {
        hideElements();
        showFilterForm();

    }

    try {
        const response = await fetch(url);
        if (!response.ok) {
            console.error('Failed to fetch campaigns', response.status);
            return;
        }

        let campaigns = await response.json();
        if (campaigns === null) {
            campaigns = [];
        }

        campaignList.innerHTML = ""; // Clear the campaign list
        let mediaContent;
        if (campaigns.length === 0) {
            campaignList.innerHTML = "<p>No campaigns found.</p>";
        } else {


            campaigns.forEach(campaign => {
                const progress = (campaign.amount_raised / campaign.target_amount) * 100;
                const formattedGoal = new Intl.NumberFormat().format(campaign.target_amount);
                const formattedRaised = new Intl.NumberFormat().format(campaign.amount_raised);
                const campaignDiv = document.createElement("div");
                campaignDiv.classList.add("campaign", "campaign-card");

                campaignDiv.innerHTML = `
                  <div class="campaign-media">
                            <img src="${campaign.media_path}" alt="Campaign Image">
                  </div>
                 <div class="campaign-content">
                          <h3>${campaign.title}</h3>
                           <p class="category">${campaign.category}</p>
                  
                           <p class="organizer"><strong>By:</strong> ${campaign.user_id}</p>
                   <p><strong>Goal:</strong> $${formattedGoal}</p>
                  <p>Amount raised: $${formattedRaised}</p>
        <div class="progress-bar">
            <div class="progress" style="width: ${progress}%;"></div>
        </div>
        <div class="actions">
             <a href="/static/donation.html?id=${campaign.campaign_id}" class="btn donate-btn">Donate</a>
        </div>
    </div>
                `;

                campaignDiv.addEventListener("click", (e) => {
                    if (!e.target.closest(".donate-btn")) {
                        window.location.href = `/static/campaign-details.html?id=${campaign.campaign_id}`;
                    }
                });

                campaignList.appendChild(campaignDiv);
            });
        }
    } catch (error) {
        console.error("Error fetching campaigns:", error);
    }
}

function hideElements() {
    document.querySelectorAll(".hide-when-filter").forEach(element => {
        element.style.display = "none";
    });
}

function showFilterForm() {
    const filterContainer = document.getElementById("filter-container");
    if (filterContainer) {
        filterContainer.style.display = "block";
    }
}

function checkUserLoginStatus() {
    const token = localStorage.getItem("token");
    const campaignMenu = document.getElementById("campaign-menu");

    if (!token) {
        campaignMenu.style.display = "none";
        return;
    }

    const decoded = jwt_decode(token);
    const currentTime = Date.now() / 1000; // Current time in seconds

    if (decoded.exp < currentTime) {
        campaignMenu.style.display = "none";
    } else {
        campaignMenu.style.display = "block";
    }
}

function categoryFilter() {
    const categoryDropdown = document.getElementById("category");
    if (!categoryDropdown) return;

    // Read the current category from the URL
    const urlParams = new URLSearchParams(window.location.search);
    const selectedCategory = urlParams.get("category") || "";
    if (selectedCategory) {
        categoryDropdown.value = selectedCategory;
    }

    // When the dropdown value changes, update the URL and reload
    categoryDropdown.addEventListener("change", function () {
        const newCategory = this.value;
        const params = new URLSearchParams(window.location.search);
        if (newCategory) {
            params.set("category", newCategory);
        } else {
            params.delete("category");
        }
        // Reload page with updated filters
        window.location.search = params.toString();
    });

    // Optionally, highlight any links that match the selected category
    const links = document.querySelectorAll("li a");
    links.forEach(link => {
        if (link.href.includes(`category=${selectedCategory}`)) {
            link.style.fontWeight = "bold";
        }
    });
}

document.getElementById("toggle-filters").addEventListener("click", function () {
    const filterSection = document.getElementById("more-filters");
    if (filterSection.style.display === "none" || filterSection.style.display === "") {
        filterSection.style.display = "block";
        this.textContent = "More Filters ▲";
    } else {
        filterSection.style.display = "none";
        this.textContent = "More Filters ▼";
    }
});

// Update filters when target amount or amount raised changes
async function updateFilters() {
    const category = document.getElementById("category").value;
    const targetAmount = document.getElementById("target-amount").value;
    const amountRaised = document.getElementById("amount-raised").value;
    const query = document.getElementById("searchInput").value.trim();
    const params = new URLSearchParams(window.location.search);


    if (query) params.set("search", query);
    else params.delete("search");

    if (category) params.set("category", category);
    else params.delete("category");

    if (targetAmount) params.set("target_amount", targetAmount);
    else params.delete("target_amount");

    if (amountRaised) params.set("amount_raised", amountRaised);
    else params.delete("amount_raised");

    window.history.pushState({}, "", `?${params.toString()}`);
    await fetchCampaigns(query, category, targetAmount, amountRaised);
}

document.getElementById("target-amount").addEventListener("change", updateFilters);
document.getElementById("amount-raised").addEventListener("change", updateFilters);
document.getElementById("category").addEventListener("change", updateFilters);
document.getElementById("searchForm").addEventListener("submit", (e) => {
    e.preventDefault();
    updateFilters();
});
window.onload = checkUserLoginStatus;
