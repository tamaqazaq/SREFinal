const token = localStorage.getItem("token"); 

async function fetchUsers() {
    try {
        let response = await fetch('/admin/users', {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        let data = await response.json();
        console.log("Users:", data);
        document.getElementById('users').innerHTML = data.map(user =>
            `<tr><td>${user.user_id}</td><td>${user.name}</td><td>${user.email}</td>
            <td class="actions">
                <button class="delete-btn" onclick="deleteUser(${user.user_id})"><i class="fas fa-trash-alt"></i> Delete</button>
            </td></tr>`
        ).join('');
    } catch (error) {
        console.error("Error fetching users:", error);
    }
}

document.addEventListener("DOMContentLoaded", async function () {
    try {
        let response = await fetch("/api/get-user-role");  // Example endpoint
        let user = await response.json();

        if (user.role !== "admin") {
            alert("Access Denied: Admins Only");
            window.location.href = "/login";
        }

        // Load user data if admin
        let usersResponse = await fetch("/admin/users");
        let users = await usersResponse.json();
        let userList = document.getElementById("user-list");
        users.forEach(user => {
            let userDiv = document.createElement("div");
            userDiv.innerText = `User: ${user.name}, Role: ${user.role}`;
            userList.appendChild(userDiv);
        });
    } catch (error) {
        console.error("Error:", error);
    }
});


async function deleteUser(userId) {
    try {
        await fetch(`/admin/users/${userId}`, {
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        fetchUsers();
    } catch (error) {
        console.error("Error deleting user:", error);
    }
}

async function fetchCampaigns() {
    try {
        let response = await fetch('/admin/campaigns', {
            headers: { 'Authorization': `Bearer ${token}` }
        });
        let data = await response.json();
        console.log("Campaigns:", data);
        document.getElementById('campaigns').innerHTML = data.map(campaign =>
            `<tr><td>${campaign.campaign_id}</td><td>${campaign.title}</td><td id="status-${campaign.campaign_id}">${campaign.status}</td>
            <td class="actions">
                <button class="delete-btn" onclick="deleteCampaign(${campaign.campaign_id})"><i class="fas fa-trash-alt"></i> Delete</button>
                <button class="edit-btn" onclick="editCampaign(${campaign.campaign_id})"><i class="fas fa-pencil-alt"></i> Edit</button>
            </td></tr>`
        ).join('');
    } catch (error) {
        console.error("Error fetching campaigns:", error);
    }
}

async function deleteCampaign(campaignId) {
    try {
        await fetch(`/admin/campaigns/${campaignId}`, {
            method: 'DELETE',
            headers: { 'Authorization': `Bearer ${token}` }
        });
        fetchCampaigns(); 
    } catch (error) {
        console.error("Error deleting campaign:", error);
    }
}

async function updateStatus(campaignId, status) {
    try {
        await fetch(`/admin/campaigns/${campaignId}/status`, {
            method: 'PUT',
            body: JSON.stringify({ status }),
            headers: { 
                'Content-Type': 'application/json',
                'Authorization': `Bearer ${token}` 
            }
        });
        document.getElementById(`status-${campaignId}`).innerText = status; 
    } catch (error) {
        console.error("Error updating campaign status:", error);
    }
}

function editCampaign(campaignId) {
    const statusCell = document.getElementById(`status-${campaignId}`);
    const currentStatus = statusCell.innerText;
    statusCell.innerHTML = `
        <select onchange="updateStatus(${campaignId}, this.value)">
            <option value="active" ${currentStatus === 'active' ? 'selected' : ''}>Active</option>
            <option value="completed" ${currentStatus === 'completed' ? 'selected' : ''}>Completed</option>
            <option value="inactive" ${currentStatus === 'inactive' ? 'selected' : ''}>Inactive</option>
        </select>
    `;
}
window.onload = () => {
    fetchUsers();
    fetchCampaigns();
};

document.addEventListener("DOMContentLoaded", async () => {
    const token = localStorage.getItem("token");

    try {
        let response = await fetch('/admin/dashboard', {
            headers: {
                'Authorization': `Bearer ${token}`
            }
        });

        if (!response.ok) {
            throw new Error('Network response was not ok');
        }

        let data = await response.json();

        document.getElementById('totalCampaigns').innerHTML = `<h2>Total Campaigns: ${data.total_campaigns}</h2>`;

        document.getElementById('totalDonations').innerHTML = `<h2>Total Donated: $${data.total_donations}</h2>`;

        document.getElementById('totalUsers').innerHTML = `<h2>Total Users: ${data.total_users}</h2>`;

        document.getElementById('campaignStats').innerHTML = `
            <p>Active: ${data.campaign_stats.active_campaigns}</p>
            <p>Inactive: ${data.campaign_stats.inactive_campaigns}</p>
            <p>Completed: ${data.campaign_stats.completed_campaigns}</p>`;

        let topCampaignsHtml = '<ul>';
        data.top_donated_campaigns.forEach(campaign => {
            topCampaignsHtml += `<li>${campaign.title}: $${campaign.total_donated}</li>`;
        });
        topCampaignsHtml += '</ul>';
        document.getElementById('topDonatedCampaigns').innerHTML = topCampaignsHtml;

        let topDonorsHtml = '<ul>';
        data.top_donors.forEach(donor => {
            topDonorsHtml += `<li>${donor.name}: $${donor.total_donated}</li>`;
        });
        topDonorsHtml += '</ul>';
        document.getElementById('topDonors').innerHTML = topDonorsHtml;

    } catch (error) {
        console.error('Error fetching dashboard data:', error);
        document.getElementById('totalCampaigns').innerHTML = `<p>Error loading data</p>`;
    }
});
